// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package categories contains per-category fixture tests for the Swift stable
// demangler. Each .txt file in this directory targets one mismatch category from
// the production divergences report.
//
// Tests are informational: they report pass/fail counts per category but do NOT
// fail the test suite when a symbol doesn't match — the point is to track
// improvement across iterations, not to enforce correctness.
//
// Run with: go test ./scheme/swift/stable/testdata/categories/
package categories

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

func TestCategoryFixtures(t *testing.T) {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}

	totalPass, totalFail := 0, 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		catName := strings.TrimSuffix(entry.Name(), ".txt")
		pass, fail := runCategoryFile(t, ctx, cat, entry.Name(), catName)
		totalPass += pass
		totalFail += fail
		pct := 0.0
		if pass+fail > 0 {
			pct = float64(pass) / float64(pass+fail) * 100
		}
		t.Logf("category %-35s pass=%d fail=%d (%.0f%%)", catName, pass, fail, pct)
	}

	t.Logf("TOTAL pass=%d fail=%d", totalPass, totalFail)
}

// TestCategoryPropertyDescriptor targets property descriptor symbols.
func TestCategoryPropertyDescriptor(t *testing.T) {
	runSingleCategory(t, "property_descriptor.txt")
}

// TestCategoryProtocolConformanceDescriptor targets protocol conformance descriptor.
func TestCategoryProtocolConformanceDescriptor(t *testing.T) {
	runSingleCategory(t, "protocol_conformance_descriptor.txt")
}

// TestCategoryMethodDescriptor targets method descriptor symbols.
func TestCategoryMethodDescriptor(t *testing.T) {
	runSingleCategory(t, "method_descriptor.txt")
}

// TestCategoryDispatchThunk targets dispatch thunk symbols.
func TestCategoryDispatchThunk(t *testing.T) {
	runSingleCategory(t, "dispatch_thunk.txt")
}

// TestCategoryEnumCase targets enum case (WC suffix) symbols.
func TestCategoryEnumCase(t *testing.T) {
	runSingleCategory(t, "enum_case.txt")
}

// TestCategoryAssociatedTypeDescriptor targets associated type descriptor symbols.
func TestCategoryAssociatedTypeDescriptor(t *testing.T) {
	runSingleCategory(t, "associated_type_descriptor.txt")
}

// TestCategoryProtocolWitnessTable targets protocol witness table symbols.
func TestCategoryProtocolWitnessTable(t *testing.T) {
	runSingleCategory(t, "protocol_witness_table.txt")
}

// TestCategoryOpaqueTypeDescriptor targets opaque type descriptor symbols.
func TestCategoryOpaqueTypeDescriptor(t *testing.T) {
	runSingleCategory(t, "opaque_type_descriptor.txt")
}

// TestCategoryNominalTypeDescriptor targets nominal type descriptor symbols.
func TestCategoryNominalTypeDescriptor(t *testing.T) {
	runSingleCategory(t, "nominal_type_descriptor.txt")
}

// TestCategoryTypeMetadataAccessor targets type metadata accessor symbols.
func TestCategoryTypeMetadataAccessor(t *testing.T) {
	runSingleCategory(t, "type_metadata_accessor.txt")
}

// TestCategoryTypeMetadata targets type metadata symbols.
func TestCategoryTypeMetadata(t *testing.T) {
	runSingleCategory(t, "type_metadata.txt")
}

func runSingleCategory(t *testing.T, filename string) {
	t.Helper()
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()
	catName := strings.TrimSuffix(filename, ".txt")
	pass, fail := runCategoryFile(t, ctx, cat, filename, catName)
	pct := 0.0
	if pass+fail > 0 {
		pct = float64(pass) / float64(pass+fail) * 100
	}
	t.Logf("category %s: pass=%d fail=%d (%.0f%%)", catName, pass, fail, pct)
}

func runCategoryFile(t *testing.T, ctx context.Context, cat *demangle.Catalog, filename, catName string) (pass, fail int) {
	t.Helper()
	fpath := filepath.Join(".", filename)
	f, err := os.Open(fpath)
	if err != nil {
		t.Errorf("open %s: %v", filename, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 256*1024), 256*1024)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ---> ", 2)
		if len(parts) != 2 {
			t.Logf("%s:%d: skip malformed line", filename, lineNum)
			continue
		}
		sym, want := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		result, err := cat.Demangle(ctx, sym, nil)
		if err != nil {
			t.Logf("%s:%d FAIL (error) %s: %v", catName, lineNum, sym, err)
			fail++
			continue
		}
		got := result.Output
		if got == want {
			pass++
		} else {
			t.Logf("%s:%d FAIL got=%q want=%q", catName, lineNum, got, want)
			fail++
		}
	}
	return pass, fail
}
