// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command swiftclose-digest reads production-divergences.txt and emits digest.md
// (≤100 lines) with current parity/round-trip stats, top-20 mismatch categories,
// last-10 git commits, and 3 suggested next items.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	reParityHeader = regexp.MustCompile(`^# parity run (\S+) — (\d+)/(\d+) pass$`)
	reRTHeader     = regexp.MustCompile(`^# \[round-trip\] run (\S+) — (\d+)/(\d+) pass$`)
	reWant         = regexp.MustCompile(`want="([^"]+)"`)
)

type runStat struct {
	ts    string
	pass  int
	total int
}

type catItem struct {
	name  string
	count int
}

func main() {
	divFile := locateDivFile()

	f, err := os.Open(divFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open divergences: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	lastParity, lastRT, lastParityLines := parseDivergences(f)

	// Fallback: if divergence file lacks the round-trip header (parity-only
	// runs don't write it), pull RT count from testdata/baselines.json so
	// the digest doesn't show "0/0". Use parity total as denominator
	// (same corpus, RT is a subset of parity-attempted symbols).
	if lastRT.total == 0 {
		if rt, ok := readBaselineRT(); ok {
			if rt.total == 0 && lastParity.total > 0 {
				rt.total = lastParity.total
			}
			lastRT = rt
		}
	}

	catCount := map[string]int{}
	errCount := 0
	for _, line := range lastParityLines {
		switch {
		case strings.HasPrefix(line, "[mismatch]"):
			if m := reWant.FindStringSubmatch(line); m != nil {
				catCount[mismatchCategory(m[1])]++
			}
		case strings.HasPrefix(line, "[error]"):
			errCount++
		}
	}

	cats := sortedCats(catCount)
	mismatchCount := len(lastParityLines) - errCount

	parityPct := pct(lastParity.pass, lastParity.total)
	rtPct := pct(lastRT.pass, lastRT.total)

	out, err := os.Create("digest.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create digest.md: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	fmt.Fprintf(w, "# Swift Production Digest\n\n")
	fmt.Fprintf(w, "**Parity**: %.2f%% (%d/%d) — %s\n", parityPct, lastParity.pass, lastParity.total, lastParity.ts)
	fmt.Fprintf(w, "**Round-trip**: %.2f%% (%d/%d) — %s\n", rtPct, lastRT.pass, lastRT.total, lastRT.ts)
	fmt.Fprintf(w, "**Failures**: %d parse-errors + %d mismatches\n\n", errCount, mismatchCount)

	fmt.Fprintf(w, "## Top-20 Mismatch Categories\n\n")
	limit := 20
	if len(cats) < limit {
		limit = len(cats)
	}
	for i := 0; i < limit; i++ {
		fmt.Fprintf(w, "- %-42s %d\n", cats[i].name, cats[i].count)
	}

	fmt.Fprintf(w, "\n## Last 10 Commits\n\n")
	for _, c := range gitLog(10) {
		fmt.Fprintf(w, "- %s\n", c)
	}

	fmt.Fprintf(w, "\n## Suggested Next 3 Items\n\n")
	for i, s := range suggest(cats) {
		fmt.Fprintf(w, "%d. %s\n", i+1, s)
	}

	w.Flush()
	fmt.Println("digest.md written.")
}

func locateDivFile() string {
	candidates := []string{
		filepath.Join("scheme", "swift", "stable", "testdata", "production", "production-divergences.txt"),
		filepath.Join("..", "..", "scheme", "swift", "stable", "testdata", "production", "production-divergences.txt"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[0]
}

func readBaselineRT() (runStat, bool) {
	candidates := []string{
		"testdata/baselines.json",
		filepath.Join("..", "..", "..", "..", "testdata", "baselines.json"),
	}
	var data []byte
	var err error
	for _, p := range candidates {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return runStat{}, false
	}
	var b struct {
		RT    int    `json:"production_roundtrip_pass"`
		Total int    `json:"production_parity_total"`
		TS    string `json:"snapshot_ts"`
	}
	if jerr := json.Unmarshal(data, &b); jerr != nil {
		return runStat{}, false
	}
	// Total may not be in baselines; default to parity total via separate read.
	total := b.Total
	if total == 0 {
		var p struct {
			Parity int `json:"production_parity_pass"`
		}
		_ = json.Unmarshal(data, &p)
		// Approximation: assume corpus total ~63757 (current). Better: read
		// from passing-roundtrip.txt line count? RT total tracked elsewhere.
		// For now: skip if unknown so we don't show fake percentage.
		total = 0
	}
	if b.RT == 0 && total == 0 {
		return runStat{}, false
	}
	return runStat{ts: b.TS, pass: b.RT, total: total}, b.RT > 0
}

func parseDivergences(f *os.File) (lastParity, lastRT runStat, lastParityLines []string) {
	var inLastParity bool
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if m := reParityHeader.FindStringSubmatch(line); m != nil {
			pass, _ := strconv.Atoi(m[2])
			total, _ := strconv.Atoi(m[3])
			lastParity = runStat{ts: m[1], pass: pass, total: total}
			inLastParity = true
			lastParityLines = nil
			continue
		}
		if m := reRTHeader.FindStringSubmatch(line); m != nil {
			pass, _ := strconv.Atoi(m[2])
			total, _ := strconv.Atoi(m[3])
			lastRT = runStat{ts: m[1], pass: pass, total: total}
			inLastParity = false
			continue
		}
		if inLastParity {
			lastParityLines = append(lastParityLines, line)
		}
	}
	return lastParity, lastRT, lastParityLines
}

func mismatchCategory(want string) string {
	for _, p := range []string{
		"property descriptor",
		"protocol conformance descriptor",
		"method descriptor",
		"associated type descriptor",
		"protocol witness table",
		"nominal type descriptor",
		"opaque type descriptor",
		"type metadata accessor",
		"type metadata",
		"enum case",
		"ObjC resilient class stub",
		"dispatch thunk",
	} {
		if strings.HasPrefix(want, p) {
			return p
		}
	}
	if strings.HasPrefix(want, "static") {
		parts := strings.SplitN(want, " ", 3)
		if len(parts) >= 2 {
			return "static " + parts[1]
		}
		return "static"
	}
	if len(want) > 50 {
		return want[:50] + "…"
	}
	return want
}

func sortedCats(catCount map[string]int) []catItem {
	cats := make([]catItem, 0, len(catCount))
	for k, v := range catCount {
		cats = append(cats, catItem{k, v})
	}
	sort.Slice(cats, func(i, j int) bool {
		if cats[i].count != cats[j].count {
			return cats[i].count > cats[j].count
		}
		return cats[i].name < cats[j].name
	})
	return cats
}

func pct(pass, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(pass) / float64(total) * 100
}

func gitLog(n int) []string {
	out, err := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--oneline").Output()
	if err != nil {
		return []string{"(git log unavailable)"}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines
}

func suggest(cats []catItem) []string {
	// Top-3 actionable printer-parity items.
	var result []string
	for _, c := range cats {
		if len(result) >= 3 {
			break
		}
		switch c.name {
		case "property descriptor":
			result = append(result, fmt.Sprintf("P1: property descriptor fix — %d mismatches", c.count))
		case "protocol conformance descriptor":
			result = append(result, fmt.Sprintf("P2: protocol conformance descriptor — %d mismatches", c.count))
		case "method descriptor":
			result = append(result, fmt.Sprintf("P3: method descriptor — %d mismatches", c.count))
		case "enum case":
			result = append(result, fmt.Sprintf("P6: enum case — %d mismatches", c.count))
		case "protocol witness table":
			result = append(result, fmt.Sprintf("P5: protocol witness table — %d mismatches", c.count))
		case "associated type descriptor":
			result = append(result, fmt.Sprintf("P-assoctype: associated type descriptor — %d mismatches", c.count))
		case "opaque type descriptor":
			result = append(result, fmt.Sprintf("P10: opaque type descriptor — %d mismatches", c.count))
		case "nominal type descriptor":
			result = append(result, fmt.Sprintf("P8: nominal type descriptor — %d mismatches", c.count))
		case "type metadata accessor":
			result = append(result, fmt.Sprintf("P7: type metadata accessor — %d mismatches", c.count))
		case "type metadata":
			result = append(result, fmt.Sprintf("P9: type metadata — %d mismatches", c.count))
		default:
			if c.count >= 10 {
				result = append(result, fmt.Sprintf("investigate: %s — %d mismatches", c.name, c.count))
			}
		}
	}
	if len(result) == 0 {
		result = append(result, "All categories < 10 — re-triage")
	}
	return result
}
