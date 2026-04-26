// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command swiftprod-extract reads a Mach-O binary (fat or regular) and prints
// each Swift symbol to stdout, one per line, in bare "$s..." form (no leading
// underscore).  It is intended to replace "xcrun nm -gU" in the corpus
// pipeline.
//
// Usage:
//
//	swiftprod-extract [--arch arm64e|arm64|x86_64] [--count] <binary>
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jelmerdehen/demangle/internal/symtab/macho"
)

func main() {
	arch := flag.String("arch", "", "arch to prefer (arm64e|arm64|x86_64)")
	count := flag.Bool("count", false, "print symbol count to stderr")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: swiftprod-extract [--arch ARCH] [--count] <binary>\n")
		os.Exit(1)
	}
	path := flag.Arg(0)
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var n int
	var walkFn func(io.ReaderAt, func(string) error) error
	if *arch != "" {
		walkFn = func(r io.ReaderAt, fn func(string) error) error {
			return macho.WalkArch(r, *arch, fn)
		}
	} else {
		walkFn = macho.Walk
	}

	err = walkFn(f, func(sym string) error {
		fmt.Println(sym)
		n++
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if *count {
		fmt.Fprintf(os.Stderr, "%d symbols\n", n)
	}
}
