# DOS-6004 — Versioning

**Status:** Draft
**Version:** 1.0
**Category:** Reference

------

# Abstract

This specification defines the versioning model of Documentation OS.

Documentation OS consists of multiple independently evolving components.

This specification defines how those versions relate to one another while preserving compatibility across the Documentation OS ecosystem.

Versioning enables Documentation OS to evolve without introducing ambiguity between specifications, repositories, Repository Profiles, Documentation Engines, and tooling.

------

# Purpose

The purpose of Versioning is to establish a consistent versioning strategy for Documentation OS.

Versioning enables:

- specification evolution;
- backward compatibility;
- implementation compatibility;
- repository migration;
- deterministic interoperability.

Documentation OS intentionally separates specification evolution from implementation evolution.

------

# Scope

Versioning applies to the following Documentation OS components:

- Specification
- Repository
- Repository Profile
- Documentation Engine
- CLI

Each component evolves independently.

------

# Versioning Principles

Documentation OS versioning follows the following principles.

## VS-1 Independent Evolution

Different Documentation OS components evolve independently.

Changing one version does not necessarily require changing another.

------

## VS-2 Explicit Compatibility

Compatibility shall be explicitly declared.

Implementations should never assume compatibility implicitly.

------

## VS-3 Stable Semantics

Minor version updates should preserve existing conceptual semantics whenever possible.

Breaking semantic changes require a major version increment.

------

## VS-4 Observable Versioning

Every implementation should expose its supported Documentation OS version.

Version information should be discoverable by both humans and tooling.

------

# Version Types

Documentation OS defines five independent version types.

```text
Documentation OS

├── Specification Version

├── Repository Version

├── Repository Profile Version

├── Documentation Engine Version

└── CLI Version
```

Each version serves a different purpose.

------

# Specification Version

The Specification Version identifies the Documentation OS specification against which compatibility is evaluated.

Examples:

```text
Documentation OS Specification

1.0

1.1

2.0
```

Specification Version defines:

- concepts;
- semantics;
- normative behavior.

All other versions reference the Specification Version.

------

# Repository Version

The Repository Version identifies the documentation state of an individual repository.

Repository Version is implementation-defined.

Typical repository evolution includes:

- documentation restructuring;
- Knowledge growth;
- Runtime completion;
- repository migration.

Repository Version does not imply Specification Version.

------

# Repository Profile Version

The Repository Profile Version identifies the version of the active Repository Profile.

Examples include:

```text
Single Repository Profile

v1

Workspace Profile

v1
```

Repository Profiles evolve independently from the Documentation OS Specification provided they preserve normative semantics.

------

# Documentation Engine Version

The Documentation Engine Version identifies the implementation version of the Documentation Engine.

Examples:

```text
Documentation Engine

0.9

1.0

1.2
```

Different Documentation Engine versions may support the same Specification Version.

------

# CLI Version

The CLI Version identifies the implementation version of the command-line interface.

CLI evolution is independent from the Documentation Engine provided externally observable behavior remains compatible.

------

# Compatibility Matrix

Documentation OS compatibility is determined primarily by Specification Version.

Example:

| Component            | Version              |
| -------------------- | -------------------- |
| Specification        | 1.0                  |
| Repository Profile   | Single Repository v1 |
| Documentation Engine | 1.4                  |
| CLI                  | 1.3                  |

The Documentation Engine and CLI may evolve independently while remaining compatible with Specification 1.0.

------

# Version Relationships

The following dependency hierarchy applies.

```text
Specification

↓

Repository Profile

↓

Documentation Engine

↓

CLI
```

Repository evolution occurs alongside these components.

Repository Version remains independent.

------

# Version Compatibility

Documentation OS defines three compatibility categories.

## Fully Compatible

Implementation supports the declared Specification Version without known deviations.

------

## Backward Compatible

Implementation introduces additional capabilities while preserving all existing normative behavior.

Existing repositories continue functioning without modification.

------

## Breaking Change

Implementation modifies or removes normative behavior.

Breaking changes require a major Specification Version increment.

Repository migration may become necessary.

------

# Version Discovery

Implementations should expose version information through appropriate mechanisms.

Typical information includes:

- Specification Version;
- Repository Profile Version;
- Documentation Engine Version;
- CLI Version.

Discovery mechanism is implementation-defined.

------

# Repository Migration

Version changes requiring repository transformation should use the Migration process defined in:

DOS-4004 — Migration.

Migration preserves repository semantics while updating repository representation.

------

# Deprecation

Documentation OS may deprecate features.

Deprecated behavior:

- remains documented;
- remains functional for at least one compatible Specification Version where practical;
- should identify replacement behavior.

Removal of deprecated behavior constitutes a breaking change.

------

# Semantic Versioning

Documentation OS recommends semantic versioning for implementations.

Typical interpretation:

| Version Component | Meaning                                          |
| ----------------- | ------------------------------------------------ |
| Major             | Breaking specification or implementation changes |
| Minor             | Backward-compatible capability additions         |
| Patch             | Bug fixes and implementation improvements        |

Implementation versioning remains implementation-defined.

------

# Compliance

A Documentation OS implementation SHALL:

- declare its supported Specification Version;
- declare its active Repository Profile;
- expose implementation version information;
- preserve compatibility according to the declared Specification Version.

------

# Non-Goals

This specification intentionally does not define:

- release schedules;
- implementation branching strategy;
- source control workflow;
- release management policy.

These concerns remain implementation-specific.

------

# References

- README
- DOS-4004 — Migration
- DOS-5004 — Documentation Engine
- DOS-6003 — Conformance
- DOS-6005 — Change Log

------

# Summary

Versioning enables Documentation OS to evolve predictably while preserving interoperability.

By separating Specification, Repository Profile, Documentation Engine, CLI, and Repository versions, Documentation OS allows independent implementation evolution without compromising semantic compatibility or repository consistency.

The Specification Version remains the authoritative compatibility reference for the entire Documentation OS ecosystem.