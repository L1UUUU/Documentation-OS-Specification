# DOS-3002 — Runtime Lifecycle

**Status:** Draft
**Version:** 1.0
**Category:** Lifecycle

------

# Abstract

This specification defines the lifecycle of Runtime within Documentation OS.

Unlike the Document Lifecycle, which applies to all managed documentation artifacts, the Runtime Lifecycle describes how a unit of work progresses from an accepted implementation effort to a synchronized terminal state, updating repository knowledge or recording an explicit no-change result.

The Runtime Lifecycle is the execution lifecycle of Documentation OS.

Its purpose is to ensure that ended Work activity is analysed, that affected Knowledge is synchronized, and that an explicit no-change result is recorded when no Knowledge update is required.

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
- Runtime completes Knowledge Synchronization, which either updates Knowledge or records an explicit no-change result.
- Knowledge MUST be synchronized before Runtime is completed.
- Runtime SHALL NOT become long-term repository storage.

------

# Runtime Lifecycle

Every Work possesses exactly one observable lifecycle state — Active or Completed — determined by its directory location under `.scratch/` (`active/` or `completed/`).

While a Work is Active, it progresses through the following conceptual workflow phases. These phases describe engineering progress; they are not separately persistable states and are not individually observable from repository structure alone.

```text
Accepted

↓

Clarified

↓

Planned

↓

Executing

↓

Implemented  (MAY be partial or absent for a non-succeeded outcome)

↓

Knowledge Impact Analysed

↓

Knowledge Synchronized

↓

Validated
```

The transition from Active to Completed occurs through the Work Close Pipeline (see DOS-3004): the Work directory moves from `active/<workstream-slug>/` to `completed/<workstream-slug>/`, which is the terminal observable state.

Each phase transition represents a meaningful engineering event.

No transition SHOULD occur implicitly.

------

# Workflow Phases and Terminal State

The following phases describe an Active Work's engineering progress. They are conceptual phases rather than separately persistable lifecycle states; an Active Work is in the Active observable state throughout, regardless of which phase it has reached.

`Completed` is the terminal observable state, reached via the Work Close Pipeline (DOS-3004).

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

For a non-succeeded outcome (`cancelled`, `superseded`, `failed`), a Work MAY enter Knowledge Impact Analysis from any earlier active phase, and implementation completion is not required (see Failure Handling; DOS-1003, DOS-3004).

------

## Knowledge Impact Analysed

The Work's implementation or execution activity has been evaluated to determine whether repository knowledge requires updating.

Typical questions include:

- Does Architecture change?
- Is a new ADR required?
- Do Standards need revision?
- Should Inbox observations be created?

Knowledge Impact Analysis determines what repository knowledge must evolve, or records an explicit no-change result where no Knowledge edits are required.

------

## Knowledge Synchronized

Required Knowledge has been updated, or an explicit no-change result has been recorded where no Knowledge edits were required.

Synchronization transfers newly created understanding from Runtime into the Knowledge domain.

At this point the repository accurately reflects the Work's completed or terminated activity, including any no-change result.

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

The Work has reached its terminal observable lifecycle state, with its terminal `outcome` (`succeeded`, `cancelled`, `superseded`, or `failed`) recorded in PRD front matter.

The workstream directory has moved from active/<slug>/ to completed/<slug>/.

Completed Runtime no longer participates in active engineering work. Ownership concludes at this point.

Core Runtime assets are preserved (immutable business content; only the terminal `outcome` is added to PRD front matter at completion, and the generated INDEX.md MAY be regenerated).

Repository knowledge now contains the enduring understanding created by the Work.

------

# Lifecycle Transitions

Workflow phases progress only forward.

```text
Accepted

↓

Clarified

↓

Planned

↓

Executing

↓

Implemented  (MAY be partial or absent for a non-succeeded outcome)

↓

Knowledge Impact Analysed

↓

Knowledge Synchronized

↓

Validated

↓

(Work Close Pipeline: active/ → completed/)

↓

Completed
```

The Active → Completed transition is the only observable state transition and is performed by the Work Close Pipeline (DOS-3004).

Implementations SHOULD avoid reversing workflow phases.

If significant additional engineering work becomes necessary after completion, a new Work SHOULD be created.

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

Ownership concludes when the Work reaches the Completed state.

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

Every Work possesses exactly one observable lifecycle state — Active or Completed — determined by its directory location under `.scratch/`.

------

## RI-2

Runtime SHALL NOT bypass Knowledge Impact Analysis.

------

## RI-3

Knowledge Synchronization SHALL precede Runtime completion.

------

## RI-4

Completed Runtime core assets SHALL be preserved (immutable business content; only the terminal `outcome` is added to PRD front matter at completion, and the generated INDEX.md MAY be regenerated).

------

## RI-5

Completed Work SHALL NOT return to Executing.

Subsequent implementation requires a new Work.

------

# Failure Handling

A Work may fail before completion.

Typical failure causes include:

- implementation cancelled;
- requirements withdrawn;
- engineering superseded.

`Completed` denotes termination, not necessarily success. A Work reaches `Completed` with a terminal `outcome` recorded in its PRD front matter. Normative `outcome` values are:

- `succeeded` — implementation objectives achieved and Knowledge synchronized;
- `cancelled` — requirements withdrawn or implementation abandoned;
- `superseded` — engineering approach replaced by another Work;
- `failed` — implementation could not be completed.

A non-`succeeded` Work MAY reach `Completed` without having achieved its implementation objectives, provided its `outcome` truthfully records the termination reason.

A Work terminated before reaching the `Implemented` phase (for example requirements withdrawn before coding began, or the approach superseded by another Work) MAY still enter the Work Close Pipeline from its current active phase: Knowledge Impact Analysis evaluates the impact of whatever activity did occur (frequently no impact), and the pipeline proceeds to Complete with a non-`succeeded` `outcome`. The workflow phases preceding `Implemented` are therefore not prerequisites for termination; only Knowledge Impact Analysis → Knowledge Synchronization → Validation → Complete is mandatory (DOS-3003, DOS-3004).

Failed Works MAY reach the Completed state.

Repository Profiles MAY define additional failure handling conventions.

Failure SHALL NOT bypass required documentation consistency operations: a non-`succeeded` Work SHALL still record its `outcome` and SHALL NOT leave the repository in an inconsistent state.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every Work follows the Runtime Lifecycle;
- lifecycle transitions are explicit;
- Knowledge Impact Analysis occurs after the Work's implementation or execution activity has ended;
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

The Runtime Lifecycle defines how every Work progresses from accepted implementation to a synchronized terminal state.

Implementation alone does not complete a Work.

A Work is complete only after:

- Knowledge has been analysed;
- Knowledge has been synchronized;
- the repository has been validated;
- Runtime has been completed.

This lifecycle ensures that every Work leaves the repository in a synchronized state, whether through updated Knowledge or an explicit no-change result.