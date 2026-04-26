// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

//go:build version_diff

package production

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var swiftVersions = []string{"5.9", "6.0", "6.1", "6.2", "6.3"}

// findSwiftDemangle returns the path to swift-demangle for the given version,
// or "" if not found.
func findSwiftDemangle(version string) string {
	candidates := []string{
		fmt.Sprintf("/opt/swift-toolchains/%s/usr/bin/swift-demangle", version),
		fmt.Sprintf("/opt/swift-toolchains/%s.0/usr/bin/swift-demangle", version),
		fmt.Sprintf("/opt/swift-toolchains/%s.1/usr/bin/swift-demangle", version),
		fmt.Sprintf("/opt/swift-toolchains/%s.2/usr/bin/swift-demangle", version),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// demangle1 runs a single swift-demangle binary on a symbol and returns its
// output, stripping the " ---> " arrow prefix that swift-demangle emits.
func demangle1(binaryPath, symbol string) string {
	out, err := exec.Command(binaryPath, symbol).Output()
	if err != nil {
		return ""
	}
	s := strings.TrimSpace(string(out))
	// swift-demangle prints: <input> ---> <output>
	const arrow = " ---> "
	if i := strings.Index(s, arrow); i >= 0 {
		return strings.TrimSpace(s[i+len(arrow):])
	}
	return s
}

// TestCrossVersionDifferential reads every *.txt corpus file, runs each
// symbol through every available Swift toolchain, and logs divergences to
// swift-version-divergences.txt.
//
// The test skips (not fails) when fewer than 2 toolchains are installed, so
// CI without the toolchains stays green.
//
// Run with:
//
//	go test -tags version_diff -v -run TestCrossVersionDifferential ./scheme/swift/stable/testdata/production/
func TestCrossVersionDifferential(t *testing.T) {
	// Discover available toolchains.
	available := map[string]string{}
	for _, v := range swiftVersions {
		if path := findSwiftDemangle(v); path != "" {
			available[v] = path
		}
	}
	if len(available) < 2 {
		t.Skipf("need ≥2 Swift toolchains in /opt/swift-toolchains; found %d: %v",
			len(available), available)
	}
	t.Logf("Available toolchains (%d):", len(available))
	for _, v := range swiftVersions {
		if p, ok := available[v]; ok {
			t.Logf("  %s -> %s", v, p)
		}
	}

	corpusDir := "corpus"
	files, err := filepath.Glob(filepath.Join(corpusDir, "*.txt"))
	if err != nil || len(files) == 0 {
		t.Skip("no corpus/*.txt files found")
	}

	type divergenceRow struct {
		file    string
		symbol  string
		outputs map[string]string // version -> output
	}

	var divergences []divergenceRow
	var totalSymbols, totalDivergent int

	for _, f := range files {
		fh, ferr := os.Open(f)
		if ferr != nil {
			t.Fatalf("open %s: %v", f, ferr)
		}
		sc := bufio.NewScanner(fh)
		sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

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
			if sym == "" {
				continue
			}
			totalSymbols++

			// Run through each available toolchain.
			outputs := make(map[string]string, len(available))
			for _, v := range swiftVersions {
				p, ok := available[v]
				if !ok {
					continue
				}
				outputs[v] = demangle1(p, sym)
			}

			// Check if any version disagrees with another.
			var first string
			diverged := false
			for _, v := range swiftVersions {
				out, ok := outputs[v]
				if !ok {
					continue
				}
				if first == "" {
					first = out
				} else if out != first {
					diverged = true
					break
				}
			}
			if diverged {
				totalDivergent++
				divergences = append(divergences, divergenceRow{
					file:    filepath.Base(f),
					symbol:  sym,
					outputs: outputs,
				})
			}
		}
		fh.Close()
		if err := sc.Err(); err != nil {
			t.Fatalf("scan %s: %v", f, err)
		}
	}

	t.Logf("Cross-version diff: %d divergences found out of %d symbols across %d toolchains",
		totalDivergent, totalSymbols, len(available))

	// Write divergences to swift-version-divergences.txt.
	if len(divergences) > 0 {
		df, ferr := os.OpenFile("swift-version-divergences.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if ferr == nil {
			fmt.Fprintf(df, "\n# version-diff run %s — %d/%d divergent, toolchains: %s\n",
				time.Now().UTC().Format(time.RFC3339),
				totalDivergent, totalSymbols,
				strings.Join(func() []string {
					var vs []string
					for _, v := range swiftVersions {
						if _, ok := available[v]; ok {
							vs = append(vs, v)
						}
					}
					return vs
				}(), " "))
			for _, d := range divergences {
				var parts []string
				for _, v := range swiftVersions {
					if out, ok := d.outputs[v]; ok {
						parts = append(parts, fmt.Sprintf("%s:%s", v, out))
					}
				}
				fmt.Fprintf(df, "%s\t%s\t%s\n", d.file, d.symbol,
					strings.Join(parts, "\t"))
			}
			df.Close()
		}
	}
}
