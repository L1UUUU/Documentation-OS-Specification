# DOS-6005 — Change Log

**Status:** Draft
**Version:** 1.0
**Category:** Reference

------

# Abstract

This specification defines the change management policy and change log format for Documentation OS.

The Change Log records the evolution of the Documentation OS Specification itself.

Unlike repository history, which records project-specific changes, the Change Log records changes to the Documentation OS standard.

Its purpose is to provide a transparent, traceable history of specification evolution.

------

# Purpose

The purpose of the Change Log is to:

- document specification evolution;
- communicate compatibility changes;
- support implementation upgrades;
- preserve historical context;
- provide a stable migration reference.

Every published Documentation OS Specification version shall include a corresponding Change Log.

------

# Scope

The Change Log applies only to the Documentation OS Specification.

It does not record:

- repository history;
- project development history;
- implementation releases;
- Documentation Engine releases;
- Repository Profile adoption.

Those changes belong to their respective version histories.

------

# Design Principles

The Change Log follows the following principles.

## CL-1 Chronological

Changes shall be recorded in chronological order.

The newest published version should appear first.

------

## CL-2 Immutable

Published Change Log entries shall not be modified retrospectively.

Corrections should be recorded as new entries.

------

## CL-3 Observable

Every published specification change should be visible through the Change Log.

No published specification change should occur silently.

------

## CL-4 Version-Based

Changes are grouped by Specification Version.

Each version forms one logical release.

------

# Change Categories

Documentation OS defines the following change categories.

```text
Specification Changes

├── Added

├── Changed

├── Deprecated

├── Removed

├── Fixed

└── Clarified
```

Each recorded change should belong to exactly one category.

------

# Added

Records newly introduced capabilities.

Examples include:

- new specifications;
- new Repository Profiles;
- additional Documentation Operations;
- new lifecycle stages.

Added entries introduce new functionality.

------

# Changed

Records modifications to existing behavior.

Examples include:

- updated requirements;
- revised lifecycle behavior;
- improved operational semantics.

Changed entries may affect implementations.

Compatibility impact should be stated explicitly.

------

# Deprecated

Records features scheduled for future removal.

Deprecated features remain supported.

Each deprecation should identify:

- affected specification;
- replacement behavior;
- expected removal version (if known).

------

# Removed

Records features removed from the specification.

Removal constitutes a breaking specification change.

Removed entries should reference:

- the version in which deprecation occurred;
- the replacement mechanism.

------

# Fixed

Records corrections that preserve specification semantics.

Examples include:

- wording corrections;
- normative inconsistencies;
- specification defects.

Fixed entries should not introduce new functionality.

------

# Clarified

Records editorial improvements.

Clarifications improve understanding without changing observable behavior.

Examples include:

- improved terminology;
- clearer examples;
- additional explanatory text.

Clarifications are non-breaking.

------

# Entry Format

Each Specification Version should use the following structure.

```text
Version

Release Date

Compatibility

Added

Changed

Deprecated

Removed

Fixed

Clarified
```

Categories with no entries may be omitted.

------

# Compatibility Statement

Each version shall include a compatibility statement.

Typical examples include:

- Fully Backward Compatible
- Partially Backward Compatible
- Breaking Release

Compatibility definitions are specified in:

DOS-6004 — Versioning.

------

# Example

Example release structure.

```text
Version 1.1

Compatibility:
Backward Compatible

Added
- New Repository Profile

Changed
- Validation behavior

Fixed
- Relationship clarification
```

This example is informative.

------

# Relationship to Versioning

The Change Log records specification history.

Versioning defines compatibility.

The two specifications complement one another.

Version numbers identify releases.

Change Log entries explain what changed.

------

# Relationship to Migration

Migration procedures should reference Change Log entries whenever repository transformation becomes necessary.

Change Log entries help implementations determine:

- whether migration is required;
- which behaviors changed;
- which Repository Profiles are affected.

Migration procedures are defined separately.

------

# Initial Release

The first Documentation OS release establishes the baseline specification.

## Documentation OS Specification 1.0

**Status:** Initial Release

**Compatibility:** Baseline Specification

### Added

- Foundation layer
- Information Model
- Knowledge Model
- Runtime Model
- Identity Model
- Relationship Model
- Single Repository Profile
- Repository Layout
- Knowledge Mapping
- Runtime Mapping
- Repository Conventions
- Document Lifecycle
- Runtime Lifecycle
- Knowledge Impact Analysis
- Work Close Pipeline
- Ownership
- Documentation Operations
- Validation
- Health
- Migration
- Documentation Testing
- Agent Entry
- Reading Strategy
- Execution Contract
- Documentation Engine
- CLI
- Reference specifications
- Appendices

This release establishes the initial Documentation OS Specification.

------

## Documentation OS Specification 1.0 — Draft Revision

**Status:** Draft

**Compatibility:** Breaking — Inbox semantics, Identity representation, and Agent Entry validation changed relative to the initial draft baseline.

### Changed

- Information Model: introduced Staging Information as a third top-level category alongside Managed Information and Repository Guidance (DOS-1001).
- Inbox: reclassified from a Knowledge Category to unmanaged Staging Information across DOS-0004, DOS-2003, DOS-2004, DOS-3001, DOS-3003, DOS-3004, DOS-5002, Appendix A, Appendix C.
- Runtime: clarified the Active/Archived Runtime distinction and unified "leave the active execution context" wording (DOS-2004).
- Identity: defined the Work metadata file (`work.yaml`), status and relationship token enumerations, PRD/Issue inheritance of Work identifiers, and archived Runtime path resolution (DOS-2005).
- Agent Entry: local entry files now require both `AGENTS.md` and a content-equivalent `CLAUDE.md`; equivalence validation upgraded from SHOULD to SHALL (DOS-5001, DOS-4002).
- Terminology SoT: DOS-6002 now references DOS-0004 for core terms instead of redefining them.
- Dependency direction: clarified the distinct roles of redefine (concept ownership) and contradict (compliance priority) in DOS-6002.
- Core Principles: renamed Principle 8 from Human Knowledge to Authored Knowledge (DOS-0003).

### Fixed

- Resolved the contradiction where Inbox items were treated as managed Artifacts with identity (DOS-0004, DOS-3001) while not requiring stable identifiers (DOS-2005).
- Aligned specification status to Draft across all specifications pending the final 1.0 release.

------

# Future Releases

Subsequent Specification Versions should append additional entries above older releases.

Historical entries shall remain unchanged.

The Change Log therefore forms the permanent history of Documentation OS evolution.

------

# Compliance

A Documentation OS Specification SHALL:

- maintain a Change Log;
- record every published Specification Version;
- classify changes consistently;
- preserve historical entries;
- identify compatibility impact.

------

# Non-Goals

This specification intentionally does not define:

- source control history;
- implementation release notes;
- project changelogs;
- Documentation Engine release cadence;
- Repository migration procedures.

These concerns remain outside the scope of the Documentation OS Change Log.

------

# References

- README
- DOS-4004 — Migration
- DOS-6003 — Conformance
- DOS-6004 — Versioning

------

# Summary

The Documentation OS Change Log provides the authoritative history of Specification evolution.

By recording every published Specification Version together with its compatibility impact and categorized changes, the Change Log enables transparent evolution of Documentation OS while preserving long-term traceability and implementation confidence.