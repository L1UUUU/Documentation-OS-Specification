# DOS-2002 — Repository Layout

**Status:** Draft
**Version:** 1.0
**Category:** Repository Profile

------

# Abstract

This specification defines the normative repository layout for the Single Repository Profile.

Unlike the Information Model, which defines conceptual domains, this specification defines how those domains are organized within a repository.

The repository layout exists to improve discoverability, consistency, automation, and AI comprehension.

It does not define documentation semantics.

Semantics are defined by the Model specifications.

------

# Purpose

The purpose of the Repository Layout is to provide a deterministic directory structure that:

- separates persistent Knowledge from Runtime;
- minimizes navigation cost;
- enables deterministic Documentation Operations;
- remains stable throughout repository evolution;
- allows AI agents to rapidly establish repository context.

A repository layout is an implementation detail.

It does not redefine the conceptual models established by Documentation OS.

------

# Design Principles

The Repository Layout follows the following principles.

## RL-1 Domain Separation

Knowledge and Runtime SHALL occupy independent repository locations.

Mixing persistent knowledge with temporary execution artifacts is prohibited.

------

## RL-2 Predictability

Every repository using the Single Repository Profile SHOULD expose an identical documentation structure.

Agents SHOULD NOT need repository-specific heuristics to locate documentation.

------

## RL-3 Minimal Top-Level Structure

The repository root SHOULD contain only a small number of Documentation OS entry points.

Additional documentation should be organized beneath dedicated directories.

------

## RL-4 Stable Navigation

Repository layout SHOULD change infrequently.

Repository evolution SHOULD occur primarily within documentation content rather than directory organization.

------

# Repository Root

The following documentation-related entries are defined at the repository root.

```text
repository/

├── AGENTS.md
├── CLAUDE.md
├── docs/
├── .scratch/
└── ...
```

Documentation OS does not constrain source code layout.

Only documentation-related structures are defined here.

------

# AGENTS.md

## Purpose

The root `AGENTS.md` is the canonical source of the Agent Entry Document for AI agents.

A content-equivalent mirror `CLAUDE.md` SHALL accompany it; both files are valid entry points (see DOS-5001 — Agent Entry).

It defines:

- repository overview;
- documentation architecture;
- execution contract;
- Documentation Operations;
- reading strategy;
- repository conventions.

The root AGENTS.md should allow an agent to understand how to navigate the repository before reading project knowledge.

------

# docs/

## Purpose

The `docs/` directory contains persistent Knowledge.

Knowledge stored here represents the long-term understanding of the repository.

The directory participates in the Knowledge lifecycle.

Runtime artifacts SHALL NOT be stored here.

------

# .scratch/

## Purpose

The `.scratch/` directory contains Runtime organized by workstream.

Runtime artifacts support active implementation.

Examples include:

- PRDs;
- generated implementation Issues;
- Handoff documents;
- temporary execution notes.

## Structure

Runtime is organized into:

```text
.scratch/
├── AGENTS.md
├── CLAUDE.md
├── active/
│   └── <workstream-slug>/
├── completed/
│   └── <workstream-slug>/
└── INDEX.md
```

`.scratch/` is the local entry point for Runtime workstreams and SHALL expose both `AGENTS.md` and a content-equivalent `CLAUDE.md` mirror, together with the generated `INDEX.md` (see DOS-5001 — Agent Entry).

Active Work resides beneath `active/`. Completed Work is preserved beneath `completed/` with its Core Runtime Assets (PRD, Issues, Handoff).

The Documentation Engine generates `INDEX.md` to provide Runtime overview and navigation.

Contents of `.scratch/active/` should disappear from active development once Work has completed.

------

# Internal Documentation

Documentation directories may contain local Agent Entry files.

Where present, both `AGENTS.md` and a content-equivalent `CLAUDE.md` SHALL be provided (see DOS-5001 — Agent Entry).

These files provide directory-specific guidance.

For example:

```text
docs/
│
├── AGENTS.md
├── CLAUDE.md
├── architecture/
├── adr/
├── standards/
└── inbox/   (Staging Information)
```

A local AGENTS.md supplements, but does not replace, the root AGENTS.md.

Agents should interpret local guidance within the scope of the corresponding directory.

------

# Documentation Hierarchy

The repository layout forms a navigation hierarchy.

```text
Repository

│

├── Root AGENTS

│

├── Knowledge

│      ├── Local AGENTS

│      └── Knowledge Categories

│

├── Staging

│      └── Inbox

│

└── Runtime

       ├── Local AGENTS

       └── Runtime Artifacts
```

Navigation flows from general guidance toward increasingly specialized documentation.

------

# Documentation Boundaries

The repository layout establishes physical boundaries.

These boundaries exist for organization.

They do not redefine ownership.

For example:

Moving an Architecture document between directories does not change its Knowledge Category.

Changing repository layout never changes documentation semantics.

------

# Reserved Directories

The following paths are reserved by the Single Repository Profile.

```text
docs/

.scratch/
```

Implementations SHALL NOT assign unrelated responsibilities to these directories.

Additional subdirectories MAY be introduced beneath them.

------

# Repository Evolution

Repository layout should remain stable.

When new capabilities are introduced:

Preferred:

```text
docs/

└── new-category/
```

Avoid:

Frequent restructuring of top-level directories.

Stable layout improves:

- discoverability;
- automation;
- documentation tooling;
- AI navigation.

------

# Compatibility

Future Repository Profiles may implement different layouts.

Examples include:

Workspace Profile

```text
workspace/

projects/

shared-docs/
```

Cloud Profile

```text
knowledge/

runtime/

agents/
```

These profiles remain compliant provided they preserve the conceptual models defined by Documentation OS.

------

# Compliance

A compliant Single Repository implementation SHALL satisfy the following requirements.

- The repository SHALL expose one Knowledge location.
- The repository SHALL expose one Runtime location.
- The repository SHALL expose a root Agent Entry.
- Knowledge and Runtime SHALL remain physically separated.
- Repository layout SHALL NOT redefine conceptual semantics.

------

# Non-Goals

This specification intentionally does not define:

- Knowledge Categories;
- Runtime organization;
- file naming conventions;
- identifier allocation;
- lifecycle behavior;
- documentation operations.

These concerns are defined by subsequent specifications.

------

# References

- DOS-1001 — Information Model
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-2001 — Single Repository Profile
- DOS-2003 — Knowledge Mapping
- DOS-2004 — Runtime Mapping

------

# Summary

The Repository Layout provides the concrete filesystem organization for the Single Repository Profile.

It establishes three stable repository entry points:

- `AGENTS.md` (with its content-equivalent mirror `CLAUDE.md`)
- `docs/`
- `.scratch/`

while preserving the conceptual separation between Knowledge and Runtime established by the Documentation OS Model layer.