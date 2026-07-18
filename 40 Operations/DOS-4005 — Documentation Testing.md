# DOS-4005 — Documentation Testing

**Status:** Draft
**Version:** 1.0
**Category:** Operations

------

# Abstract

This specification defines Documentation Testing within Documentation OS.

Documentation Testing verifies that a Documentation OS implementation behaves correctly according to the specifications.

Unlike Validation, which evaluates the state of a repository, Documentation Testing evaluates the behavior of the Documentation Engine itself.

Its objective is to ensure that Documentation Operations remain deterministic, reliable, and specification-compliant.

------

# Purpose

The purpose of Documentation Testing is to verify that Documentation OS implementations correctly execute the specifications.

Documentation Testing answers questions such as:

- Does the Documentation Engine allocate identifiers correctly?
- Does Validation detect required failures?
- Does Migration preserve repository semantics?
- Are Documentation Operations deterministic?

Documentation Testing validates implementation behavior rather than repository state.

------

# Scope

Documentation Testing applies to every Documentation Engine implementation.

Typical testing targets include:

- Documentation Operations
- Repository Profile implementations
- Validation
- Migration
- Health
- Lifecycle operations

Documentation Testing does not evaluate engineering quality or software functionality outside Documentation OS.

------

# Design Principles

Documentation Testing follows the following principles.

## DT-1 Specification-Driven

Every test SHALL verify behavior defined by a Documentation OS specification.

Tests should never depend upon undocumented implementation details.

------

## DT-2 Deterministic

Executing identical tests against identical repository state SHALL produce identical results.

Tests SHALL avoid timing-dependent behavior.

------

## DT-3 Isolated

Each test should execute independently.

One failing test SHALL NOT invalidate unrelated tests.

------

## DT-4 Repeatable

Tests should be safely repeatable.

Repeated execution SHALL NOT modify repository semantics.

------

## DT-5 Observable

Test results SHALL clearly identify:

- tested specification;
- expected behavior;
- observed behavior;
- pass/fail status.

------

# Testing Categories

Documentation OS defines six normative testing categories.

```text
Documentation Testing

├── Operation Tests

├── Validation Tests

├── Lifecycle Tests

├── Migration Tests

├── Repository Profile Tests

└── Integration Tests
```

Implementations may define additional testing categories.

------

# Operation Tests

Operation Tests verify individual Documentation Operations.

Examples include:

- identifier allocation;
- Generate Issue numbering, idempotency, conflict detection, concurrency, and
  transactional INDEX behavior;
- generated content updates;
- relationship synchronization;
- complete operations (active→completed move);
- .scratch/INDEX.md regeneration determinism.

Each operation should be tested independently.

------

# Validation Tests

Validation Tests verify that Validation correctly detects both valid and invalid repository states.

Typical cases include:

- duplicate identities;
- invalid lifecycle transitions;
- broken relationships;
- inconsistent generated artifacts.

Validation Tests should include both positive and negative scenarios.

------

# Lifecycle Tests

Lifecycle Tests verify lifecycle behavior.

Typical examples include:

- valid lifecycle progression;
- prohibited transitions;
- Work Close Pipeline ordering (including 4-phase sequence: Knowledge Synchronization → Validation → Complete → Cleanup);
- ownership transitions.

Lifecycle behavior should remain specification compliant.

------

# Migration Tests

Migration Tests verify repository evolution.

Typical checks include:

- semantic preservation;
- identity preservation;
- relationship preservation;
- successful post-migration Validation.

Migration Tests should verify recoverability where supported.

------

# Repository Profile Tests

Repository Profile Tests verify compliance with the active Repository Profile.

Examples include:

- required repository entry points;
- required directory layout;
- required documentation locations;
- profile-specific conventions;
- .scratch structure assertions (.scratch/AGENTS.md + .scratch/CLAUDE.md + .scratch/active/ + .scratch/completed/ + .scratch/INDEX.md + each Work containing PRD.md, issues/NN-<slug>.md, HANDOFF.md).

Repository Profile Tests should remain independent of implementation language.

------

# Integration Tests

Integration Tests verify cooperation between multiple Documentation Operations.

Typical scenarios include:

```text
Synchronize

↓

Validate

↓

Complete (active→completed + INDEX regeneration)

↓

Health
```

Integration Tests ensure that independent operations function correctly as a complete Documentation OS implementation.

------

# Test Fixtures

Documentation Testing should use deterministic repository fixtures.

Fixtures should represent realistic repository states.

Typical fixtures include:

- empty repository;
- healthy repository;
- invalid repository;
- partially synchronized repository;
- migrated repository;
- Runtime Profile conformance fixtures:
  - repository with active Work(s);
  - repository with completed Work(s) (verifying Core Assets preserved in completed/<slug>/);
  - repository with regenerated INDEX.md (verifying reproducibility);
  - repository violating Runtime Profile (negative: missing PRD/HANDOFF, incorrect WORK-NNNN/work.yaml usage, presence of archive/);
  - normative conformance scenarios — see Normative Conformance Scenarios below.

Fixtures should remain version-controlled.

------

# Normative Conformance Scenarios

The following scenarios SHALL be covered by the conformance tests required in the Compliance section. Each scenario follows the Expected Results structure: initial repository state, operation under test, expected repository state, and expected diagnostics.

## Scenario 1 — Active Work with empty `issues/` is valid

- Initial state: a repository with one active Work under `.scratch/active/<slug>/` containing `PRD.md`, an empty `issues/` directory, and `HANDOFF.md`.
- Operation: run Validation.
- Expected state: repository unchanged (Validation is read-only).
- Expected diagnostics: Passed — an active Work MAY have an empty `issues/` directory.

## Scenario 2 — Work entering Complete with empty `issues/` is invalid

- Initial state: an active Work whose `issues/` directory is empty.
- Operation: invoke the Complete operation.
- Expected state: the Work remains in `active/`; no directory movement occurs.
- Expected diagnostics: Failure — the `issues/` directory SHALL contain at least one Issue file before Complete.

## Scenario 3 — `generate work` leaves `.scratch/INDEX.md` consistent

- Initial state: a repository with N existing Works; `.scratch/INDEX.md` reflects them.
- Operation: run `generate work <new-slug>` for a new, globally unique slug.
- Expected state: `.scratch/active/<new-slug>/` is created with `PRD.md`, `issues/`, and `HANDOFF.md`; `.scratch/INDEX.md` is regenerated and lists the new Work under "Active Works".
- Expected diagnostics: Success; INDEX.md is reproducible from the contents of `active/` and `completed/`.

## Scenario 4 — `generate work` rolls back on mid-creation failure

- Initial state: a repository in which a creation step cannot succeed (for example the requested slug already exists under `active/` or `completed/`).
- Operation: run `generate work <slug>`.
- Expected state: no partial `.scratch/active/<slug>/` remains; `.scratch/INDEX.md` is identical to its pre-operation state.
- Expected diagnostics: Failure; the report states that the workspace and the INDEX regeneration were rolled back.

## Scenario 5 — Complete succeeds, Cleanup fails; Work stays Completed and Cleanup is retriable

- Initial state: an active Work that satisfies all Complete preconditions.
- Operation: invoke Complete; the Complete stage succeeds (the Work moves to `completed/` and its terminal `outcome` is recorded) but the Cleanup stage fails.
- Expected state: the Work resides in `.scratch/completed/<slug>/` with Core Assets preserved and a legal `outcome` in PRD front matter; `.scratch/INDEX.md` is not yet regenerated.
- Re-operation: re-invoke Complete; execution resumes from the Cleanup stage only (no repeated directory movement); `.scratch/INDEX.md` is regenerated.
- Expected diagnostics: the first invocation returns Failure with a recovery instruction to retry only the idempotent Cleanup stage; the re-invocation returns Success.

## Scenario 6 — A failed or interrupted pipeline stage is retried idempotently

- Initial state: a Work in the Work Close Pipeline where an earlier stage (for example Validation) has failed.
- Operation: re-invoke the pipeline after the cause of failure has been resolved.
- Expected state: the Work progresses from the failed stage; no duplicated Knowledge edits and no inconsistent intermediate state are introduced.
- Expected diagnostics: Success on retry; the resulting repository state is equivalent to a single uninterrupted run.

## Scenario 7 — `completed/` location corresponds to the Completed terminal state

- Initial state: a Work that has reached Completed.
- Operation: run Validation / inspect `.scratch/`.
- Expected state: the Work directory is located under `.scratch/completed/<slug>/` (not `active/`); the PRD declares a legal `outcome`; Core Runtime Assets are preserved.
- Expected diagnostics: Validation passes; the Runtime location-presence check confirms the Completed terminal state.

## Scenario 8 — Active Work declaring an `outcome` is invalid

- Initial state: an active Work whose PRD front matter declares `outcome: succeeded`.
- Operation: run Validation.
- Expected state: repository unchanged (Validation is read-only).
- Expected diagnostics: Failure — an active Work SHALL NOT declare an `outcome`.

## Scenario 9 — `outcome=succeeded` with a non-terminal Issue fails Complete preflight

- Initial state: an active Work with at least one Issue, one of which has `status: in-progress`; the PRD carries no `outcome`.
- Operation: invoke the Complete operation with `outcome=succeeded`.
- Expected state: the Work remains in `active/`; no `outcome` is written; no directory movement occurs.
- Expected diagnostics: Failure — preflight rejected a `succeeded` outcome because a non-terminal Issue remains.

## Scenario 10 — A `cancelled` Work reaches Completed through the standard pipeline

- Initial state: an active Work whose requirements were withdrawn before its implementation objectives were achieved; the PRD carries no `outcome`; Knowledge Impact Analysis produced no Knowledge impact; the no-change Knowledge Synchronization stage has completed (recording that no Knowledge edits were required); repository Validation has passed; all Complete preconditions are satisfied.
- Operation: invoke the Complete operation with `outcome=cancelled`.
- Expected state: the Work moves to `.scratch/completed/<slug>/`; the PRD front matter records `outcome: cancelled`; Core Runtime Assets are preserved; `.scratch/INDEX.md` is regenerated (Cleanup) and surfaces the `cancelled` outcome.
- Expected diagnostics: Success.

## Scenario 11 — Complete rolls back if directory movement fails after writing `outcome`

- Initial state: an active Work satisfying all Complete preconditions.
- Operation: invoke the Complete operation with `outcome=succeeded`; the directory movement step fails.
- Expected state: the Work remains valid under `.scratch/active/<slug>/`; the PRD carries no `outcome` (the write was rolled back); any Ephemeral cleanup is rolled back where required.
- Expected diagnostics: Failure with recovery guidance; the Work remains valid for retry.

## Scenario 12 — INDEX surfaces the terminal `outcome` for Completed Works

- Initial state: a repository with one Completed Work whose PRD declares `outcome: superseded`.
- Operation: regenerate `.scratch/INDEX.md`.
- Expected state: INDEX lists the Completed Work with its `outcome` (`superseded`), alongside its slug, PRD path, HANDOFF path, and Issues.
- Expected diagnostics: Success; INDEX is reproducible from repository state alone.

## Scenario 13 — Concurrent ADR drafts receive distinct final numbers at integration

- Initial state: two parallel Works each carry a draft ADR using a non-final placeholder name; neither pre-allocates a final `ADR-NNNN`.
- Operation: integrate both Works; the Documentation Engine allocates final numbers against the integration target state.
- Expected state: each ADR receives a distinct, atomically allocated `ADR-NNNN`; any collision from concurrent integration is resolved by re-numbering the later integration and regenerating its managed references.
- Expected diagnostics: Success; no duplicate or reused ADR numbers.

## Scenario 14 — ADR draft allocation, reference regeneration, and Validation gate

- Initial state: a Work carries a draft ADR using a non-final placeholder name; Identity Validation treats it as exempt.
- Operation: integrate the Work; the Documentation Engine allocates a final ADR-NNNN, regenerates managed references, then runs Identity Validation and Relationship Validation.
- Expected state: the artifact carries its final identifier; all managed references resolve; both validations pass.
- Expected diagnostics: Success. If either validation fails after allocation, the integration is rolled back or corrected before retrying.

## Scenario 15 — Generate Issue allocates the next Work-local number without reusing gaps

- Initial state: one Active Work has an empty `issues/` directory; a second
  Active Work contains `01-first.md` and `03-third.md`.
- Operation: invoke Generate Issue once for each Work with valid, distinct
  inputs.
- Expected state: the first Work receives `01-<slug>.md`; the second receives
  `04-<slug>.md`; each file contains the supplied title, DOS lifecycle status,
  and body; `.scratch/INDEX.md` lists both new Issues.
- Expected diagnostics: Success; each result exposes its number, name, path,
  and `created=true`.

## Scenario 16 — Generate Issue retry is idempotent and conflicting reuse is rejected

- Initial state: an Active Work already contains an Issue created from a known
  Work slug, Issue slug, title, status, and body; INDEX is consistent.
- Operation: repeat the identical logical input, then repeat the same Issue slug
  with a different title, status, or body.
- Expected state: the identical retry neither writes nor allocates another
  number and returns the existing Issue; the conflicting retry leaves the Issue
  and INDEX byte-for-byte unchanged.
- Expected diagnostics: the identical retry succeeds with `created=false`; the
  conflicting retry returns a deterministic conflict.

## Scenario 17 — Concurrent Generate Issue calls allocate safely

- Initial state: an Active Work with a valid Runtime structure and an empty
  `issues/` directory.
- Operation: concurrently invoke Generate Issue with multiple distinct valid
  slugs, including at least two concurrent invocations of one identical logical
  input.
- Expected state: each distinct input has exactly one Issue, all allocated
  numbers are unique, the identical calls converge on one Issue, and INDEX is
  consistent with the final `issues/` directory.
- Expected diagnostics: all non-conflicting calls succeed; exactly one result
  for each logical input reports `created=true`, and identical followers report
  `created=false`.

## Scenario 18 — Generate Issue rejects Completed or invalid Work state

- Initial state: one Work exists only under `.scratch/completed/`; another slug
  is missing; a third appears under both `active/` and `completed/`.
- Operation: invoke Generate Issue for each slug.
- Expected state: no Issue or INDEX change occurs in any case.
- Expected diagnostics: deterministic failures distinguish immutable Completed
  Work, missing Work, and inconsistent repository state.

## Scenario 19 — Generate Issue and INDEX publication are transactional

- Initial state: an Active Work with a consistent INDEX; arrange for INDEX
  publication to fail after the Issue has been staged or published.
- Operation: invoke Generate Issue, then remove the failure and retry the same
  logical input.
- Expected state: after the returned failure, neither the new Issue nor an INDEX
  change remains; after retry, exactly one Issue exists and INDEX lists it.
- Expected diagnostics: the first invocation reports transaction failure and
  rollback; the retry succeeds without a duplicate number.

## Scenario 20 — Generate Issue validates its bounded input contract

- Initial state: a valid Active Work, plus a valid Active Work whose greatest
  Issue number is `99`.
- Operation: separately request an invalid Issue slug, an unknown status, an
  empty title, an empty body, and a new Issue after `99`.
- Expected state: repository state and INDEX remain unchanged for every request.
- Expected diagnostics: each request fails deterministically before mutation;
  the capacity case reports that no Issue number is available.

------

# Expected Results

Every Documentation Test SHALL define:

- initial repository state;
- operation under test;
- expected repository state;
- expected diagnostics.

Tests SHALL NOT rely upon manual interpretation.

------

# Test Reports

Documentation Testing should produce structured reports.

Typical report sections include:

- executed tests;
- passed tests;
- failed tests;
- execution duration;
- implementation version;
- specification version.

Report format is implementation-defined.

------

# Relationship to Validation

Documentation Testing and Validation serve different purposes.

| Validation                 | Documentation Testing                     |
| -------------------------- | ----------------------------------------- |
| Tests repository state     | Tests Documentation Engine behavior       |
| Repository-focused         | Implementation-focused                    |
| Evaluates one repository   | Evaluates Documentation OS implementation |
| Produces validation report | Produces test report                      |

Validation may itself become the subject of Documentation Testing.

------

# Relationship to Health

Health evaluates repository sustainability.

Documentation Testing evaluates implementation correctness.

A Documentation Engine may successfully pass Documentation Testing while operating on an unhealthy repository.

Likewise, a healthy repository does not guarantee a correct Documentation Engine implementation.

These concerns remain intentionally independent.

------

# Automation

Documentation Testing should support automated execution.

Typical execution points include:

- implementation development;
- release verification;
- Repository Profile changes;
- Documentation OS upgrades;
- continuous integration.

Execution policy is implementation-defined.

------

# Compliance

A Documentation Engine claiming Documentation OS compliance SHALL provide tests capable of verifying:

- Documentation Operations;
- Validation;
- Lifecycle behavior;
- Migration;
- Repository Profile implementation;
- operation integration.
- Generate Issue behavior, including all normative scenarios defined above.

Implementations may provide additional tests provided they remain consistent with Documentation OS specifications.

------

# Non-Goals

Documentation Testing intentionally does not evaluate:

- application functionality;
- source code quality;
- software performance;
- engineering decisions;
- architectural correctness.

These concerns belong to software engineering rather than Documentation OS.

------

# References

- DOS-2001 — Single Repository Profile
- DOS-2004 — Runtime Mapping
- DOS-3001 — Document Lifecycle
- DOS-3002 — Runtime Lifecycle
- DOS-4001 — Documentation Operations
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration
- DOS-6003 — Conformance

------

# Summary

Documentation Testing verifies that Documentation OS implementations faithfully execute the specifications.

Unlike Validation, which evaluates repository correctness, Documentation Testing evaluates Documentation Engine behavior.

Together with Documentation Operations, Validation, Health, and Migration, Documentation Testing completes the operational foundation of Documentation OS, ensuring that compliant implementations remain deterministic, reliable, and specification-driven.
