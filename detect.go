// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"regexp"
	"sort"
	"strings"
)

// DetectOptions configures how Catalog.Detect ranks candidates.
//
// AmbiguityWindow (default 5): if the top candidate's confidence is
// within this many points of the runner-up, Catalog.Demangle returns
// ErrAmbiguous instead of picking arbitrarily. Set TieBreak to
// override.
//
// Strict (default false): any tie at any distance — including the top
// being N points ahead of the runner-up where N ≤ AmbiguityWindow —
// becomes ErrAmbiguous. Equivalent to AmbiguityWindow = ∞ for the
// tie-break decision.
type DetectOptions struct {
	MaxCandidates    int
	MinConfidence    int
	IncludeWeak      bool
	SchemeHintFamily string
	Strict           bool
	TieBreak         TieBreakPolicy
	AmbiguityWindow  int
	// SchemeSpecific boosts: forwarded by Catalog.Demangle when it
	// passes Options through. Populated from Options.SchemeSpecific.
	Boosts map[string]string
}

// TieBreakPolicy governs the exact-tie case (two or more candidates
// with the same top score and AmbiguityWindow == 0).
type TieBreakPolicy int

const (
	// PickHighest — current non-deterministic behaviour; caller does
	// not care which of two equally-scored schemes wins.
	PickHighest TieBreakPolicy = iota
	// PickAlphabetical — deterministic across runs + hosts. Default
	// for CI and bench runs.
	PickAlphabetical
	// ReturnError — equivalent to Strict for exact ties.
	ReturnError
)

// defaultAmbiguityWindow is applied when DetectOptions.AmbiguityWindow
// is zero (explicit zero is allowed; callers who want no window set
// a negative value or use a distinct config).
const defaultAmbiguityWindow = 5

// Detect runs every registered scheme's Sniff predicate, applies
// catalog boosts + scheme negatives, and returns ranked candidates.
func (c *Catalog) Detect(input string, opts DetectOptions) []Candidate {
	c.mu.RLock()
	schemes := make([]Scheme, 0, len(c.schemes))
	for _, s := range c.schemes {
		schemes = append(schemes, s)
	}
	boostMap := c.boosts
	c.mu.RUnlock()

	maxCand := opts.MaxCandidates
	if maxCand <= 0 {
		maxCand = 5
	}
	minConf := opts.MinConfidence
	if minConf <= 0 && !opts.IncludeWeak {
		minConf = 30
	}

	out := make([]Candidate, 0, len(schemes))
	for _, s := range schemes {
		info := s.Info()
		if opts.SchemeHintFamily != "" && info.Family != opts.SchemeHintFamily {
			// Hint is a soft preference: we still run the sniff but
			// drop the score by a large amount if the family mismatches.
			// Zero means skip for now; revisit when we have a real
			// hint consumer.
		}

		conf, ok := s.Sniff(input)
		var signals, negatives []string
		if ok && conf > 0 {
			signals = append(signals, "sniff")
		}

		// Negatives: deduct penalties.
		for _, neg := range info.Negatives {
			if matchNegative(neg, input) {
				conf -= neg.Penalty
				negatives = append(negatives, neg.Pattern)
			}
		}

		// Catalog-level boost.
		if boostMap != nil {
			if signalsBySignal, ok := boostMap[info.Name]; ok {
				for signal, delta := range signalsBySignal {
					if opts.Boosts != nil {
						if _, set := opts.Boosts[signal]; set {
							conf += delta
							signals = append(signals, "boost:"+signal)
						}
					}
				}
			}
		}

		if conf < 0 {
			conf = 0
		}
		if conf < minConf && !opts.IncludeWeak {
			continue
		}

		cand := Candidate{
			Scheme:     info.Name,
			Confidence: conf,
			Signals:    signals,
			Negatives:  negatives,
		}
		if conf < 50 {
			cand.Diagnostic = "weak match"
		}
		out = append(out, cand)
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Confidence != out[j].Confidence {
			return out[i].Confidence > out[j].Confidence
		}
		// Alphabetical tiebreak for determinism inside same-score runs.
		return out[i].Scheme < out[j].Scheme
	})

	if len(out) > maxCand {
		out = out[:maxCand]
	}
	return out
}

// DetectionBoosts is the config shape accepted by WithBoostsMap /
// WithBoostsFile. Keys are scheme names; inner keys are signal names
// checked against Options.SchemeSpecific at call time.
type DetectionBoosts map[string]map[string]int

func matchNegative(n Negative, input string) bool {
	switch n.Kind {
	case NegContains:
		return strings.Contains(input, n.Pattern)
	case NegRegex:
		re, err := regexp.Compile(n.Pattern)
		if err != nil {
			return false
		}
		return re.MatchString(input)
	case NegBodyShape:
		// BodyShape is a catalog-level hook; schemes can register
		// named matchers via RegisterBodyShape. For v1 we treat it
		// as a no-op — schemes that need body-shape negatives use
		// NegRegex instead.
		return false
	}
	return false
}
