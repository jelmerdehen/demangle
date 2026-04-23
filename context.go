// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

// Context is what schemes consume at call time. Passed via
// Options.Context. Accessed by ops via Lookup(key), Reader(), and
// scheme-specific extension interfaces (see SymbolicResolver pattern
// in docs/writing-a-scheme.md).
//
// Implementations MUST be safe for concurrent use by multiple
// goroutines. DemangleBatch dispatches Lookup from many goroutines
// simultaneously.
//
// Metadata: the returned map is owned by the Context. Callers MUST
// NOT mutate it after the Context has been passed to any Scheme
// method. Implementations that need defensive safety can return a
// copy, but the default is "don't mutate after handoff".
type Context interface {
	Kind() string
	SHA256() string // "" for live callbacks that have no blob identity
	Metadata() map[string]string
	Lookup(key string) (value string, ok bool)
	// Reader is an optional fallback for schemes that need the raw
	// blob. Callback contexts return ErrUnsupported.
	Reader() (io.ReadCloser, error)
}

// ErrUnsupported is returned by Context implementations whose Reader
// method doesn't apply (e.g. CallbackContext).
var errUnsupported = errors.New("demangle: Reader unsupported for this Context")

// CallbackContext adapts a pure function into a Context. The Fn field
// itself must be goroutine-safe — the wrapper adds no synchronisation
// around it. Use SyncContext to wrap a non-safe inner value.
type CallbackContext struct {
	KindName string
	Meta     map[string]string // see Context.Metadata — must not mutate after handoff
	Fn       func(key string) (string, bool)
}

func (c *CallbackContext) Kind() string                { return c.KindName }
func (c *CallbackContext) SHA256() string              { return "" }
func (c *CallbackContext) Metadata() map[string]string { return c.Meta }
func (c *CallbackContext) Lookup(k string) (string, bool) {
	if c.Fn == nil {
		return "", false
	}
	return c.Fn(k)
}
func (c *CallbackContext) Reader() (io.ReadCloser, error) { return nil, errUnsupported }

// Scheme-specific Context extensions — richer signatures for typed
// keys when string is awkward.
//
// Example — Swift symbolic-reference resolver (bytes 0x01..0x0c +
// uint32 offset into runtime metadata). Lives at
// scheme/swift/common/resolver.go:
//
//   type SymbolicResolver interface {
//       demangle.Context                                          // kind + metadata + generic Lookup
//       ResolveSymbolic(ctx context.Context,
//                       tag byte, offset uint32) (*demangle.Node, error)
//   }
//
// The Swift scheme calls RequireContext(opts, "swift_symref") then
// type-asserts to SymbolicResolver for the typed call. Avoids
// stringly-encoded keys. Pattern reusable by any scheme whose context
// key is richer than string.

// SyncContext wraps a (potentially non-safe) inner Context with a
// single mutex around Lookup / Metadata / Reader. Use when the caller
// can't guarantee the inner is goroutine-safe.
func SyncContext(inner Context) Context {
	if inner == nil {
		return nil
	}
	return &syncContext{inner: inner}
}

type syncContext struct {
	mu    sync.Mutex
	inner Context
}

func (s *syncContext) Kind() string { return s.inner.Kind() }
func (s *syncContext) SHA256() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inner.SHA256()
}
func (s *syncContext) Metadata() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Defensive copy so callers can't mutate the inner map while
	// another goroutine is reading it.
	m := s.inner.Metadata()
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
func (s *syncContext) Lookup(k string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inner.Lookup(k)
}
func (s *syncContext) Reader() (io.ReadCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.inner.Reader()
}

// RequireContext is the boilerplate every context-consuming scheme
// writes. Returns the Context if non-nil and of the requested kind;
// otherwise a structured ErrNeedsContext.
func RequireContext(opts Options, kind string) (Context, error) {
	if opts.Context == nil {
		return nil, &Error{Kind: ErrNeedsContext, Expected: kind, Offset: -1}
	}
	if opts.Context.Kind() != kind {
		return nil, &Error{Kind: ErrNeedsContext, Expected: kind, Got: opts.Context.Kind(), Offset: -1}
	}
	return opts.Context, nil
}

// ContextInfo is the inspectable metadata for a stored Context in a
// ContextStore. Returned by ContextStore.List.
type ContextInfo struct {
	Kind         string
	SHA256       string
	ByteSize     int64
	UploadedTS   time.Time
	LastAccessTS time.Time
	Metadata     map[string]string
}

// ContextStore is the SQLite-backed persistence layer for
// blob-identity contexts (ProGuard maps, JS source maps). Live
// callback contexts do NOT use this store; they live in caller memory.
type ContextStore interface {
	Put(ctx context.Context, kind string, blob []byte, meta map[string]string) (sha256 string, err error)
	Get(ctx context.Context, sha256 string) (Context, error)
	List(ctx context.Context, kind string) ([]ContextInfo, error)
	Delete(ctx context.Context, sha256 string) error
	Touch(ctx context.Context, sha256 string) error // updates LastAccessTS
	Close() error
}
