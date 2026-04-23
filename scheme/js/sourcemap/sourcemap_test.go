// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package sourcemap_test

import (
	"context"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/js/sourcemap"
)

// A hand-built tiny source map. The `mappings` string encodes these
// segments (line 0):
//
//	genCol=0  srcIdx=0  origLine=0 origCol=0  nameIdx=0  (hello)
//	genCol=5  srcIdx=0  origLine=1 origCol=2  nameIdx=1  (world)
const tinyMap = `{
  "version": 3,
  "file": "out.js",
  "sources": ["src/main.js"],
  "names": ["hello", "world"],
  "mappings": "AAAAA,KACEC"
}`

func newCatalog(t *testing.T) (*demangle.Catalog, demangle.Context) {
	t.Helper()
	c := demangle.NewCatalog()
	c.Register(sourcemap.Scheme{})
	store := demangle.InMemoryContextStore()
	sha, err := store.Put(context.Background(), "js_source_map", []byte(tinyMap), nil)
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	ctx, err := store.Get(context.Background(), sha)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	return c, ctx
}

func TestSourceMapLookup(t *testing.T) {
	t.Parallel()
	cat, mapCtx := newCatalog(t)

	// Query exactly at the first mapping.
	r, err := cat.Demangle(context.Background(), "0:0", &demangle.Options{Context: mapCtx})
	if err != nil {
		t.Fatalf("demangle 0:0: %v", err)
	}
	if r.Output != "hello" {
		t.Fatalf("output = %q, want %q", r.Output, "hello")
	}
	if r.Annotations["js.has_name_mapping"] != "true" {
		t.Fatalf("has_name_mapping flag wrong: %+v", r.Annotations)
	}
	if r.Annotations["js.original_source"] != "src/main.js" {
		t.Fatalf("original_source = %q", r.Annotations["js.original_source"])
	}
}

func TestSourceMapMissingContext(t *testing.T) {
	t.Parallel()
	c := demangle.NewCatalog()
	c.Register(sourcemap.Scheme{})
	_, err := c.Demangle(context.Background(), "0:5", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSourceMapBadCoord(t *testing.T) {
	t.Parallel()
	cat, mapCtx := newCatalog(t)
	_, err := cat.Demangle(context.Background(), "not-a-coord", &demangle.Options{Context: mapCtx})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSourceMapNoMapping(t *testing.T) {
	t.Parallel()
	cat, mapCtx := newCatalog(t)
	// Line 42 has no data in the tiny map.
	_, err := cat.Demangle(context.Background(), "42:0", &demangle.Options{Context: mapCtx})
	if err == nil {
		t.Fatalf("expected no-mapping error")
	}
}

func FuzzSourceMap(f *testing.F) {
	seeds := []string{"0:0", "0:5", "1:0", "42:0", "", "not-a-coord", "1,2"}
	for _, s := range seeds {
		f.Add(s)
	}
	// Build a catalog + map context inline (newCatalog needs *testing.T).
	cat := demangle.NewCatalog()
	cat.Register(sourcemap.Scheme{})
	store := demangle.InMemoryContextStore()
	sha, _ := store.Put(context.Background(), "js_source_map", []byte(tinyMap), nil)
	mapCtx, _ := store.Get(context.Background(), sha)
	f.Fuzz(func(t *testing.T, in string) {
		if len(in) > 256 {
			t.Skip()
		}
		_, _ = cat.Demangle(context.Background(), in, &demangle.Options{Context: mapCtx})
	})
}

func TestSourceMapSniffAcceptsCoord(t *testing.T) {
	t.Parallel()
	s := sourcemap.Scheme{}
	for _, c := range []struct {
		in      string
		wantHit bool
	}{
		{"0:5", true},
		{"123:456", true},
		{"1,2", true},
		{"not-a-coord", false},
		{"", false},
	} {
		c := c
		t.Run(c.in, func(t *testing.T) {
			_, ok := s.Sniff(c.in)
			if ok != c.wantHit {
				t.Fatalf("sniff = %v, want %v", ok, c.wantHit)
			}
		})
	}
}
