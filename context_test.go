// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/jelmerdehen/demangle"
)

func TestCallbackContext(t *testing.T) {
	t.Parallel()
	ctx := &demangle.CallbackContext{
		KindName: "test",
		Meta:     map[string]string{"a": "b"},
		Fn: func(k string) (string, bool) {
			if k == "hello" {
				return "world", true
			}
			return "", false
		},
	}
	if ctx.Kind() != "test" {
		t.Fatalf("kind = %q", ctx.Kind())
	}
	if ctx.SHA256() != "" {
		t.Fatalf("callback ctx should have empty sha256, got %q", ctx.SHA256())
	}
	if v, ok := ctx.Lookup("hello"); !ok || v != "world" {
		t.Fatalf("lookup = %q, %v", v, ok)
	}
	if _, ok := ctx.Lookup("missing"); ok {
		t.Fatalf("lookup(missing) unexpectedly ok")
	}
	if ctx.Metadata()["a"] != "b" {
		t.Fatalf("metadata = %+v", ctx.Metadata())
	}
	if _, err := ctx.Reader(); err == nil {
		t.Fatalf("callback Reader should error")
	}
}

func TestCallbackContextNilFn(t *testing.T) {
	t.Parallel()
	ctx := &demangle.CallbackContext{KindName: "empty"}
	if _, ok := ctx.Lookup("x"); ok {
		t.Fatalf("nil Fn should return ok=false")
	}
}

func TestSyncContextConcurrent(t *testing.T) {
	t.Parallel()
	// Inner callback is deliberately not safe for concurrent use —
	// SyncContext must serialise.
	var counter int
	inner := &demangle.CallbackContext{
		KindName: "k",
		Meta:     map[string]string{"m": "v"},
		Fn: func(k string) (string, bool) {
			counter++ // unguarded on purpose
			return k + "!", true
		},
	}
	wrapped := demangle.SyncContext(inner)
	const N = 50
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = wrapped.Lookup("x")
			_ = wrapped.Metadata()
		}()
	}
	wg.Wait()
	if counter != N {
		t.Fatalf("counter = %d want %d (race or lost updates)", counter, N)
	}
}

func TestSyncContextNilInner(t *testing.T) {
	t.Parallel()
	if demangle.SyncContext(nil) != nil {
		t.Fatalf("SyncContext(nil) should return nil")
	}
}

func TestSyncContextMetadataIsolation(t *testing.T) {
	t.Parallel()
	inner := &demangle.CallbackContext{
		KindName: "k",
		Meta:     map[string]string{"a": "1"},
	}
	wrapped := demangle.SyncContext(inner)
	m := wrapped.Metadata()
	m["a"] = "mutated"
	if inner.Meta["a"] != "1" {
		t.Fatalf("SyncContext.Metadata leaked inner map: inner now %q", inner.Meta["a"])
	}
}

func TestRequireContextMissingReturnsErrNeedsContext(t *testing.T) {
	t.Parallel()
	_, err := demangle.RequireContext(demangle.Options{}, "proguard_map")
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrNeedsContext {
		t.Fatalf("err = %v, want ErrNeedsContext", err)
	}
	if e.Expected != "proguard_map" {
		t.Fatalf("Expected = %q", e.Expected)
	}
}

func TestRequireContextWrongKind(t *testing.T) {
	t.Parallel()
	ctx := &demangle.CallbackContext{KindName: "other"}
	_, err := demangle.RequireContext(demangle.Options{Context: ctx}, "proguard_map")
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrNeedsContext {
		t.Fatalf("err = %v, want ErrNeedsContext", err)
	}
	if e.Got != "other" {
		t.Fatalf("Got = %q want 'other'", e.Got)
	}
}

func TestRequireContextMatchReturnsCtx(t *testing.T) {
	t.Parallel()
	ctx := &demangle.CallbackContext{KindName: "proguard_map"}
	got, err := demangle.RequireContext(demangle.Options{Context: ctx}, "proguard_map")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got != ctx {
		t.Fatalf("returned ctx differs")
	}
}
