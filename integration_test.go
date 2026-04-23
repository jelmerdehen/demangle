// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/internal/testscheme"
)

// Stage 0 integration test: a trivial Scheme + Mangler registered on
// a hermetic Catalog round-trips cleanly through Demangle + Mangle +
// Detect + DemangleBatch.

func newTestCatalog(t *testing.T) *demangle.Catalog {
	t.Helper()
	cat := demangle.NewCatalog()
	cat.Register(testscheme.Scheme{})
	return cat
}

func TestStage0_DemangleRoundTrip(t *testing.T) {
	t.Parallel()
	cat := newTestCatalog(t)
	ctx := context.Background()

	r, err := cat.Demangle(ctx, "Xhello", nil)
	if err != nil {
		t.Fatalf("demangle: %v", err)
	}
	if r.Scheme != "testscheme" {
		t.Fatalf("scheme = %q want testscheme", r.Scheme)
	}
	if r.Output != "hello" {
		t.Fatalf("output = %q want hello", r.Output)
	}
	if r.Tree == nil || r.Tree.Text != "hello" {
		t.Fatalf("tree = %+v", r.Tree)
	}

	// Mangle back.
	back, err := cat.Mangle(ctx, "testscheme", r.Tree, nil)
	if err != nil {
		t.Fatalf("mangle: %v", err)
	}
	if back.Output != "Xhello" {
		t.Fatalf("remangle = %q want Xhello", back.Output)
	}
}

func TestStage0_Detect(t *testing.T) {
	t.Parallel()
	cat := newTestCatalog(t)

	cands := cat.Detect("Xtest", demangle.DetectOptions{})
	if len(cands) != 1 {
		t.Fatalf("candidates = %d want 1", len(cands))
	}
	if cands[0].Scheme != "testscheme" {
		t.Fatalf("top = %q", cands[0].Scheme)
	}
	if cands[0].Confidence != 90 {
		t.Fatalf("confidence = %d", cands[0].Confidence)
	}
}

func TestStage0_UnrecognisedInput(t *testing.T) {
	t.Parallel()
	cat := newTestCatalog(t)
	_, err := cat.Demangle(context.Background(), "no-prefix", nil)
	var e *demangle.Error
	if !errors.As(err, &e) {
		t.Fatalf("err = %v, want *Error", err)
	}
	if e.Kind != demangle.ErrUnrecognisedInput {
		t.Fatalf("kind = %v", e.Kind)
	}
}

func TestStage0_MangleNotInvertible(t *testing.T) {
	t.Parallel()
	// Register a scheme that does NOT implement Mangler.
	cat := demangle.NewCatalog()
	cat.Register(noMangler{})
	_, err := cat.Mangle(context.Background(), "nomangler", &demangle.Node{}, nil)
	var e *demangle.Error
	if !errors.As(err, &e) || e.Kind != demangle.ErrNotInvertible {
		t.Fatalf("err = %v, want ErrNotInvertible", err)
	}
}

func TestStage0_DemangleBatch(t *testing.T) {
	t.Parallel()
	cat := newTestCatalog(t)

	inputs := []string{"Xalpha", "Xbeta", "no", "Xgamma"}
	in := make(chan demangle.BatchRequest, len(inputs))
	out := make(chan demangle.BatchResponse, len(inputs))
	for i, s := range inputs {
		in <- demangle.BatchRequest{ID: uint64(i), Input: s}
	}
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var summary demangle.BatchSummary
	done := make(chan struct{})
	go func() {
		summary = cat.DemangleBatch(ctx, in, out, demangle.BatchOptions{Workers: 2})
		close(done)
	}()

	got := map[uint64]demangle.BatchResponse{}
	for resp := range out {
		got[resp.ID] = resp
	}
	<-done

	if summary.Processed != len(inputs) {
		t.Fatalf("processed = %d want %d", summary.Processed, len(inputs))
	}
	if summary.Succeeded != 3 {
		t.Fatalf("succeeded = %d want 3", summary.Succeeded)
	}
	if summary.Unrecognised != 1 {
		t.Fatalf("unrecognised = %d want 1", summary.Unrecognised)
	}
	if got[0].Result == nil || got[0].Result.Output != "alpha" {
		t.Fatalf("id 0 result = %+v", got[0].Result)
	}
	if got[2].Err == nil || got[2].Err.Kind != demangle.ErrUnrecognisedInput {
		t.Fatalf("id 2 err = %+v", got[2].Err)
	}
}

func TestStage0_WalkAndWalkFunc(t *testing.T) {
	t.Parallel()
	root := &demangle.Node{Scheme: "testscheme", Text: "root", Children: []*demangle.Node{
		{Scheme: "testscheme", Text: "child1"},
		{Scheme: "testscheme", Text: "child2", Children: []*demangle.Node{
			{Scheme: "testscheme", Text: "grandchild"},
		}},
	}}

	var visited []string
	err := demangle.WalkFunc(root, func(n *demangle.Node) (bool, error) {
		visited = append(visited, n.Text)
		return true, nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	want := []string{"root", "child1", "child2", "grandchild"}
	if len(visited) != len(want) {
		t.Fatalf("visited = %v want %v", visited, want)
	}
	for i := range want {
		if visited[i] != want[i] {
			t.Fatalf("visited[%d] = %q want %q", i, visited[i], want[i])
		}
	}
}

// noMangler is a Scheme that does NOT implement Mangler; used to
// verify Catalog.Mangle returns ErrNotInvertible cleanly.
type noMangler struct{}

func (noMangler) Info() demangle.Info {
	return demangle.Info{Name: "nomangler", Family: "test", Version: "0", MangleFidelity: demangle.None}
}
func (noMangler) Capabilities() demangle.Capabilities { return demangle.Capabilities{} }
func (noMangler) Sniff(string) (int, bool)            { return 0, false }
func (noMangler) Demangle(context.Context, string, demangle.Options) (*demangle.Result, error) {
	return nil, demangle.WrongScheme("nomangler", "")
}
