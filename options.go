// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

// Options is the per-call tuning for Demangle / Mangle.
//
// NOTE: there is no Deadline field. Deadlines ride context.Context.
// Callers use context.WithTimeout or context.WithDeadline and pass the
// resulting ctx to Demangle / Mangle / DemangleBatch.
type Options struct {
	Simplified                    bool
	SynthesizeSugar               bool
	QualifyEntities               bool
	DisplayGenericSpecialisations bool
	DisplayThunks                 bool
	ReturnTree                    bool
	VerifyRoundTrip               bool
	AllowLegacy                   bool

	// BestEffortMangle — consumer opt-in for BestEffort-fidelity
	// schemes. Without this, Mangle on a known-lossy input returns
	// ErrNotInvertible instead of a partial Result.
	BestEffortMangle bool

	// Trace — populate Result.Trace with per-op audit info.
	Trace bool

	// SchemeSpecific — per-scheme string knobs that don't fit above.
	// Keys are namespaced by scheme (e.g. "js.allow_deobfuscation").
	SchemeSpecific map[string]string

	// Context — scheme-specific runtime state (ProGuard map, JS source
	// map, Swift symbolic resolver). See context.go.
	Context Context
}

// Warning is a structured non-fatal note attached to a Result.
type Warning struct {
	Code    string // stable identifier, e.g. "swift.symref.unresolved"
	Message string
	OpName  string // optional — primitive op that emitted the warning
}

// TraceEntry is one step in a per-call audit trail (Options.Trace).
type TraceEntry struct {
	Op       string
	InLen    int
	OutLen   int
	DurationNanos int64
	Note     string
}
