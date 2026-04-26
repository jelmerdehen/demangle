// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command swiftc-corpus manages and checks the swiftc Swift symbol corpus.
//
// Usage:
//
//	swiftc-corpus --corpus <path>       path to corpus.txt
//	swiftc-corpus --regenerate          wipe + rebuild corpus via "make regenerate"
//	swiftc-corpus --check               assert committed corpus matches regeneration output
//	swiftc-corpus --diff <symbol>       print demangle / remangle / oracle output for one symbol
//	swiftc-corpus --feature <NN>        restrict --check or --diff to symbols from src/NN_*.swift
//	swiftc-corpus --divergences <path>  path to divergences file (used with --check)
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jelmerdehen/demangle"
	_ "github.com/jelmerdehen/demangle/scheme/swift/all"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

const (
	oracleBin     = "/usr/lib/swift/bin/swift-demangle"
	corpusRelPath = "scheme/swift/stable/testdata/swiftc/corpus.txt"
	divRelPath    = "scheme/swift/stable/testdata/swiftc/swiftc-oracle-divergences.txt"
	testdataDir   = "scheme/swift/stable/testdata/swiftc"
	moduleName    = "github.com/jelmerdehen/demangle"
)

func main() {
	fs := flag.NewFlagSet("swiftc-corpus", flag.ExitOnError)
	corpusFlag := fs.String("corpus", "", "path to corpus.txt (default: auto-located relative to repo root)")
	divergencesFlag := fs.String("divergences", "", "path to swiftc-oracle-divergences.txt")
	regenerate := fs.Bool("regenerate", false, "wipe + rebuild corpus: run 'make regenerate' in testdata dir")
	check := fs.Bool("check", false, "assert committed corpus matches current regeneration output")
	diff := fs.String("diff", "", "print demangle/remangle/oracle output for one symbol")
	feature := fs.String("feature", "", "restrict to symbols from src/NN_*.swift (e.g. '01')")

	_ = fs.Parse(os.Args[1:])

	repoRoot, err := findRepoRoot()
	if err != nil {
		fatalf("cannot locate repo root: %v", err)
	}

	corpusPath := *corpusFlag
	if corpusPath == "" {
		corpusPath = filepath.Join(repoRoot, corpusRelPath)
	}

	divPath := *divergencesFlag
	if divPath == "" {
		divPath = filepath.Join(repoRoot, divRelPath)
	}

	testdata := filepath.Join(repoRoot, testdataDir)

	switch {
	case *regenerate:
		if err := runRegenerate(testdata); err != nil {
			fatalf("regenerate: %v", err)
		}
	case *check:
		if err := runCheck(testdata, corpusPath, *feature); err != nil {
			fmt.Fprintf(os.Stderr, "check: %v\n", err)
			os.Exit(1)
		}
	case *diff != "":
		if err := runDiff(*diff, corpusPath, *feature); err != nil {
			fatalf("diff: %v", err)
		}
	default:
		fs.Usage()
		os.Exit(2)
	}
}

// findRepoRoot walks up from os.Executable() until it finds a go.mod
// declaring module github.com/jelmerdehen/demangle.
func findRepoRoot() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	for {
		candidate := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(candidate); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(line) == "module "+moduleName {
					return dir, nil
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback: check working directory and its parents
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine working directory: %w", err)
	}
	dir = wd
	for {
		candidate := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(candidate); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(line) == "module "+moduleName {
					return dir, nil
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no go.mod with module %q found above %q or %q", moduleName, exe, wd)
}

// runRegenerate runs "make regenerate" in the testdata/swiftc directory.
func runRegenerate(testdata string) error {
	cmd := exec.Command("make", "regenerate")
	cmd.Dir = testdata
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Fprintf(os.Stderr, "regenerate: running make regenerate in %s\n", testdata)
	return cmd.Run()
}

// runCheck regenerates corpus into a temp dir and diffs against the
// committed corpus.txt. Prints the first 20 differing lines; exits 1 on drift.
func runCheck(testdata, corpusPath, featureNN string) error {
	tmpDir, err := os.MkdirTemp("", "swiftc-corpus-check-*")
	if err != nil {
		return fmt.Errorf("mktemp: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpCorpus := filepath.Join(tmpDir, "corpus.txt")
	tmpRaw := filepath.Join(tmpDir, "corpus-raw.txt")
	tmpOut := filepath.Join(tmpDir, "out")

	// Run make oracle (or regenerate) with overridden output paths.
	// We use regenerate but redirect CORPUS, CORPUS_RAW, OUT_DIR.
	cmd := exec.Command("make", "regenerate",
		"OUT_DIR="+tmpOut,
		"CORPUS_RAW="+tmpRaw,
		"CORPUS="+tmpCorpus,
	)
	cmd.Dir = testdata
	cmd.Stdout = os.Stderr // progress to stderr
	cmd.Stderr = os.Stderr
	fmt.Fprintf(os.Stderr, "check: regenerating into %s\n", tmpDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make regenerate: %w", err)
	}

	// Load both corpus files; filter by feature if requested.
	committed, err := loadCorpusLines(corpusPath, featureNN)
	if err != nil {
		return fmt.Errorf("load committed corpus: %w", err)
	}
	fresh, err := loadCorpusLines(tmpCorpus, featureNN)
	if err != nil {
		return fmt.Errorf("load regenerated corpus: %w", err)
	}

	diffs := diffLines(committed, fresh)
	if len(diffs) == 0 {
		fmt.Println("check: OK — committed corpus matches regeneration output")
		return nil
	}

	fmt.Fprintf(os.Stderr, "check: DRIFT — %d differing lines (showing first 20):\n", len(diffs))
	limit := 20
	if len(diffs) < limit {
		limit = len(diffs)
	}
	for _, d := range diffs[:limit] {
		fmt.Fprintln(os.Stderr, d)
	}
	return fmt.Errorf("corpus drift: %d lines differ", len(diffs))
}

// runDiff prints three demangle outputs for a symbol:
//  1. Our Demangle output
//  2. Our Remangle(Demangle(sym)).Output
//  3. swift-demangle oracle output
func runDiff(symbol, corpusPath, featureNN string) error {
	ctx := context.Background()

	// If featureNN is set, verify the symbol appears in that feature's corpus.
	if featureNN != "" {
		lines, err := loadCorpusLines(corpusPath, featureNN)
		if err != nil {
			return fmt.Errorf("load corpus: %w", err)
		}
		found := false
		for _, line := range lines {
			mangled, _, _ := parseCorpusLine(line)
			if mangled == symbol {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "warning: symbol %q not found in feature %q corpus\n", symbol, featureNN)
		}
	}

	// Build a catalog with the stable scheme.
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	// 1. Demangle output.
	demangleOut := "(error)"
	var tree *demangle.Node
	r1, err := cat.Demangle(ctx, symbol, &demangle.Options{ReturnTree: true})
	if err != nil {
		demangleOut = fmt.Sprintf("(demangle error: %v)", err)
	} else {
		demangleOut = r1.Output
		tree = r1.Tree
	}

	// 2. Remangle output (Remangle(Demangle(sym))).
	remangleOut := "(no tree — demangle failed)"
	if tree != nil {
		r2, err := cat.Mangle(ctx, "swift-stable", tree, nil)
		if err != nil {
			remangleOut = fmt.Sprintf("(remangle error: %v)", err)
		} else {
			remangleOut = r2.Output
		}
	}

	// 3. Oracle output (swift-demangle).
	oracleOut := runOracle(symbol)

	fmt.Printf("symbol:   %s\n", symbol)
	fmt.Printf("demangle: %s\n", demangleOut)
	fmt.Printf("remangle: %s\n", remangleOut)
	fmt.Printf("oracle:   %s\n", oracleOut)
	return nil
}

// runOracle invokes swift-demangle and returns the demangled string,
// or an error message if the binary is unavailable.
func runOracle(symbol string) string {
	cmd := exec.Command(oracleBin, symbol)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("(oracle error: %v)", err)
	}
	// swift-demangle outputs: "<sym> ---> <demangled>\n"
	line := strings.TrimSpace(string(out))
	if idx := strings.Index(line, " ---> "); idx >= 0 {
		return strings.TrimSpace(line[idx+len(" ---> "):])
	}
	// If no arrow found, swift-demangle returned the symbol unchanged.
	if line == symbol {
		return "(oracle: identity — not demangled)"
	}
	return line
}

// loadCorpusLines reads corpus.txt, strips comment/blank lines,
// and optionally filters to a feature prefix (e.g. "01" → "BasicTypes").
func loadCorpusLines(path, featureNN string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Derive module name filter from feature number if given.
	moduleFilt := ""
	if featureNN != "" {
		moduleFilt = featureModuleName(featureNN, filepath.Dir(path))
	}

	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if moduleFilt != "" {
			_, demangled, ok := parseCorpusLine(line)
			if !ok {
				continue
			}
			if !strings.HasPrefix(demangled, moduleFilt+".") && demangled != moduleFilt {
				continue
			}
		}
		lines = append(lines, line)
	}
	return lines, sc.Err()
}

// parseCorpusLine splits "mangled ---> demangled" into its parts.
// Returns (mangled, demangled, true) or ("", "", false) on bad format.
func parseCorpusLine(line string) (mangled, demangled string, ok bool) {
	const sep = " ---> "
	idx := strings.Index(line, sep)
	if idx < 0 {
		return "", "", false
	}
	return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+len(sep):]), true
}

// featureModuleName derives the Swift module name from a feature number.
// It looks in testdata/swiftc/src/ for a file matching NN_*.swift, then
// applies the same transform as the Makefile module-name function:
//
//	strip leading digits + underscore, replace _x → X, capitalise first letter.
//
// If the src directory cannot be read, a best-effort derivation is used.
func featureModuleName(nn, testdataSwiftcDir string) string {
	srcDir := filepath.Join(testdataSwiftcDir, "src")
	entries, err := os.ReadDir(srcDir)
	if err == nil {
		for _, e := range entries {
			name := e.Name()
			if !strings.HasSuffix(name, ".swift") {
				continue
			}
			if !strings.HasPrefix(name, nn+"_") {
				continue
			}
			base := strings.TrimSuffix(name, ".swift")
			return makefileModuleName(base)
		}
	}
	// Fallback: no source directory found.
	return nn
}

// makefileModuleName applies the Makefile module-name transformation:
//
//	strip leading "NN_", replace "_x" → "X", capitalise first letter.
func makefileModuleName(base string) string {
	// Strip leading digits and underscore (e.g. "01_basic_types" → "basic_types").
	rest := base
	for i, ch := range base {
		if ch == '_' {
			rest = base[i+1:]
			break
		}
		if ch < '0' || ch > '9' {
			break
		}
	}
	// Replace _x → X (camelCase) for each underscore-letter pair.
	var sb strings.Builder
	capitaliseNext := true
	for _, ch := range rest {
		if ch == '_' {
			capitaliseNext = true
			continue
		}
		if capitaliseNext {
			ch = toUpper(ch)
			capitaliseNext = false
		}
		sb.WriteRune(ch)
	}
	return sb.String()
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

// diffLines returns lines that appear in one slice but not the other,
// as a list of "- <line>" / "+ <line>" pairs.
func diffLines(committed, fresh []string) []string {
	// Build lookup sets.
	committedSet := make(map[string]struct{}, len(committed))
	for _, l := range committed {
		committedSet[l] = struct{}{}
	}
	freshSet := make(map[string]struct{}, len(fresh))
	for _, l := range fresh {
		freshSet[l] = struct{}{}
	}

	var diffs []string
	for _, l := range committed {
		if _, ok := freshSet[l]; !ok {
			diffs = append(diffs, "- "+l)
		}
	}
	for _, l := range fresh {
		if _, ok := committedSet[l]; !ok {
			diffs = append(diffs, "+ "+l)
		}
	}
	return diffs
}

// fatalf prints an error to stderr and exits 1.
func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "swiftc-corpus: "+format+"\n", args...)
	os.Exit(1)
}

// Ensure the stable package's Remangle function is referenced so the
// import is not pruned — we use it indirectly via cat.Mangle.
var _ = stable.Scheme{}
