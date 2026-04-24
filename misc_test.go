// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle_test

import (
	"io"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
)

func TestKindCategoryString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		c    demangle.KindCategory
		want string
	}{
		{demangle.KindCatFunction, "function"},
		{demangle.KindCatMethod, "method"},
		{demangle.KindCatConstructor, "constructor"},
		{demangle.KindCatDestructor, "destructor"},
		{demangle.KindCatType, "type"},
		{demangle.KindCatNamespace, "namespace"},
		{demangle.KindCatModule, "module"},
		{demangle.KindCatTemplate, "template"},
		{demangle.KindCatOperator, "operator"},
		{demangle.KindCatLiteral, "literal"},
		{demangle.KindCatThunk, "thunk"},
		{demangle.KindCatAccessor, "accessor"},
		{demangle.KindCatVariable, "variable"},
		{demangle.KindCatParameter, "parameter"},
		{demangle.KindCatVTable, "vtable"},
		{demangle.KindCatTypeInfo, "typeinfo"},
		{demangle.KindCatClosure, "closure"},
		{demangle.KindCatMacro, "macro"},
		{demangle.KindCatOther, "other"},
		{demangle.KindCategory(999), "unknown"},
	}
	for _, c := range cases {
		if got := c.c.String(); got != c.want {
			t.Errorf("KindCategory(%d).String() = %q want %q", c.c, got, c.want)
		}
	}
}

func TestStabilityString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s    demangle.Stability
		want string
	}{
		{demangle.Stable, "stable"},
		{demangle.Experimental, "experimental"},
		{demangle.Deprecated, "deprecated"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Stability(%d).String() = %q want %q", c.s, got, c.want)
		}
	}
}

func TestFidelityString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		f    demangle.Fidelity
		want string
	}{
		{demangle.Exact, "exact"},
		{demangle.Canonical, "canonical"},
		{demangle.BestEffort, "best-effort"},
		{demangle.None, "none"},
	}
	for _, c := range cases {
		if got := c.f.String(); got != c.want {
			t.Errorf("Fidelity(%d).String() = %q want %q", c.f, got, c.want)
		}
	}
}

func TestSyncContextRealReader(t *testing.T) {
	t.Parallel()
	// Wrap a callback context in SyncContext and exercise Kind/SHA256/Reader
	// on the syncContext path (Reader returns ErrUnsupported because inner
	// is a callback).
	inner := &demangle.CallbackContext{KindName: "kcb"}
	wrapped := demangle.SyncContext(inner)
	if wrapped.Kind() != "kcb" {
		t.Errorf("Kind = %q", wrapped.Kind())
	}
	if wrapped.SHA256() != "" {
		t.Errorf("SHA256 = %q want empty (callback)", wrapped.SHA256())
	}
	r, err := wrapped.Reader()
	if err == nil {
		_ = r.(io.Closer).Close()
		t.Fatal("Reader should error on callback-backed context")
	}
}

// stopVisitor refuses to descend; Enter / Leave both return errors to
// verify Walk's error-propagation paths (currently partially covered).
type stopVisitor struct {
	enterErr, leaveErr error
	entered, left      int
}

func (v *stopVisitor) Enter(n *demangle.Node) (bool, error) {
	v.entered++
	return false, v.enterErr
}
func (v *stopVisitor) Leave(n *demangle.Node) error {
	v.left++
	return v.leaveErr
}

func TestWalkNilRootIsNoOp(t *testing.T) {
	t.Parallel()
	if err := demangle.Walk(nil, &stopVisitor{}); err != nil {
		t.Fatalf("nil root: %v", err)
	}
}

func TestWalkEnterError(t *testing.T) {
	t.Parallel()
	boom := testErr("boom")
	v := &stopVisitor{enterErr: boom}
	root := &demangle.Node{Text: "r", Children: []*demangle.Node{{Text: "c"}}}
	err := demangle.Walk(root, v)
	if err != boom {
		t.Fatalf("err = %v want boom", err)
	}
	if v.entered != 1 || v.left != 0 {
		t.Fatalf("Enter should have been called once with no Leave; got enter=%d leave=%d", v.entered, v.left)
	}
}

func TestWalkLeaveError(t *testing.T) {
	t.Parallel()
	boom := testErr("boom")
	v := &stopVisitor{leaveErr: boom}
	root := &demangle.Node{Text: "r"}
	err := demangle.Walk(root, v)
	if err != boom {
		t.Fatalf("err = %v want boom", err)
	}
}

func TestWalkSkipsNilChildren(t *testing.T) {
	t.Parallel()
	root := &demangle.Node{
		Text: "r",
		Children: []*demangle.Node{
			nil,
			{Text: "c1"},
			nil,
			{Text: "c2"},
		},
	}
	var visited []string
	err := demangle.WalkFunc(root, func(n *demangle.Node) (bool, error) {
		visited = append(visited, n.Text)
		return true, nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	want := []string{"r", "c1", "c2"}
	if len(visited) != len(want) {
		t.Fatalf("visited = %v want %v", visited, want)
	}
}

type testErr string

func (e testErr) Error() string { return string(e) }

func TestStabilityStringUnknown(t *testing.T) {
	t.Parallel()
	if got := demangle.Stability(999).String(); got != "unknown" {
		t.Fatalf("unknown Stability = %q want 'unknown'", got)
	}
}

func TestFidelityStringUnknown(t *testing.T) {
	t.Parallel()
	if got := demangle.Fidelity(999).String(); got != "unknown" {
		t.Fatalf("unknown Fidelity = %q want 'unknown'", got)
	}
}

func TestErrKindStringUnknown(t *testing.T) {
	t.Parallel()
	if got := demangle.ErrKind(999).String(); got != "unknown" {
		t.Fatalf("unknown ErrKind = %q", got)
	}
}

func TestCatalogRegisterPanicsOnDuplicate(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Register should panic on duplicate name")
		}
	}()
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "dup", prefix: "A", conf: 80})
	cat.Register(prefixScheme{name: "dup", prefix: "B", conf: 80})
}

func TestAmbiguousErrorMessage(t *testing.T) {
	t.Parallel()
	// Trigger Catalog.Demangle ambiguity path and verify the wrapped
	// AmbiguousError.Error() formatter includes the candidate list.
	cat := demangle.NewCatalog()
	cat.Register(prefixScheme{name: "alpha", prefix: "P", conf: 80})
	cat.Register(prefixScheme{name: "beta", prefix: "P", conf: 79})
	_, err := cat.Demangle(nil, "Phello", nil)
	if err == nil {
		t.Fatal("expected ambiguous error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "alpha") || !strings.Contains(msg, "beta") {
		t.Fatalf("ambiguous error %q should list candidates alpha + beta", msg)
	}
}
