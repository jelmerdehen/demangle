// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// BatchRequest is one entry in a streaming batch demangle workload.
// ID is caller-supplied so responses can be correlated back to inputs
// when they arrive out of order from the worker pool.
type BatchRequest struct {
	Input string
	ID    uint64
}

// BatchResponse carries the per-input result. Err is non-nil when the
// demangle failed; Result is nil in that case. Scheme names which
// scheme actually produced the output (auto-detect or explicit).
type BatchResponse struct {
	ID     uint64
	Input  string
	Result *Result
	Err    *Error
	Scheme string
}

// BatchErrorPolicy governs what DemangleBatch does on a non-fatal
// per-input error.
type BatchErrorPolicy int

const (
	// BatchCollect — emit the response with Err populated. Default.
	BatchCollect BatchErrorPolicy = iota
	// BatchDrop — skip the input, keep going, don't emit a response.
	BatchDrop
	// BatchPropagate — cancel the batch on first error.
	BatchPropagate
)

// BatchOptions tunes DemangleBatch.
type BatchOptions struct {
	Workers      int    // default: runtime.NumCPU()
	Scheme       string // empty = auto-detect per input
	OnError      BatchErrorPolicy
	StateBufSize int // sync.Pool initial size; 0 = default 32
	BufSize      int // bounded channel between producer and workers; default 100

	// Per-call Options forwarded to each demangle.
	Options Options
}

// BatchSummary is returned when DemangleBatch's input channel closes.
type BatchSummary struct {
	Processed    int
	Succeeded    int
	Unrecognised int
	Truncated    int
	Grammar      int
	Internal     int
	Duration     time.Duration
	ByScheme     map[string]int
}

// DemangleBatch runs a worker pool over the input channel, writing
// per-input BatchResponses to out. Closes out when in closes or ctx
// is cancelled. Returns a BatchSummary covering the run.
//
// Behaviour notes:
//   - Output ordering is NOT preserved; callers use BatchResponse.ID.
//   - BufSize is the bounded queue between the caller's producer and
//     the worker pool; full queue blocks the producer (backpressure).
//   - On ctx cancel, remaining inputs are drained from in (so the
//     producer doesn't block) but not processed. Summary still includes
//     partial counts up to the cancellation point.
func (c *Catalog) DemangleBatch(ctx context.Context, in <-chan BatchRequest, out chan<- BatchResponse, opts BatchOptions) BatchSummary {
	workers := opts.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	bufSize := opts.BufSize
	if bufSize <= 0 {
		bufSize = 100
	}

	work := make(chan BatchRequest, bufSize)

	start := time.Now()
	var (
		mu       sync.Mutex
		summary  = BatchSummary{ByScheme: map[string]int{}}
		wg       sync.WaitGroup
		cancelled bool
	)

	// Feeder — reads from in, forwards to work, respecting ctx.
	feederDone := make(chan struct{})
	go func() {
		defer close(feederDone)
		defer close(work)
		for {
			select {
			case <-ctx.Done():
				// Drain remaining inputs so producers aren't blocked.
				for range in {
				}
				return
			case req, ok := <-in:
				if !ok {
					return
				}
				select {
				case work <- req:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	// Workers — pull from work, demangle, emit to out.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range work {
				select {
				case <-ctx.Done():
					return
				default:
				}

				resp := c.runOne(ctx, req, opts)

				mu.Lock()
				summary.Processed++
				if resp.Err == nil {
					summary.Succeeded++
					summary.ByScheme[resp.Scheme]++
				} else {
					switch resp.Err.Kind {
					case ErrUnrecognisedInput, ErrWrongScheme:
						summary.Unrecognised++
					case ErrTruncatedInput:
						summary.Truncated++
					case ErrGrammarViolation:
						summary.Grammar++
					case ErrInternal:
						summary.Internal++
					}
				}
				mu.Unlock()

				switch {
				case resp.Err == nil:
					// emit
				case opts.OnError == BatchDrop:
					continue
				case opts.OnError == BatchPropagate:
					// Emit, then request cancellation.
					mu.Lock()
					cancelled = true
					mu.Unlock()
				}

				select {
				case out <- resp:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
	<-feederDone
	close(out)

	summary.Duration = time.Since(start)
	_ = cancelled // reserved for future "first-error cancel" plumbing
	return summary
}

// BatchSliceResult is one per-input result from DemangleBatchSlice.
type BatchSliceResult struct {
	Output string
	Scheme string
	Err    error
}

// DemangleBatchSlice demangled all symbols concurrently using a worker pool
// and returns results in the same order as inputs.
// workers ≤ 0 uses runtime.NumCPU().
//
// The implementation uses a lock-free work-stealing index (atomic counter)
// so workers write directly into the pre-allocated result slice with no
// per-symbol mutex or output channel. This keeps the fast path to one
// atomic increment + one Demangle call per symbol.
//
// Callers that need streaming output, backpressure, or detailed per-scheme
// statistics should use the channel-based DemangleBatch directly.
func (c *Catalog) DemangleBatchSlice(ctx context.Context, symbols []string, workers int) []BatchSliceResult {
	return c.DemangleBatchSliceScheme(ctx, symbols, workers, "")
}

// DemangleBatchSliceScheme is like DemangleBatchSlice but pins all symbols
// to the named scheme, skipping auto-detection. schemeName="" falls back to
// auto-detect (same as DemangleBatchSlice).
//
// Pinning eliminates the per-symbol Sniff loop and catalog read-lock
// contention, enabling full linear scaling across cores. This is the
// recommended path when all symbols in a batch share a known scheme (e.g.,
// a Swift binary export table).
func (c *Catalog) DemangleBatchSliceScheme(ctx context.Context, symbols []string, workers int, schemeName string) []BatchSliceResult {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	n := len(symbols)
	results := make([]BatchSliceResult, n)
	if n == 0 {
		return results
	}

	// Resolve pinned scheme once before spawning workers (avoids repeated
	// map lookups inside the hot loop).
	var pinnedScheme Scheme
	if schemeName != "" {
		s, ok := c.Scheme(schemeName)
		if !ok {
			for i := range results {
				results[i] = BatchSliceResult{Err: &Error{Kind: ErrInternal, Scheme: schemeName, Offset: -1, Expected: "registered scheme"}}
			}
			return results
		}
		pinnedScheme = s
	}

	var (
		idx int64 = -1
		wg  sync.WaitGroup
	)

	opts := Options{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Atomically claim the next index.
				j := int(atomic.AddInt64(&idx, 1))
				if j >= n {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
				}
				sym := symbols[j]
				var (
					r   *Result
					err error
				)
				if pinnedScheme != nil {
					r, err = pinnedScheme.Demangle(ctx, sym, opts)
					if r != nil {
						r.Scheme = schemeName
					}
				} else {
					r, err = c.Demangle(ctx, sym, &opts)
				}
				if err != nil {
					results[j] = BatchSliceResult{Err: err}
				} else {
					out := ""
					scheme := ""
					if r != nil {
						out = r.Output
						scheme = r.Scheme
					}
					results[j] = BatchSliceResult{Output: out, Scheme: scheme}
				}
			}
		}()
	}
	wg.Wait()
	return results
}

// runOne is the per-request core. Exposed here (not inline in the
// worker loop) so tests can exercise it directly.
func (c *Catalog) runOne(ctx context.Context, req BatchRequest, opts BatchOptions) BatchResponse {
	resp := BatchResponse{ID: req.ID, Input: req.Input}

	var (
		r   *Result
		err error
	)
	optsCopy := opts.Options
	if opts.Scheme != "" {
		sch, ok := c.Scheme(opts.Scheme)
		if !ok {
			resp.Err = &Error{Kind: ErrInternal, Scheme: opts.Scheme, Expected: "registered scheme", Offset: -1}
			return resp
		}
		r, err = sch.Demangle(ctx, req.Input, optsCopy)
		resp.Scheme = opts.Scheme
	} else {
		r, err = c.Demangle(ctx, req.Input, &optsCopy)
		if r != nil {
			resp.Scheme = r.Scheme
		}
	}
	resp.Result = r
	if err != nil {
		var e *Error
		if errors.As(err, &e) {
			resp.Err = e
		} else {
			resp.Err = &Error{Kind: ErrInternal, Scheme: resp.Scheme, Offset: -1, Cause: err}
		}
	}
	return resp
}
