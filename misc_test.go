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
