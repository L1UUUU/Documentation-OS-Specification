# DOS-1004 — Identity Model

**Status:** Stable
**Version:** 1.0
**Category:** Model

------

# Abstract

This specification defines the Identity Model of Documentation OS.

Identity provides stable, implementation-independent identifiers for documentation artifacts throughout their lifecycle.

Identity enables documentation to evolve without breaking references caused by file renames, directory restructuring, or repository profile changes.

Repository location is mutable.

Identity is not.

------

# Purpose

The purpose of the Identity Model is to establish a stable identity system for every significant documentation artifact.

Identity enables:

- stable cross-references;
- long-term traceability;
- deterministic tooling;
- repository migration;
- profile-independent document management.

Identity is a property of the artifact itself rather than its storage location.

------

# Design Goals

The Identity Model is designed to satisfy the following goals.

## ID-1 Stability

Identity should remain unchanged throughout an artifact's lifetime.

Renaming a document must not change its identity.

Moving a document must not change its identity.

Changing repository layout must not change its identity.

------

## ID-2 Uniqueness

Every managed documentation artifact shall possess exactly one unique identity.

Duplicate identities are prohibited.

------

## ID-3 Human Readability

Identities should remain understandable by humans.

Developers should be able to recognize the approximate category of an artifact simply by reading its identifier.

------

## ID-4 Tool Friendliness

Identities should be deterministic and easy for tooling to process.

Identity allocation should never depend on execution timing or random generation.

------

## ID-5 Profile Independence

Identity belongs to the conceptual model.

Repository Profiles decide where an artifact is stored.

Profiles never redefine identity.

------

# Identity

An Identity is the permanent identifier assigned to a documentation artifact.

Identity is independent of:

- repository location;
- filename;
- directory structure;
- repository profile;
- implementation language.

Identity remains constant until the artifact is retired.

------

# Managed Artifacts

Documentation OS requires stable identities only for managed artifacts.

Examples include:

- Architecture documents
- ADRs
- Standards
- Runtime Work
- PRDs
- Issue collections

Derived artifacts such as generated indexes do not require independent identities.

------

# Identity Lifecycle

Identity is allocated exactly once.

```text
Create Artifact
        │
        ▼
Allocate Identity
        │
        ▼
Artifact Evolves
        │
        ▼
Archive / Retire
```

Identity shall never be reassigned.

Once retired, an identity shall never be reused.

------

# Identity Persistence

Identity survives:

- file renames;
- directory relocation;
- repository restructuring;
- repository profile migration.

Identity does not survive artifact deletion.

Deleted identities remain retired permanently.

------

# Identity Format

Documentation OS intentionally separates identity semantics from identity syntax.

This specification defines only the semantic requirements.

Repository Profiles define the concrete syntax.

For example, a Repository Profile may specify:

```text
ADR-0001

ARCH-0003

STD-0007

WORK-0015
```

Another profile may adopt a different naming convention while preserving identical semantics.

------

# Identity Allocation

Identity allocation shall satisfy the following properties.

## Deterministic

Allocation should produce identical results under identical repository state.

------

## Monotonic

New identities should not invalidate existing identities.

Existing identifiers remain stable indefinitely.

------

## Category Aware

Identity should communicate the conceptual category of the artifact.

The category should be derivable from the identifier.

------

## Repository Scoped

Identity uniqueness is guaranteed within one Documentation OS repository.

Cross-repository uniqueness is outside the scope of this specification.

------

# Identity Ownership

Identity belongs to the artifact.

Identity does not belong to:

- filenames;
- directories;
- repository layout;
- implementation tooling.

Implementations may move artifacts without affecting identity.

------

# Identity References

All long-lived references should target identities rather than filenames whenever possible.

For example:

Preferred:

```text
ADR-0007
```

Less Preferred:

```text
docs/adr/0007-runtime-model.md
```

Identity-based references remain valid after repository restructuring.

------

# Identity and Repository Profiles

Repository Profiles specify:

- filename conventions;
- numbering conventions;
- storage layout.

They do not redefine identity semantics.

For example:

Single Repository Profile may map:

```text
ADR-0007
```

to

```text
docs/adr/0007-runtime-model.md
```

A future Workspace Profile may map the same identity elsewhere without affecting documentation semantics.

------

# Identity and Relationships

Relationships should reference stable identities.

Identity therefore forms the foundation of:

- cross-document references;
- reverse references;
- traceability;
- impact analysis.

Relationship semantics are defined separately.

------

# Compliance

A Documentation OS implementation SHALL satisfy the following requirements.

- Every managed artifact SHALL possess exactly one identity.
- Identity SHALL remain stable throughout the artifact lifecycle.
- Identity SHALL be independent of repository layout.
- Identity SHALL NOT be reused.
- Repository Profiles SHALL preserve identity semantics.

------

# Non-Goals

This specification intentionally does not define:

- filename syntax;
- numbering format;
- directory naming;
- repository conventions;
- document templates.

These concerns belong to Repository Profiles.

------

# References

- DOS-1001 — Information Model
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-1005 — Relationship Model
- DOS-2001 — Single Repository Profile
- DOS-2002 — Repository Layout

------

# Summary

Identity provides a stable, implementation-independent identifier for managed documentation artifacts.

Identity enables Documentation OS to preserve long-term traceability while allowing repository layouts, filenames, and repository profiles to evolve independently.

Identity belongs to the artifact.

Repository structure belongs to the profile.

The two concerns remain intentionally separated.