// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// healthState tracks service liveness + request counters so the
// /healthz and /metrics handlers can report them.
type healthState struct {
	startedAt time.Time
	requests  atomic.Uint64
	errors    atomic.Uint64
	bytesIn   atomic.Uint64
	bytesOut  atomic.Uint64
	// perScheme request counts; populated inside the service wrapper.
	perScheme sync.Map // key: scheme name, value: *atomic.Uint64
}

// startHealthEndpoint spins up a plain-HTTP server on addr serving
// /healthz (liveness probe) and /metrics (Prometheus text format).
// Runs until ctx is cancelled (callers close the server on shutdown).
func startHealthEndpoint(addr string, h *healthState) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "ok\nuptime=%.0fs\nrequests=%d\nerrors=%d\n",
			time.Since(h.startedAt).Seconds(),
			h.requests.Load(), h.errors.Load())
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "# HELP demangle_uptime_seconds Seconds since service start.\n")
		fmt.Fprintf(w, "# TYPE demangle_uptime_seconds gauge\n")
		fmt.Fprintf(w, "demangle_uptime_seconds %.0f\n", time.Since(h.startedAt).Seconds())

		fmt.Fprintf(w, "# HELP demangle_requests_total Total demangle requests served.\n")
		fmt.Fprintf(w, "# TYPE demangle_requests_total counter\n")
		fmt.Fprintf(w, "demangle_requests_total %d\n", h.requests.Load())

		fmt.Fprintf(w, "# HELP demangle_errors_total Total demangle requests that returned an error.\n")
		fmt.Fprintf(w, "# TYPE demangle_errors_total counter\n")
		fmt.Fprintf(w, "demangle_errors_total %d\n", h.errors.Load())

		fmt.Fprintf(w, "# HELP demangle_bytes_in_total Total input bytes processed.\n")
		fmt.Fprintf(w, "# TYPE demangle_bytes_in_total counter\n")
		fmt.Fprintf(w, "demangle_bytes_in_total %d\n", h.bytesIn.Load())

		fmt.Fprintf(w, "# HELP demangle_bytes_out_total Total output bytes produced.\n")
		fmt.Fprintf(w, "# TYPE demangle_bytes_out_total counter\n")
		fmt.Fprintf(w, "demangle_bytes_out_total %d\n", h.bytesOut.Load())

		fmt.Fprintf(w, "# HELP demangle_goroutines Number of running goroutines.\n")
		fmt.Fprintf(w, "# TYPE demangle_goroutines gauge\n")
		fmt.Fprintf(w, "demangle_goroutines %d\n", runtime.NumGoroutine())

		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Fprintf(w, "# HELP demangle_heap_bytes Current heap allocation in bytes.\n")
		fmt.Fprintf(w, "# TYPE demangle_heap_bytes gauge\n")
		fmt.Fprintf(w, "demangle_heap_bytes %d\n", ms.HeapAlloc)

		h.perScheme.Range(func(k, v any) bool {
			fmt.Fprintf(w, "demangle_requests_by_scheme_total{scheme=%q} %d\n",
				k.(string), v.(*atomic.Uint64).Load())
			return true
		})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Ready = process up + library initialised. Both are true by
		// the time this handler is reachable.
		fmt.Fprint(w, "ok")
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		_ = srv.ListenAndServe()
	}()
	return srv
}
