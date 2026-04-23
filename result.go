// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

// Result is what every Demangle / Mangle call returns on success.
// Fields are populated opportunistically; Tree is nil unless
// Options.ReturnTree (or the scheme always fills it).
type Result struct {
	Scheme      string
	Input       string
	Output      string
	Tree        *Node
	Confidence  int // 0–100; populated when auto-detected, else 0
	Annotations map[string]string

	// Partial is true when BestEffortMangle produced an output that
	// could not fully reconstruct the original input. Callers that
	// need exact round-trip check this before trusting Output.
	Partial bool

	// Warnings is structured non-fatal feedback; may be nil.
	Warnings []Warning

	// LostInfo enumerates, in human-readable form, what the scheme
	// could not reconstruct (for Partial == true).
	LostInfo []string

	// Trace is populated only when Options.Trace is set. Nil otherwise.
	Trace []TraceEntry
}

// Candidate is one ranked detection hit.
type Candidate struct {
	Scheme     string
	Confidence int
	Signals    []string // which positive sniffs fired
	Negatives  []string // which negatives fired (reducing score)
	Diagnostic string   // optional weak-match explanation
}
