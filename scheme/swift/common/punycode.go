// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Punycode codec — Apple's RFC-3492 variant used by the Swift $e
// embedded-ABI mangling.
//
// Key differences from standard RFC 3492:
//   - Custom alphabet: digits 0–25 → 'a'–'z', digits 26–35 → 'A'–'J'
//     (standard uses a–z then 0–9).
//   - isValidUnicodeScalar accepts 0xD800–0xD87F for non-symbol ASCII
//     remapping (used by encodePunycodeUTF8 with mapNonSymbolChars).
//
// Reference: swift/lib/Demangling/Punycode.cpp (Apache-2.0 / LLVM-exception).
package common

import (
	"errors"
	"unicode/utf8"
)

// RFC 3492 §5 parameters — identical to Apple's implementation.
const (
	pBase        = 36
	pTmin        = 1
	pTmax        = 26
	pSkew        = 38
	pDamp        = 700
	pInitialBias = 72
	pInitialN    = uint32(128)
	pDelimiter   = '_'
)

var (
	errPunycodeOverflow    = errors.New("punycode: integer overflow")
	errPunycodeInvalidChar = errors.New("punycode: invalid encoded character")
	errPunycodeInvalidCP   = errors.New("punycode: decoded code point is basic (< 0x80)")
	errPunycodeInvalidUTF8 = errors.New("punycode: input is not valid UTF-8")
	errPunycodeInvalidScalar = errors.New("punycode: invalid unicode scalar")
	errPunycodeSurrogate   = errors.New("punycode: surrogate codepoint rejected")
)

// digitValue maps a digit (0–35) to Apple's custom alphabet character.
// Digits 0–25 → 'a'–'z'; digits 26–35 → 'A'–'J'.
func digitValue(digit int) byte {
	if digit < 26 {
		return byte('a' + digit)
	}
	return byte('A' - 26 + digit)
}

// digitIndex maps a character back to its digit value, or -1.
func digitIndex(c byte) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a')
	}
	if c >= 'A' && c <= 'J' {
		return int(c-'A') + 26
	}
	return -1
}

// isValidUnicodeScalar mirrors Apple's check (accepts 0xD800–0xD87F for
// non-symbol ASCII remapping).
func isValidUnicodeScalar(s uint32) bool {
	return s < 0xD880 || (s >= 0xE000 && s <= 0x10FFFF)
}

// adaptBias is the RFC 3492 §6.1 bias adaptation function.
func adaptBias(delta, numPoints int, firstTime bool) int {
	if firstTime {
		delta /= pDamp
	} else {
		delta /= 2
	}
	delta += delta / numPoints
	k := 0
	for delta > ((pBase-pTmin)*pTmax)/2 {
		delta /= pBase - pTmin
		k += pBase
	}
	return k + ((pBase-pTmin+1)*delta)/(delta+pSkew)
}

// PunycodeDecode decodes an Apple-variant Punycode string to a UTF-8 string.
// Returns ("", nil) for empty input.
func PunycodeDecode(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	// Locate last delimiter.
	codePoints := make([]uint32, 0, len(s))
	rest := s

	lastDelim := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == pDelimiter {
			lastDelim = i
			break
		}
	}

	if lastDelim >= 0 {
		// Copy ASCII prefix up to (not including) last delimiter.
		for i := 0; i < lastDelim; i++ {
			c := s[i]
			if c > 0x7F {
				return "", errPunycodeInvalidChar
			}
			codePoints = append(codePoints, uint32(c))
		}
		rest = s[lastDelim+1:]
	}

	n := pInitialN
	i := 0
	bias := pInitialBias

	for len(rest) > 0 {
		oldi := i
		w := 1
		for k := pBase; ; k += pBase {
			if len(rest) == 0 {
				return "", errPunycodeInvalidChar
			}
			cp := rest[0]
			rest = rest[1:]
			digit := digitIndex(cp)
			if digit < 0 {
				return "", errPunycodeInvalidChar
			}
			// Overflow check: digit * w + i
			if digit > (1<<31-1-i)/w {
				return "", errPunycodeOverflow
			}
			i += digit * w
			t := k - bias
			if t < pTmin {
				t = pTmin
			} else if t > pTmax {
				t = pTmax
			}
			if digit < t {
				break
			}
			if w > (1<<31-1)/(pBase-t) {
				return "", errPunycodeOverflow
			}
			w *= pBase - t
		}
		bias = adaptBias(i-oldi, len(codePoints)+1, oldi == 0)
		// n += i / (len(codePoints) + 1)
		inc := i / (len(codePoints) + 1)
		if inc > (1<<31-1)-int(n) {
			return "", errPunycodeOverflow
		}
		n += uint32(inc)
		i %= len(codePoints) + 1
		if n < 0x80 {
			return "", errPunycodeInvalidCP
		}
		if n >= 0xD800 && n <= 0xDFFF {
			return "", errPunycodeSurrogate
		}
		// Insert n at position i.
		codePoints = append(codePoints, 0)
		copy(codePoints[i+1:], codePoints[i:])
		codePoints[i] = n
		i++
	}

	return codePointsToUTF8(codePoints)
}

// codePointsToUTF8 re-encodes a []uint32 slice to a UTF-8 string.
// Mirrors Apple's encodeToUTF8: values in 0xD800–0xD87F are mapped
// back to their ASCII originals (c -= 0xD800).
func codePointsToUTF8(cps []uint32) (string, error) {
	buf := make([]byte, 0, len(cps)*3)
	var tmp [4]byte
	for _, s := range cps {
		if !isValidUnicodeScalar(s) {
			return "", errPunycodeInvalidScalar
		}
		if s >= 0xD800 && s < 0xD880 {
			s -= 0xD800
		}
		n := utf8.EncodeRune(tmp[:], rune(s))
		buf = append(buf, tmp[:n]...)
	}
	return string(buf), nil
}

// PunycodeEncode encodes a UTF-8 string using Apple's Punycode variant.
// Returns ("", nil) for empty input.
// Pure-ASCII strings with no non-symbol characters are returned with a
// trailing delimiter (e.g. "abc" → "abc_") — matching Apple's encoder
// behaviour when b > 0.
func PunycodeEncode(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if !utf8.ValidString(s) {
		return "", errPunycodeInvalidUTF8
	}

	// Decode UTF-8 to []uint32 code points.
	codePoints := make([]uint32, 0, len(s))
	for _, r := range s {
		if r >= 0xD800 && r <= 0xDFFF {
			return "", errPunycodeSurrogate
		}
		cp := uint32(r)
		if !isValidUnicodeScalar(cp) {
			return "", errPunycodeInvalidScalar
		}
		codePoints = append(codePoints, cp)
	}

	return encodePunycode(codePoints)
}

// encodePunycode is the core encoder operating on []uint32 code points.
func encodePunycode(codePoints []uint32) (string, error) {
	out := make([]byte, 0, len(codePoints)+8)

	n := pInitialN
	delta := 0
	bias := pInitialBias

	// Copy all basic (ASCII) code points to output; count them as h and b.
	h := 0
	for _, c := range codePoints {
		if c >= 0xD800 && c <= 0xDFFF {
			return "", errPunycodeSurrogate
		}
		if c < 0x80 {
			h++
			out = append(out, byte(c))
		}
		if !isValidUnicodeScalar(c) {
			return "", errPunycodeInvalidScalar
		}
	}
	b := h
	if b > 0 {
		out = append(out, pDelimiter)
	}

	for h < len(codePoints) {
		// Find minimum code point >= n.
		m := uint32(0x10FFFF)
		for _, cp := range codePoints {
			if cp >= n && cp < m {
				m = cp
			}
		}
		// delta += (m - n) * (h + 1)
		mMinusN := int(m - n)
		if mMinusN > (1<<31-1-delta)/(h+1) {
			return "", errPunycodeOverflow
		}
		delta += mMinusN * (h + 1)
		n = m

		for _, c := range codePoints {
			if c < n {
				if delta == 1<<31-1 {
					return "", errPunycodeOverflow
				}
				delta++
			}
			if c == n {
				q := delta
				for k := pBase; ; k += pBase {
					t := k - bias
					if t < pTmin {
						t = pTmin
					} else if t > pTmax {
						t = pTmax
					}
					if q < t {
						break
					}
					out = append(out, digitValue(t+((q-t)%(pBase-t))))
					q = (q - t) / (pBase - t)
				}
				out = append(out, digitValue(q))
				bias = adaptBias(delta, h+1, h == b)
				delta = 0
				h++
			}
		}
		delta++
		n++
	}
	return string(out), nil
}
