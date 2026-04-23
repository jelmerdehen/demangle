// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Package sourcemap implements Source Map V3 reverse lookup. Given a
// JS source map blob (uploaded via ContextStore, kind js_source_map)
// and a (line, column) coordinate, it returns the original source +
// line + column + identifier name when the map records one at that
// position.
//
// Input format for Catalog.Demangle is a string "line:column" OR the
// JSON object {"line":L,"column":C}. Callers that hold both the
// minified symbol and the coordinates typically build the string
// form for brevity.
//
// Fidelity None — the relationship is coordinate-keyed lookup rather
// than bijective text transform; there's nothing to round-trip.
//
// Context caching: the parsed map (VLQ-decoded segments + sources +
// names) is cached per sha256 so a 50 MB bundle is parsed once per
// process lifetime.
package sourcemap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jelmerdehen/demangle"
)

// strings import kept; helper.

const (
	KindLocation int32 = iota + 1
)

const contextKind = "js_source_map"

type Scheme struct{}

var info = demangle.Info{
	Name:            "js-sourcemap",
	Family:          "js",
	Version:         "V3",
	Description:     "Source Map V3 reverse lookup. Requires a js_source_map Context.",
	Stability:       demangle.Stable,
	MangleFidelity:  demangle.None,
	RequiresContext: []string{contextKind},
}

var caps = demangle.Capabilities{
	MaxInputBytes: 1024,
	KindNames: map[int32]string{
		KindLocation: "Location",
	},
	KindCategories: map[int32]demangle.KindCategory{
		KindLocation: demangle.KindCatOther,
	},
}

func (Scheme) Info() demangle.Info                 { return info }
func (Scheme) Capabilities() demangle.Capabilities { return caps }

// Sniff is weak — the input is a coordinate "L:C" plus a stored
// Context; there's no characteristic prefix. We return a low
// confidence for "digit:digit" or "digit,digit" shapes so Detect
// surfaces the scheme as a possible candidate.
func (Scheme) Sniff(s string) (int, bool) {
	if _, _, ok := parseCoord(s); ok {
		return 30, true
	}
	return 0, false
}

func (Scheme) Demangle(_ context.Context, in string, opts demangle.Options) (*demangle.Result, error) {
	line, col, ok := parseCoord(in)
	if !ok {
		return nil, &demangle.Error{
			Kind: demangle.ErrGrammarViolation, Scheme: "js-sourcemap",
			Offset: -1, Expected: `"line:column" coordinate`, Got: in,
		}
	}
	ctx, err := demangle.RequireContext(opts, contextKind)
	if err != nil {
		return nil, err
	}
	idx, err := mapFromContext(ctx)
	if err != nil {
		return nil, err
	}
	seg := idx.lookup(line, col)
	if seg == nil {
		return nil, &demangle.Error{
			Kind: demangle.ErrUnrecognisedInput, Scheme: "js-sourcemap",
			Offset: -1, Expected: "mapping at coordinate",
			Got: fmt.Sprintf("line %d col %d", line, col),
		}
	}
	source := ""
	if seg.sourceIdx >= 0 && seg.sourceIdx < len(idx.sources) {
		source = idx.sources[seg.sourceIdx]
	}
	name := ""
	hasName := seg.nameIdx >= 0 && seg.nameIdx < len(idx.names)
	if hasName {
		name = idx.names[seg.nameIdx]
	}
	// sourcesContent lookup — the original source line at origLine,
	// when the map inlined the source text.
	originalLine := ""
	if seg.sourceIdx >= 0 && seg.sourceIdx < len(idx.sourcesContent) {
		content := idx.sourcesContent[seg.sourceIdx]
		if content != "" && seg.origLine >= 0 {
			lines := strings.Split(content, "\n")
			if seg.origLine < len(lines) {
				originalLine = lines[seg.origLine]
			}
		}
	}
	displayLoc := fmt.Sprintf("%s:%d:%d", source, seg.origLine+1, seg.origCol)
	output := displayLoc
	if hasName {
		output = name
	}
	attrs := map[string]string{
		"js.original_source":  source,
		"js.original_name":    name,
		"js.has_name_mapping": boolStr(hasName),
		"js.map_sha256":       idx.sha,
	}
	if originalLine != "" {
		attrs["js.original_line"] = originalLine
	}
	return &demangle.Result{
		Scheme: "js-sourcemap",
		Input:  in,
		Output: output,
		Tree: &demangle.Node{
			Scheme: "js-sourcemap", Kind: KindLocation,
			Text:  displayLoc,
			Attrs: attrs,
		},
		Annotations: attrs,
	}, nil
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func parseCoord(s string) (line, col int, ok bool) {
	// Accept "L:C" or "L,C".
	for _, sep := range []byte{':', ','} {
		if i := strings.IndexByte(s, sep); i > 0 {
			l, errL := strconv.Atoi(strings.TrimSpace(s[:i]))
			c, errC := strconv.Atoi(strings.TrimSpace(s[i+1:]))
			if errL == nil && errC == nil && l >= 0 && c >= 0 {
				return l, c, true
			}
		}
	}
	return 0, 0, false
}

// --- parsed map + index ------------------------------------------

type segment struct {
	genCol     int
	sourceIdx  int
	origLine   int
	origCol    int
	nameIdx    int
	hasSource  bool
	hasOrig    bool
	hasName    bool
}

type parsed struct {
	sha            string
	sources        []string
	sourcesContent []string
	names          []string
	// lineSegments[L] is the slice of segments for generated-line L,
	// sorted by genCol. Binary search for lookup.
	lineSegments [][]segment
}

// lookup finds the largest-genCol segment on line <= col.
func (p *parsed) lookup(line, col int) *segment {
	if line < 0 || line >= len(p.lineSegments) {
		return nil
	}
	segs := p.lineSegments[line]
	if len(segs) == 0 {
		return nil
	}
	// Binary search for largest genCol ≤ col.
	i := sort.Search(len(segs), func(i int) bool {
		return segs[i].genCol > col
	})
	if i == 0 {
		return nil
	}
	return &segs[i-1]
}

// --- JSON envelope ------------------------------------------------

type envelope struct {
	Version        int      `json:"version"`
	File           string   `json:"file"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
	Names          []string `json:"names"`
	Mappings       string   `json:"mappings"`
	// Sections-based indexed maps: future work. For now we fail
	// gracefully on their presence.
	Sections json.RawMessage `json:"sections"`
}

func parseMap(data []byte, sha string) (*parsed, error) {
	var e envelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("source map: json: %w", err)
	}
	if e.Version != 3 && e.Version != 0 {
		return nil, fmt.Errorf("source map: unsupported version %d", e.Version)
	}
	if len(e.Sections) > 0 {
		return nil, fmt.Errorf("source map: indexed sections maps not yet supported")
	}
	p := &parsed{
		sha:            sha,
		sources:        e.Sources,
		sourcesContent: e.SourcesContent,
		names:          e.Names,
	}
	lines := strings.Split(e.Mappings, ";")
	p.lineSegments = make([][]segment, len(lines))

	var (
		curSource = 0
		curLine   = 0
		curCol    = 0
		curName   = 0
	)
	for i, line := range lines {
		if line == "" {
			continue
		}
		// Each line resets genCol to absolute 0; source/origLine/origCol/
		// name deltas carry over.
		genCol := 0
		parts := strings.Split(line, ",")
		segs := make([]segment, 0, len(parts))
		for _, part := range parts {
			if part == "" {
				continue
			}
			seg, err := decodeSegment(part, &genCol, &curSource, &curLine, &curCol, &curName)
			if err != nil {
				return nil, fmt.Errorf("source map line %d: %w", i, err)
			}
			segs = append(segs, seg)
		}
		// Sort for binary search (usually already sorted but tolerant).
		sort.SliceStable(segs, func(i, j int) bool { return segs[i].genCol < segs[j].genCol })
		p.lineSegments[i] = segs
	}
	return p, nil
}

func decodeSegment(s string, genCol, curSource, curLine, curCol, curName *int) (segment, error) {
	var vals []int
	pos := 0
	for pos < len(s) {
		v, next, err := decodeVLQ(s, pos)
		if err != nil {
			return segment{}, err
		}
		vals = append(vals, v)
		pos = next
	}
	seg := segment{sourceIdx: -1, origLine: -1, origCol: -1, nameIdx: -1}
	switch len(vals) {
	case 1:
		*genCol += vals[0]
		seg.genCol = *genCol
	case 4:
		*genCol += vals[0]
		*curSource += vals[1]
		*curLine += vals[2]
		*curCol += vals[3]
		seg.genCol = *genCol
		seg.sourceIdx = *curSource
		seg.origLine = *curLine
		seg.origCol = *curCol
		seg.hasSource = true
		seg.hasOrig = true
	case 5:
		*genCol += vals[0]
		*curSource += vals[1]
		*curLine += vals[2]
		*curCol += vals[3]
		*curName += vals[4]
		seg.genCol = *genCol
		seg.sourceIdx = *curSource
		seg.origLine = *curLine
		seg.origCol = *curCol
		seg.nameIdx = *curName
		seg.hasSource = true
		seg.hasOrig = true
		seg.hasName = true
	default:
		return segment{}, fmt.Errorf("invalid segment field count %d", len(vals))
	}
	return seg, nil
}

// --- context caching ---------------------------------------------

type blobAccessor interface {
	Blob() []byte
}

var (
	cacheMu sync.Mutex
	cache   = map[string]*parsed{}
)

func mapFromContext(c demangle.Context) (*parsed, error) {
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
	p, err := parseMap(data, sha)
	if err != nil {
		return nil, err
	}
	if sha != "" {
		cacheMu.Lock()
		cache[sha] = p
		cacheMu.Unlock()
	}
	return p, nil
}

func init() {
	demangle.Default.Register(Scheme{})
}
