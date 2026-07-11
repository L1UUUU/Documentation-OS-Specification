# DOS-2002 — Repository Layout

**Status:** Stable
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

It should never redefine the conceptual models established by Documentation OS.

------

# Design Principles

The Repository Layout follows the following principles.

## RL-1 Domain Separation

Knowledge and Runtime shall occupy independent repository locations.

Mixing persistent knowledge with temporary execution artifacts is prohibited.

------

## RL-2 Predictability

Every repository using the Single Repository Profile should expose an identical documentation structure.

Agents should not need repository-specific heuristics to locate documentation.

------

## RL-3 Minimal Top-Level Structure

The repository root should contain only a small number of Documentation OS entry points.

Additional documentation should be organized beneath dedicated directories.

------

## RL-4 Stable Navigation

Repository layout should change infrequently.

Repository evolution should occur primarily within documentation content rather than directory organization.

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

Runtime artifacts shall not be stored here.

------

# .scratch/

## Purpose

The `.scratch/` directory contains Runtime.

Runtime artifacts support active implementation.

Examples include:

- clarified requirements;
- PRDs;
- generated implementation Issues;
- temporary execution notes.

Contents of `.scratch/` should disappear from active development once Work has completed.

------

# Internal Documentation

Documentation directories may contain local Agent Entry files (`AGENTS.md`, optionally mirrored by a local `CLAUDE.md`).

These files provide directory-specific guidance.

For example:

```text
docs/
│
├── AGENTS.md
├── architecture/
├── adr/
├── standards/
└── inbox/
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

Implementations shall not assign unrelated responsibilities to these directories.

Additional subdirectories may be introduced beneath them.

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

A compliant Single Repository implementation shall satisfy the following requirements.

- The repository shall expose one Knowledge location.
- The repository shall expose one Runtime location.
- The repository shall expose a root Agent Entry.
- Knowledge and Runtime shall remain physically separated.
- Repository layout shall not redefine conceptual semantics.

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