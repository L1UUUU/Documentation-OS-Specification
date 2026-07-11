# DOS-3002 — Runtime Lifecycle

**Status:** Stable
**Version:** 1.0
**Category:** Lifecycle

------

# Abstract

This specification defines the lifecycle of Runtime within Documentation OS.

Unlike the Document Lifecycle, which applies to all managed documentation artifacts, the Runtime Lifecycle describes how a unit of work evolves from an accepted implementation effort into completed repository knowledge.

The Runtime Lifecycle is the execution lifecycle of Documentation OS.

Its purpose is not merely to complete implementation, but to ensure that implementation permanently improves the repository.

------

# Purpose

The Runtime Lifecycle establishes a deterministic lifecycle for every Work.

Every Work should progress through the same conceptual stages regardless of:

- programming language;
- engineering methodology;
- AI tooling;
- repository implementation.

The lifecycle exists to guarantee that temporary execution artifacts eventually leave the active execution context while preserving any knowledge created during implementation.

------

# Scope

This specification applies to every Work managed by the Runtime domain.

Examples include:

- feature development;
- bug fixes;
- refactoring;
- documentation improvements;
- engineering maintenance.

The Runtime Lifecycle begins after implementation work has been accepted.

Observations that have not yet become implementation work belong to the Inbox lifecycle and are outside the scope of this specification.

------

# Lifecycle Principles

The Runtime Lifecycle follows the principles established by:

- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles

In particular:

- Runtime is temporary.
- Runtime produces Knowledge.
- Knowledge must be synchronized before Runtime is archived.
- Runtime shall not become long-term repository storage.

------

# Runtime Lifecycle

Every Work progresses through the following conceptual lifecycle.

```text
Accepted

↓

Clarified

↓

Planned

↓

Executing

↓

Implemented

↓

Knowledge Impact Analysed

↓

Knowledge Synchronized

↓

Validated

↓

Archived

↓

Closed
```

Each transition represents a meaningful engineering event.

No transition should occur implicitly.

------

# State Definitions

## Accepted

A repository observation has been accepted as implementation work.

Runtime begins at this point.

The Work receives:

- ownership;
- identity;
- execution context.

------

## Clarified

The implementation objective has been clarified.

Clarification establishes:

- scope;
- constraints;
- assumptions;
- success criteria.

Clarification provides the execution boundary for the Work.

------

## Planned

Execution planning has completed.

Planning may produce artifacts such as:

- implementation plans;
- PRDs;
- execution strategies;
- engineering breakdowns.

Planning prepares implementation.

Planning is not implementation.

------

## Executing

Engineering work is actively being performed.

Typical activities include:

- coding;
- testing;
- refactoring;
- documentation updates;
- validation during implementation.

Runtime is expected to change frequently during this state.

------

## Implemented

The implementation objective has been completed.

Implementation completion indicates that the engineering work has finished.

It does **not** indicate that the Work has completed.

Additional lifecycle stages remain mandatory.

------

## Knowledge Impact Analysed

The completed implementation has been evaluated to determine whether repository knowledge requires updating.

Typical questions include:

- Does Architecture change?
- Is a new ADR required?
- Do Standards need revision?
- Should Inbox observations be created?

Knowledge Impact Analysis determines what repository knowledge must evolve.

------

## Knowledge Synchronized

Required Knowledge has been updated.

Synchronization transfers newly created understanding from Runtime into the Knowledge domain.

At this point the repository accurately reflects the completed implementation.

------

## Validated

Documentation Operations verify repository consistency.

Validation confirms:

- identifiers;
- references;
- generated artifacts;
- documentation consistency;
- repository integrity.

Validation ensures that synchronization has not introduced structural inconsistencies.

------

## Archived

Runtime artifacts become historical execution records.

Archived Runtime no longer participates in active engineering work.

Repository knowledge now contains the enduring understanding created by the Work.

------

## Closed

The Work has completed successfully.

Closure confirms that:

- implementation is complete;
- Knowledge has been synchronized;
- repository validation has succeeded;
- Runtime has been archived.

Only at this point is the Work considered finished.

------

# Lifecycle Transitions

The Runtime Lifecycle permits only forward progression.

```text
Accepted

↓

Clarified

↓

Planned

↓

Executing

↓

Implemented

↓

Knowledge Impact Analysed

↓

Knowledge Synchronized

↓

Validated

↓

Archived

↓

Closed
```

Implementations should avoid reversing lifecycle states.

If significant additional engineering work becomes necessary after closure, a new Work should be created.

------

# Runtime Artifacts

Runtime artifacts participate in the lifecycle together with their Work.

Examples include:

- clarified requirements;
- implementation plans;
- execution tasks;
- temporary notes.

Artifacts may evolve while the Work remains active.

They become immutable after archival.

------

# Runtime Ownership

Ownership remains attached to the Work throughout the lifecycle.

Ownership includes responsibility for:

- maintaining Runtime artifacts;
- initiating Knowledge Impact Analysis;
- synchronizing repository knowledge;
- completing validation;
- archiving Runtime.

Ownership concludes when the Work reaches the Closed state.

------

# Runtime and Knowledge

Runtime continuously consumes Knowledge.

Knowledge should guide implementation.

Implementation should improve Knowledge.

The relationship is therefore cyclical.

```text
Knowledge

↓

Runtime

↓

Implementation

↓

Knowledge Synchronization

↓

Knowledge
```

This continuous feedback loop forms the central operating model of Documentation OS.

------

# Runtime Invariants

The following invariants shall always remain true.

## RI-1

Every Work possesses exactly one lifecycle state.

------

## RI-2

Runtime shall not bypass Knowledge Impact Analysis.

------

## RI-3

Knowledge Synchronization shall precede Runtime archival.

------

## RI-4

Archived Runtime shall remain immutable.

------

## RI-5

Closed Work shall not return to Executing.

Subsequent implementation requires a new Work.

------

# Failure Handling

A Work may fail before completion.

Typical failure causes include:

- implementation cancelled;
- requirements withdrawn;
- engineering superseded.

Failed Works may be archived without reaching the Closed state.

Repository Profiles may define additional failure handling conventions.

Failure shall not bypass required documentation consistency operations.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every Work follows the Runtime Lifecycle;
- lifecycle transitions are explicit;
- Knowledge Impact Analysis occurs after implementation;
- Knowledge Synchronization precedes archival;
- archived Runtime remains immutable.

------

# Non-Goals

This specification intentionally does not define:

- implementation methodology;
- project management workflow;
- sprint planning;
- issue tracking systems;
- engineering estimation.

These concerns remain outside Documentation OS.

------

# References

- DOS-1003 — Runtime Model
- DOS-2004 — Runtime Mapping
- DOS-3001 — Document Lifecycle
- DOS-3003 — Knowledge Impact Analysis
- DOS-3004 — Work Close Pipeline
- DOS-4002 — Validation

------

# Summary

The Runtime Lifecycle defines how every Work evolves from accepted implementation to completed repository knowledge.

Implementation alone does not complete a Work.

A Work is complete only after:

- Knowledge has been analysed;
- Knowledge has been synchronized;
- the repository has been validated;
- Runtime has been archived.

This lifecycle ensures that every implementation leaves the repository with greater long-term understanding than before the Work began.