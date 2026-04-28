package stable_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

func TestLockedState(t *testing.T) {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	cases := []struct{ sym, want string }{
		{"_$s10Foundation11LockedStateVAAytRszlE4lockyyF", "(extension in Foundation):Foundation.LockedState<A where A == ()>.lock() -> ()"},
		{"_$s10Foundation11LockedStateVAAytRszlE6unlockyyF", "(extension in Foundation):Foundation.LockedState<A where A == ()>.unlock() -> ()"},
	}
	for _, c := range cases {
		r, err := cat.Demangle(context.Background(), c.sym, nil)
		if err != nil { t.Errorf("%s: %v", c.sym, err); continue }
		if r.Output != c.want { t.Errorf("\nGOT:  %s\nWANT: %s", r.Output, c.want) }
	}
}
