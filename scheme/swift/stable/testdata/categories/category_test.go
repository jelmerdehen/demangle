// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package categories contains per-category fixture tests for the Swift stable
// demangler. Each .txt file in this directory targets one mismatch category from
// the production divergences report.
//
// Bootstrap-or-enforce semantics: for each <cat>.txt, the test looks for a
// `passing-<cat>.txt` snapshot file. If absent (first run), the test logs
// pass/fail informationally and writes the current pass-set as the bootstrap
// snapshot. If present, the test ENFORCES: any symbol in the snapshot that no
// longer passes is a regression (`t.Errorf`).
//
// Run with: go test ./scheme/swift/stable/testdata/categories/
package categories

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// loadCategorySnapshot reads passing-<cat>.txt and returns the set of
// symbols that previously passed. Returns (empty, false) if the file
// does not exist (bootstrap mode).
func loadCategorySnapshot(catName string) (map[string]struct{}, bool) {
	path := "passing-" + catName + ".txt"
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	out := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		out[s] = struct{}{}
	}
	return out, true
}

// writeCategorySnapshot writes the bootstrap snapshot for a category.
func writeCategorySnapshot(catName string, passed map[string]struct{}) error {
	path := "passing-" + catName + ".txt"
	keys := make([]string, 0, len(passed))
	for k := range passed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("# bootstrap snapshot — first run; subsequent runs enforce non-regression\n")
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

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
		// Skip snapshot files written by category bootstrapping.
		if strings.HasPrefix(entry.Name(), "passing-") {
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

// TestProtocolWitnessTableFixtureStrictGate verifies all protocol_witness_table.txt entries pass.
func TestProtocolWitnessTableFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "protocol_witness_table.txt")
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

// TestCategoryInitConstraint targets ufC init own-generic-where-clause symbols.
func TestCategoryInitConstraint(t *testing.T) {
	runSingleCategory(t, "init_constraint.txt")
}

// TestCategoryExtensionInExtension targets nested Foundation extension symbols.
func TestCategoryExtensionInExtension(t *testing.T) {
	runSingleCategory(t, "extension_in_extension.txt")
}

// TestCategorySelfReturnAttributed targets FormatStyle.Attributed builder-pattern symbols.
func TestCategorySelfReturnAttributed(t *testing.T) {
	runSingleCategory(t, "self_return_attributed.txt")
}

// TestSelfReturnAttributedFixtureStrictGate verifies all self_return_attributed.txt entries pass.
func TestSelfReturnAttributedFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "self_return_attributed.txt")
}

// TestExtensionInExtensionFixtureStrictGate verifies the PC fix:
// (extension in Foundation):(extension in Foundation):... double-prefix symbols.
func TestExtensionInExtensionFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "extension_in_extension.txt")
}

// TestProtocolConformanceDescriptorFixtureStrictGate verifies the PB-3 fix:
// sADRzrlMc conditional-conformance symbols must all produce the correct
// "< where A: Swift.Protocol>" prefix.
func TestProtocolConformanceDescriptorFixtureStrictGate(t *testing.T) {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()
	cases := []struct{ sym, want string }{
		{
			"_$ss18ReversedCollectionVyxGs20LazySequenceProtocolssADRzrlMc",
			"protocol conformance descriptor for < where A: Swift.LazySequenceProtocol> Swift.ReversedCollection<A> : Swift.LazySequenceProtocol in Swift",
		},
		{
			"_$ss5SliceVyxGs20LazySequenceProtocolssADRzrlMc",
			"protocol conformance descriptor for < where A: Swift.LazySequenceProtocol> Swift.Slice<A> : Swift.LazySequenceProtocol in Swift",
		},
	}
	for _, tc := range cases {
		r, err := cat.Demangle(ctx, tc.sym, nil)
		if err != nil {
			t.Errorf("Demangle(%q) error: %v", tc.sym, err)
			continue
		}
		if r.Output != tc.want {
			t.Errorf("Demangle(%q)\n  got:  %q\n  want: %q", tc.sym, r.Output, tc.want)
		}
	}
}

// runStrictGateFile reads a category fixture file and requires every entry to pass.
// Use for buckets where 100% of entries are expected to demangle correctly.
func runStrictGateFile(t *testing.T, filename string) {
	t.Helper()
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()
	f, err := os.Open(filepath.Join(".", filename))
	if err != nil {
		t.Fatalf("open %s: %v", filename, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 256*1024), 256*1024)
	lineNum, pass, fail := 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ---> ", 2)
		if len(parts) != 2 {
			continue
		}
		sym, want := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		r, err2 := cat.Demangle(ctx, sym, nil)
		if err2 != nil {
			t.Errorf("%s:%d: Demangle(%q) error: %v", filename, lineNum, sym, err2)
			fail++
			continue
		}
		if r.Output != want {
			t.Errorf("%s:%d: Demangle(%q)\n  got:  %q\n  want: %q", filename, lineNum, sym, r.Output, want)
			fail++
			continue
		}
		pass++
	}
	t.Logf("strict gate %s: %d/%d pass", filename, pass, pass+fail)
}

// TestInitConstraintFixtureStrictGate verifies all init_constraint.txt entries pass.
func TestInitConstraintFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "init_constraint.txt")
}

// TestAssociatedTypeDescriptorFixtureStrictGate verifies all associated_type_descriptor.txt entries pass.
func TestAssociatedTypeDescriptorFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "associated_type_descriptor.txt")
}

// TestPropertyDescriptorFixtureStrictGate verifies all property_descriptor.txt entries pass.
func TestPropertyDescriptorFixtureStrictGate(t *testing.T) {
	runStrictGateFile(t, "property_descriptor.txt")
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

	currentPass := map[string]struct{}{}
	snapshot, hasSnapshot := loadCategorySnapshot(catName)

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
			fail++
			continue
		}
		got := result.Output
		if got == want {
			pass++
			currentPass[sym] = struct{}{}
		} else {
			fail++
		}
	}

	// Bootstrap-or-enforce against passing-<cat>.txt snapshot.
	if !hasSnapshot {
		if err := writeCategorySnapshot(catName, currentPass); err != nil {
			t.Logf("bootstrap snapshot write %s: %v", catName, err)
		} else {
			t.Logf("category %s: bootstrap snapshot written (%d symbols)", catName, len(currentPass))
		}
		return pass, fail
	}
	// Enforce: every symbol in the snapshot must still pass.
	var missing []string
	for s := range snapshot {
		if _, ok := currentPass[s]; !ok {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		shown := missing
		if len(shown) > 5 {
			shown = shown[:5]
		}
		t.Errorf("category %s REGRESSION: %d symbol(s) in snapshot no longer pass:\n  %s",
			catName, len(missing), strings.Join(shown, "\n  "))
		if len(missing) > 5 {
			t.Errorf("  ... (%d more — see passing-%s.txt)", len(missing)-5, catName)
		}
	}
	_ = fmt.Sprintf // ensure fmt import retained in case of future use
	return pass, fail
}
