// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command demangle is the CLI frontend to the library. Subcommands:
//
//	demangle <input> [--scheme NAME] [--simplified] [--tree] [--json]
//	mangle --scheme NAME <tree-json|display-string>
//	detect <input> [--strict] [--window POINTS]
//	batch --corpus FILE [--scheme NAME] [--workers N]
//	scheme list [--family FAMILY]
//	scheme show <name>
//	context upload <path> --kind KIND [--meta k=v]
//	context list [--kind KIND]
//	context delete <sha256>
//	fuzz --scheme NAME --duration DURATION
//	version
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jelmerdehen/demangle"

	// Register every bundled in-process scheme. Subprocess adapters
	// (js/obfuscated) require explicit opt-in and are not imported here.
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "demangle":
		err = runDemangle(args)
	case "mangle":
		err = runMangle(args)
	case "detect":
		err = runDetect(args)
	case "batch":
		err = runBatch(args)
	case "scheme":
		err = runScheme(args)
	case "context":
		err = runContext(args)
	case "fuzz":
		err = runFuzz(args)
	case "catalog":
		err = runCatalog(args)
	case "version":
		err = runVersion(args)
	case "help", "-h", "--help":
		usage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %q\n", cmd)
		usage(os.Stderr)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `demangle — CLI frontend to github.com/jelmerdehen/demangle

Commands:
  demangle <input>                 auto-detect + demangle
  mangle --scheme NAME <input>     mangle (only for schemes implementing Mangler)
  detect <input>                   print ranked candidates
  batch --corpus FILE              streaming demangle over a fixture file
  scheme list                      list every registered scheme
  scheme show <name>               print scheme Info + Capabilities
  context upload <path> --kind K   upload a context blob (ProGuard map, JS source map)
  context list                     list stored contexts
  context delete <sha256>          delete a context by sha256
  fuzz --scheme NAME               convenience wrapper over "go test -fuzz"
  catalog stats                    summary: scheme count, family / fidelity / stability breakdowns
  version                          print library build info

See docs/ for details.
`)
}

// ---- demangle ----------------------------------------------------

func runDemangle(args []string) error {
	fs := flag.NewFlagSet("demangle", flag.ExitOnError)
	scheme := fs.String("scheme", "", "scheme name; empty = auto-detect")
	simplified := fs.Bool("simplified", false, "simplified output (scheme-specific)")
	tree := fs.Bool("tree", false, "include AST in output")
	asJSON := fs.Bool("json", false, "print full Result as JSON")
	contextSHA := fs.String("context-sha", "", "sha256 of a stored Context to attach (see `demangle context list`)")
	qualify := fs.Bool("qualify", true, "qualify names with module prefix (scheme-specific)")
	sugar := fs.Bool("sugar", true, "synthesise sugared types ([T], T?, [K:V]) where applicable")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("demangle: expected exactly one <input> argument")
	}
	input := fs.Arg(0)

	opts := &demangle.Options{
		Simplified:      *simplified,
		ReturnTree:      *tree,
		QualifyEntities: *qualify,
		SynthesizeSugar: *sugar,
	}
	if *contextSHA != "" {
		store, err := openDefaultContextStore()
		if err != nil {
			return err
		}
		defer store.Close()
		ctx, err := store.Get(context.Background(), *contextSHA)
		if err != nil {
			return fmt.Errorf("load context %s: %w", *contextSHA, err)
		}
		opts.Context = ctx
	}

	var (
		r   *demangle.Result
		err error
	)
	if *scheme != "" {
		s, ok := demangle.Default.Scheme(*scheme)
		if !ok {
			return fmt.Errorf("unknown scheme: %q", *scheme)
		}
		r, err = s.Demangle(context.Background(), input, *opts)
	} else {
		r, err = demangle.Default.Demangle(context.Background(), input, opts)
	}
	if err != nil {
		return err
	}

	if *asJSON {
		return json.NewEncoder(os.Stdout).Encode(r)
	}
	fmt.Println(r.Output)
	if *tree && r.Tree != nil {
		printTree(os.Stdout, r.Tree, 0)
	}
	return nil
}

func openDefaultContextStore() (demangle.ContextStore, error) {
	path := os.Getenv("DEMANGLE_CONTEXT_DB")
	if path == "" {
		path = filepath.Join(os.TempDir(), "demangle-contexts.db")
	}
	return demangle.OpenContextStore(path)
}

func printTree(w io.Writer, n *demangle.Node, depth int) {
	if n == nil {
		return
	}
	fmt.Fprintf(w, "%s%s(%d) text=%q index=%d\n",
		strings.Repeat("  ", depth), n.Scheme, n.Kind, n.Text, n.Index)
	for _, c := range n.Children {
		printTree(w, c, depth+1)
	}
}

// ---- mangle -------------------------------------------------------

func runMangle(args []string) error {
	fs := flag.NewFlagSet("mangle", flag.ExitOnError)
	scheme := fs.String("scheme", "", "scheme name (required)")
	_ = fs.Parse(args)

	if *scheme == "" {
		return errors.New("mangle: --scheme is required")
	}
	if fs.NArg() != 1 {
		return errors.New("mangle: expected exactly one <tree-json> argument")
	}

	var tree demangle.Node
	if err := json.Unmarshal([]byte(fs.Arg(0)), &tree); err != nil {
		return fmt.Errorf("parse tree json: %w", err)
	}
	r, err := demangle.Default.Mangle(context.Background(), *scheme, &tree, nil)
	if err != nil {
		return err
	}
	fmt.Println(r.Output)
	return nil
}

// ---- detect -------------------------------------------------------

func runDetect(args []string) error {
	fs := flag.NewFlagSet("detect", flag.ExitOnError)
	strict := fs.Bool("strict", false, "fail on any ambiguity")
	window := fs.Int("window", 5, "ambiguity window in points")
	maxCand := fs.Int("max", 5, "max candidates to print")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("detect: expected exactly one <input> argument")
	}
	cands := demangle.Default.Detect(fs.Arg(0), demangle.DetectOptions{
		Strict:          *strict,
		AmbiguityWindow: *window,
		MaxCandidates:   *maxCand,
		IncludeWeak:     true,
	})
	if len(cands) == 0 {
		return errors.New("no candidates")
	}
	for _, c := range cands {
		fmt.Printf("%-24s %3d  signals=%v negatives=%v %s\n",
			c.Scheme, c.Confidence, c.Signals, c.Negatives, c.Diagnostic)
	}
	return nil
}

// ---- batch --------------------------------------------------------

func runBatch(args []string) error {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)
	corpus := fs.String("corpus", "", `path to newline-delimited input file; "-" for stdin`)
	scheme := fs.String("scheme", "", "scheme name; empty = auto-detect per input")
	workers := fs.Int("workers", 0, "worker count; 0 = NumCPU")
	format := fs.String("format", "text", `output format: text | jsonl`)
	onlyOK := fs.Bool("only-ok", false, "suppress rows that errored")
	_ = fs.Parse(args)

	if *corpus == "" {
		return errors.New("batch: --corpus FILE (or - for stdin) is required")
	}
	var f *os.File
	if *corpus == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(*corpus)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	in := make(chan demangle.BatchRequest, 256)
	out := make(chan demangle.BatchResponse, 256)
	done := make(chan demangle.BatchSummary, 1)
	ctx := context.Background()

	go func() {
		summary := demangle.Default.DemangleBatch(ctx, in, out, demangle.BatchOptions{
			Workers: *workers,
			Scheme:  *scheme,
		})
		done <- summary
	}()

	go func() {
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
		var id uint64
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			in <- demangle.BatchRequest{ID: id, Input: line}
			id++
		}
		close(in)
	}()

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	enc := json.NewEncoder(bw)
	for resp := range out {
		if resp.Err != nil {
			if *onlyOK {
				continue
			}
			switch *format {
			case "jsonl":
				_ = enc.Encode(struct {
					ID     uint64 `json:"id"`
					Input  string `json:"input"`
					Error  string `json:"error"`
				}{resp.ID, resp.Input, resp.Err.Kind.String()})
			default:
				fmt.Fprintf(bw, "%d\t%s\tERROR\t%s\n", resp.ID, resp.Input, resp.Err.Kind.String())
			}
			continue
		}
		switch *format {
		case "jsonl":
			_ = enc.Encode(struct {
				ID     uint64 `json:"id"`
				Input  string `json:"input"`
				Scheme string `json:"scheme"`
				Output string `json:"output"`
			}{resp.ID, resp.Input, resp.Scheme, resp.Result.Output})
		default:
			fmt.Fprintf(bw, "%d\t%s\t%s\t%s\n", resp.ID, resp.Input, resp.Scheme, resp.Result.Output)
		}
	}
	summary := <-done
	fmt.Fprintf(os.Stderr,
		"processed=%d succeeded=%d unrecognised=%d truncated=%d grammar=%d internal=%d duration=%s\n",
		summary.Processed, summary.Succeeded, summary.Unrecognised,
		summary.Truncated, summary.Grammar, summary.Internal, summary.Duration)
	return nil
}

// ---- scheme -------------------------------------------------------

func runScheme(args []string) error {
	if len(args) == 0 {
		return errors.New("scheme: expected subcommand list|show")
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "list":
		return runSchemeList(rest)
	case "show":
		return runSchemeShow(rest)
	default:
		return fmt.Errorf("scheme: unknown subcommand %q", sub)
	}
}

func runSchemeList(args []string) error {
	fs := flag.NewFlagSet("scheme list", flag.ExitOnError)
	family := fs.String("family", "", "filter by family")
	_ = fs.Parse(args)

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	for _, info := range demangle.Default.Schemes() {
		if *family != "" && info.Family != *family {
			continue
		}
		fmt.Fprintf(bw, "%-24s %-8s %-14s %-10s %s\n",
			info.Name, info.Family, info.Version, info.Stability, info.MangleFidelity)
	}
	return nil
}

func runSchemeShow(args []string) error {
	if len(args) != 1 {
		return errors.New("scheme show: expected exactly one <name> argument")
	}
	name := args[0]
	s, ok := demangle.Default.Scheme(name)
	if !ok {
		return fmt.Errorf("unknown scheme: %q", name)
	}
	info := s.Info()
	caps := s.Capabilities()
	fmt.Printf("name:            %s\n", info.Name)
	fmt.Printf("family:          %s\n", info.Family)
	fmt.Printf("version:         %s\n", info.Version)
	fmt.Printf("description:     %s\n", info.Description)
	fmt.Printf("stability:       %s\n", info.Stability)
	fmt.Printf("mangle fidelity: %s\n", info.MangleFidelity)
	if len(info.RequiresContext) > 0 {
		fmt.Printf("requires:        %v\n", info.RequiresContext)
	}
	if _, isMangler := s.(demangle.Mangler); isMangler {
		fmt.Printf("implements:      Mangler\n")
	} else {
		fmt.Printf("implements:      Scheme (demangle-only)\n")
	}
	if caps.MaxInputBytes > 0 {
		fmt.Printf("max input:       %d bytes\n", caps.MaxInputBytes)
	}
	if len(caps.KindNames) > 0 {
		fmt.Printf("kinds:           %d\n", len(caps.KindNames))
	}
	if len(info.Negatives) > 0 {
		fmt.Printf("negatives:\n")
		for _, n := range info.Negatives {
			fmt.Printf("  - kind=%d pattern=%q penalty=%d\n", n.Kind, n.Pattern, n.Penalty)
		}
	}
	if len(info.KnownLossy) > 0 {
		fmt.Printf("known-lossy:\n")
		for _, k := range info.KnownLossy {
			fmt.Printf("  - %s — %s\n", k.Pattern, k.Reason)
		}
	}
	return nil
}

// ---- context ------------------------------------------------------

func runContext(args []string) error {
	if len(args) == 0 {
		return errors.New("context: expected subcommand upload|list|delete")
	}
	store, err := openDefaultContextStore()
	if err != nil {
		return err
	}
	defer store.Close()

	sub := args[0]
	rest := args[1:]
	switch sub {
	case "upload":
		return runContextUpload(store, rest)
	case "list":
		return runContextList(store, rest)
	case "delete":
		return runContextDelete(store, rest)
	default:
		return fmt.Errorf("context: unknown subcommand %q", sub)
	}
}

func runContextUpload(store demangle.ContextStore, args []string) error {
	fs := flag.NewFlagSet("context upload", flag.ExitOnError)
	kind := fs.String("kind", "", "context kind (required)")
	var meta multiFlag
	fs.Var(&meta, "meta", "metadata key=value (repeatable)")
	_ = fs.Parse(args)

	if *kind == "" {
		return errors.New("context upload: --kind is required")
	}
	if fs.NArg() != 1 {
		return errors.New("context upload: expected exactly one <path> argument")
	}
	blob, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	m := map[string]string{}
	for _, kv := range meta {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --meta %q, want key=value", kv)
		}
		m[parts[0]] = parts[1]
	}

	sha, err := store.Put(context.Background(), *kind, blob, m)
	if err != nil {
		return err
	}
	fmt.Println(sha)
	return nil
}

func runContextList(store demangle.ContextStore, args []string) error {
	fs := flag.NewFlagSet("context list", flag.ExitOnError)
	kind := fs.String("kind", "", "filter by kind")
	_ = fs.Parse(args)

	list, err := store.List(context.Background(), *kind)
	if err != nil {
		return err
	}
	for _, info := range list {
		fmt.Printf("%s  %-18s %8d bytes  uploaded=%s  last-access=%s\n",
			info.SHA256, info.Kind, info.ByteSize,
			info.UploadedTS.Format(time.RFC3339),
			info.LastAccessTS.Format(time.RFC3339))
	}
	return nil
}

func runContextDelete(store demangle.ContextStore, args []string) error {
	if len(args) != 1 {
		return errors.New("context delete: expected exactly one <sha256> argument")
	}
	return store.Delete(context.Background(), args[0])
}

// ---- fuzz ---------------------------------------------------------

func runFuzz(args []string) error {
	fs := flag.NewFlagSet("fuzz", flag.ExitOnError)
	scheme := fs.String("scheme", "", "scheme name (required)")
	duration := fs.Duration("duration", 5*time.Minute, "fuzz duration")
	_ = fs.Parse(args)

	if *scheme == "" {
		return errors.New("fuzz: --scheme is required")
	}
	// This is a sugar command — it shells out to `go test -fuzz`
	// inside the scheme's subpackage. Implementation deferred to
	// Stage 0.5 when the first native adapter arrives.
	return fmt.Errorf("fuzz: deferred to Stage 0.5 (scheme=%s duration=%s)", *scheme, *duration)
}

// ---- catalog ------------------------------------------------------

func runCatalog(args []string) error {
	if len(args) == 0 {
		return errors.New("catalog: expected subcommand stats")
	}
	switch args[0] {
	case "stats":
		return runCatalogStats(args[1:])
	default:
		return fmt.Errorf("catalog: unknown subcommand %q", args[0])
	}
}

func runCatalogStats(_ []string) error {
	schemes := demangle.Default.Schemes()
	byFidelity := map[string]int{}
	byFamily := map[string]int{}
	byStability := map[string]int{}
	manglers := 0
	for _, info := range schemes {
		byFidelity[info.MangleFidelity.String()]++
		byFamily[info.Family]++
		byStability[info.Stability.String()]++
		s, _ := demangle.Default.Scheme(info.Name)
		if _, ok := s.(demangle.Mangler); ok {
			manglers++
		}
	}
	fmt.Printf("schemes:       %d\n", len(schemes))
	fmt.Printf("manglers:      %d (%.0f%%)\n", manglers, 100*float64(manglers)/float64(len(schemes)))
	fmt.Printf("\nby family:\n")
	for k, v := range byFamily {
		fmt.Printf("  %-8s %d\n", k, v)
	}
	fmt.Printf("\nby fidelity:\n")
	for k, v := range byFidelity {
		fmt.Printf("  %-12s %d\n", k, v)
	}
	fmt.Printf("\nby stability:\n")
	for k, v := range byStability {
		fmt.Printf("  %-14s %d\n", k, v)
	}
	return nil
}

// ---- version ------------------------------------------------------

func runVersion(_ []string) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("demangle: build info unavailable")
		return nil
	}
	fmt.Printf("module:   %s\n", info.Main.Path)
	fmt.Printf("version:  %s\n", info.Main.Version)
	fmt.Printf("go:       %s\n", info.GoVersion)
	return nil
}

// ---- flag helpers -------------------------------------------------

type multiFlag []string

func (m *multiFlag) String() string     { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }
