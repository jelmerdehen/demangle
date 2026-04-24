// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"context"
	"path/filepath"
	"testing"
)

func TestInMemoryContextStore_PutGetListDelete(t *testing.T) {
	t.Parallel()
	store := InMemoryContextStore()
	defer store.Close()
	runStoreSuite(t, store)
}

func TestSqliteContextStore_PutGetListDelete(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "contexts.db")
	store, err := OpenContextStore(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer store.Close()
	runStoreSuite(t, store)
}

func runStoreSuite(t *testing.T, store ContextStore) {
	t.Helper()
	ctx := context.Background()

	// Put
	blob := []byte("com.example.Foo -> a:\n    void bar(int) -> b\n")
	sha, err := store.Put(ctx, "proguard_map", blob, map[string]string{"app": "test"})
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	if sha == "" {
		t.Fatalf("empty sha256")
	}

	// Idempotent put
	sha2, err := store.Put(ctx, "proguard_map", blob, map[string]string{"app": "test"})
	if err != nil {
		t.Fatalf("re-put: %v", err)
	}
	if sha2 != sha {
		t.Fatalf("sha mismatch after re-put: %q vs %q", sha2, sha)
	}

	// Get
	cc, err := store.Get(ctx, sha)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if cc.Kind() != "proguard_map" {
		t.Fatalf("kind = %q", cc.Kind())
	}
	if cc.SHA256() != sha {
		t.Fatalf("sha = %q want %q", cc.SHA256(), sha)
	}
	if v, ok := cc.Metadata()["app"]; !ok || v != "test" {
		t.Fatalf("metadata app = %q ok=%v", v, ok)
	}
	if b, ok := cc.(interface{ Blob() []byte }); ok {
		if string(b.Blob()) != string(blob) {
			t.Fatalf("blob mismatch")
		}
	}
	// Lookup defaults to metadata scan; "app" is present, "bogus" isn't.
	if v, ok := cc.Lookup("app"); !ok || v != "test" {
		t.Fatalf("Lookup(app) = %q,%v", v, ok)
	}
	if _, ok := cc.Lookup("bogus"); ok {
		t.Fatalf("Lookup(bogus) should be false")
	}
	// Reader streams the blob bytes back.
	rc, err := cc.Reader()
	if err != nil {
		t.Fatalf("Reader: %v", err)
	}
	buf := make([]byte, len(blob))
	n, _ := rc.Read(buf)
	if err := rc.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if n != len(blob) || string(buf) != string(blob) {
		t.Fatalf("Read got %d bytes: %q", n, string(buf[:n]))
	}

	// List
	list, err := store.List(ctx, "proguard_map")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("list len = %d want 1", len(list))
	}
	if list[0].SHA256 != sha {
		t.Fatalf("list sha = %q", list[0].SHA256)
	}

	// Touch → last_access_ts updates
	if err := store.Touch(ctx, sha); err != nil {
		t.Fatalf("touch: %v", err)
	}

	// Delete
	if err := store.Delete(ctx, sha); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := store.Get(ctx, sha); err == nil {
		t.Fatalf("expected error after delete")
	}
}
