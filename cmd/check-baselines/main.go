// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen
//
// check-baselines is the absolute-count aggregate ratchet enforcer.
//
// Reads `testdata/baselines.json` (committed) + `.snapshot-cache` (fresh
// from `make smoke`). Fails non-zero if any committed absolute count
// is HIGHER than the cached current count — i.e. an aggregate
// regression. Reports percentages live but never gates on them.
//
// Aggregate ratchet means: the corpus may grow, refactors may add or
// move passes, but the total count of passing symbols cannot drop.
// Per-symbol regression detection lives in `cmd/snapshot-pass-set
// --mode=check` — this tool covers the orthogonal aggregate axis.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
	TS               int64               `json:"ts"`
	HeadSHA          string              `json:"head_sha"`
	ParityPassSet    []string            `json:"parity_pass_set"`
	RoundtripPassSet []string            `json:"roundtrip_pass_set"`
	CategoryPassSets map[string][]string `json:"category_pass_sets"`
	Counts           map[string]int      `json:"counts"`
}

func main() {
	repoRoot := flag.String("repo", "/data/p/demangle", "repo root")
	flag.Parse()

	bp := filepath.Join(*repoRoot, "testdata/baselines.json")
	cp := filepath.Join(*repoRoot, ".snapshot-cache")

	bsRaw, err := os.ReadFile(bp)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("ratchet: no baselines.json yet — first run, skipping (run `make snapshot` to bootstrap)")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "read baselines.json: %v\n", err)
		os.Exit(2)
	}
	var bs baselines
	if err := json.Unmarshal(bsRaw, &bs); err != nil {
		fmt.Fprintf(os.Stderr, "parse baselines.json: %v\n", err)
		os.Exit(2)
	}

	cRaw, err := os.ReadFile(cp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ratchet: .snapshot-cache missing — run `make smoke` first to populate")
		os.Exit(2)
	}
	var c cacheFile
	if err := json.Unmarshal(cRaw, &c); err != nil {
		fmt.Fprintf(os.Stderr, "parse .snapshot-cache: %v\n", err)
		os.Exit(2)
	}

	// Cache freshness: ≤1 hour old.
	age := time.Now().Unix() - c.TS
	if age > 3600 {
		fmt.Fprintf(os.Stderr, "ratchet: .snapshot-cache stale (%d s old, max 3600) — re-run `make smoke`\n", age)
		os.Exit(2)
	}

	// Compare absolute counts. Cache always wins (cannot regress vs
	// committed baseline); cache may equal or exceed.
	type axis struct {
		name     string
		baseline int
		current  int
	}
	axes := []axis{
		{"production_parity_pass", bs.ProductionParityPass, c.Counts["production_parity_pass"]},
		{"production_roundtrip_pass", bs.ProductionRoundTripPass, c.Counts["production_roundtrip_pass"]},
	}

	var fail []string
	for _, a := range axes {
		if a.current < a.baseline {
			fail = append(fail, fmt.Sprintf("%s: baseline=%d current=%d (drop %d)", a.name, a.baseline, a.current, a.baseline-a.current))
		}
	}

	if len(fail) > 0 {
		fmt.Fprintln(os.Stderr, "RATCHET FAILURE: aggregate counts dropped vs committed baseline")
		for _, f := range fail {
			fmt.Fprintln(os.Stderr, "  -", f)
		}
		fmt.Fprintln(os.Stderr, "\nIf intentional (foundational refactor), add BREAK_OK footer to commit.")
		os.Exit(1)
	}

	// Report live percentages from cached counts (informational).
	parityPct := 100.0
	if bs.ProductionParityPass > 0 {
		parityPct = float64(c.Counts["production_parity_pass"]) / float64(bs.ProductionParityPass) * 100
	}
	rtPct := 100.0
	if bs.ProductionRoundTripPass > 0 {
		rtPct = float64(c.Counts["production_roundtrip_pass"]) / float64(bs.ProductionRoundTripPass) * 100
	}
	fmt.Printf("ratchet: parity=%d (%.1f %% of baseline %d) / roundtrip=%d (%.1f %% of baseline %d)\n",
		c.Counts["production_parity_pass"], parityPct, bs.ProductionParityPass,
		c.Counts["production_roundtrip_pass"], rtPct, bs.ProductionRoundTripPass)
}
