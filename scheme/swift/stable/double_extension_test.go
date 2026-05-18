// SPDX-License-Identifier: Apache-2.0 OR MIT

package stable

import "testing"

// Canonical double-extension conformance-descriptor body (the symbol
// drives plans/double-extension-grammar.md):
//
//	Foundation.Measurement<A>< where A: __C.NSDimension>
//	  .FormatStyle< where A == __C.NSUnitInformationStorage>.ByteCount
//	  : Foundation.FormatStyle in Foundation
const deCanonicalBody = "10Foundation11MeasurementVAASo11NSDimensionCRbzrlE11FormatStyleVAASo24NSUnitInformationStorageCRszrlE9ByteCountVyx__GAafAMc"

func TestScanStructuralE(t *testing.T) {
	// The two structural 'E' terminators in the canonical body sit right
	// after the "rl" of each layer's generic signature; an 'E' inside an
	// identifier body must not be matched.
	e1 := scanStructuralE(deCanonicalBody, 0)
	if e1 < 0 || deCanonicalBody[e1] != 'E' {
		t.Fatalf("scanStructuralE first: got %d", e1)
	}
	if got := deCanonicalBody[e1-2 : e1]; got != "rl" {
		t.Fatalf("first E not after gensig 'rl': preceded by %q", got)
	}
	e2 := scanStructuralE(deCanonicalBody, e1+1)
	if e2 < 0 || deCanonicalBody[e2] != 'E' {
		t.Fatalf("scanStructuralE second: got %d", e2)
	}
	if got := deCanonicalBody[e2-2 : e2]; got != "rl" {
		t.Fatalf("second E not after gensig 'rl': preceded by %q", got)
	}
	// No third structural 'E' — the tail is yx__GAafAMc.
	if e3 := scanStructuralE(deCanonicalBody, e2+1); e3 >= 0 {
		t.Fatalf("scanStructuralE third: want -1, got %d", e3)
	}
	// 'E' inside an identifier body must be skipped.
	if got := scanStructuralE("7ElementV", 0); got >= 0 {
		t.Fatalf("identifier-body E matched at %d", got)
	}
}

func TestParseNominalDecl(t *testing.T) {
	name, kind, next, ok := parseNominalDecl("11FormatStyleV9ByteCount", 0)
	if !ok || name != "FormatStyle" || kind != 'V' {
		t.Fatalf("parseNominalDecl: ok=%v name=%q kind=%c", ok, name, kind)
	}
	if next != len("11FormatStyleV") {
		t.Fatalf("parseNominalDecl next: got %d", next)
	}
	if _, _, _, ok := parseNominalDecl("yx__G", 0); ok {
		t.Fatalf("parseNominalDecl matched a non-nominal tail")
	}
}

func TestParseExtLayerModuleRef(t *testing.T) {
	mods := []string{"Foundation"}
	// 'AA' back-reference resolves substitution index 0.
	module, next, ok := parseExtLayerModuleRef("AASo11NSDimensionC", 0, mods)
	if !ok || module != "Foundation" || next != 2 {
		t.Fatalf("AA ref: ok=%v module=%q next=%d", ok, module, next)
	}
	// Literal length-prefixed module.
	module, next, ok = parseExtLayerModuleRef("10FoundationE", 0, mods)
	if !ok || module != "Foundation" || next != 12 {
		t.Fatalf("literal mod: ok=%v module=%q next=%d", ok, module, next)
	}
	// Out-of-range back-reference declines.
	if _, _, ok := parseExtLayerModuleRef("AZ", 0, mods); ok {
		t.Fatalf("out-of-range ref accepted")
	}
}
