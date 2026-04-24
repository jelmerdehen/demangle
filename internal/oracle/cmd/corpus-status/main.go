// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command corpus-status shows the status of every Apple corpus fixture
// against our Swift stable demangler and the oracle binary.
//
// Usage:
//
//	go run ./internal/oracle/cmd/corpus-status/ [-corpus PATH]
//
// Output: a tab-separated 4-column table sorted by category
// (MISMATCH first, then GRAMMAR, TRUNCATED, UNSUPPORTED, ORACLE_ERR, MATCH last),
// followed by a summary line.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/jelmerdehen/demangle"
	stable "github.com/jelmerdehen/demangle/scheme/swift/stable"
)

const (
	statusMatch      = "MATCH"
	statusMismatch   = "MISMATCH"
	statusUnsupported = "UNSUPPORTED"
	statusGrammar    = "GRAMMAR"
	statusTruncated  = "TRUNCATED"
	statusOracleErr  = "ORACLE_ERR"
)

// categoryOrder defines sort priority (lower index = earlier in output).
var categoryOrder = map[string]int{
	statusMismatch:    0,
	statusGrammar:     1,
	statusTruncated:   2,
	statusUnsupported: 3,
	statusOracleErr:   4,
	statusMatch:       5,
}

type row struct {
	status  string
	mangled string
	ours    string
	theirs  string
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	// file is .../internal/oracle/cmd/corpus-status/main.go
	// walk up 4 levels to repo root
	dir := filepath.Dir(file)
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}
	return dir
}

func oracleOutput(bin, mangled string) (string, error) {
	cmd := exec.Command(bin, "-expand", mangled)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// swift-demangle outputs lines like:
	//   Demangling for $s...
	//   kind=Global
	//     kind=...
	//   $s... ---> Decoded.Name
	// We want the last "---> " line's RHS, or the whole trimmed output if
	// the arrow format is absent.
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if idx := strings.Index(lines[i], " ---> "); idx >= 0 {
			return strings.TrimSpace(lines[i][idx+6:]), nil
		}
	}
	return strings.TrimSpace(string(out)), nil
}

func classifyErr(err error) string {
	var de *demangle.Error
	if errors.As(err, &de) {
		switch de.Kind {
		case demangle.ErrUnsupported:
			return statusUnsupported
		case demangle.ErrTruncatedInput:
			return statusTruncated
		case demangle.ErrGrammarViolation:
			return statusGrammar
		}
	}
	return statusGrammar
}

func main() {
	var corpusFlag string
	flag.StringVar(&corpusFlag, "corpus", "", "path to corpus file (default: scheme/swift/stable/testdata/apple/manglings.txt under repo root)")
	flag.Parse()

	corpusPath := corpusFlag
	if corpusPath == "" {
		corpusPath = filepath.Join(repoRoot(), "scheme", "swift", "stable", "testdata", "apple", "manglings.txt")
	}

	f, err := os.Open(corpusPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "corpus-status: open corpus: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	oracleBin := "/usr/lib/swift/bin/swift-demangle"
	oracleAvailable := true
	if _, err := os.Stat(oracleBin); err != nil {
		oracleAvailable = false
	}

	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()

	var rows []row

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		parts := strings.SplitN(line, " ---> ", 2)
		if len(parts) != 2 {
			continue
		}
		mangled := strings.TrimSpace(parts[0])
		if !strings.HasPrefix(mangled, "$s") {
			continue
		}

		// Our demangler.
		got, deErr := cat.Demangle(ctx, mangled, nil)
		var ours string
		var status string
		if deErr != nil {
			ours = "<error: " + deErr.Error() + ">"
			status = classifyErr(deErr)
		} else {
			ours = got.Output
		}

		// Oracle.
		var theirs string
		if !oracleAvailable {
			theirs = "(oracle unavailable)"
			if status == "" {
				// We succeeded; mark as MATCH tentatively but no oracle to compare
				status = statusOracleErr
			}
		} else {
			oracleOut, oErr := oracleOutput(oracleBin, mangled)
			if oErr != nil {
				theirs = "<oracle-error: " + oErr.Error() + ">"
				if status == "" {
					status = statusOracleErr
				}
			} else {
				theirs = oracleOut
				if status == "" {
					// deErr was nil
					if ours == theirs {
						status = statusMatch
					} else {
						status = statusMismatch
					}
				}
			}
		}

		rows = append(rows, row{
			status:  status,
			mangled: truncate(mangled, 60),
			ours:    truncate(ours, 80),
			theirs:  truncate(theirs, 80),
		})
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "corpus-status: scanner: %v\n", err)
		os.Exit(1)
	}

	// Sort: by category priority, then mangled for stability.
	sort.SliceStable(rows, func(i, j int) bool {
		pi := categoryOrder[rows[i].status]
		pj := categoryOrder[rows[j].status]
		if pi != pj {
			return pi < pj
		}
		return rows[i].mangled < rows[j].mangled
	})

	// Print table.
	fmt.Printf("STATUS\tMANGLED\tOURS\tTHEIRS\n")
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.status]++
		fmt.Printf("%s\t%s\t%s\t%s\n", r.status, r.mangled, r.ours, r.theirs)
	}

	total := len(rows)
	fmt.Printf("# total: %d, match: %d, mismatch: %d, unsupported: %d, grammar: %d\n",
		total,
		counts[statusMatch],
		counts[statusMismatch],
		counts[statusUnsupported],
		counts[statusGrammar],
	)
}
