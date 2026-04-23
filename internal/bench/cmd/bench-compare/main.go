// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command bench-compare diffs two benchmark runs and fails on
// regressions beyond a configurable threshold.
//
// Input is go-test -bench output (stdout). Old + new files are
// parsed; per-benchmark ns/op deltas computed; non-zero exit if any
// tracked benchmark regresses more than --threshold percent.
//
// Wired into CI per the v5 plan: store testdata/baselines.bench as
// the tracked baseline; each PR runs make bench JSON=build/bench.txt
// then invokes this tool.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type sample struct {
	name   string
	nsOp   float64
	allocs float64
}

func main() {
	oldPath := flag.String("old", "testdata/baselines.bench", "baseline bench output path")
	newPath := flag.String("new", "", "new bench output path (required)")
	threshold := flag.Float64("threshold", 10.0, "percent regression that fails the run")
	flag.Parse()

	if *newPath == "" {
		fmt.Fprintln(os.Stderr, "bench-compare: --new FILE is required")
		os.Exit(2)
	}

	oldSamples, err := parseFile(*oldPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read old: %v\n", err)
		os.Exit(1)
	}
	newSamples, err := parseFile(*newPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read new: %v\n", err)
		os.Exit(1)
	}

	oldBy := map[string]sample{}
	for _, s := range oldSamples {
		oldBy[s.name] = s
	}

	failed := 0
	fmt.Printf("%-60s  %10s  %10s  %7s\n", "benchmark", "old ns/op", "new ns/op", "delta")
	for _, n := range newSamples {
		o, ok := oldBy[n.name]
		if !ok {
			fmt.Printf("%-60s  %10s  %10.0f  %7s\n", n.name, "(new)", n.nsOp, "")
			continue
		}
		deltaPct := 0.0
		if o.nsOp > 0 {
			deltaPct = 100 * (n.nsOp - o.nsOp) / o.nsOp
		}
		marker := "  ok"
		if deltaPct > *threshold {
			marker = "REGRESS"
			failed++
		} else if deltaPct < -5 {
			marker = " IMPROV"
		}
		fmt.Printf("%-60s  %10.0f  %10.0f  %6.1f%%  %s\n",
			n.name, o.nsOp, n.nsOp, deltaPct, marker)
	}
	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\n%d benchmark(s) regressed by > %.1f %%\n", failed, *threshold)
		os.Exit(1)
	}
}

func parseFile(path string) ([]sample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f)
}

func parse(r io.Reader) ([]sample, error) {
	var out []sample
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// fields[0] = BenchmarkName-N (strip -N suffix)
		name := strings.TrimRight(strings.TrimRightFunc(fields[0], isDigit), "-")
		// Find the "ns/op" column.
		var nsOp, allocs float64
		for i := 1; i < len(fields)-1; i++ {
			if fields[i+1] == "ns/op" {
				v, err := strconv.ParseFloat(fields[i], 64)
				if err == nil {
					nsOp = v
				}
			}
			if fields[i+1] == "allocs/op" {
				v, err := strconv.ParseFloat(fields[i], 64)
				if err == nil {
					allocs = v
				}
			}
		}
		if nsOp > 0 {
			out = append(out, sample{name: name, nsOp: nsOp, allocs: allocs})
		}
	}
	return out, sc.Err()
}

func isDigit(r rune) bool { return r >= '0' && r <= '9' }
