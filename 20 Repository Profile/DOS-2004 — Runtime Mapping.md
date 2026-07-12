# DOS-2004 — Runtime Mapping

**Status:** Draft
**Version:** 1.0
**Category:** Repository Profile

------

# Abstract

This specification defines how the conceptual Runtime Model is mapped into concrete Runtime structures within the Single Repository Profile.

The Runtime Model defines the semantics of temporary execution information.

This specification defines how those concepts are represented inside a repository.

Unlike the Knowledge domain, Runtime exists solely to support execution.

Its organization optimizes implementation efficiency rather than long-term preservation.

------

# Purpose

The purpose of Runtime Mapping is to organize temporary execution artifacts into a deterministic structure that:

- supports engineering execution;
- supports AI collaboration;
- minimizes long-term repository pollution;
- enables deterministic Work lifecycle management.

Runtime organization should maximize execution efficiency while minimizing permanent repository complexity.

------

# Design Principles

Runtime Mapping follows the following principles.

## RM-1 Runtime Exists for Execution

Runtime exists to complete Work.

It does not exist to preserve repository knowledge.

Every Runtime artifact should contribute directly to implementation.

------

## RM-2 Runtime Is Disposable

Active Runtime SHALL leave the active execution context after Work completion.

It may then be completed (Core Runtime Assets preserved), discarded, or transformed into Knowledge.

Long-lived active Runtime indicates incomplete lifecycle management.

------

## RM-3 Runtime Is Self-Contained

A Runtime workspace should contain everything necessary to execute one Work.

An agent should not need to infer missing planning information.

------

## RM-4 Runtime Produces Knowledge

Runtime itself is temporary.

The understanding created during Runtime becomes permanent Knowledge through Knowledge Synchronization.

------

# Runtime Structure

The Single Repository Profile represents Runtime beneath:

```text
.scratch/
```

The internal organization of `.scratch/` is determined by Runtime responsibilities rather than document types.

A Runtime implementation consists of one or more independent Works.

------

# Runtime Responsibilities

Each Work maintains the following Core Runtime Assets.

```text
Work

├── PRD

├── Issues

├── Handoff

└── Ephemeral Runtime Content
```

These are the canonical Runtime components.

------

## PRD

### Purpose

The PRD (Product Requirements Document) is the canonical entry point for a Work.

It combines clarified requirements and implementation plan into a single authoritative source.

A PRD serves as:

- confirmed scope and constraints;
- implementation strategy and execution phases;
- the Work's primary identity and coordination document.

### Structure

The PRD resides at:

```text
.scratch/active/<workstream-slug>/PRD.md
```

Upon completion, it moves to:

```text
.scratch/completed/<workstream-slug>/PRD.md
```

### Relationships

Work-level relationships are declared in the PRD's front matter YAML block.

This establishes the Work's connections to other Works and Knowledge artifacts.

When a relationship targets a Runtime asset, the reference uses a Work-scoped path. A Work-scoped path has the form `<workstream-slug>/<file>` (for example `<workstream-slug>/PRD.md`, `<workstream-slug>/issues/01-<slug>.md`, or `<workstream-slug>/HANDOFF.md`); it does not include an `active/` or `completed/` segment. The Documentation Engine resolves a Work-scoped path to the Work's current physical location (`active/` or `completed/`) at lookup time, consistent with DOS-1004 and DOS-1005.

------

## Issues

### Purpose

Issues represent execution tasks—the concrete engineering work that implements the PRD.

Each issue corresponds to one implementation task or subtask.

### Structure

Issues reside within the Work's issues directory:

```text
.scratch/active/<workstream-slug>/issues/NN-<slug>.md
```

Where `NN` is a sequential number and `<slug>` is a kebab-case descriptor.

The issue's status and execution details are recorded within the issue file content.

Issues SHALL declare a `status` field in front matter. Normative status values are: `open`, `in-progress`, `done`, `blocked`. The `title` field is optional; when absent, the INDEX falls back to the issue filename. Relationship declarations remain optional.

Minimal front matter example:

```yaml
---
status: in-progress
title: <human-readable issue title>
---
```

------

## Handoff

### Purpose

Handoff documents capture execution context for transfer across:

- agents
- sessions
- execution phases

They preserve the state necessary for another agent or session to continue the Work seamlessly.

### Structure

Handoff resides at:

```text
.scratch/active/<workstream-slug>/HANDOFF.md
```

Handoff content is contextual and may include:

- current execution state;
- pending decisions;
- intermediate results;
- context for continuation.

------

## Ephemeral Runtime Content

### Purpose

Ephemeral Runtime Content captures temporary information useful during implementation but not intended for long-term preservation.

Examples include:

- debugging observations;
- exploratory experiments;
- discarded implementation ideas;
- temporary references.

### Lifecycle

Ephemeral Content may be:

- cleaned up after Work completion;
- promoted into Knowledge if it contains lasting value;
- discarded when no longer relevant.

Unlike Core Runtime Assets, Ephemeral Content is not preserved in completed Work.

------

# Runtime Independence

Each Work should remain independent.

Multiple Runtime Works may coexist.

One Work should not require another Work's temporary artifacts to understand its execution context.

Shared understanding belongs in the Knowledge domain rather than Runtime.

------

# Runtime Lifecycle

Runtime evolves according to the Runtime Lifecycle specification.

Typical evolution:

```text
Observation

↓

Clarification

↓

Implementation Plan

↓

Execution

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Complete

↓

Close
```

Runtime Mapping defines representation.

Lifecycle specifications define behavior.

------

# Runtime and Inbox

Inbox is not part of Runtime.

Inbox is Staging Information, not a Knowledge Category and not Runtime.

Inbox stores unresolved repository observations.

Runtime begins only after an observation has been accepted as implementation work.

The conceptual transition is:

```text
Observation

↓

Inbox

↓

Clarification

↓

Runtime
```

This separation prevents repository observations from becoming implementation artifacts prematurely.

------

# Runtime and Knowledge

Runtime continuously consumes Knowledge.

Examples include:

- Architecture;
- Standards;
- ADR.

Runtime may produce new Knowledge through implementation.

However, Runtime SHALL NOT become the repository's long-term memory.

------

# Repository Representation

The Single Repository Profile defines the normative Runtime structure as:

```text
.scratch/

├── active/
│   └── <workstream-slug>/
│       ├── PRD.md
│       ├── issues/
│       │   ├── NN-<slug>.md
│       │   └── ...
│       └── HANDOFF.md
│
├── completed/
│   └── <workstream-slug>/
│       ├── PRD.md
│       ├── issues/
│       │   ├── NN-<slug>.md
│       │   └── ...
│       └── HANDOFF.md
│
└── INDEX.md
```

This is the normative Single Repository Profile structure.

A `<workstream-slug>` is a lowercase kebab-case identifier that uniquely names a Work within active/ and completed/.

### Core Runtime Assets Preservation

Upon Work completion, the Core Runtime Assets (PRD, Issues, Handoff) are preserved when the Work transitions from active/ to completed/.

Only Ephemeral Runtime Content may be cleaned up.

### INDEX.md

The `.scratch/INDEX.md` file is generated by the Documentation Engine to provide Runtime overview and navigation.

INDEX.md SHALL list, under "Active Works" and "Completed Works" sections, each Work's workstream slug, PRD path, and HANDOFF path, and SHALL list each Issue as `NN-<slug>.md [status] <title>`. Paths MAY include an `active/` or `completed/` segment; INDEX.md is a derived artifact and is regenerated as Work state changes.

------

# Compliance

A compliant Single Repository implementation SHALL satisfy the following requirements.

- Runtime SHALL exist beneath the Runtime domain.
- Runtime SHALL organize execution around Work.
- Runtime SHALL separate planning from execution.
- Runtime SHALL remain temporary.
- Runtime SHALL support Knowledge Synchronization before completion.

------

# Non-Goals

This specification intentionally does not define:

- PRD templates;
- task estimation;
- agile methodologies;
- implementation workflows.

These concerns belong to engineering process rather than Documentation OS.

------

# References

- DOS-1003 — Runtime Model
- DOS-2001 — Single Repository Profile
- DOS-3002 — Runtime Lifecycle
- DOS-3003 — Knowledge Impact Analysis
- DOS-3004 — Work Close Pipeline

------

# Summary

Runtime Mapping defines how temporary engineering execution is represented within the Single Repository Profile.

Runtime is organized around independent Works.

Each Work contains:

- PRD (canonical entry, combining clarified requirements and implementation plan);
- Issues (execution tasks);
- Handoff (cross-context transfer);
- Ephemeral Runtime Content (temporary, disposable).

Runtime exists solely to complete implementation and produce new repository knowledge before leaving the active execution domain.