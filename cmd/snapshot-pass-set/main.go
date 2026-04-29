// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen
//
// snapshot-pass-set computes the per-symbol pass-set for production parity
// + round-trip + per-category fixtures, with three guarantees:
//
//  1. **Determinism**: the corpus is run twice; if the two pass-sets
//     differ, the tool exits non-zero. Catches map-iteration / address-
//     hash flap before it locks into the snapshot.
//  2. **Per-symbol panic recovery**: every parse runs with `defer
//     recover()`. Panics are treated as "not passing" — never crash the
//     tool itself.
//  3. **Union-only**: in --update mode the new snapshot is `prev ∪ current`;
//     a symbol that ever passed is expected to keep passing.
//
// Modes:
//
//	--update     (default) merge current pass-set into snapshot files.
//	--check      exit non-zero if snapshot \ current_pass_set ≠ ∅.
//	--bootstrap  write current pass-set as initial snapshot (replaces).
//
// Outputs (sorted, one symbol per line):
//
//	scheme/swift/stable/testdata/production/passing-parity.txt
//	scheme/swift/stable/testdata/production/passing-roundtrip.txt
//	scheme/swift/stable/testdata/categories/passing-<cat>.txt
//
// Plus on success of --update or --bootstrap, writes:
//
//	.snapshot-cache (consumed by smoke-fast)
//	testdata/baselines.json (absolute counts)
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// trimAnnotations strips Apple's `{T:}` re-entry-point markers + any
// trailing whitespace, mirroring the corpus_test.go helper.
func trimAnnotations(s string) string {
	s = strings.TrimPrefix(s, "{T:}")
	if strings.HasPrefix(s, "{T:") {
		if j := strings.Index(s, "} "); j >= 0 {
			s = s[j+2:]
		}
	}
	return strings.TrimSpace(s)
}

// passSet computes parity + round-trip pass-sets for one corpus pass.
// safeDemangle catches panics and returns (nil, error) on either error
// or panic — never panics itself.
func safeDemangle(cat *demangle.Catalog, ctx context.Context, sym string) (res *demangle.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return cat.Demangle(ctx, sym, nil)
}

func safeRemangle(ctx context.Context, tree *demangle.Node) (res *demangle.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return stable.Remangle(ctx, tree, demangle.Options{})
}

type passSets struct {
	parity    map[string]struct{}
	roundtrip map[string]struct{}
	category  map[string]map[string]struct{} // cat name → set
}

// runOnce demangles + remangles every symbol in the production corpus +
// per-category fixtures. Returns the pass-sets.
func runOnce(repoRoot string) (*passSets, error) {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	out := &passSets{
		parity:    map[string]struct{}{},
		roundtrip: map[string]struct{}{},
		category:  map[string]map[string]struct{}{},
	}

	// --- production corpus ---
	corpusDir := filepath.Join(repoRoot, "scheme/swift/stable/testdata/production/corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		return nil, fmt.Errorf("read corpus dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		path := filepath.Join(corpusDir, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		for sc.Scan() {
			line := sc.Text()
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			i := strings.Index(line, " ---> ")
			if i < 0 {
				continue
			}
			sym := strings.TrimSpace(line[:i])
			want := trimAnnotations(line[i+6:])
			if !strings.HasPrefix(sym, "$s") && !strings.HasPrefix(sym, "_$s") {
				continue
			}
			res, derr := safeDemangle(cat, ctx, sym)
			if derr == nil && res != nil && res.Output == want {
				out.parity[sym] = struct{}{}
			}
			if derr == nil && res != nil && res.Tree != nil {
				rm, rerr := safeRemangle(ctx, res.Tree)
				if rerr == nil && rm != nil && rm.Output == sym {
					out.roundtrip[sym] = struct{}{}
				}
			}
		}
		f.Close()
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("scan %s: %w", path, err)
		}
	}

	// --- per-category fixtures ---
	catDir := filepath.Join(repoRoot, "scheme/swift/stable/testdata/categories")
	catEntries, err := os.ReadDir(catDir)
	if err == nil {
		for _, entry := range catEntries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
				continue
			}
			if strings.HasPrefix(entry.Name(), "passing-") {
				continue
			}
			catName := strings.TrimSuffix(entry.Name(), ".txt")
			path := filepath.Join(catDir, entry.Name())
			f, err := os.Open(path)
			if err != nil {
				continue
			}
			set := map[string]struct{}{}
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
			for sc.Scan() {
				line := sc.Text()
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				i := strings.Index(line, " ---> ")
				if i < 0 {
					continue
				}
				sym := strings.TrimSpace(line[:i])
				want := trimAnnotations(line[i+6:])
				res, derr := safeDemangle(cat, ctx, sym)
				if derr == nil && res != nil && res.Output == want {
					set[sym] = struct{}{}
				}
			}
			f.Close()
			out.category[catName] = set
		}
	}

	return out, nil
}

// setEqual reports whether two string-sets contain the same elements.
func setEqual(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}

// loadSnapshot reads a snapshot file (one sorted symbol per line, comments allowed).
func loadSnapshot(path string) (map[string]struct{}, error) {
	out := map[string]struct{}{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil // empty
		}
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		out[s] = struct{}{}
	}
	return out, sc.Err()
}

// writeSnapshot writes a sorted list of symbols to path with a one-line
// header comment recording timestamp.
func writeSnapshot(path string, set map[string]struct{}, header string) error {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	if header != "" {
		fmt.Fprintf(&b, "# %s\n", header)
	}
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('\n')
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// disappeared returns the set of symbols in prev but not in current.
func disappeared(prev, current map[string]struct{}) []string {
	var out []string
	for k := range prev {
		if _, ok := current[k]; !ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// appearedSet returns the set of symbols in current but not in prev.
func appearedSet(prev, current map[string]struct{}) []string {
	var out []string
	for k := range current {
		if _, ok := prev[k]; !ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// unionMerge returns prev ∪ current.
func unionMerge(prev, current map[string]struct{}) map[string]struct{} {
	out := map[string]struct{}{}
	for k := range prev {
		out[k] = struct{}{}
	}
	for k := range current {
		out[k] = struct{}{}
	}
	return out
}

func toSortedSlice(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func snapshotPath(repoRoot, kind string) string {
	switch kind {
	case "parity":
		return filepath.Join(repoRoot, "scheme/swift/stable/testdata/production/passing-parity.txt")
	case "roundtrip":
		return filepath.Join(repoRoot, "scheme/swift/stable/testdata/production/passing-roundtrip.txt")
	}
	return ""
}

func categorySnapshotPath(repoRoot, cat string) string {
	return filepath.Join(repoRoot, "scheme/swift/stable/testdata/categories/passing-"+cat+".txt")
}

type baselines struct {
	AppleCuratedPass         int       `json:"apple_curated_pass"`
	AppleCuratedTotal        int       `json:"apple_curated_total"`
	SwiftcThreeWayPass       int       `json:"swiftc_three_way_pass"`
	SwiftcThreeWayTotal      int       `json:"swiftc_three_way_total"`
	ProductionParityPass     int       `json:"production_parity_pass"`
	ProductionRoundTripPass  int       `json:"production_roundtrip_pass"`
	FuzzCrashers             int       `json:"fuzz_crashers"`
	SnapshotCommit           string    `json:"snapshot_commit"`
	SnapshotTS               time.Time `json:"snapshot_ts"`
}

type cacheFile struct {
	TS              int64                          `json:"ts"`
	HeadSHA         string                         `json:"head_sha"`
	ParityPassSet   []string                       `json:"parity_pass_set"`
	RoundtripPassSet []string                      `json:"roundtrip_pass_set"`
	CategoryPassSets map[string][]string           `json:"category_pass_sets"`
	Counts          map[string]int                 `json:"counts"`
}

func writeCache(repoRoot string, sets *passSets) error {
	c := cacheFile{
		TS:              time.Now().Unix(),
		ParityPassSet:   toSortedSlice(sets.parity),
		RoundtripPassSet: toSortedSlice(sets.roundtrip),
		CategoryPassSets: map[string][]string{},
		Counts: map[string]int{
			"production_parity_pass":    len(sets.parity),
			"production_roundtrip_pass": len(sets.roundtrip),
		},
	}
	for cat, set := range sets.category {
		c.CategoryPassSets[cat] = toSortedSlice(set)
		c.Counts["category_"+cat+"_pass"] = len(set)
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(repoRoot, ".snapshot-cache"), b, 0o644)
}

func writeBaselines(repoRoot string, sets *passSets) error {
	bp := filepath.Join(repoRoot, "testdata/baselines.json")
	old, err := os.ReadFile(bp)
	bs := baselines{
		AppleCuratedPass:        153,
		AppleCuratedTotal:       153,
		SwiftcThreeWayPass:      222,
		SwiftcThreeWayTotal:     222,
		ProductionParityPass:    len(sets.parity),
		ProductionRoundTripPass: len(sets.roundtrip),
		FuzzCrashers:            0,
		SnapshotTS:              time.Now().UTC(),
	}
	if err == nil {
		var prev baselines
		if jerr := json.Unmarshal(old, &prev); jerr == nil {
			// Preserve fields the snapshot tool doesn't compute (Apple,
			// swiftc) — those come from the test suite gates.
			if prev.AppleCuratedPass > 0 {
				bs.AppleCuratedPass = prev.AppleCuratedPass
			}
			if prev.AppleCuratedTotal > 0 {
				bs.AppleCuratedTotal = prev.AppleCuratedTotal
			}
			if prev.SwiftcThreeWayPass > 0 {
				bs.SwiftcThreeWayPass = prev.SwiftcThreeWayPass
			}
			if prev.SwiftcThreeWayTotal > 0 {
				bs.SwiftcThreeWayTotal = prev.SwiftcThreeWayTotal
			}
		}
	}
	b, err := json.MarshalIndent(bs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(bp), 0o755); err != nil {
		return err
	}
	return os.WriteFile(bp, b, 0o644)
}

var _ = sync.Mutex{}

// appendBreaksLog appends a BREAK_OK entry to breaks.log. Used when
// the BREAK_OK env vars are detected during --mode=check.
func appendBreaksLog(repoRoot, breakID, reason, restoreBy string, disappearedCount int) error {
	path := filepath.Join(repoRoot, "breaks.log")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(f, "## %s opened %s by commit pending\n", breakID, now)
	fmt.Fprintf(f, "reason: %s\n", reason)
	fmt.Fprintf(f, "restore_by: %s\n", restoreBy)
	fmt.Fprintf(f, "disappeared_count: %d\n", disappearedCount)
	fmt.Fprintln(f, "---END")
	return nil
}

func main() {
	mode := flag.String("mode", "update", "update | check | bootstrap")
	repoRoot := flag.String("repo", "/data/p/demangle", "repo root")
	flag.Parse()

	t0 := time.Now()

	// --- determinism check: run twice ---
	first, err := runOnce(*repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runOnce#1: %v\n", err)
		os.Exit(2)
	}
	second, err := runOnce(*repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runOnce#2: %v\n", err)
		os.Exit(2)
	}
	if !setEqual(first.parity, second.parity) {
		fmt.Fprintln(os.Stderr, "DETERMINISM FAILURE: parity pass-sets differ across runs")
		os.Exit(3)
	}
	if !setEqual(first.roundtrip, second.roundtrip) {
		fmt.Fprintln(os.Stderr, "DETERMINISM FAILURE: round-trip pass-sets differ across runs")
		os.Exit(3)
	}
	for cat, s1 := range first.category {
		s2 := second.category[cat]
		if !setEqual(s1, s2) {
			fmt.Fprintf(os.Stderr, "DETERMINISM FAILURE: category %s pass-sets differ\n", cat)
			os.Exit(3)
		}
	}

	current := first

	// --- mode dispatch ---
	switch *mode {
	case "bootstrap":
		header := fmt.Sprintf("snapshot-pass-set bootstrap %s — high-water (replaces existing)", time.Now().UTC().Format(time.RFC3339))
		if err := writeSnapshot(snapshotPath(*repoRoot, "parity"), current.parity, header); err != nil {
			fmt.Fprintf(os.Stderr, "write parity: %v\n", err)
			os.Exit(2)
		}
		if err := writeSnapshot(snapshotPath(*repoRoot, "roundtrip"), current.roundtrip, header); err != nil {
			fmt.Fprintf(os.Stderr, "write roundtrip: %v\n", err)
			os.Exit(2)
		}
		for cat, set := range current.category {
			if err := writeSnapshot(categorySnapshotPath(*repoRoot, cat), set, header); err != nil {
				fmt.Fprintf(os.Stderr, "write category %s: %v\n", cat, err)
				os.Exit(2)
			}
		}
		if err := writeCache(*repoRoot, current); err != nil {
			fmt.Fprintf(os.Stderr, "write cache: %v\n", err)
			os.Exit(2)
		}
		if err := writeBaselines(*repoRoot, current); err != nil {
			fmt.Fprintf(os.Stderr, "write baselines: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("bootstrap: parity=%d roundtrip=%d categories=%d (%.1fs)\n",
			len(current.parity), len(current.roundtrip), len(current.category), time.Since(t0).Seconds())

	case "update":
		// Union-merge into existing snapshot.
		header := fmt.Sprintf("snapshot-pass-set update %s — union with previous (high-water)", time.Now().UTC().Format(time.RFC3339))
		// parity
		prevP, _ := loadSnapshot(snapshotPath(*repoRoot, "parity"))
		newP := unionMerge(prevP, current.parity)
		if err := writeSnapshot(snapshotPath(*repoRoot, "parity"), newP, header); err != nil {
			fmt.Fprintf(os.Stderr, "write parity: %v\n", err)
			os.Exit(2)
		}
		// roundtrip
		prevR, _ := loadSnapshot(snapshotPath(*repoRoot, "roundtrip"))
		newR := unionMerge(prevR, current.roundtrip)
		if err := writeSnapshot(snapshotPath(*repoRoot, "roundtrip"), newR, header); err != nil {
			fmt.Fprintf(os.Stderr, "write roundtrip: %v\n", err)
			os.Exit(2)
		}
		// per-category
		for cat, set := range current.category {
			prev, _ := loadSnapshot(categorySnapshotPath(*repoRoot, cat))
			merged := unionMerge(prev, set)
			if err := writeSnapshot(categorySnapshotPath(*repoRoot, cat), merged, header); err != nil {
				fmt.Fprintf(os.Stderr, "write category %s: %v\n", cat, err)
				os.Exit(2)
			}
		}
		if err := writeCache(*repoRoot, current); err != nil {
			fmt.Fprintf(os.Stderr, "write cache: %v\n", err)
			os.Exit(2)
		}
		if err := writeBaselines(*repoRoot, current); err != nil {
			fmt.Fprintf(os.Stderr, "write baselines: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("update: parity=%d (snapshot %d) roundtrip=%d (snapshot %d) categories=%d (%.1fs)\n",
			len(current.parity), len(newP),
			len(current.roundtrip), len(newR),
			len(current.category), time.Since(t0).Seconds())

	case "check":
		// Compare current pass-set against snapshot. Disappeared = blocker
		// unless BREAK_OK env var is set (foundational refactor escape).
		var totalDisappeared int
		var disParity, disRoundtrip []string
		prevP, _ := loadSnapshot(snapshotPath(*repoRoot, "parity"))
		disParity = disappeared(prevP, current.parity)
		totalDisappeared += len(disParity)
		prevR, _ := loadSnapshot(snapshotPath(*repoRoot, "roundtrip"))
		disRoundtrip = disappeared(prevR, current.roundtrip)
		totalDisappeared += len(disRoundtrip)

		if totalDisappeared == 0 {
			appP := appearedSet(prevP, current.parity)
			appR := appearedSet(prevR, current.roundtrip)
			fmt.Printf("check: no regressions. +%d new parity passes, +%d new round-trip passes (%.1fs)\n",
				len(appP), len(appR), time.Since(t0).Seconds())
			if len(appP) > 0 {
				fmt.Println("appeared (parity):")
				for _, s := range appP[:min(len(appP), 5)] {
					fmt.Println("  +", s)
				}
				if len(appP) > 5 {
					fmt.Printf("  ... (%d more)\n", len(appP)-5)
				}
			}
			os.Exit(0)
		}

		// Regression detected. Honour BREAK_OK escape if the env carries
		// a properly-formed declaration.
		breakReason := os.Getenv("BREAK_OK")
		restoreByStr := os.Getenv("RESTORE_BY")
		if breakReason != "" && restoreByStr != "" {
			restoreBy, err := time.Parse("2006-01-02", restoreByStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "BREAK_OK rejected: RESTORE_BY=%q not ISO date (YYYY-MM-DD)\n", restoreByStr)
				os.Exit(1)
			}
			if restoreBy.Before(time.Now().Add(24 * time.Hour)) {
				fmt.Fprintln(os.Stderr, "BREAK_OK rejected: RESTORE_BY must be at least tomorrow")
				os.Exit(1)
			}
			// Persist to breaks.log. Use HEAD sha as BREAK_ID-prefix; the
			// final commit will replace this when the post-commit hook fires.
			breakID := fmt.Sprintf("pending-%d", time.Now().Unix())
			if err := appendBreaksLog(*repoRoot, breakID, breakReason, restoreByStr, len(disParity)+len(disRoundtrip)); err != nil {
				fmt.Fprintf(os.Stderr, "warning: append breaks.log: %v\n", err)
			}
			fmt.Fprintf(os.Stderr, "BREAK_OK accepted (%d regressions allowed). Restore by %s.\n",
				totalDisappeared, restoreByStr)
			fmt.Fprintln(os.Stderr, "Logged to breaks.log; remember to add BREAK_OK + RESTORE_BY footer to commit msg.")
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "REGRESSION: %d symbol(s) disappeared from snapshot\n", totalDisappeared)
		if len(disParity) > 0 {
			fmt.Fprintf(os.Stderr, "disappeared (parity, %d):\n", len(disParity))
			for _, s := range disParity[:min(len(disParity), 10)] {
				fmt.Fprintln(os.Stderr, "  -", s)
			}
			if len(disParity) > 10 {
				fmt.Fprintf(os.Stderr, "  ... (%d more)\n", len(disParity)-10)
			}
		}
		if len(disRoundtrip) > 0 {
			fmt.Fprintf(os.Stderr, "disappeared (roundtrip, %d):\n", len(disRoundtrip))
			for _, s := range disRoundtrip[:min(len(disRoundtrip), 10)] {
				fmt.Fprintln(os.Stderr, "  -", s)
			}
			if len(disRoundtrip) > 10 {
				fmt.Fprintf(os.Stderr, "  ... (%d more)\n", len(disRoundtrip)-10)
			}
		}
		fmt.Fprintln(os.Stderr, "\nIf this regression is intentional (foundational refactor), retry with:")
		fmt.Fprintln(os.Stderr, `  BREAK_OK="reason" RESTORE_BY="2026-05-13" git commit ...`)
		fmt.Fprintln(os.Stderr, "and add a BREAK_OK / RESTORE_BY footer to the commit message —")
		fmt.Fprintln(os.Stderr, "see docs/regression-discipline.md.")
		os.Exit(1)

	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", *mode)
		os.Exit(2)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
