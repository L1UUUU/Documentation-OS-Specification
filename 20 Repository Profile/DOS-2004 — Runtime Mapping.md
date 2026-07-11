# DOS-2004 — Runtime Mapping

**Status:** Stable
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

Runtime should disappear naturally after Work completion.

Long-lived Runtime indicates incomplete lifecycle management.

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

Each Work maintains the following conceptual responsibilities.

```text
Work

├── Clarified Requirements

├── Implementation Plan

├── Execution Tasks

└── Temporary Notes
```

These responsibilities are conceptual.

Repository layout is an implementation choice.

------

# Clarified Requirements

## Purpose

Clarified Requirements describe the agreed implementation objective.

They are produced after requirement clarification.

Clarified Requirements should represent:

- confirmed scope;
- confirmed constraints;
- accepted assumptions;
- implementation boundaries.

Clarified Requirements become the authoritative execution input.

------

# Implementation Plan

## Purpose

Implementation Plans translate clarified requirements into executable engineering work.

Typical examples include:

- PRDs;
- implementation strategy;
- execution phases.

Implementation Plans exist only for the duration of the Work.

They are Runtime artifacts.

------

# Execution Tasks

## Purpose

Execution Tasks decompose implementation into concrete engineering work.

Typical examples include:

- generated Issues;
- implementation checklist;
- engineering subtasks.

Documentation OS intentionally does not prescribe issue tracking methodology.

Execution Tasks merely provide execution guidance.

------

# Temporary Notes

## Purpose

Temporary Notes capture information useful during implementation but not intended to become permanent repository knowledge.

Examples include:

- debugging observations;
- exploratory experiments;
- discarded implementation ideas;
- temporary references.

Temporary Notes should either:

- disappear;
- or produce Knowledge.

They should not accumulate indefinitely.

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

Archive

↓

Close
```

Runtime Mapping defines representation.

Lifecycle specifications define behavior.

------

# Runtime and Inbox

Inbox is not part of Runtime.

Inbox belongs to the Knowledge domain.

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

However, Runtime shall never become the repository's long-term memory.

------

# Repository Representation

A typical Runtime representation may resemble:

```text
.scratch/

└── WORK-0007/

    ├── requirements/

    ├── plan/

    ├── tasks/

    └── notes/
```

This example illustrates one valid implementation.

Repository Profiles may evolve internal Runtime organization while preserving identical Runtime semantics.

------

# Compliance

A compliant Single Repository implementation SHALL satisfy the following requirements.

- Runtime SHALL exist beneath the Runtime domain.
- Runtime SHALL organize execution around Work.
- Runtime SHALL separate planning from execution.
- Runtime SHALL remain temporary.
- Runtime SHALL support Knowledge Synchronization before archival.

------

# Non-Goals

This specification intentionally does not define:

- PRD templates;
- Issue formats;
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

- Clarified Requirements;
- Implementation Plan;
- Execution Tasks;
- Temporary Notes.

Runtime exists solely to complete implementation and produce new repository knowledge before leaving the active execution domain.