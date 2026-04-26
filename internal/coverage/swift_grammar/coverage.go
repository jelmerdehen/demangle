// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package swift_grammar provides coverage analysis for the Swift stable demangler.
// It categorizes parse results across a production corpus and identifies grammar
// gaps — positions where the parser consumes a valid prefix but cannot parse
// the remainder because a grammar rule is not yet implemented.
//
// The analysis is "failure-coverage": it groups every error by its normalised
// message prefix so callers can see exactly which grammar paths are never
// exercised by a fully-passing parse.
package swift_grammar

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

// Result holds the outcome of demangling a single symbol.
type Result struct {
	Symbol   string
	Fixture  string // basename of the source file
	Output   string // non-empty on success
	Err      error
	Category string // "ok", "grammar-gap", "wrong-scheme", "parse-error", "other"
}

// Report aggregates coverage statistics across an entire corpus.
type Report struct {
	// Total is the number of symbols attempted.
	Total int
	// OK is the number of symbols that demangled without error.
	OK int
	// GapByMsg maps normalised error message → count.
	GapByMsg map[string]int
	// FirstFixture maps normalised error message → basename of the first file
	// that triggered it.
	FirstFixture map[string]string
	// FirstSymbol maps normalised error message → the first symbol that
	// triggered it (useful for reproducing the failure).
	FirstSymbol map[string]string
	// LastFixture maps normalised error message → basename of the last file
	// that triggered it.
	LastFixture map[string]string
}

func newReport() *Report {
	return &Report{
		GapByMsg:     make(map[string]int),
		FirstFixture: make(map[string]string),
		FirstSymbol:  make(map[string]string),
		LastFixture:  make(map[string]string),
	}
}

// AnalyzeCorpusDir reads all *.txt files in dir, demangles each symbol, and
// returns a Report.  Lines that start with '#' or are empty are skipped.
// Lines in "symbol ---> expected" format are parsed; the symbol column is used.
func AnalyzeCorpusDir(dir string) (*Report, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, fmt.Errorf("glob %s: %w", dir, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no *.txt files found in %s", dir)
	}

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	r := newReport()
	for _, fpath := range files {
		fname := filepath.Base(fpath)
		fh, err := os.Open(fpath)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", fpath, err)
		}
		if err := analyzeFile(ctx, cat, fh, fname, r); err != nil {
			fh.Close()
			return nil, fmt.Errorf("analyze %s: %w", fname, err)
		}
		fh.Close()
	}
	return r, nil
}

// AnalyzeReader processes symbols from an already-opened reader (useful in
// tests with synthetic input).  fname is used only for attribution in the
// report.
func AnalyzeReader(ctx context.Context, cat *demangle.Catalog, rd io.Reader, fname string, r *Report) error {
	return analyzeFile(ctx, cat, rd, fname, r)
}

// NewCatalog returns a Catalog pre-loaded with the Swift stable scheme.
func NewCatalog() *demangle.Catalog {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	return cat
}

func analyzeFile(ctx context.Context, cat *demangle.Catalog, rd io.Reader, fname string, r *Report) error {
	sc := bufio.NewScanner(rd)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Support "symbol ---> expected" format as well as bare symbol lines.
		sym := line
		if i := strings.Index(line, " ---> "); i >= 0 {
			sym = strings.TrimSpace(line[:i])
		}
		if sym == "" {
			continue
		}

		r.Total++
		res, demErr := cat.Demangle(ctx, sym, nil)
		if demErr == nil {
			r.OK++
			_ = res
			continue
		}

		category := classifyError(demErr)
		msg := normalizeError(demErr, category)

		r.GapByMsg[msg]++
		if _, seen := r.FirstFixture[msg]; !seen {
			r.FirstFixture[msg] = fname
			r.FirstSymbol[msg] = sym
		}
		r.LastFixture[msg] = fname
	}
	return sc.Err()
}

// classifyError maps a demangle error to one of our category strings.
func classifyError(err error) string {
	var de *demangle.Error
	if !errors.As(err, &de) {
		return "other"
	}
	switch de.Kind {
	case demangle.ErrWrongScheme, demangle.ErrUnrecognisedInput:
		return "wrong-scheme"
	case demangle.ErrUnsupported:
		return "grammar-gap"
	case demangle.ErrGrammarViolation, demangle.ErrTruncatedInput:
		return "parse-error"
	default:
		return "other"
	}
}

// normalizeError strips position- and symbol-specific parts from error messages
// so identical grammar gaps collapse into a single bucket.
func normalizeError(err error, category string) string {
	var de *demangle.Error
	if !errors.As(err, &de) {
		// Non-structured error: use the full message.
		return truncate(err.Error(), 120)
	}

	switch category {
	case "grammar-gap":
		// "swift-stable: unsupported (expected end of input (grammar feature
		//  not yet supported), got -) at offset 42 near \"…\""
		// → keep only the Expected field; it already identifies the rule.
		exp := de.Expected
		if exp == "" {
			exp = de.Kind.String()
		}
		return "grammar-gap: " + exp

	case "parse-error":
		// Keep Kind + Expected; strip offset/window which are symbol-specific.
		base := de.Kind.String()
		if de.Expected != "" {
			base += " (expected " + de.Expected + ")"
		}
		return "parse-error: " + base

	case "wrong-scheme":
		return "wrong-scheme"

	default:
		return "other: " + truncate(err.Error(), 80)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

// GapEntry is one row in the sorted gap table.
type GapEntry struct {
	Msg          string
	Count        int
	FirstFixture string
	FirstSymbol  string
	LastFixture  string
}

// SortedGaps returns gap entries sorted descending by count.
func (r *Report) SortedGaps() []GapEntry {
	entries := make([]GapEntry, 0, len(r.GapByMsg))
	for msg, count := range r.GapByMsg {
		entries = append(entries, GapEntry{
			Msg:          msg,
			Count:        count,
			FirstFixture: r.FirstFixture[msg],
			FirstSymbol:  r.FirstSymbol[msg],
			LastFixture:  r.LastFixture[msg],
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Msg < entries[j].Msg
	})
	return entries
}

// PrintReport writes a human-readable coverage report to w.
func PrintReport(r *Report, w io.Writer) {
	total := r.Total
	if total == 0 {
		fmt.Fprintln(w, "Coverage Report: 0 symbols — nothing to report")
		return
	}
	failCount := total - r.OK
	fmt.Fprintf(w, "Coverage Report\n")
	fmt.Fprintf(w, "  Total:  %d\n", total)
	fmt.Fprintf(w, "  OK:     %d (%.2f%%)\n", r.OK, 100*float64(r.OK)/float64(total))
	fmt.Fprintf(w, "  Failed: %d (%.2f%%)\n", failCount, 100*float64(failCount)/float64(total))
	fmt.Fprintf(w, "\nFailure categories (by message, descending count):\n")
	fmt.Fprintf(w, "  %-7s  %-50s  %-26s  %-26s  %s\n",
		"count", "message", "first-fixture", "last-fixture", "first-symbol")
	fmt.Fprintf(w, "  %s\n", strings.Repeat("-", 160))

	entries := r.SortedGaps()
	for i, e := range entries {
		if i >= 50 {
			fmt.Fprintf(w, "  … (%d more categories)\n", len(entries)-50)
			break
		}
		// Truncate long fields for display.
		msg := e.Msg
		if len(msg) > 48 {
			msg = msg[:45] + "…"
		}
		sym := e.FirstSymbol
		if len(sym) > 60 {
			sym = sym[:57] + "…"
		}
		fmt.Fprintf(w, "  %7d  %-50s  %-26s  %-26s  %s\n",
			e.Count, msg, e.FirstFixture, e.LastFixture, sym)
	}
}

// PrintSummaryLine emits a single-line summary suitable for a log file.
func PrintSummaryLine(r *Report, w io.Writer) {
	total := r.Total
	if total == 0 {
		fmt.Fprintln(w, "total=0 ok=0 fail=0 ok_pct=100.00")
		return
	}
	failCount := total - r.OK
	fmt.Fprintf(w, "total=%d ok=%d fail=%d ok_pct=%.2f\n",
		total, r.OK, failCount, 100*float64(r.OK)/float64(total))
}
