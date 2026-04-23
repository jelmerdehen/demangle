// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"context"
	"sort"
	"sync"
)

// CatalogOption tunes a Catalog at construction.
type CatalogOption func(*Catalog)

// WithContextStore installs a custom ContextStore. Defaults to
// in-memory (InMemoryContextStore) when omitted.
func WithContextStore(s ContextStore) CatalogOption {
	return func(c *Catalog) { c.contexts = s }
}

// WithBoosts installs detection-time boosts. See DetectionBoosts and
// DetectOptions.Boosts.
func WithBoosts(b DetectionBoosts) CatalogOption {
	return func(c *Catalog) { c.boosts = b }
}

// WithMaxInputBytes sets the catalog-level default input-size cap
// used when a scheme declares Capabilities.MaxInputBytes == 0. Zero
// means "use package default" (DefaultMaxInputBytes).
func WithMaxInputBytes(n int) CatalogOption {
	return func(c *Catalog) { c.maxInputBytes = n }
}

// DefaultMaxInputBytes is the fallback cap applied when neither the
// scheme nor the catalog overrides it.
const DefaultMaxInputBytes = 64 * 1024

// Catalog holds a set of registered Schemes plus an optional
// ContextStore for uploaded blob-identity contexts.
type Catalog struct {
	mu            sync.RWMutex
	schemes       map[string]Scheme
	contexts      ContextStore
	boosts        DetectionBoosts
	maxInputBytes int
}

// NewCatalog builds a fresh Catalog. Schemes are not registered
// automatically — use Register, or blank-import a scheme/<family>
// subpackage to populate Default.
func NewCatalog(opts ...CatalogOption) *Catalog {
	c := &Catalog{schemes: map[string]Scheme{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Default is the package-level catalog populated by blank imports.
// Tests should construct their own catalog via NewCatalog; Default is
// intended for CLI, oracle harness, and integration tests.
var Default = NewCatalog()

// Register adds a Scheme. Panics on duplicate names — schemes are
// shared singletons; overwriting would leak state across tests.
func (c *Catalog) Register(s Scheme) {
	name := s.Info().Name
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, dup := c.schemes[name]; dup {
		panic("demangle: duplicate scheme registration: " + name)
	}
	if c.schemes == nil {
		c.schemes = map[string]Scheme{}
	}
	c.schemes[name] = s
}

// Schemes returns the Info for every registered scheme, sorted by
// name for deterministic output.
func (c *Catalog) Schemes() []Info {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Info, 0, len(c.schemes))
	for _, s := range c.schemes {
		out = append(out, s.Info())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Scheme fetches a registered scheme by name.
func (c *Catalog) Scheme(name string) (Scheme, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.schemes[name]
	return s, ok
}

// ContextStore returns the catalog's uploaded-context store (or nil
// if none was attached at construction).
func (c *Catalog) ContextStore() ContextStore { return c.contexts }

// Demangle auto-detects the scheme and dispatches. See DetectOptions
// for how tie-breaking works on close candidates.
func (c *Catalog) Demangle(ctx context.Context, input string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	// Promote SchemeSpecific into detect-time Boosts so a consumer
	// can signal "apple_binary: true" → boost swift-stable, etc.
	detOpts := DetectOptions{Boosts: opts.SchemeSpecific}
	candidates := c.Detect(input, detOpts)
	if len(candidates) == 0 {
		return nil, &Error{Kind: ErrUnrecognisedInput, Offset: -1}
	}

	// Tie-break check.
	top := candidates[0]
	if len(candidates) > 1 {
		runnerUp := candidates[1]
		window := detOpts.AmbiguityWindow
		if window == 0 {
			window = defaultAmbiguityWindow
		}
		ambiguous := (top.Confidence-runnerUp.Confidence) < window && top.Confidence > 0
		if detOpts.Strict && top.Confidence == runnerUp.Confidence {
			ambiguous = true
		}
		if ambiguous {
			// Apply tie-break policy for exact ties (runnerUp score
			// equals top).
			if top.Confidence == runnerUp.Confidence {
				switch detOpts.TieBreak {
				case PickAlphabetical:
					// candidates is already sorted alphabetically
					// within a score band; top is the first — keep it.
				case ReturnError:
					return nil, ambiguousError(candidates)
				default: // PickHighest
					if detOpts.Strict {
						return nil, ambiguousError(candidates)
					}
				}
			} else {
				// Within window but not exact tie: always ambiguous
				// unless explicitly suppressed (no option for that
				// yet — design intent is "ambiguous surfaces").
				return nil, ambiguousError(candidates)
			}
		}
	}

	sch, ok := c.Scheme(top.Scheme)
	if !ok {
		return nil, &Error{Kind: ErrInternal, Scheme: top.Scheme, Offset: -1}
	}

	// Input-size cap.
	if max := c.effectiveMaxInputBytes(sch); max > 0 && len(input) > max {
		return nil, &Error{Kind: ErrInputTooLarge, Scheme: top.Scheme, Offset: -1}
	}

	r, err := sch.Demangle(ctx, input, *opts)
	if r != nil && r.Confidence == 0 {
		r.Confidence = top.Confidence
	}
	return r, err
}

// Mangle dispatches to the named scheme's Mangler implementation.
// Returns ErrNotInvertible when the scheme does not implement Mangler.
func (c *Catalog) Mangle(ctx context.Context, schemeName string, tree *Node, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	s, ok := c.Scheme(schemeName)
	if !ok {
		return nil, &Error{Kind: ErrInternal, Scheme: schemeName, Offset: -1, Expected: "registered scheme"}
	}
	m, ok := s.(Mangler)
	if !ok {
		return nil, &Error{Kind: ErrNotInvertible, Scheme: schemeName, Offset: -1}
	}
	return m.Mangle(ctx, tree, *opts)
}

// effectiveMaxInputBytes resolves the MaxInputBytes precedence chain:
// scheme override → catalog default → package default.
func (c *Catalog) effectiveMaxInputBytes(s Scheme) int {
	if n := s.Capabilities().MaxInputBytes; n > 0 {
		return n
	}
	if c.maxInputBytes > 0 {
		return c.maxInputBytes
	}
	return DefaultMaxInputBytes
}

func ambiguousError(cands []Candidate) error {
	return &Error{
		Kind:   ErrAmbiguous,
		Offset: -1,
		Cause:  &AmbiguousError{Candidates: cands},
	}
}
