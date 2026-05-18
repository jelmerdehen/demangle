// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package stable

import (
	"context"
	"strings"
	"testing"

	"github.com/jelmerdehen/demangle"
)

// TestDecodeEntitySignatureSpan exercises decodeEntitySignatureSpan, the
// pure pre-`F`/`FZ` span decoder built by plans/entity-signature-parser.md
// P1. The decoder is unit-only this fire (P2 wires it into rendering),
// so this test is its sole exerciser.
//
// Each case builds a real parser over the full mangled body (so the
// span's substitutions resolve), locates the span, and asserts the
// decoded label list plus the result/arg byte-ranges. Spans are
// trace-verified against the Apple --expand tree (kodo swift-demangle).
func TestDecodeEntitySignatureSpan(t *testing.T) {
	type want struct {
		labels  []string
		result  string   // bytes of the result type
		args    []string // bytes of each arg element
		notOK   bool
		reasonC string // substring expected in reason when notOK
	}
	cases := []struct {
		name string
		body string // full mangled body (no $s prefix)
		span string // the captured pre-terminal span, must be a substring
		want want
	}{
		{
			// Foundation.AffineTransform.scale(x:y:) -> ()
			// labels x,y; result `y` (empty tuple); arg tuple of two CGFloat.
			name: "scale-two-labels-void-result",
			body: "10Foundation15AffineTransformV5scale1x1yy12CoreGraphics7CGFloatV_AItF",
			span: "1x1yy12CoreGraphics7CGFloatV_AIt",
			want: want{
				labels: []string{"x", "y"},
				result: "y",
				args:   []string{"12CoreGraphics7CGFloatV", "AI"},
			},
		},
		{
			// (extension in Foundation):Swift.StringProtocol.canBeConverted(to:) -> Swift.Bool
			// label `to`; result `Sb`; single tuple-wrapped arg.
			name: "canBeConverted-one-label-typed-result",
			body: "Sy10FoundationE14canBeConverted2toSbSSAAE8EncodingV_tF",
			span: "2toSbSSAAE8EncodingV_t",
			want: want{
				labels: []string{"to"},
				result: "Sb",
				args:   []string{"SSAAE8EncodingV"},
			},
		},
		{
			// (extension in Foundation):Foundation._BridgedNSError.hash(into:) -> ()
			// label `into`; result `y` (empty tuple); single inout-Hasher arg.
			name: "hash-into-inout-arg-void-result",
			body: "10Foundation15_BridgedNSErrorPAAE4hash4intoys6HasherVz_tF",
			span: "4intoys6HasherVz_t",
			want: want{
				labels: []string{"into"},
				result: "y",
				args:   []string{"s6HasherVz"},
			},
		},
		{
			// Foundation.AttributedString.replaceAttributes(_:with:) -> ()
			// FirstElementMarker `_` + label `with`; result `y`; two args.
			name: "replaceAttributes-first-elem-marker",
			body: "10Foundation16AttributedStringV17replaceAttributes_4withyAA18AttributeContainerV_AGtF",
			span: "_4withyAA18AttributeContainerV_AGt",
			want: want{
				labels: []string{"", "with"},
				result: "y",
				args:   []string{"AA18AttributeContainerV", "AG"},
			},
		},
		{
			// Foundation.InflectionRule.init(morphology:) — the result
			// type is a substitution ref (`AcA`) directly followed by an
			// identifier that belongs to the *next* element. A byte-only
			// scanner cannot tell `AcA` (complete result) from
			// `AcA10MorphologyV` (ref-qualified nominal): it takes the
			// greedy reading and the decoder then reports a clean
			// failure rather than a wrong split. This is the documented
			// P1 substitution-boundary limitation — P2 resolves it with
			// the live mid-symbol substitution table.
			name: "init-subst-ref-boundary-out-of-scope",
			body: "10Foundation14InflectionRuleO10morphologyAcA10MorphologyV_tcfC",
			span: "10morphologyAcA10MorphologyV_t",
			want: want{
				notOK:   true,
				reasonC: "unparseable",
			},
		},
		{
			// DateComponents.ISO8601FormatStyle.dateSeparator(_:) — the
			// result type contains a `0`-prefixed word-substitution
			// identifier (`A0B0VADV…`) which the structural byte scanner
			// does not decode. Decoder must cleanly report a reason.
			name: "dateSeparator-word-sub-ident-out-of-scope",
			body: "10Foundation14DateComponentsV18ISO8601FormatStyleV13dateSeparatoryAeA0B0VADV0bH0OF",
			span: "yAeA0B0VADV0bH0O",
			want: want{
				notOK:   true,
				reasonC: "unparseable",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			off := strings.Index(tc.body, tc.span)
			if off < 0 {
				t.Fatalf("span %q not found in body %q", tc.span, tc.body)
			}
			p := &parser{s: tc.body, origin: "_$s" + tc.body, schemeName: "swift.stable"}
			got := p.decodeEntitySignatureSpan(off, off+len(tc.span))

			if tc.want.notOK {
				if got.ok {
					t.Fatalf("expected decode failure, got ok=%+v", got)
				}
				if tc.want.reasonC != "" && !strings.Contains(got.reason, tc.want.reasonC) {
					t.Fatalf("reason %q does not contain %q", got.reason, tc.want.reasonC)
				}
				return
			}
			if !got.ok {
				t.Fatalf("decode failed: reason=%q", got.reason)
			}

			if len(got.labels) != len(tc.want.labels) {
				t.Fatalf("label count: got %d %q want %d %q",
					len(got.labels), got.labels, len(tc.want.labels), tc.want.labels)
			}
			for i := range got.labels {
				if got.labels[i] != tc.want.labels[i] {
					t.Errorf("label[%d]: got %q want %q", i, got.labels[i], tc.want.labels[i])
				}
			}

			gotResult := tc.body[got.resultStart:got.resultEnd]
			if gotResult != tc.want.result {
				t.Errorf("result bytes: got %q want %q", gotResult, tc.want.result)
			}

			if len(got.argRanges) != len(tc.want.args) {
				t.Fatalf("arg count: got %d want %d (%q)",
					len(got.argRanges), len(tc.want.args), tc.want.args)
			}
			for i, r := range got.argRanges {
				gotArg := tc.body[r[0]:r[1]]
				if gotArg != tc.want.args[i] {
					t.Errorf("arg[%d] bytes: got %q want %q", i, gotArg, tc.want.args[i])
				}
			}

			// Cross-check: arg count == label count, an invariant of the
			// decoder's contract.
			if len(got.labels) != len(got.argRanges) {
				t.Errorf("invariant: labels=%d argRanges=%d", len(got.labels), len(got.argRanges))
			}
		})
	}
}

// TestDecodeEntitySignatureSpanRejects confirms the decoder cleanly
// reports a reason (rather than panicking or mis-decoding) for spans
// outside the P1 subset.
func TestDecodeEntitySignatureSpanRejects(t *testing.T) {
	cases := []struct {
		name string
		body string
		span string
	}{
		{"empty-span", "abc", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &parser{s: tc.body, origin: "_$s" + tc.body, schemeName: "swift.stable"}
			off := 0
			got := p.decodeEntitySignatureSpan(off, off+len(tc.span))
			if got.ok {
				t.Fatalf("expected failure for %q, got ok", tc.name)
			}
			if got.reason == "" {
				t.Fatalf("expected non-empty reason for %q", tc.name)
			}
		})
	}
}

// TestVerboseFunctionLiteralRender exercises the P2 wiring of
// decodeEntitySignatureSpan into fpVerboseFunctionText: a stdlib-host
// extension function whose parameter and result types are all literal
// (literal nominal / bare stdlib / module-ref-qualified nominal) is
// rendered in Apple's verbose `(label: type, …) -> result` form. The
// expected strings are byte-exact against kodo `xcrun swift-demangle`.
// See plans/entity-signature-parser.md P2.
func TestVerboseFunctionLiteralRender(t *testing.T) {
	cases := []struct {
		name string
		sym  string
		want string
	}{
		{
			// single labelled param, bare-stdlib result (P2 regression
			// guard — was the function-verbose-form P2/CKL slice).
			name: "single-label-bare-result",
			sym:  "_$sSy10FoundationE14canBeConverted2toSbSSAAE8EncodingV_tF",
			want: "(extension in Foundation):Swift.StringProtocol.canBeConverted(to: (extension in Foundation):Swift.String.Encoding) -> Swift.Bool",
		},
		{
			// two labelled params, module-ref-qualified optional result
			// (`AA4DataVSg` → Foundation.Data?) — the new P2 gain.
			name: "multi-label-ref-qualified-result",
			sym:  "_$sSS10FoundationE4data5using20allowLossyConversionAA4DataVSgSSAAE8EncodingV_SbtF",
			want: "(extension in Foundation):Swift.String.data(using: (extension in Foundation):Swift.String.Encoding, allowLossyConversion: Swift.Bool) -> Foundation.Data?",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Scheme{}.Demangle(context.Background(), tc.sym, demangle.Options{})
			if err != nil {
				t.Fatalf("demangle %q: %v", tc.sym, err)
			}
			if res.Output != tc.want {
				t.Errorf("output mismatch\n got: %s\nwant: %s", res.Output, tc.want)
			}
		})
	}
}

// TestVerbosePlainHostRender exercises fpVerbosePlainHostText, the P3
// renderer for plain module-qualified nominal-host functions and
// initializers (the bucket the stdlib-host `S<letter>` candidate
// detector does not catch). It re-parses the entity from offset 0 so
// the substitution table is naturally populated from the host chain,
// resolving host-chain back-refs (`AE`, `A2C`) in arg/result types.
// The expected string is byte-exact against kodo `xcrun swift-demangle`.
// See plans/entity-signature-parser.md P3.
func TestVerbosePlainHostRender(t *testing.T) {
	cases := []struct {
		name string
		sym  string
		want string
	}{
		{
			// nested module-host init; word-sub label `09dependentB0`
			// → dependentMorphology; result `AE` and arg `A2CSg` are
			// host-chain substitution back-refs (the new P3 gain).
			name: "nested-host-init-hostchain-refs",
			sym:  "_$s10Foundation10MorphologyV7PronounV7pronoun10morphology09dependentB0AESS_A2CSgtcfC",
			want: "Foundation.Morphology.Pronoun.init(pronoun: Swift.String, morphology: Foundation.Morphology, dependentMorphology: Foundation.Morphology?) -> Foundation.Morphology.Pronoun",
		},
		{
			// non-Foundation/Swift module host: the renderer is module-
			// gated and must decline, leaving the simplified label-only
			// form (the production corpus baseline for SwiftUI).
			name: "non-verbose-module-declines",
			sym:  "_$s7SwiftUI16GlassButtonStyleVyAcA0C0VcfC",
			want: "GlassButtonStyle.init(_:)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Scheme{}.Demangle(context.Background(), tc.sym, demangle.Options{})
			if err != nil {
				t.Fatalf("demangle %q: %v", tc.sym, err)
			}
			if res.Output != tc.want {
				t.Errorf("output mismatch\n got: %s\nwant: %s", res.Output, tc.want)
			}
		})
	}
}
