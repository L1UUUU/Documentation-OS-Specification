# DOS-3004 — Work Close Pipeline

**Status:** Draft
**Version:** 1.0
**Category:** Lifecycle

------

# Abstract

This specification defines the Work Close Pipeline of Documentation OS.

The Work Close Pipeline is the mandatory sequence of operations required to bring a Work to its terminal state.

Implementation completion alone does not complete a Work.

A Work reaches the **Completed** terminal state only after repository Knowledge has been synchronized, repository consistency has been validated, and the Complete stage has moved the Work to `completed/` and released lifecycle ownership. The Cleanup stage then finalizes repository state; it is mandatory but may be retried independently without affecting the already-established Completed state.

The Work Close Pipeline guarantees that every completed Work permanently improves repository quality.

------

# Purpose

The purpose of the Work Close Pipeline is to provide a deterministic completion process for every Work.

Without a standardized closing procedure, repositories gradually accumulate:

- outdated Runtime artifacts;
- unsynchronized documentation;
- inconsistent references;
- stale planning documents.

The Work Close Pipeline prevents these forms of repository decay.

------

# Scope

This specification applies to every Runtime Work managed by Documentation OS.

The pipeline begins after a Work's implementation or execution activity has ended — achieved, abandoned, or superseded. A Work terminated before implementation began still enters the pipeline, with Knowledge Impact Analysis typically producing no Knowledge impact (DOS-3002, DOS-3003). A Work reaches its **Completed** terminal lifecycle state at the Complete stage; the subsequent Cleanup stage finalizes repository state and is mandatory but independently retriable.

------

# Lifecycle Position

The Work Close Pipeline occupies the final phase of the Runtime Lifecycle. Implementation and Knowledge Impact Analysis are preceding lifecycle activity that establishes the context for the pipeline; they are not pipeline stages. The pipeline itself consists of the four stages defined in Pipeline Overview.

Preceding activity:

```text
Implementation

↓

Knowledge Impact Analysis
```

Pipeline stages:

```text
Knowledge Synchronization

↓

Validation

↓

Complete

↓

Cleanup
```

Every completed Work SHALL pass through Knowledge Impact Analysis and all four pipeline stages; Implementation MAY be partial or absent for a Work terminated before its objectives were reached (DOS-3002, DOS-3003).

Stages SHALL NOT be skipped.

The abstract Complete stage maps to active/<workstream-slug>/ → completed/<workstream-slug>/ directory movement in the Single Repository implementation.

------

# Pipeline Overview

The Work Close Pipeline consists of four sequential stages.

| Stage                     | Responsibility                |
| ------------------------- | ----------------------------- |
| Knowledge Synchronization | Update persistent Knowledge   |
| Validation                | Verify repository consistency |
| Complete                  | Complete Runtime artifacts     |
| Cleanup                   | Finalize Runtime state        |

Each stage depends upon successful completion of the previous stage.

------

# Stage 1 — Knowledge Synchronization

## Purpose

Synchronize the repository's persistent Knowledge using the output of Knowledge Impact Analysis.

Typical activities include:

- updating Architecture;
- creating or updating ADRs;
- updating Standards.

Knowledge Synchronization modifies the Knowledge domain.

Unresolved concerns discovered during synchronization are recorded as Inbox staging items; these do not modify Knowledge.

Runtime remains active during this stage.

------

# Stage 2 — Validation

## Purpose

Verify repository consistency.

Validation includes, where applicable:

- identity integrity;
- reference integrity;
- generated artifacts;
- managed regions;
- documentation consistency.

Validation confirms repository correctness before Runtime leaves the active domain.

If validation fails, the Work SHALL NOT proceed to Complete.

------

# Stage 3 — Complete

## Purpose

Move Runtime from active execution into completed state.

The Complete stage receives the Work's terminal `outcome` as an input to the Complete Operation — an engineering judgement supplied by the Agent or caller (DOS-4001). It performs the atomic transition of the Work to its final lifecycle state.

Deterministic actions:

1. Preflight — verify Core Runtime Assets exist (PRD.md, at least one `issues/NN-<slug>.md`, HANDOFF.md); verify a legal `outcome` value was supplied; and when `outcome = succeeded`, verify no Issue remains `open`, `in-progress`, or `blocked`. A non-`succeeded` Work MAY leave non-terminal Issues. Preflight failure aborts the stage with the Work unchanged under `active/`.
2. Record the supplied terminal `outcome` (`succeeded`, `cancelled`, `superseded`, or `failed`) in PRD front matter (see DOS-3002);
3. Clean Ephemeral Runtime Content only (temporary notes, planning artifacts);
4. Atomically move active/<workstream-slug>/ → completed/<workstream-slug>/;
5. Preserve Core Runtime Assets unchanged — apart from the deterministic addition of the terminal `outcome` to PRD front matter in step 2, Core Runtime Assets SHALL remain unchanged (immutable business content; generated INDEX.md MAY be regenerated).

Steps 2–4 form a single transaction: if the directory movement fails, the `outcome` write and the Ephemeral cleanup SHALL be rolled back so that the Work remains valid under `active/` (an active Work SHALL NOT carry an `outcome`).

Note: HANDOFF.md is generated at Work creation (DOS-2004) and is therefore always present; preflight verifies integrity rather than first-time presence. For a non-`succeeded` Work, the recorded `outcome` truthfully reflects that the implementation objectives were not achieved; repository consistency verification SHALL NOT be bypassed regardless of outcome.

Stage 3 owns the Work directory movement. Apart from the deterministic addition of the terminal `outcome` to PRD front matter, Core Runtime Assets SHALL remain unchanged by completion.

The directory movement to completed/<workstream-slug>/ represents the Work's entry into the Completed lifecycle state—the terminal state of the Work lifecycle. Ownership is released at this point.

Repository Knowledge becomes the authoritative project memory.

------

# Stage 4 — Cleanup

## Purpose

Finalize Runtime after completion.

Cleanup operates after the Work has already reached its Completed terminal state. Cleanup is idempotent and MAY be retried independently if it fails; failure does not affect the Work's already-established Completed state.

Typical cleanup activities include:

- regenerating .scratch/INDEX.md (generated summaries MAY be refreshed);
- cleaning up external temporary state and caches;
- updating repository indexes.

Cleanup SHALL NOT move the Work directory; directory movement is owned by Stage 3 — Complete.

Cleanup prepares the repository for subsequent Work.

Cleanup SHALL NOT modify synchronized Knowledge.

------

# Pipeline Invariants

The following invariants SHALL always remain true.

## WC-1

Knowledge Synchronization precedes Validation.

------

## WC-2

Validation precedes Complete.

------

## WC-3

Complete precedes Cleanup.

------
## WC-4

Pipeline stages SHALL reach successful completion at most once for each Work. Failed or interrupted attempts MAY be retried idempotently.

------

## WC-5

Pipeline stages execute in order.

Reordering stages is prohibited.

------

# Failure Handling

A pipeline stage can fail.

Examples include:

- validation errors;
- synchronization conflicts;
- documentation inconsistencies;
- repository operation failures.

If a stage fails before Complete:

- the Work SHALL remain in active/;
- subsequent stages SHALL NOT execute;
- failure SHALL be explicitly reported.

If a stage fails after Complete (the Work directory has already moved to completed/):

- the Work SHALL remain in completed/;
- the Work has already reached its Completed terminal state;
- after recovery, the remaining Cleanup stage MAY be retried idempotently;
- failure SHALL be explicitly reported.

In both cases the failure is recoverable, and the pipeline MAY resume after the failure has been resolved.

------

# Idempotency

Pipeline stages SHOULD be designed to support safe re-execution.

Re-running a partially completed Work Close Pipeline SHOULD NOT introduce inconsistent repository state.

Documentation Operations SHOULD therefore favor deterministic behavior and avoid destructive side effects.

------

# Documentation Operations

Documentation Operations MAY automate portions of the Work Close Pipeline.

Typical operations include:

- synchronizing generated indexes;
- validating references;
- refreshing managed regions;
- completing Runtime;
- updating repository summaries.

Documentation Operations maintain structural consistency.

Engineering decisions remain the responsibility of humans or AI agents.

------

# Repository Effects

Successful completion of the Work Close Pipeline guarantees:

- Runtime has left the active execution domain;
- repository Knowledge reflects the completed or terminated Work activity, including an explicit no-change result where applicable;
- repository consistency has been verified;
- lifecycle ownership has concluded.

The repository is therefore immediately ready for the next Work.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every completed Work executes the Work Close Pipeline;
- pipeline stages occur in the defined order;
- Runtime is not completed before Knowledge Synchronization;
- repository validation succeeds before Work completion;
- Work ownership concludes when the Complete stage moves the Work to `completed/`; the subsequent Cleanup stage is mandatory and may be retried independently.

------

# Non-Goals

This specification intentionally does not define:

- implementation methodology;
- engineering review process;
- CI/CD pipelines;
- version control workflow;
- deployment procedures.

These concerns remain outside the scope of Documentation OS.

------

# References

- DOS-3002 — Runtime Lifecycle
- DOS-3003 — Knowledge Impact Analysis
- DOS-3005 — Ownership
- DOS-4001 — Documentation Operations
- DOS-4002 — Validation

------

# Summary

The Work Close Pipeline defines the mandatory completion process for every Runtime Work.

A Work is not complete when implementation ends.

A Work reaches the Completed terminal state when:

- Knowledge has been synchronized;
- repository validation has succeeded;
- the Complete stage has moved the Work to `completed/` and released Ownership.

The subsequent Cleanup stage finalizes repository state (regenerating `.scratch/INDEX.md`); it is mandatory but idempotent and may be retried independently without affecting the already-established Completed state.

This pipeline ensures that every completed Work leaves the repository in a consistent, maintainable, and knowledge-rich state.