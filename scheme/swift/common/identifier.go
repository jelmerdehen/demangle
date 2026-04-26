// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import "unicode"

// IsValidSwiftIdentifier reports whether s is a valid Swift identifier.
// Rules per the Swift ABI mangling spec and Unicode.swift:
//   - Empty string → false.
//   - First rune: isXIDStart(r) or r == '_' or r == '$'.
//   - Subsequent runes: isXIDContinue(r) or r == '$'.
//   - No surrogate codepoints (U+D800–U+DFFF) anywhere.
//
// XID_Start ≈ unicode.L ∪ unicode.Nl ∪ {'_', '$'}
// XID_Continue ≈ XID_Start ∪ unicode.Mn ∪ unicode.Mc ∪ unicode.Nd ∪ unicode.Pc
func IsValidSwiftIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r >= 0xD800 && r <= 0xDFFF {
			return false
		}
		if i == 0 {
			// XID_Start: letters (L), letter-numbers (Nl), underscore, dollar.
			if r == '_' || r == '$' {
				continue
			}
			if !unicode.In(r, unicode.L, unicode.Nl) {
				return false
			}
		} else {
			// XID_Continue: XID_Start ∪ non-spacing marks (Mn) ∪ spacing
			// combining marks (Mc) ∪ decimal digits (Nd) ∪ connector
			// punctuation (Pc, which includes '_').
			if r == '_' || r == '$' {
				continue
			}
			if !unicode.In(r, unicode.L, unicode.Nl, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc) {
				return false
			}
		}
	}
	return true
}
