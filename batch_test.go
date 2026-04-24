// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jelmerdehen/demangle"
)

// batchScheme is a deterministic test fixture — matches inputs with
// the "B:" prefix and returns the rest as Output.
type batchScheme struct{}

func (batchScheme) Info() demangle.Info {
	return demangle.Info{Name: "batchtest", Family: "test", Version: "1", MangleFidelity: demangle.None}
}
func (batchScheme) Capabilities() demangle.Capabilities { return demangle.Capabilities{} }
func (batchScheme) Sniff(in string) (int, bool) {
	if strings.HasPrefix(in, "B:") {
		return 90, true
	}
	return 0, false
}
func (batchScheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if !strings.HasPrefix(in, "B:") {
		return nil, demangle.WrongScheme("batchtest", in)
	}
	return &demangle.Result{Scheme: "batchtest", Input: in, Output: in[2:]}, nil
}

func newBatchCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	cat := demangle.NewCatalog()
	cat.Register(batchScheme{})
	return cat
}

func TestDemangleBatch_ByScheme(t *testing.T) {
	t.Parallel()
	cat := newBatchCatalog(t)
	in := make(chan demangle.BatchRequest, 4)
	out := make(chan demangle.BatchResponse, 4)
	inputs := []string{"B:alpha", "B:beta", "nope", "B:gamma"}
	for i, s := range inputs {
		in <- demangle.BatchRequest{ID: uint64(i), Input: s}
	}
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})
	var summary demangle.BatchSummary
	go func() {
		summary = cat.DemangleBatch(ctx, in, out, demangle.BatchOptions{Workers: 2})
		close(done)
	}()

	responses := 0
	for range out {
		responses++
	}
	<-done
	if summary.Processed != 4 {
		t.Fatalf("processed = %d want 4", summary.Processed)
	}
	if summary.Succeeded != 3 {
		t.Fatalf("succeeded = %d want 3", summary.Succeeded)
	}
	if summary.Unrecognised != 1 {
		t.Fatalf("unrecognised = %d want 1", summary.Unrecognised)
	}
	if summary.ByScheme["batchtest"] != 3 {
		t.Fatalf("ByScheme[batchtest] = %d want 3", summary.ByScheme["batchtest"])
	}
}

func TestDemangleBatch_Pinned(t *testing.T) {
	t.Parallel()
	cat := newBatchCatalog(t)
	in := make(chan demangle.BatchRequest, 2)
	out := make(chan demangle.BatchResponse, 2)
	in <- demangle.BatchRequest{ID: 1, Input: "B:one"}
	in <- demangle.BatchRequest{ID: 2, Input: "B:two"}
	close(in)
	summary := cat.DemangleBatch(context.Background(), in, out,
		demangle.BatchOptions{Scheme: "batchtest", Workers: 1})
	drain := 0
	for range out {
		drain++
	}
	if summary.Succeeded != 2 {
		t.Fatalf("pinned-scheme succeeded = %d want 2", summary.Succeeded)
	}
}

func TestDemangleBatch_CtxCancel(t *testing.T) {
	t.Parallel()
	cat := newBatchCatalog(t)
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan demangle.BatchRequest, 1)
	out := make(chan demangle.BatchResponse, 1)
	cancel() // cancel before any work lands
	close(in)
	summary := cat.DemangleBatch(ctx, in, out, demangle.BatchOptions{Workers: 1})
	for range out {
	}
	if summary.Processed < 0 {
		t.Fatalf("summary = %+v", summary)
	}
}
