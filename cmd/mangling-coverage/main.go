// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// mangling-coverage cross-references Apple's documented Swift mangling
// grammar against our parser + remangler implementation.
//
// METHODOLOGY v1 — APPROXIMATE (grep heuristic).
// Produces FALSE POSITIVES (matches in comments / unreachable branches)
// and FALSE NEGATIVES (rules dispatched via tables without literal
// operator strings). Use as a directional indicator, not ground truth.
// v2 will replace this with parser-instrumented dispatch logging.
//
// Reads:
//   /data/apps/swift/swiftlang/swift/docs/ABI/Mangling.rst
//   /data/apps/swift/swiftlang/swift/docs/ABI/OldMangling.rst
// Greps:
//   scheme/swift/stable/stable.go
//   scheme/swift/stable/remangler.go
//   scheme/swift/common/printer.go
//
// Emits docs/mangling-coverage.md (≤200 lines).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

const disclaimer = `> **Methodology v1**: grep-based string match against ` + "`stable.go`" + ` +
> ` + "`remangler.go`" + ` + ` + "`printer.go`" + `. Produces FALSE POSITIVES (matches in
> comments / unreachable branches) and FALSE NEGATIVES (rules dispatched
> via tables without literal operator strings). Use as a directional
> indicator, not ground truth. v2 will replace this with parser-
> instrumented dispatch logging.
`

// production captures one grammar rule from .rst.
type production struct {
	rule     string   // LHS — e.g. "global", "type"
	op       string   // RHS terminal — e.g. "Ma", "MP", "$s"
	src      string   // source file (Mangling.rst | OldMangling.rst)
	line     int
	rawline  string
}

// reProduction matches a typical grammar line like:
//
//	global ::= type 'Ma'
//	mangled-name ::= '$s' global
//
// Captures: rule (group 1), terminal-token (group 2).
//
// Some lines have multiple terminals or non-terminal RHS only; we record
// each terminal we see for a rule.
var reProduction = regexp.MustCompile(`^\s*([a-z][a-zA-Z0-9_-]*)\s*::=\s*(.*)$`)
var reTerminal = regexp.MustCompile(`'([^']{1,16})'`)

func parseRST(path string) ([]production, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src := strings.TrimPrefix(path, "/data/apps/swift/swiftlang/swift/docs/ABI/")
	var out []production
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		// Skip non-grammar lines: must contain `::=`.
		if !strings.Contains(line, "::=") {
			continue
		}
		m := reProduction.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		rule := m[1]
		rhs := m[2]
		// Strip RST comments (`//` style).
		if i := strings.Index(rhs, "//"); i >= 0 {
			rhs = rhs[:i]
		}
		// Find each terminal token in the RHS.
		terms := reTerminal.FindAllStringSubmatch(rhs, -1)
		if len(terms) == 0 {
			// No literal terminals — record rule with empty op for context.
			out = append(out, production{rule: rule, op: "", src: src, line: lineNum, rawline: strings.TrimSpace(line)})
			continue
		}
		for _, t := range terms {
			out = append(out, production{rule: rule, op: t[1], src: src, line: lineNum, rawline: strings.TrimSpace(line)})
		}
	}
	return out, sc.Err()
}

// loadFile reads the entire file content for grep checks.
func loadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// hasOpRef reports whether the haystack contains a string literal mentioning the
// given mangling operator. We accept several common Go syntaxes:
//
//	"Ma" / 'M' / 'a' (split)
//	`Ma` (raw string)
//	== 'M' && next == 'a' (sequential byte check)
//
// To stay simple + conservative we look for occurrences as quoted
// substrings or as a literal pattern after a single-letter check.
func hasOpRef(haystack, op string) bool {
	if op == "" {
		return false
	}
	// Single-byte op: check `'X'` exists in source.
	if len(op) == 1 {
		return strings.Contains(haystack, "'"+op+"'")
	}
	// Multi-byte op: look for double-quoted substring "Ma", or
	// sequential byte match like 'M', 'a', or backtick raw `Ma`.
	if strings.Contains(haystack, `"`+op+`"`) {
		return true
	}
	if strings.Contains(haystack, "`"+op+"`") {
		return true
	}
	// Sequential byte literals: 'M' followed by 'a' on the same/adjacent
	// lines is the most common parser shape. A loose check: look for the
	// exact pattern 'M', 'a' in close proximity.
	first := op[0]
	rest := op[1:]
	// Search for 'X' literal anchors then check the rest is mentioned
	// near it as literal byte refs.
	anchor := "'" + string(first) + "'"
	idx := 0
	for {
		pos := strings.Index(haystack[idx:], anchor)
		if pos < 0 {
			break
		}
		abs := idx + pos
		// Look forward up to 200 chars for the next byte literal(s).
		window := haystack[abs:min(abs+200, len(haystack))]
		ok := true
		for _, c := range []byte(rest) {
			if !strings.Contains(window, "'"+string(c)+"'") {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
		idx = abs + 1
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type coverage struct {
	rule       string
	op         string
	src        string
	line       int
	parser     bool // mentioned in stable.go
	remangler  bool // mentioned in remangler.go
	printer    bool // mentioned in printer.go
}

func (c coverage) status() string {
	switch {
	case c.parser && c.remangler && c.printer:
		return "✓ covered"
	case c.parser || c.remangler || c.printer:
		return "~ partial"
	default:
		return "✗ missing"
	}
}

func main() {
	repoRoot := flag.String("repo", "/data/p/demangle", "demangle repository root")
	rstStable := flag.String("rst-stable", "/data/apps/swift/swiftlang/swift/docs/ABI/Mangling.rst", "Mangling.rst path")
	rstOld := flag.String("rst-old", "/data/apps/swift/swiftlang/swift/docs/ABI/OldMangling.rst", "OldMangling.rst path")
	out := flag.String("o", "docs/mangling-coverage.md", "output path (relative to repo)")
	flag.Parse()

	prods, err := parseRST(*rstStable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseRST stable: %v\n", err)
		os.Exit(2)
	}
	old, err := parseRST(*rstOld)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseRST old: %v\n", err)
		os.Exit(2)
	}
	prods = append(prods, old...)

	parserSrc, err := loadFile(*repoRoot + "/scheme/swift/stable/stable.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load stable.go: %v\n", err)
		os.Exit(2)
	}
	remanglerSrc, err := loadFile(*repoRoot + "/scheme/swift/stable/remangler.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load remangler.go: %v\n", err)
		os.Exit(2)
	}
	printerSrc, err := loadFile(*repoRoot + "/scheme/swift/common/printer.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load printer.go: %v\n", err)
		os.Exit(2)
	}

	// Dedupe by (rule, op) — many .rst lines reference the same operator.
	seen := make(map[string]coverage)
	for _, p := range prods {
		key := p.rule + "|" + p.op
		if _, exists := seen[key]; exists {
			continue
		}
		c := coverage{rule: p.rule, op: p.op, src: p.src, line: p.line}
		if p.op != "" {
			c.parser = hasOpRef(parserSrc, p.op)
			c.remangler = hasOpRef(remanglerSrc, p.op)
			c.printer = hasOpRef(printerSrc, p.op)
		}
		seen[key] = c
	}

	// Sort by rule, then op.
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Tally.
	var total, covered, partial, missing int
	for _, k := range keys {
		c := seen[k]
		if c.op == "" {
			continue
		}
		total++
		switch {
		case c.parser && c.remangler && c.printer:
			covered++
		case c.parser || c.remangler || c.printer:
			partial++
		default:
			missing++
		}
	}

	// Emit markdown.
	var b strings.Builder
	b.WriteString("# Mangling.rst Coverage\n\n")
	b.WriteString(disclaimer)
	b.WriteString("\n## Summary\n\n")
	pct := float64(covered) / float64(total) * 100
	fmt.Fprintf(&b, "- Total operator productions: **%d**\n", total)
	fmt.Fprintf(&b, "- ✓ covered (parser + remangler + printer): **%d (%.1f %%)**\n", covered, pct)
	fmt.Fprintf(&b, "- ~ partial (some-but-not-all): **%d**\n", partial)
	fmt.Fprintf(&b, "- ✗ missing (none mentioned): **%d**\n\n", missing)

	b.WriteString("## Missing operators\n\n")
	b.WriteString("These appear in Mangling.rst but no literal mention in our parser /\nremangler / printer source. May be false negatives if dispatched via\ntable.\n\n")
	b.WriteString("| rule | op | src |\n|---|---|---|\n")
	missingCount := 0
	for _, k := range keys {
		c := seen[k]
		if c.op == "" || c.parser || c.remangler || c.printer {
			continue
		}
		fmt.Fprintf(&b, "| `%s` | `%s` | %s:%d |\n", c.rule, c.op, c.src, c.line)
		missingCount++
		if missingCount >= 80 {
			fmt.Fprintf(&b, "| ... | ... | (truncated; %d more)| \n", missing-missingCount)
			break
		}
	}

	b.WriteString("\n## Partial operators\n\n")
	b.WriteString("Mentioned in some but not all of (parser, remangler, printer).\n\n")
	b.WriteString("| rule | op | parser | remangler | printer |\n|---|---|---|---|---|\n")
	partialCount := 0
	for _, k := range keys {
		c := seen[k]
		if c.op == "" {
			continue
		}
		full := c.parser && c.remangler && c.printer
		none := !c.parser && !c.remangler && !c.printer
		if full || none {
			continue
		}
		mark := func(b bool) string {
			if b {
				return "✓"
			}
			return ""
		}
		fmt.Fprintf(&b, "| `%s` | `%s` | %s | %s | %s |\n",
			c.rule, c.op, mark(c.parser), mark(c.remangler), mark(c.printer))
		partialCount++
		if partialCount >= 60 {
			fmt.Fprintf(&b, "| ... | ... | ... | ... | (truncated; %d more) |\n", partial-partialCount)
			break
		}
	}

	if err := os.WriteFile(*repoRoot+"/"+*out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", *out, err)
		os.Exit(2)
	}
	fmt.Printf("coverage: total=%d covered=%d (%.1f %%) partial=%d missing=%d → %s\n",
		total, covered, pct, partial, missing, *out)
}
