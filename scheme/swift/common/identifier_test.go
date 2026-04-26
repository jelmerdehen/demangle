// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import "testing"

func TestIsValidSwiftIdentifier(t *testing.T) {
	t.Parallel()

	valid := []struct {
		name  string
		input string
	}{
		// Simple ASCII identifiers.
		{"simple lowercase", "foo"},
		{"simple uppercase", "Bar"},
		{"mixed case", "myVariable"},
		{"all caps", "CONSTANT"},
		{"single letter", "x"},
		// Underscore variants.
		{"leading underscore", "_foo"},
		{"double underscore prefix", "__bar"},
		{"underscore only", "_"},
		{"underscore in middle", "foo_bar"},
		// Dollar sign.
		{"dollar prefix", "$foo"},
		{"dollar only", "$"},
		// With digits (not at start).
		{"trailing digit", "foo1"},
		{"digits in middle", "a1b2c3"},
		// Non-ASCII: Japanese.
		{"japanese hiragana", "てすと"},
		{"japanese katakana", "テスト"},
		{"japanese kanji", "変数"},
		// Non-ASCII: Latin extended.
		{"latin extended e acute", "café"},
		{"greek letters", "αβγ"},
		// Longer identifier.
		{"long ascii identifier", "thisIsAVeryLongIdentifierName"},
		// Letter-number (Nl) category.
		{"roman numeral I (Nl)", "Ⅰfoo"},
	}

	invalid := []struct {
		name  string
		input string
	}{
		// Empty string.
		{"empty string", ""},
		// Starting with a digit.
		{"starts with digit", "1foo"},
		{"starts with zero", "0x"},
		{"all digits", "123"},
		// Surrogate codepoints as CESU-8 byte sequences.
		// Go's range over a string produces utf8.RuneError for invalid UTF-8,
		// so these are rejected as invalid characters (RuneError = U+FFFD is
		// not a valid identifier start/continue character).
		// The IsValidSwiftIdentifier surrogate guard provides defence-in-depth
		// when rune values are constructed directly (e.g., from ABI byte slices).
		{"U+D800 CESU-8", string([]byte{0xed, 0xa0, 0x80})},        // U+D800
		{"U+DC00 CESU-8", string([]byte{0xed, 0xb0, 0x80})},        // U+DC00
		{"U+DFFF CESU-8", string([]byte{0xed, 0xbf, 0xbf})},        // U+DFFF
		{"U+D87F CESU-8", string([]byte{0xed, 0xa1, 0xbf})},        // U+D87F
		{"U+DBFF CESU-8", string([]byte{0xed, 0xaf, 0xbf})},        // U+DBFF
		// Surrogate after valid start.
		{"dollar then U+D800 CESU-8", "$" + string([]byte{0xed, 0xa0, 0x80})},
		// Starting with punctuation/symbols that are not '_' or '$'.
		{"starts with hash", "#foo"},
		{"starts with at sign", "@foo"},
		{"starts with hyphen", "-foo"},
		{"starts with space", " foo"},
		// Whitespace / control characters.
		{"contains space", "foo bar"},
		{"contains newline", "foo\nbar"},
		{"contains tab", "foo\tbar"},
		// Emoji (not letters, not XID_Start).
		{"emoji start", "😀foo"},
		// Symbol character.
		{"starts with plus", "+foo"},
		// Dot (not allowed by XID rules).
		{"starts with dot", ".foo"},
	}

	for _, tc := range valid {
		tc := tc
		t.Run("valid/"+tc.name, func(t *testing.T) {
			t.Parallel()
			if !IsValidSwiftIdentifier(tc.input) {
				t.Errorf("IsValidSwiftIdentifier(%q) = false, want true", tc.input)
			}
		})
	}

	for _, tc := range invalid {
		tc := tc
		t.Run("invalid/"+tc.name, func(t *testing.T) {
			t.Parallel()
			if IsValidSwiftIdentifier(tc.input) {
				t.Errorf("IsValidSwiftIdentifier(%q) = true, want false", tc.input)
			}
		})
	}
}
