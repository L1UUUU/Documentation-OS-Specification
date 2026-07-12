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

Every published Documentation OS Specification version includes a corresponding Change Log.

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

Changes SHALL be recorded in chronological order.

The newest published version SHOULD appear first.

------

## CL-2 Immutable

Published Change Log entries SHALL NOT be modified retrospectively.

Corrections SHOULD be recorded as new entries.

------

## CL-3 Observable

Every published specification change SHOULD be visible through the Change Log.

No published specification change SHOULD occur silently.

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

Each recorded change SHOULD belong to exactly one category.

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

Compatibility impact SHOULD be stated explicitly.

------

# Deprecated

Records features scheduled for future removal.

Deprecated features remain supported.

Each deprecation SHOULD identify:

- affected specification;
- replacement behavior;
- expected removal version (if known).

------

# Removed

Records features removed from the specification.

Removal constitutes a breaking specification change.

Removed entries SHOULD reference:

- the version in which deprecation occurred;
- the replacement mechanism.

------

# Fixed

Records corrections that preserve specification semantics.

Examples include:

- wording corrections;
- normative inconsistencies;
- specification defects.

Fixed entries SHOULD NOT introduce new functionality.

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

Each Specification Version SHOULD use the following structure.

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

Categories with no entries MAY be omitted.

------

# Compatibility Statement

Each version SHALL include a compatibility statement.

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

Migration procedures SHOULD reference Change Log entries whenever repository transformation becomes necessary.

Change Log entries help implementations determine:

- whether migration is required;
- which behaviors changed;
- which Repository Profiles are affected.

Migration procedures are defined separately.

------

# Initial Release

The first Documentation OS release establishes the baseline specification.

## Documentation OS Specification 1.0 — v7 Revision

**Status:** Draft

**Compatibility:** Draft-breaking refinement — Work state model consolidated to two observable states; Work completion semantics closed; no conceptual architecture change.

### Changed

- Work lifecycle state model consolidated: Accepted…Validated are now workflow phases (conceptual, non-persistable); the only observable lifecycle states are Active and Completed, determined by directory location. RI-1 narrowed accordingly (DOS-3002).
- Work completion semantics unified to the Complete stage as the Completed terminal state; Cleanup is mandatory post-completion maintenance, idempotent and independently retried, not a precondition for the Completed state (DOS-3004, DOS-0004, DOS-1003, DOS-5003).
- Knowledge lifecycle state diagram corrected: `Active ↔ Updated` replaced by a directed evolution cycle (Active → Updated → Validated → Active), resolving the conflict with "Validation MUST follow" and the missing post-Validation return path (DOS-3001).
- Generate Work operation and `generate work` CLI command now regenerate `.scratch/INDEX.md` as part of creation, with INDEX regeneration included in rollback (DOS-4001, DOS-5005).
- Documentation Testing lifecycle text corrected from "5-phase sequence" to "4-phase sequence" (DOS-4005).

### Removed

- Residual `Close` stage in the DOS-2004 Runtime Lifecycle diagram, left over from the v6 Close-stage removal (DOS-2004).

### Fixed

- WC invariant numbering renumbered consecutively (WC-1/2/3/4/5) after the v6 WC-4 removal left a gap (DOS-3004).
- HANDOFF validation rule reworded from the negative "HANDOFF.md is missing" to the positive "each Work contains a HANDOFF.md" for consistency with the other mandatory structure rules (DOS-4002).
- Completed Runtime document wording changed from "retain their established identities" to "retain their stable Work-scoped addresses and workstream slug", aligning with the Identity Model (Runtime has no global Identity) (DOS-3001).

### Clarified

- The v6 entry's "two states (Active / Completed)" now reflects the consolidated observable-state model (two normative repository-observable Work states).
- Complete command partial-success behavior: when Complete succeeds and Cleanup fails, the Work is already Completed; the command returns Failure with a recovery instruction to retry only Cleanup (DOS-5005).
- Pending Conformance test fixtures enumerated for the final Conformance Review (DOS-4005).

------

## Documentation OS Specification 1.0 — v6 Revision

**Status:** Draft

**Compatibility:** Draft-breaking refinement — lifecycle simplification (Completed/Closed merged into a single terminal state); no conceptual architecture change.

### Removed

- `Closed` lifecycle state and the Close pipeline stage: a Work now has two states (Active / Completed), and `Completed` (directory in `completed/`) is the terminal state (DOS-3002, DOS-3004).
- WC-4 invariant ("Cleanup precedes Close"), made redundant by the Close stage removal (DOS-3004).

### Changed

- Work Close Pipeline consolidated from five stages to four: Knowledge Synchronization → Validation → Complete → Cleanup (DOS-3004, DOS-0004).
- Ownership now concludes when the Work reaches the Completed state (Complete stage directory move), no longer at a separate Closed state (DOS-3002, DOS-3005).
- Failure handling after Complete rewritten: the Work is already terminal; the remaining Cleanup stage MAY be retried idempotently (DOS-3004).
- WC-5 rewritten from "execute exactly once" to "reach successful completion at most once; failed or interrupted attempts MAY be retried idempotently" (DOS-3004).
- Complete Operation / Complete Command redefined as orchestrating two lifecycle stages (Complete: move directory + release Ownership; Cleanup: regenerate INDEX), with INDEX regeneration belonging to the Cleanup stage (DOS-4001, DOS-5005).
- Runtime Completion no longer requires INDEX regeneration as a precondition; Completion is reached upon the Complete stage directory move (DOS-1003).
- Document Lifecycle split into two chains: Knowledge (Created → Active ↔ Updated → Validated → Archived → Retired) and Runtime document (Created → Active → Updated → Validated → Completed) (DOS-3001).
- Runtime definition harmonized to "Active Runtime exists only while Work remains active" across DOS-0003, DOS-0004, DOS-1001, DOS-1003, and README.
- Documentation Operation terminology list in DOS-0004 aligned to the six normative categories (Generate / Synchronize / Validate / Complete / Migrate / Inspect).

### Added

- Generate Work sub-operation under the Generate category: creates the `active/<slug>/` workspace, PRD template, empty `issues/` directory, and empty `HANDOFF.md`; verifies workstream slug uniqueness across active+completed; rolls back on failure (DOS-4001).
- `generate work` CLI command and Generate command category (DOS-5005).
- `Archived` state definition for Knowledge documents (DOS-3001).

### Fixed

- Removed the `issues/ (issue definitions, if any)` wording in DOS-5003; `issues/` is a mandatory Core Runtime Asset containing at least one Issue before Complete.

### Clarified

- Completed Runtime core assets are preserved as immutable historical records (DOS-0004, DOS-1001).
- The Completed terminal state and Ownership release point are documented at the Complete stage (DOS-3004).

------

## Documentation OS Specification 1.0 — v5 Revision

**Status:** Draft

**Compatibility:** Draft-breaking refinement — consistency closure; no architecture change, but Work completion requirements were tightened (HANDOFF at creation, at least one Issue before Complete).

### Fixed

- Resolved the HANDOFF required-time contradiction: HANDOFF.md is now generated at Work creation (DOS-2004, DOS-3002), so Validation treats a missing HANDOFF.md as an error (DOS-4002) instead of a Close-stage-only escalation.
- Reconciled the INDEX contract with HANDOFF presence (DOS-2004, DOS-5004): every Work always exposes a HANDOFF path.
- Added the missing Validation/Cleanup stages to the Work Close Pipeline diagrams in DOS-0004 and DOS-2004.
- Reclassified the v4 compatibility label from "Backward Compatible" to "Draft-breaking refinement" (DOS-6005).

### Changed

- issues/ SHALL contain at least one NN-<slug>.md before the Complete stage (DOS-4002, DOS-3004); an empty issues/ directory is permitted only while a Work is active.
- Renamed the PRD from the Work's "primary identity" to "primary definition" (DOS-2004).
- Split the long-lived reference rule into Knowledge (identifiers) and Runtime (workstream slugs / Work-scoped paths) (DOS-2005).

### Clarified

- Narrowed "Runtime exists only while active" to "Active Runtime ..." (DOS-1003).
- Aligned the Identity Model abstract/purpose/summary to "managed Knowledge artifacts" (DOS-1004).
- Upgraded the informative "Relationships should never outlive ..." to the normative "SHALL NOT outlive ..." (DOS-1005).

------

## Documentation OS Specification 1.0 — v4 Revision

**Status:** Draft

**Compatibility:** Draft-breaking refinement — no conceptual architecture change, but repository conformance requirements were strengthened (mandatory issue front matter and status vocabulary, docs/ entry points, Work internal-structure validation, INDEX fields).

### Fixed

- Resolved duplicate directory-move wording between Complete and Cleanup in DOS-3004.
- Removed the circular definition in Runtime Completion (DOS-1003).
- Eliminated residual archival terminology across DOS-3002, DOS-2004, and DOS-3004.
- Reconciled the Artifact and Identity definition conflict between DOS-0004 and DOS-2005.
- Aligned the Work starting point in DOS-1003 with DOS-3002.
- Added the missing Cleanup stage to the Appendix C lifecycle.
- Removed the inaccurate "Runtime continually disappears" wording.
- Removed Inbox from the Knowledge Relationships in error.
- Removed the misleading "directory names may differ" statement.

### Changed

- Rewrote Runtime Completion as a set of verifiable conditions.
- Promoted the four docs/ entry documents to normative.
- Clarified that Work-scoped paths do not carry a status prefix.

### Added

- Defined the minimal issue front-matter contract (status and title) and the issue status vocabulary.
- Defined the minimal INDEX.md generation contract (including issue status).
- Added Validation checks for internal Work structure and docs/ entry points.
- Added a correction flow whereby a Completed Runtime can repair historical errors through a new Work.

### Clarified

- Performed a full-text normative sentence audit: lowercase shall was either upgraded by semantics or rewritten as declarative statements.

------

## Documentation OS Specification 1.0 — v3 Revision

**Status:** Draft

**Compatibility:** Breaking — Runtime architecture reverted from v2 identifier-based model back to directory-based workstreams.

### Removed

- `WORK-NNNN` identifier scheme for Works.
- `work.yaml` metadata file for Work identity and status.
- Abstract directory model `requirements/plan/tasks/notes` under Work.
- `.scratch/archive/WORK-NNNN/` path pattern.
- PRD/Issue Work-ID inheritance convention.
- Mandatory Runtime front matter (identity/status/relationships).

### Changed

- `.scratch/archive/` renamed to `.scratch/completed/`.
- Work identification changed from `WORK-NNNN` to `<workstream-slug>`.
- "Archived Runtime immutable" refined to "Completed Runtime core assets preserved".
- Work state expression changed from status field to directory location (active/ vs completed/).
- Runtime content organization simplified from five abstract directories to three Core Runtime Assets.

### Added

- Core Runtime Assets model: `PRD.md` + `issues/NN-<slug>.md` + `HANDOFF.md`.
- `.scratch/INDEX.md` auto-generated navigation (Engine-generated).
- Work-scoped path addressing for PRD/Issues/Handoff (no global Identity required).
- Core Runtime Assets preservation rule upon Work completion.

### Clarified

- Staging Information added to DOS-0004 Information definition as a third top-level category.
- Artifact terminology harmonized: "Staging items" (not artifacts), "Inbox staging area", "Repository Guidance documents" (not artifacts).
- Agent Entry requires dual-file `AGENTS.md` + `CLAUDE.md` with SHALL-level equivalence validation.
- Dependency rules clarified as bidirectional: higher layers extend without redefining; lower layers shall not contradict higher-layer normative requirements (DOS-6002, Appendix A, README).
- Authored Knowledge principle renamed from Human Knowledge (DOS-0003).

------

## Documentation OS Specification 1.0 — v2 Revision

**Status:** Superseded by v3

**Compatibility:** Breaking — Inbox semantics, Identity representation, and Agent Entry validation changed relative to the initial draft baseline.

### Changed

- Information Model: introduced Staging Information as a third top-level category alongside Managed Information and Repository Guidance (DOS-1001).
- Inbox: reclassified from a Knowledge Category to unmanaged Staging Information across DOS-0004, DOS-2003, DOS-2004, DOS-3001, DOS-3003, DOS-3004, DOS-5002, Appendix A, Appendix C.
- Runtime: clarified the Active/Archived Runtime distinction and unified "leave the active execution context" wording (DOS-2004).
- Agent Entry: local entry files now require both `AGENTS.md` and a content-equivalent `CLAUDE.md`; equivalence validation upgraded from SHOULD to SHALL (DOS-5001, DOS-4002).
- Terminology SoT: DOS-6002 now references DOS-0004 for core terms instead of redefining them.
- Dependency direction: clarified the distinct roles of redefine (concept ownership) and contradict (compliance priority) in DOS-6002.
- Core Principles: renamed Principle 8 from Human Knowledge to Authored Knowledge (DOS-0003).

### Fixed

- Resolved the contradiction where Inbox items were treated as managed Artifacts with identity (DOS-0004, DOS-3001) while not requiring stable identifiers (DOS-2005).
- Aligned specification status to Draft across all specifications pending the final 1.0 release.

------

## Documentation OS Specification 1.0 — Baseline

**Status:** Draft (pending final release)

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

# Future Releases

Subsequent Specification Versions SHOULD append additional entries above older releases.

Historical entries SHALL remain unchanged.

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