// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package proguard implements map-assisted reverse lookup for
// ProGuard/R8 obfuscated JVM symbols.
//
// Map format (standard ProGuard/R8 emit):
//
//	com.example.Foo -> a:
//	    void bar(int) -> b
//	    int mField -> c
//	    123:456:com.example.Foo create() -> a
//	com.example.Bar -> b:
//	    com.example.Foo create() -> a
//
// Input shapes handled:
//
//	a                  → com.example.Foo                    (class only)
//	a.b                → com.example.Foo.bar                (class.member)
//	a.b(int)           → com.example.Foo.bar(int)           (class.member + arg sig)
//
// Requires a Context of kind "proguard_map" — upload the map via
// ContextStore first. A scheme-private per-sha256 cache keeps the
// parsed index alive across calls so a 50 MB map is parsed once.
//
// MangleFidelity is Exact: the map is bijective per (class, obfName,
// signature) triple. Mangle (original → obfuscated) is implemented
// and uses the same parsed index in reverse.
package proguard

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"
	"sync"

	"github.com/jelmerdehen/demangle"
)

const (
	KindClass  int32 = iota + 1
	KindMember
)

const contextKind = "proguard_map"

type Scheme struct{}

var info = demangle.Info{
	Name:            "proguard-map",
	Family:          "java",
	Version:         "proguard-any",
	Description:     "Map-assisted reverse lookup for ProGuard/R8 obfuscated names.",
	Stability:       demangle.Stable,
	MangleFidelity:  demangle.Exact,
	RequiresContext: []string{contextKind},
	// ProGuard-obfuscated identifiers are short/lowercase; collide with
	// many plain identifiers. Knock confidence down when native-mangle
	// sigils appear.
	Negatives: []demangle.Negative{
		{Kind: demangle.NegContains, Pattern: "_$s", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "_Z", Penalty: 100},
		{Kind: demangle.NegContains, Pattern: "Java_", Penalty: 60},
	},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 4 * 1024,
	KindNames: map[int32]string{
		KindClass:  "Class",
		KindMember: "Member",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindClass:  demangle.KindCatNamespace,
		KindMember: demangle.KindCatMethod,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff is a weak signal — obfuscated names share shape with plain
// short identifiers. We return 40 for any short lowercase dotted name
// so auto-detect surfaces this scheme as a "possible" candidate; the
// real gate is presence of a proguard_map Context.
func (Scheme) Sniff(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	if strings.ContainsAny(s, "$_") {
		return 0, false
	}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return 0, false
		}
	}
	// Any length ≤ 8 with all lowercase/dots/parens → weak match.
	if len(s) <= 40 {
		return 40, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
	ctx, err := demangle.RequireContext(opts, contextKind)
	if err != nil {
		return nil, err
	}
	idx, err := indexFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Split input: (class)[.member][(args)]
	argSig := ""
	head := in
	if i := strings.IndexByte(in, '('); i >= 0 {
		head = in[:i]
		argSig = in[i:]
	}

	// Try class-only lookup first.
	if cls, ok := idx.classByObf[head]; ok {
		return okResult(in, cls.original+argSig, KindClass, cls.original, "", idx.sha), nil
	}

	// Split into class + member on last '.'.
	dot := strings.LastIndexByte(head, '.')
	if dot < 0 {
		return nil, &demangle.Error{
			Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
			Offset: -1, Expected: "class-or-class.member in map",
			Got: in,
		}
	}
	clsObf := head[:dot]
	memObf := head[dot+1:]
	cls, ok := idx.classByObf[clsObf]
	if !ok {
		return nil, &demangle.Error{
			Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
			Offset: -1, Expected: "obfuscated class in map", Got: clsObf,
		}
	}
	members, ok := cls.byObf[memObf]
	if !ok || len(members) == 0 {
		return nil, &demangle.Error{
			Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
			Offset: -1, Expected: "obfuscated member on " + clsObf, Got: memObf,
		}
	}
	m := members[0]
	display := cls.original + "." + m.original + argSig
	ambiguous := ""
	if len(members) > 1 {
		ambiguous = "yes"
	}
	// Tree.Text holds the bare member name so Mangle can look it up
	// in cls.byOrig directly. Result.Output carries the display form.
	r := &demangle.Result{
		Scheme: "proguard-map", Input: in, Output: display,
		Tree: &demangle.Node{
			Scheme: "proguard-map", Kind: KindMember, Text: m.original,
			Attrs: map[string]string{
				"proguard.map_sha256": idx.sha,
				"proguard.class":      cls.original,
				"proguard.member":     m.original,
			},
		},
		Annotations: map[string]string{
			"proguard.map_sha256": idx.sha,
			"proguard.class":      cls.original,
			"proguard.member":     m.original,
		},
	}
	if ambiguous != "" {
		r.Annotations["proguard.overloads"] = ambiguous
	}
	return r, nil
}

func (Scheme) Mangle(_ context.Context, tree *demangle.Node, opts demangle.Options) (*demangle.Result, error) {
	if tree == nil {
		return nil, demangle.GrammarViolation("proguard-map", "", -1, "non-nil tree")
	}
	c, err := demangle.RequireContext(opts, contextKind)
	if err != nil {
		return nil, err
	}
	idx, err := indexFromContext(c)
	if err != nil {
		return nil, err
	}

	switch tree.Kind {
	case KindClass:
		cls, ok := idx.classByOrig[tree.Text]
		if !ok {
			return nil, &demangle.Error{
				Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
				Offset: -1, Expected: "original class in map", Got: tree.Text,
			}
		}
		return &demangle.Result{Scheme: "proguard-map", Output: cls.obf, Tree: tree}, nil
	case KindMember:
		origClass := tree.Attrs["proguard.class"]
		cls, ok := idx.classByOrig[origClass]
		if !ok {
			return nil, &demangle.Error{
				Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
				Offset: -1, Expected: "original class in map", Got: origClass,
			}
		}
		members, ok := cls.byOrig[tree.Text]
		if !ok || len(members) == 0 {
			return nil, &demangle.Error{
				Kind: demangle.ErrUnrecognisedInput, Scheme: "proguard-map",
				Offset: -1, Expected: "original member on " + origClass, Got: tree.Text,
			}
		}
		return &demangle.Result{
			Scheme: "proguard-map",
			Output: cls.obf + "." + members[0].obf,
			Tree:   tree,
		}, nil
	}
	return nil, demangle.GrammarViolation("proguard-map", "", -1, "Class or Member node")
}

func okResult(in, out string, kind int32, cls, member, sha string) *demangle.Result {
	attrs := map[string]string{"proguard.map_sha256": sha, "proguard.class": cls}
	if member != "" {
		attrs["proguard.member"] = member
	}
	node := &demangle.Node{Scheme: "proguard-map", Kind: kind, Text: out, Attrs: attrs}
	return &demangle.Result{
		Scheme: "proguard-map", Input: in, Output: out,
		Tree: node, Annotations: attrs,
	}
}

// --- parser + index --------------------------------------------------

type classEntry struct {
	original string
	obf      string
	byObf    map[string][]memberEntry // overload → N entries
	byOrig   map[string][]memberEntry
}

type memberEntry struct {
	original string // method name (no arg sig)
	obf      string
	argSig   string // "int, long" when present; "" for fields
	returnT  string // map declares "return-type name(args) -> obf"
}

type index struct {
	sha         string
	classByObf  map[string]*classEntry
	classByOrig map[string]*classEntry
}

// blobAccessor lets the parser read the raw bytes from a Context that
// exposes Blob() — this is our sqlite-backed blobContext (which stores
// the map verbatim).
type blobAccessor interface {
	Blob() []byte
}

var (
	cacheMu sync.Mutex
	cache   = map[string]*index{}
)

func indexFromContext(c demangle.Context) (*index, error) {
	sha := c.SHA256()
	if sha != "" {
		cacheMu.Lock()
		if cached, ok := cache[sha]; ok {
			cacheMu.Unlock()
			return cached, nil
		}
		cacheMu.Unlock()
	}

	var data []byte
	if b, ok := c.(blobAccessor); ok {
		data = b.Blob()
	} else {
		rc, err := c.Reader()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		data, err = io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
	}
	idx, err := parseMap(data, sha)
	if err != nil {
		return nil, err
	}
	if sha != "" {
		cacheMu.Lock()
		cache[sha] = idx
		cacheMu.Unlock()
	}
	return idx, nil
}

func parseMap(data []byte, sha string) (*index, error) {
	idx := &index{
		sha:         sha,
		classByObf:  map[string]*classEntry{},
		classByOrig: map[string]*classEntry{},
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	// Big buffer — some maps have very long single lines (generics,
	// nested types); avoid the default 64KiB limit.
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var current *classEntry
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !(line[0] == ' ' || line[0] == '\t') {
			// Class header: "orig -> obf:"
			cls, err := parseClassHeader(line)
			if err != nil {
				return nil, err
			}
			idx.classByObf[cls.obf] = cls
			idx.classByOrig[cls.original] = cls
			current = cls
			continue
		}
		// Indented — member of current class.
		if current == nil {
			// Ignore orphan member lines (shouldn't happen in valid
			// maps; be lenient on input).
			continue
		}
		if err := parseMemberLine(strings.TrimLeft(line, " \t"), current); err != nil {
			return nil, err
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return idx, nil
}

func parseClassHeader(line string) (*classEntry, error) {
	// "orig -> obf:"
	arrow := strings.Index(line, " -> ")
	if arrow < 0 {
		return nil, &demangle.Error{
			Kind: demangle.ErrGrammarViolation, Scheme: "proguard-map",
			Offset: -1, Expected: "' -> ' in class header", Got: line,
			Window: snippet(line),
		}
	}
	orig := line[:arrow]
	rest := line[arrow+4:]
	if !strings.HasSuffix(rest, ":") {
		return nil, &demangle.Error{
			Kind: demangle.ErrGrammarViolation, Scheme: "proguard-map",
			Offset: -1, Expected: "':' ending class header", Got: rest,
		}
	}
	obf := strings.TrimSuffix(rest, ":")
	return &classEntry{
		original: orig,
		obf:      obf,
		byObf:    map[string][]memberEntry{},
		byOrig:   map[string][]memberEntry{},
	}, nil
}

// parseMemberLine accepts:
//
//	return-type name -> obf
//	return-type name(args) -> obf
//	line1:line2:return-type name(args) -> obf
//	line1:line2:orig-line1:orig-line2:return-type name(args) -> obf
func parseMemberLine(line string, cls *classEntry) error {
	// Strip optional leading line-range prefix: 1-4 colon-separated
	// numeric fields before the real type/name.
	line = stripLinePrefix(line)

	arrow := strings.Index(line, " -> ")
	if arrow < 0 {
		// Silently ignore malformed member lines — R8 occasionally
		// emits annotations we don't care about.
		return nil
	}
	lhs := line[:arrow]
	obf := strings.TrimSpace(line[arrow+4:])

	// lhs = "return-type name" or "return-type name(args)"
	space := strings.IndexByte(lhs, ' ')
	if space < 0 {
		return nil
	}
	retT := lhs[:space]
	rest := lhs[space+1:]
	name := rest
	args := ""
	if lp := strings.IndexByte(rest, '('); lp >= 0 {
		rp := strings.LastIndexByte(rest, ')')
		if rp < lp {
			return nil
		}
		name = rest[:lp]
		args = rest[lp+1 : rp]
	}
	m := memberEntry{original: name, obf: obf, argSig: args, returnT: retT}
	cls.byObf[obf] = append(cls.byObf[obf], m)
	cls.byOrig[name] = append(cls.byOrig[name], m)
	return nil
}

func stripLinePrefix(line string) string {
	// Consume up to 4 leading numeric:numeric patterns.
	for i := 0; i < 4; i++ {
		j := 0
		for j < len(line) && line[j] >= '0' && line[j] <= '9' {
			j++
		}
		if j == 0 || j >= len(line) || line[j] != ':' {
			return line
		}
		line = line[j+1:]
	}
	return line
}

func snippet(s string) string {
	if len(s) > 40 {
		return s[:40]
	}
	return s
}

func init() {
	demangle.Default.Register(Scheme{})
}
