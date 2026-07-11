# DOS-5001 — Agent Entry

**Status:** Draft
**Version:** 1.0
**Category:** Runtime

------

# Abstract

This specification defines how AI agents enter and initialize a Documentation OS repository.

Agent Entry establishes the standard initialization procedure that every Documentation OS-compliant agent should follow before performing repository analysis or implementation work.

The objective of Agent Entry is to ensure that every agent begins with a consistent understanding of the repository, regardless of the underlying AI model or implementation.

------

# Purpose

The purpose of Agent Entry is to provide a deterministic repository initialization process.

Without a standardized entry procedure, AI agents often:

- read excessive documentation;
- miss critical repository conventions;
- duplicate deterministic work;
- produce inconsistent behavior.

Agent Entry ensures that every agent starts from the same repository understanding.

------

# Scope

Agent Entry applies to every AI agent interacting with a Documentation OS repository.

Examples include:

- implementation agents;
- review agents;
- planning agents;
- documentation agents;
- orchestration agents.

Repository Profiles may provide additional guidance.

They shall not redefine the Agent Entry lifecycle.

------

# Design Principles

Agent Entry follows the following principles.

## AE-1 Repository First

Agents begin with repository guidance rather than repository exploration.

The repository should explain itself before requiring discovery.

------

## AE-2 Deterministic Initialization

Two agents entering the same repository should establish equivalent repository understanding.

Initialization should not depend upon random exploration.

------

## AE-3 Minimal Required Reading

Agent Entry defines the minimum required documentation.

Agents may continue reading beyond this minimum when necessary.

------

## AE-4 Repository Guidance Over Heuristics

Repository conventions should be learned from repository guidance rather than inferred.

------

# Entry Point

Every Single Repository SHALL expose one Agent Entry Document at the repository root.

The Single Repository Profile uses `AGENTS.md` as the canonical source of the Agent Entry Document.

A content-equivalent mirror named `CLAUDE.md` SHALL accompany `AGENTS.md` at the repository root.

Both files are valid entry points; an agent MAY begin from either.

```text
Repository

│

├── AGENTS.md    (canonical source)

└── CLAUDE.md    (content-equivalent mirror)
```

`AGENTS.md` and `CLAUDE.md` SHALL remain content-equivalent.

Any change applied to one SHALL be reflected in the other.

The recommended synchronization mechanism is a symbolic link (`CLAUDE.md` → `AGENTS.md`).

Because symbolic links are not reliably preserved across all platforms and version-control clients, content-equivalence SHALL NOT depend solely on link behavior.

Documentation Engines SHALL verify equivalence through Validation.

------

# Agent Initialization

A compliant agent should initialize using the following sequence.

```text
Repository

↓

Read the Agent Entry Document

↓

Understand Repository Profile

↓

Understand Documentation Architecture

↓

Locate Relevant Knowledge

↓

Begin Task
```

The initialization sequence establishes repository context before execution begins.

------

# AGENTS.md Responsibilities

The root AGENTS.md should define:

- repository overview;
- Documentation OS version;
- active Repository Profile;
- documentation architecture;
- required Documentation Operations;
- repository conventions;
- recommended reading strategy.

AGENTS.md should not duplicate repository Knowledge.

Instead, it should direct agents toward authoritative sources.

------

# Required Repository Understanding

Before executing a task, an agent should understand:

- repository documentation architecture;
- Knowledge and Runtime separation;
- repository conventions;
- Documentation Operations;
- Work lifecycle.

Agents should avoid implementation before establishing this understanding.

------

# Local Agent Guidance

Repository subdirectories may define local Agent Entry files.

Where a subdirectory defines a local Agent Entry file, it SHALL provide both `AGENTS.md` and a content-equivalent `CLAUDE.md` mirror.

Local `AGENTS.md` and `CLAUDE.md` SHALL remain content-equivalent under the same rule that governs the repository root.

Local guidance supplements repository guidance.

Typical hierarchy:

```text
Root AGENTS.md

↓

docs/AGENTS.md

↓

architecture/AGENTS.md
```

Each level narrows repository context.

Lower-level guidance shall not contradict higher-level guidance.

------

# Repository Navigation

After initialization, agents should navigate through documented repository relationships.

Preferred navigation:

```text
AGENTS

↓

Knowledge

↓

Related Knowledge

↓

Runtime (if required)

↓

Implementation
```

Repository-wide search should supplement documented navigation rather than replace it.

------

# Documentation Operations

Agents should invoke Documentation Operations whenever deterministic repository maintenance is required.

Examples include:

- Validation;
- identifier allocation;
- generated artifact updates;
- Runtime archival.

Agents should avoid manually reproducing deterministic operations.

------

# Entry Completion

Agent Entry completes when the agent has established sufficient repository understanding to begin engineering work.

Completion does not require reading the entire repository.

Completion requires understanding repository organization.

------

# Failure Handling

Agent Entry may fail when:

- the Agent Entry Document is missing (neither `AGENTS.md` nor `CLAUDE.md` is present);
- repository structure is inconsistent;
- Repository Profile cannot be determined.

Implementations should report initialization failures explicitly.

------

# Compliance

A Documentation OS-compliant agent SHALL:

- begin with the Agent Entry Document;
- understand the active Repository Profile;
- understand Documentation Operations;
- establish repository context before implementation;
- use repository guidance before repository exploration.

------

# Non-Goals

Agent Entry intentionally does not define:

- reasoning strategy;
- implementation methodology;
- prompt design;
- AI model behavior;
- engineering workflow.

These concerns remain implementation-specific.

------

# References

- DOS-2001 — Single Repository Profile
- DOS-2002 — Repository Layout
- DOS-4001 — Documentation Operations
- DOS-5002 — Reading Strategy
- DOS-5003 — Execution Contract
- DOS-5004 — Documentation Engine

------

# Summary

Agent Entry defines the standardized initialization process for Documentation OS repositories.

Every compliant agent begins by reading the repository's Agent Entry Document (`AGENTS.md`, mirrored as `CLAUDE.md`), establishing repository understanding, learning repository conventions, and locating authoritative Knowledge before performing engineering work.

This ensures deterministic and consistent behavior across different AI agents and implementations.