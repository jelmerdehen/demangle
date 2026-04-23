// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Source Map V3 VLQ base64 codec.
//
// Encoding: each signed integer is encoded as one or more base64
// digits (64-char alphabet "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/").
// The high bit of each digit is a continuation flag; the LSB of the
// first digit's payload is a sign bit (zig-zag). Subsequent digits
// contribute 5 more payload bits each.
package sourcemap

import "fmt"

const (
	vlqShift       = 5
	vlqMask        = 1<<5 - 1
	vlqContMask    = 1 << 5
	vlqBase64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var vlqDecodeTable [256]int8

func init() {
	for i := range vlqDecodeTable {
		vlqDecodeTable[i] = -1
	}
	for i, c := range vlqBase64Chars {
		vlqDecodeTable[c] = int8(i)
	}
}

// decodeVLQ reads one VLQ integer from s starting at i. Returns
// (value, nextIndex, err).
func decodeVLQ(s string, i int) (int, int, error) {
	var (
		result uint32
		shift  uint
		more   bool
		first  = true
		sign   uint32
		pos    = i
	)
	for {
		if pos >= len(s) {
			return 0, pos, fmt.Errorf("vlq: unexpected end at %d", pos)
		}
		v := vlqDecodeTable[s[pos]]
		if v < 0 {
			return 0, pos, fmt.Errorf("vlq: invalid base64 char %q at %d", s[pos], pos)
		}
		pos++
		more = uint8(v)&vlqContMask != 0
		payload := uint32(uint8(v) & vlqMask)
		if first {
			sign = payload & 1
			payload >>= 1
			result |= payload
			shift = 4
			first = false
		} else {
			result |= payload << shift
			shift += vlqShift
		}
		if !more {
			break
		}
	}
	v := int(result)
	if sign == 1 {
		v = -v
	}
	return v, pos, nil
}
