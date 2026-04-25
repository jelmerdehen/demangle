// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package scala3 decodes Scala 3 TASTy JVM name mangling.
//
// Scala 3 (Dotty) encodes structural information into JVM class and
// method names using `$`-separated segments. The encoding is defined in
// dotty/tools/dotc/core/NameKinds.scala. Key patterns:
//
//   - Package objects:       `Foo$package$`     → `Foo (package object)`
//   - Top-level objects:     `Foo$`             → `Foo (object)`
//   - Inner / nested:        `Outer$$Inner`     → `Outer.Inner`
//   - Anonymous classes:     `$$anon$1`         → `<anon #1>`
//   - Anonymous functions:   `$$anonfun$1`      → `<anonfun #1>`
//   - Local methods:         `foo$1`            → `foo#1`
//   - Lifted lambdas:        `$lzy$foo`         → `(lazy) foo`
//   - Inline forwarders:     `foo$default$1`    → `foo (default #1)`
//   - Specialized variants:  `foo$mcI$sp`       → `foo (spec mcI)`
//   - Super accessors:       `foo$super$Bar`    → `foo (super Bar)`
//   - Protected accessors:   `foo$access$0`     → `foo (access #0)`
//   - Trait setters:         `foo$_$eq`         → `foo_= (setter)`
//   - Initializers:          `foo$init$`        → `foo.<init>`
//   - Hex-encoded chars:     `$u0041`           → `A`
//   - Scala 2 compat ops:    `$plus`, `$eq`, …  handled as fallback
//
// The scheme is heuristic (MangleFidelity = BestEffort) because the
// encoding is lossy: structural information is encoded but package
// namespaces come from the dotted class name, not the mangled suffix.
package scala3

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/jelmerdehen/demangle"
)

// Kind constants for Scala 3 AST nodes.
const (
	KindObject  int32 = iota + 1 // top-level object or companion object
	KindClass                    // class, trait, or case class
	KindMethod                   // method or function
	KindLambda                   // anonymous function / lambda
	KindAnon                     // anonymous class
	KindPackage                  // package object
	KindAccessor                 // accessor / setter
	KindInit                     // initializer ($init$)
)

// Scheme implements demangle.Scheme for Scala 3 TASTy mangled names.
type Scheme struct{}

var info = demangle.Info{
	Name:           "scala3",
	Family:         "java",
	Version:        "scala-3.x",
	Description:    "Scala 3 TASTy JVM name mangling (NameKinds).",
	Stability:      demangle.Stable,
	MangleFidelity: demangle.BestEffort,
	Negatives: []demangle.Negative{
		// Itanium C++ starts with _Z
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
		// Swift stable starts with _$s
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		// Rust v0 starts with _R
		{Kind: demangle.NegContains, Pattern: "_R", Penalty: 100},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 8 * 1024,
	KindNames: map[int32]string{
		KindObject:   "Object",
		KindClass:    "Class",
		KindMethod:   "Method",
		KindLambda:   "Lambda",
		KindAnon:     "Anon",
		KindPackage:  "Package",
		KindAccessor: "Accessor",
		KindInit:     "Init",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindObject:   demangle.KindCatModule,
		KindClass:    demangle.KindCatType,
		KindMethod:   demangle.KindCatMethod,
		KindLambda:   demangle.KindCatClosure,
		KindAnon:     demangle.KindCatType,
		KindPackage:  demangle.KindCatModule,
		KindAccessor: demangle.KindCatAccessor,
		KindInit:     demangle.KindCatConstructor,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff returns confidence ~70 for clear Scala 3 patterns and avoids
// false positives on plain Java names (which may have single `$`).
//
// Strong signals (score 75):
//   - `$package$`  — package object marker
//   - `$$anon$`    — anonymous class
//   - `$$anonfun$` — anonymous function
//   - `$lzy$`      — lazy field
//
// Medium signals (score 65):
//   - ends with `$$` or `$package` — companion / package object
//   - `$default$N` — default argument
//   - `$super$`    — super accessor
//   - `$init$`     — initializer
//   - `$mcI$sp` / `$mcD$sp` / similar — specialised method
//   - `$u` followed by 4 hex digits — unicode escape
//
// Weak signal (score 45):
//   - three or more `$`-separated segments in a dot-free Java-style name
//
// Returns false (no match) for names with no `$` at all.
func (Scheme) Sniff(s string) (int, bool) {
	if !strings.Contains(s, "$") {
		return 0, false
	}

	// Strong signals
	if strings.Contains(s, "$package$") {
		return 75, true
	}
	if strings.Contains(s, "$$anon$") || strings.Contains(s, "$$anonfun$") {
		return 75, true
	}
	if strings.Contains(s, "$lzy$") || strings.Contains(s, "$lzyINIT$") {
		return 75, true
	}
	if strings.Contains(s, "$default$") {
		return 70, true
	}
	if strings.Contains(s, "$super$") {
		return 70, true
	}
	if strings.Contains(s, "$init$") {
		return 70, true
	}
	if strings.Contains(s, "$access$") {
		return 70, true
	}
	// Specialised method: $mcX$sp where X is a type code letter or pair
	if idx := strings.Index(s, "$mc"); idx >= 0 && strings.Contains(s[idx:], "$sp") {
		return 70, true
	}
	// Unicode escape $uXXXX
	if hasUnicodeEscape(s) {
		return 70, true
	}
	// Trailing $$ (double-dollar inner class / companion)
	if strings.Contains(s, "$$") {
		return 65, true
	}
	// Trailing $package (without the closing $, seen in some bytecode tools)
	if strings.HasSuffix(s, "$package") {
		return 65, true
	}
	// Trailing $ alone — top-level object companion
	if strings.HasSuffix(s, "$") {
		return 60, true
	}
	// $extension — extension method
	if strings.HasSuffix(s, "$extension") {
		return 65, true
	}
	// $adapted — adapted forwarder
	if strings.HasSuffix(s, "$adapted") {
		return 65, true
	}
	// Trait setter: $_ $eq pattern
	if strings.Contains(s, "$_$") {
		return 65, true
	}
	// Local disambiguator: name$N where N is all digits
	if idx := strings.LastIndex(s, "$"); idx >= 0 {
		if isAllDigits(s[idx+1:]) {
			return 60, true
		}
	}
	// Multiple $-separated segments that suggest structural encoding:
	// e.g. com$example$Foo$$method$1 (>= 2 dollars, all segments ident-like)
	dollars := strings.Count(s, "$")
	if dollars >= 2 {
		parts := strings.Split(s, "$")
		allIdent := true
		for _, p := range parts {
			if p == "" {
				continue
			}
			if !isIdentLike(p) {
				allIdent = false
				break
			}
		}
		if allIdent {
			return 45, true
		}
	}

	return 0, false
}

// hasUnicodeEscape returns true if s contains a $uXXXX pattern.
func hasUnicodeEscape(s string) bool {
	for i := 0; i+6 <= len(s); i++ {
		if s[i] != '$' || s[i+1] != 'u' {
			continue
		}
		if isHex4(s[i+2 : i+6]) {
			return true
		}
	}
	return false
}

func isHex4(s string) bool {
	if len(s) != 4 {
		return false
	}
	for _, c := range s {
		if !unicode.Is(unicode.ASCII_Hex_Digit, c) {
			return false
		}
	}
	return true
}

// Demangle decodes a Scala 3 mangled name into a human-readable string.
// The algorithm:
//  1. Split on `.` to get the dotted package path; decode each segment.
//  2. Within each segment, decode `$`-based annotations in order of
//     longest-first priority (to avoid substring matches).
//  3. Decode any remaining `$uXXXX` unicode escapes.
func (Scheme) Demangle(_ context.Context, in string, _ demangle.Options) (*demangle.Result, error) {
	if !strings.Contains(in, "$") {
		return nil, demangle.WrongScheme("scala3", in)
	}

	decoded, kind, ok := decode(in)
	if !ok {
		return nil, demangle.WrongScheme("scala3", in)
	}

	tree := &demangle.Node{
		Scheme: "scala3",
		Kind:   kind,
		Text:   decoded,
	}
	return &demangle.Result{
		Scheme: "scala3",
		Input:  in,
		Output: decoded,
		Tree:   tree,
	}, nil
}

// decode is the core decoding function. Returns (humanReadable, kind, ok).
// ok is false when no Scala 3 pattern was matched.
func decode(s string) (string, int32, bool) {
	// The name may be a dotted binary name like "com/example/Foo$Bar" or
	// "com.example.Foo$Bar". We decode the last component (after the last
	// slash or dot that precedes a `$`) but preserve the package prefix.
	//
	// Strategy: split on the last `.` or `/` that is NOT inside a `$…$`
	// sequence, decode the tail, and reassemble.

	// Normalize: replace `/` with `.` for package separators.
	s = strings.ReplaceAll(s, "/", ".")

	// Split at the last `.` to get class vs member names.
	// E.g. "com.example.Foo$bar$1" → prefix="com.example", tail="Foo$bar$1"
	var prefix, tail string
	if i := strings.LastIndex(s, "."); i >= 0 {
		prefix = s[:i]
		tail = s[i+1:]
	} else {
		tail = s
	}

	decoded, kind, matched := decodeSegment(tail)
	if !matched {
		return "", 0, false
	}

	if prefix != "" {
		// Decode the prefix segments too (they may contain $-encoded chars)
		parts := strings.Split(prefix, ".")
		for i, p := range parts {
			if d, _, _ := decodeSegment(p); d != "" {
				parts[i] = d
			}
		}
		return strings.Join(parts, ".") + "." + decoded, kind, true
	}
	return decoded, kind, true
}

// decodeSegment decodes a single JVM class/method name segment (no dots).
// Returns (decoded, kind, matched). matched is true iff at least one
// Scala 3 marker was found.
func decodeSegment(s string) (string, int32, bool) {
	if !strings.Contains(s, "$") {
		return s, KindClass, false
	}

	matched := false
	kind := KindClass

	// --- Step 1: check for strong suffix / infix markers ---

	// $package$ — package object
	if strings.HasSuffix(s, "$package$") {
		base := strings.TrimSuffix(s, "$package$")
		base = decodePathBase(base)
		return fmt.Sprintf("%s (package object)", base), KindPackage, true
	}
	if strings.HasSuffix(s, "$package") {
		base := strings.TrimSuffix(s, "$package")
		base = decodePathBase(base)
		return fmt.Sprintf("%s (package object)", base), KindPackage, true
	}

	// $init$ — initializer
	if strings.Contains(s, "$init$") {
		base := strings.ReplaceAll(s, "$init$", ".<init>")
		base = decodePathBase(base)
		return base, KindInit, true
	}

	// $$anonfun$N — anonymous function (must check before $$anon$)
	if idx := strings.Index(s, "$$anonfun$"); idx >= 0 {
		rest := s[idx+len("$$anonfun$"):]
		num := extractLeadingDigits(rest)
		base := s[:idx]
		base = decodePathBase(base)
		if base != "" {
			return fmt.Sprintf("%s.<anonfun #%s>", base, num), KindLambda, true
		}
		return fmt.Sprintf("<anonfun #%s>", num), KindLambda, true
	}

	// $$anon$N — anonymous class
	if idx := strings.Index(s, "$$anon$"); idx >= 0 {
		rest := s[idx+len("$$anon$"):]
		num := extractLeadingDigits(rest)
		base := s[:idx]
		base = decodePathBase(base)
		if base != "" {
			return fmt.Sprintf("%s.<anon #%s>", base, num), KindAnon, true
		}
		return fmt.Sprintf("<anon #%s>", num), KindAnon, true
	}

	// $lzyINIT$ before $lzy$ (longer prefix wins)
	if idx := strings.Index(s, "$lzyINIT$"); idx >= 0 {
		lazyName := s[idx+len("$lzyINIT$"):]
		base := s[:idx]
		base = decodePathBase(base)
		lazyName = decodeUnicode(lazyName)
		if base != "" {
			return fmt.Sprintf("%s.(lazy init) %s", base, lazyName), KindMethod, true
		}
		return fmt.Sprintf("(lazy init) %s", lazyName), KindMethod, true
	}

	// $lzy$name — lazy field initializer
	if idx := strings.Index(s, "$lzy$"); idx >= 0 {
		lazyName := s[idx+len("$lzy$"):]
		base := s[:idx]
		base = decodePathBase(base)
		lazyName = decodeUnicode(lazyName)
		if base != "" {
			return fmt.Sprintf("%s.(lazy) %s", base, lazyName), KindMethod, true
		}
		return fmt.Sprintf("(lazy) %s", lazyName), KindMethod, true
	}

	// $super$ClassName — super accessor
	if idx := strings.Index(s, "$super$"); idx >= 0 {
		superClass := s[idx+len("$super$"):]
		base := s[:idx]
		base = decodePathBase(base)
		superClass = decodeUnicode(superClass)
		return fmt.Sprintf("%s (super %s)", base, superClass), KindAccessor, true
	}

	// $default$N — default argument forwarder
	if idx := strings.Index(s, "$default$"); idx >= 0 {
		rest := s[idx+len("$default$"):]
		num := extractLeadingDigits(rest)
		base := s[:idx]
		base = decodePathBase(base)
		return fmt.Sprintf("%s (default #%s)", base, num), KindMethod, true
	}

	// $access$N — protected accessor
	if idx := strings.Index(s, "$access$"); idx >= 0 {
		rest := s[idx+len("$access$"):]
		num := extractLeadingDigits(rest)
		base := s[:idx]
		base = decodePathBase(base)
		return fmt.Sprintf("%s (access #%s)", base, num), KindAccessor, true
	}

	// $mcXX$sp — specialised method variant
	// Pattern: $mc<typeCode>$sp
	if idx := strings.Index(s, "$mc"); idx >= 0 {
		rest := s[idx+3:]
		if j := strings.Index(rest, "$sp"); j >= 0 {
			typeCode := rest[:j]
			base := s[:idx]
			suffix := rest[j+3:]
			base = decodePathBase(base)
			result := fmt.Sprintf("%s (spec %s)", base, typeCode)
			if suffix != "" {
				result += suffix
			}
			return result, KindMethod, true
		}
	}

	// $extension — extension method marker (Scala 3)
	if strings.HasSuffix(s, "$extension") {
		base := strings.TrimSuffix(s, "$extension")
		base = decodePathBase(base)
		return fmt.Sprintf("%s (extension)", base), KindMethod, true
	}

	// $adapted — adapted forwarder
	if strings.HasSuffix(s, "$adapted") {
		base := strings.TrimSuffix(s, "$adapted")
		base = decodePathBase(base)
		return fmt.Sprintf("%s (adapted)", base), KindMethod, true
	}

	// foo$_$eq — trait setter (Scala encodes `foo_=` as `foo$_$eq`)
	// Pattern: base$_$eq
	if strings.HasSuffix(s, "$_$eq") {
		base := strings.TrimSuffix(s, "$_$eq")
		base = decodePathBase(base)
		return fmt.Sprintf("%s_= (setter)", base), KindAccessor, true
	}

	// trailing $ — companion object / top-level object
	// Avoid matching $$ suffix — that is a nested class marker.
	if strings.HasSuffix(s, "$") && !strings.HasSuffix(s, "$$") {
		base := strings.TrimSuffix(s, "$")
		// Nested: Outer$$Inner$ style — decode both parts.
		if strings.Contains(base, "$$") {
			parts := strings.SplitN(base, "$$", 2)
			outer := decodePathBase(parts[0])
			inner := decodeUnicode(parts[1])
			return fmt.Sprintf("%s.%s (object)", outer, inner), KindObject, true
		}
		base = decodePathBase(base)
		return fmt.Sprintf("%s (object)", base), KindObject, true
	}

	// $$ — inner class / nested type separator
	if strings.Contains(s, "$$") {
		parts := strings.SplitN(s, "$$", 2)
		outer := decodePathBase(parts[0])
		inner, innerKind, innerMatched := decodeSegment(parts[1])
		if !innerMatched {
			inner = decodeUnicode(parts[1])
		}
		kind = innerKind
		matched = true
		return fmt.Sprintf("%s.%s", outer, inner), kind, matched
	}

	// --- Step 2: local method disambiguation: base$N ---
	// Pattern: identifier$digits at end — local disambiguator.
	if idx := strings.LastIndex(s, "$"); idx >= 0 {
		maybeNum := s[idx+1:]
		if isAllDigits(maybeNum) {
			base := s[:idx]
			base = decodePathBase(base)
			return fmt.Sprintf("%s#%s", base, maybeNum), KindMethod, true
		}
	}

	// --- Step 3: unicode escapes ($uXXXX) must be checked before
	// the path-split to avoid splitting escape sequences.
	if hasUnicodeEscape(s) {
		return decodeUnicode(s), KindClass, true
	}

	// --- Step 4: package-path encoding Foo$bar$baz (dot-less Java name) ---
	// Only apply if all $-separated segments look like valid identifier parts.
	dollars := strings.Count(s, "$")
	if dollars >= 1 {
		parts := strings.Split(s, "$")
		allIdent := true
		for _, p := range parts {
			if p == "" {
				continue
			}
			if !isIdentLike(p) {
				allIdent = false
				break
			}
		}
		if allIdent && len(parts) >= 2 {
			var out []string
			for _, p := range parts {
				if p == "" {
					continue
				}
				out = append(out, p)
			}
			matched = true
			return strings.Join(out, "."), kind, matched
		}
	}

	return s, KindClass, false
}

// decodePathBase decodes a `$`-separated path base string.
// `com$example$Bar` → `com.example.Bar`.
// If the string contains no `$`, returns it unchanged.
// Unicode escapes (`$uXXXX`) are decoded before the `$` split so that
// `$u0041BC` → `ABC` rather than splitting into bad segments.
func decodePathBase(s string) string {
	if !strings.Contains(s, "$") {
		return s
	}
	// First decode any $uXXXX escapes in-place so the split below
	// doesn't break them.
	s = decodeUnicode(s)
	if !strings.Contains(s, "$") {
		return s
	}
	// Split remaining `$` as path separators, skipping empty segments
	// that arise from consecutive `$$`.
	parts := strings.Split(s, "$")
	var out []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return s
	}
	return strings.Join(out, ".")
}

// decodeUnicode replaces $uXXXX sequences with their Unicode characters.
func decodeUnicode(s string) string {
	if !strings.Contains(s, "$u") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		if i+6 <= len(s) && s[i] == '$' && s[i+1] == 'u' && isHex4(s[i+2:i+6]) {
			cp, err := strconv.ParseInt(s[i+2:i+6], 16, 32)
			if err == nil {
				b.WriteRune(rune(cp))
				i += 6
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// extractLeadingDigits returns the leading digit characters from s.
func extractLeadingDigits(s string) string {
	for i, c := range s {
		if c < '0' || c > '9' {
			return s[:i]
		}
	}
	return s
}

// isAllDigits returns true if s is non-empty and all bytes are ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isIdentLike returns true if s looks like a valid Java/Scala identifier segment.
// Allows letters, digits, and underscore.
func isIdentLike(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return false
		}
	}
	return true
}

func init() {
	demangle.Default.Register(Scheme{})
}
