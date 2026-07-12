# DOS-4002 — Validation

**Status:** Draft
**Version:** 1.0
**Category:** Operations

------

# Abstract

This specification defines Validation within Documentation OS.

Validation is the deterministic process of verifying the structural correctness and internal consistency of a Documentation OS repository.

Validation protects repository integrity.

It does not evaluate engineering quality.

Validation determines whether a repository conforms to Documentation OS specifications.

It does not determine whether engineering decisions are correct.

------

# Purpose

The purpose of Validation is to ensure that repository structure remains consistent as the project evolves.

Validation detects deterministic problems before they become long-term documentation defects.

Typical validation objectives include:

- structural correctness;
- reference integrity;
- lifecycle consistency;
- identifier integrity;
- generated content consistency.

Validation enables repositories to evolve safely.

------

# Scope

Validation applies to every managed artifact defined by Documentation OS.

Validation evaluates:

- repository structure;
- managed metadata;
- relationships;
- lifecycle state;
- generated artifacts.

Validation intentionally excludes engineering judgement.

------

# Design Principles

Validation follows the following principles.

## VL-1 Deterministic

Validation shall always produce identical results for identical repository state.

Validation shall never depend on:

- execution time;
- machine environment;
- user identity;
- execution order.

------

## VL-2 Read-Only

Validation shall not modify repository contents.

Repository modification belongs to Documentation Operations such as Synchronize or Migrate.

Validation only reports repository state.

------

## VL-3 Explicit

Every validation failure shall identify:

- the affected artifact;
- the violated rule;
- the reason for failure.

Validation should never produce ambiguous results.

------

## VL-4 Specification Driven

Validation verifies compliance with Documentation OS specifications.

Repository conventions not defined by Documentation OS are outside the scope of Validation.

------

# Validation Categories

Documentation OS defines six normative validation categories.

```text
Validation

├── Identity

├── Relationships

├── Lifecycle

├── Generated Content

├── Repository Structure

└── Managed Regions
```

Implementations may introduce additional validation categories.

------

# Identity Validation

Identity Validation applies exclusively to Knowledge artifacts (ARCH, ADR, STD) that possess identity-managed front matter.

Identity Validation verifies:

- identity uniqueness;
- identity stability;
- identifier format (where defined by the active Repository Profile);
- reserved identifier reuse.

Examples of failures include:

- duplicate identifiers;
- reused retired identifiers;
- missing identities on Knowledge artifacts.

------

# Relationship Validation

Relationship Validation verifies:

- referenced artifacts are resolvable (identity references resolve to docs/; workstream slug references resolve to .scratch/active|completed/<slug>/);
- relationship targets are valid;
- relationship metadata is structurally correct for Knowledge artifacts (Runtime Work relationships are defined in PRD front matter; Issues and HANDOFF.md without structured relationship metadata are exempt from structural metadata checks).

Relationship Validation does not evaluate semantic correctness.

------

# Lifecycle Validation

Lifecycle Validation verifies:

- valid lifecycle states;
- permitted lifecycle transitions;
- required lifecycle ordering.

For Knowledge artifacts: lifecycle state is determined by front matter status metadata.

For Runtime Work: lifecycle state is determined by directory location (active/ or completed/) — valid states correspond to valid directory locations under .scratch/.

Examples include:

- Runtime completed before Knowledge Synchronization;
- invalid lifecycle transitions;
- Work directory not located under active/ or completed/ (location-presence validation).

------

# Generated Content Validation

Generated Content Validation verifies:

- generated artifacts are synchronized;
- managed metadata is current (for Knowledge artifacts; Runtime metadata validation excluded following work.yaml removal);
- generated indexes are consistent (including .scratch/INDEX.md consistency and reproducibility);
- .scratch/INDEX.md exists and reflects the current state of active/ and completed/ workstreams.

Generated content should always be reproducible.

------

# Repository Structure Validation

Repository Structure Validation verifies compliance with the active Repository Profile.

Examples include:

- required entry points exist;
- reserved directories are present (including .scratch/active/, .scratch/completed/, and the requirement that .scratch/INDEX.md exists);
- required repository guidance files exist;
- Agent Entry mirrors (`AGENTS.md`, `CLAUDE.md`) remain content-equivalent at every scope that requires a mirror.

Repository Profiles define profile-specific requirements.

------

# Managed Region Validation

Managed Region Validation verifies:

- managed regions remain structurally intact;
- generated regions are identifiable;
- protected regions have not been corrupted.

Validation does not compare engineering content.

It verifies structural integrity only.

------

# Validation Results

Validation produces one of the following outcomes.

## Passed

No validation failures were detected.

Repository structure complies with Documentation OS.

------

## Passed With Warnings

Repository remains structurally valid.

Non-blocking observations are reported.

Warnings should not prevent lifecycle progression unless required by repository policy.

------

## Failed

One or more mandatory validation rules have been violated.

Repository consistency cannot be guaranteed.

Subsequent lifecycle stages requiring successful validation shall not continue until failures are resolved.

------

# Validation During Lifecycle

Validation occurs at multiple lifecycle stages.

Typical examples include:

```text
Knowledge Synchronization

↓

Validation

↓

Complete
```

Additional validation may occur:

- before migration;
- before release;
- before repository profile upgrades;
- during Documentation Testing.

------

# Validation Reports

Validation should produce deterministic reports.

A report should include:

- validation category;
- affected artifact;
- violated rule;
- severity;
- suggested remediation.

Report format is implementation-defined.

Report semantics are standardized by this specification.

------

# Documentation Operations

Validation is a Documentation Operation.

Documentation Engines shall expose Validation as a first-class operation.

Validation may be executed:

- manually;
- automatically;
- during lifecycle transitions.

Validation execution policy is implementation-defined.

------

# Failure Handling

Validation failures shall be explicit.

Documentation Engines shall avoid partial success reporting when repository correctness cannot be guaranteed.

Validation failures should remain reproducible.

------

# Compliance

A compliant Documentation Engine SHALL validate:

- identities (Knowledge artifacts only);
- relationships (referenced artifacts resolvable; structural metadata checks apply to Knowledge front matter);
- lifecycle consistency (location-presence for Runtime Work; metadata-presence for Knowledge artifacts);
- repository structure (including .scratch/active/, .scratch/completed/, .scratch/INDEX.md);
- generated content (including .scratch/INDEX.md consistency);
- managed regions.

Additional validation rules may be introduced provided they remain compatible with Documentation OS.

------

# Non-Goals

Validation intentionally does not evaluate:

- architecture quality;
- engineering decisions;
- implementation correctness;
- coding style;
- software quality.

These concerns belong to engineering review rather than Documentation OS.

------

# References

- DOS-0003 — Core Principles
- DOS-3001 — Document Lifecycle
- DOS-3002 — Runtime Lifecycle
- DOS-4001 — Documentation Operations
- DOS-4003 — Health
- DOS-4005 — Documentation Testing

------

# Summary

Validation is the deterministic verification mechanism of Documentation OS.

Validation ensures that a repository remains structurally correct throughout its lifecycle.

Validation protects repository integrity.

It does not evaluate engineering judgement.

By separating structural correctness from engineering quality, Documentation OS enables reliable automation while preserving human and AI responsibility for engineering decisions.