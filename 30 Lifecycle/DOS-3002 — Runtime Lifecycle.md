# DOS-3002 — Runtime Lifecycle

**Status:** Draft
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
- Knowledge MUST be synchronized before Runtime is completed.
- Runtime SHALL NOT become long-term repository storage.

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

Completed

↓

Closed
```

Each transition represents a meaningful engineering event.

No transition SHOULD occur implicitly.

------

# State Definitions

## Accepted

A repository observation has been accepted as implementation work.

Runtime begins at this point.

The Work receives:

- ownership;
- workstream slug and active/<slug>/ path (Runtime location, not global identity);
- execution context;
- an initial HANDOFF.md (MAY be empty), generated when the workspace is created.

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

## Completed

Runtime artifacts become historical execution records.

Completed Runtime no longer participates in active engineering work.

The workstream directory moves from active/<slug>/ to completed/<slug>/.

Core Runtime assets are preserved (immutable business content; generated INDEX.md MAY be regenerated).

Repository knowledge now contains the enduring understanding created by the Work.

------

## Closed

The Work has completed successfully.

Closure confirms that:

- implementation is complete;
- Knowledge has been synchronized;
- repository validation has succeeded;
- Runtime has been completed (moved to completed/<slug>/).

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

Completed

↓

Closed
```

Implementations SHOULD avoid reversing lifecycle states.

If significant additional engineering work becomes necessary after closure, a new Work SHOULD be created.

------

# Runtime Artifacts

Runtime artifacts participate in the lifecycle together with their Work.

Core Runtime Assets (preserved upon completion):

- PRD.md;
- issues/*.md;
- HANDOFF.md.

Ephemeral Runtime Content (MAY be cleaned up):

- clarified requirements;
- implementation plans;
- execution tasks;
- temporary notes.

Artifacts may evolve while the Work remains active.

Core assets are preserved upon completion (immutable business content; generated INDEX.md MAY be regenerated).

------

# Runtime Ownership

Ownership remains attached to the Work throughout the lifecycle.

Ownership includes responsibility for:

- maintaining Runtime artifacts;
- initiating Knowledge Impact Analysis;
- synchronizing repository knowledge;
- completing validation;
- completing Runtime.

Ownership concludes when the Work reaches the Closed state.

Completed Runtime core assets are immutable. To correct a historical error in a completed Work, do not edit the completed Work; instead create a new Work whose relationships point to the original completed Work. Only INDEX.md MAY be regenerated; completed core content SHALL NOT be modified.

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

The following invariants SHALL always remain true.

## RI-1

Every Work possesses exactly one lifecycle state.

------

## RI-2

Runtime SHALL NOT bypass Knowledge Impact Analysis.

------

## RI-3

Knowledge Synchronization SHALL precede Runtime completion.

------

## RI-4

Completed Runtime core assets SHALL be preserved (immutable business content; the generated INDEX.md MAY be regenerated).

------

## RI-5

Closed Work SHALL NOT return to Executing.

Subsequent implementation requires a new Work.

------

# Failure Handling

A Work may fail before completion.

Typical failure causes include:

- implementation cancelled;
- requirements withdrawn;
- engineering superseded.

Failed Works MAY be completed without reaching the Closed state.

Repository Profiles MAY define additional failure handling conventions.

Failure SHALL NOT bypass required documentation consistency operations.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every Work follows the Runtime Lifecycle;
- lifecycle transitions are explicit;
- Knowledge Impact Analysis occurs after implementation;
- Knowledge Synchronization precedes completion;
- completed Runtime core assets remain preserved.

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
- Runtime has been completed.

This lifecycle ensures that every implementation leaves the repository with greater long-term understanding than before the Work began.