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
- DOS-3001 — Document Lifecycle
- DOS-3002 — Runtime Lifecycle
- DOS-4001 — Documentation Operations
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration

------

# Summary

Documentation Testing verifies that Documentation OS implementations faithfully execute the specifications.

Unlike Validation, which evaluates repository correctness, Documentation Testing evaluates Documentation Engine behavior.

Together with Documentation Operations, Validation, Health, and Migration, Documentation Testing completes the operational foundation of Documentation OS, ensuring that compliant implementations remain deterministic, reliable, and specification-driven.