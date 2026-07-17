// Package engine provides deterministic Documentation OS repository operations.
//
// The package implements Documentation OS Specification 1.0 Draft revision 12
// and Single Repository Profile 1.0. Callers retain responsibility for
// engineering decisions, including the explicit KnowledgeImpact passed to
// Synchronize and the terminal outcome passed to Complete.
//
// This release is an RC contract. Its exported API is intended for consumer
// validation and may receive breaking corrections before v0.1.0. Starting with
// v0.1.0, compatibility follows semantic versioning within the v0 major series.
// Machine integrations should use VersionInfo for compatibility negotiation and
// ErrorCodeOf for error mapping instead of parsing error messages.
package engine
