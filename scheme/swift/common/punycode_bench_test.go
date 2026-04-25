// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package common

import (
	"fmt"
	"testing"
)

// benchCorpus is a 1000-entry mix of ASCII, Unicode, CJK, and emoji
// identifiers representative of real Swift mangled names.
var benchCorpus []string

func init() {
	const total = 1000
	benchCorpus = make([]string, 0, total)

	// ~250 ASCII-only identifiers.
	for i := 0; i < 250; i++ {
		benchCorpus = append(benchCorpus, fmt.Sprintf("ident%d", i))
	}
	// ~250 Latin extended (accented chars).
	latin := []string{"café", "résumé", "naïve", "über", "façade", "piñata", "jalapeño"}
	for i := 0; len(benchCorpus) < 500; i++ {
		benchCorpus = append(benchCorpus, latin[i%len(latin)]+fmt.Sprintf("%d", i))
	}
	// ~250 CJK.
	cjk := []string{"中文", "日本語", "한국어", "数学", "物理", "化学", "天文"}
	for i := 0; len(benchCorpus) < 750; i++ {
		benchCorpus = append(benchCorpus, cjk[i%len(cjk)]+fmt.Sprintf("%d", i))
	}
	// ~250 emoji / multi-codepoint.
	emoji := []string{"\U0001F600", "\U0001F680", "\U0001F4BB", "\U0001F525", "\U0001F98A"}
	for i := 0; len(benchCorpus) < total; i++ {
		benchCorpus = append(benchCorpus, emoji[i%len(emoji)]+fmt.Sprintf("func%d", i))
	}
}

// preEncoded holds Punycode-encoded versions of benchCorpus for the
// decode benchmark.
var preEncoded []string

func init() {
	preEncoded = make([]string, len(benchCorpus))
	for i, s := range benchCorpus {
		enc, err := PunycodeEncode(s)
		if err != nil {
			panic("bench init encode: " + err.Error())
		}
		preEncoded[i] = enc
	}
}

func BenchmarkPunycodeEncode(b *testing.B) {
	b.ReportAllocs()
	n := len(benchCorpus)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PunycodeEncode(benchCorpus[i%n])
	}
}

func BenchmarkPunycodeDecode(b *testing.B) {
	b.ReportAllocs()
	n := len(preEncoded)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PunycodeDecode(preEncoded[i%n])
	}
}
