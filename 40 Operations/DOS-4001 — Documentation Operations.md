# DOS-4001 — Documentation Operations

**Status:** Draft
**Version:** 1.0
**Category:** Operations

------

# Abstract

This specification defines the Documentation Operations of Documentation OS.

Documentation Operations are deterministic repository operations responsible for maintaining documentation consistency.

Documentation Operations never create engineering knowledge.

Instead, they maintain, validate, synchronize, and evolve the documentation system according to the specifications defined by Documentation OS.

This specification defines the operational contract that every Documentation Engine SHALL implement.

------

# Purpose

The purpose of Documentation Operations is to separate deterministic repository maintenance from engineering reasoning.

Engineering reasoning belongs to:

- humans;
- AI agents.

Deterministic maintenance belongs to the Documentation Engine.

This separation allows Documentation OS to automate repository consistency without automating engineering judgement.

------

# Scope

Documentation Operations are responsible for repository maintenance.

Typical responsibilities include:

- identifier management;
- relationship maintenance;
- generated content;
- validation;
- lifecycle transitions;
- migration;
- repository health.

Documentation Operations do not perform:

- architectural design;
- engineering decisions;
- implementation planning;
- repository reasoning.

------

# Design Principles

Documentation Operations follow the following principles.

## DO-1 Deterministic

Every operation SHALL produce identical results when executed against identical repository state.

Operations SHALL avoid non-deterministic behavior.

------

## DO-2 Idempotent

Operations should support safe repeated execution.

Executing the same operation multiple times should not introduce inconsistent repository state.

------

## DO-3 Non-Destructive

Operations SHALL preserve human-authored knowledge.

Managed content may be regenerated.

Human-authored content outside managed regions SHALL remain untouched.

------

## DO-4 Specification Driven

Operations SHALL implement Documentation OS specifications.

Repository behavior SHALL NOT depend upon undocumented implementation decisions.

------

## DO-5 Repository Local

Operations SHALL derive their behavior entirely from repository contents.

Persistent external state SHALL NOT be required.

------

# Operation Categories

Documentation Operations are divided into six categories.

```text
Documentation Operations

├── Generate

├── Synchronize

├── Validate

├── Complete

├── Migrate

└── Inspect
```

Each category has independent responsibilities.

------

# Generate

## Purpose

Generate deterministic repository artifacts.

Typical examples include:

- indexes (including `.scratch/INDEX.md`, which provides navigation and statistics for active and completed Works);
- summaries;
- managed metadata;
- navigation files;
- compatibility files.

Generated artifacts SHALL always be reproducible.

## Generate Work

Generate Work is responsible for creating a new Work workspace.

Operations include:

- creating the `active/<slug>/` directory structure;
- generating a PRD template;
- creating an empty `issues/` directory;
- creating an empty `HANDOFF.md` file;
- verifying that the workstream slug is globally unique across both `active/` and `completed/`;
- regenerating `.scratch/INDEX.md` to reflect the newly created Work;
- rolling back all created artifacts (including the INDEX regeneration) if creation fails.

Generate Work belongs to the Generate category and SHALL NOT introduce a new operation category.

## Generate Issue

Generate Issue appends one caller-authored Issue to an existing Active Work.
The caller supplies the Work slug, Issue slug, title, lifecycle status, and
Markdown body. The Documentation Engine supplies all deterministic repository
representation, including numbering, filename construction, front matter, and
INDEX maintenance. Generate Issue belongs to the Generate category and SHALL
NOT introduce a new operation category.

### Input Contract

The Work slug and Issue slug SHALL satisfy the Single Repository Profile slug
rules. The title SHALL contain non-whitespace content and SHALL be representable
as one front-matter scalar line. The body SHALL contain non-whitespace Markdown
content. The status SHALL be exactly one of `open`, `in-progress`, `done`,
`blocked`, `cancelled`, or `superseded`.

The caller SHALL NOT supply an Issue number or destination path. Workflow or
triage labels defined by a particular repository are outside this operation's
status vocabulary.

### Active Work Precondition

Generate Issue SHALL operate only on
`.scratch/active/<workstream-slug>/`. If the Work exists only under
`.scratch/completed/`, the operation SHALL reject the request because Completed
Core Runtime Assets are immutable. If the Work is missing, appears in both
locations, is not a directory, or lacks the required Active Work structure, the
operation SHALL fail without modifying repository state.

### Number and Name Allocation

While providing repository-scoped exclusive serialization across allocation
and publication, the operation SHALL inspect the target Work's existing Issue
filenames. It SHALL
allocate one greater than the greatest existing `NN`; an empty `issues/`
directory receives `01`. Deleted-number gaps SHALL NOT be reused. The resulting
name is `NN-<slug>.md`. Allocation SHALL fail before mutation if the next number
would exceed `99`, if existing Issue numbers are duplicated or malformed, or if
the requested slug is already associated with conflicting content.

Concurrent Generate Issue calls for the same Work SHALL be serialized across
number inspection, Issue publication, and INDEX update. Distinct successful
requests SHALL receive distinct numbers. Concurrent identical requests SHALL
converge on one Issue: one call may report creation and the others report an
idempotent existing result.

### Idempotency and Conflict

The operation's logical input is the validated tuple of Work slug, Issue slug,
title, status, and Markdown body. When an Active Work already contains the same
Issue slug whose parsed `title` and `status` equal the validated input and whose
Markdown body equals the supplied body, a retry
SHALL return that Issue's existing name, number, and path and report that no new
artifact was created. It SHALL NOT allocate another number or rewrite the
Issue. Before returning, it SHALL verify that `.scratch/INDEX.md` lists the
existing Issue and transactionally regenerate INDEX when an earlier
interruption left it stale.

When the same Issue slug exists with a different title, status, or body, the
operation SHALL return a conflict and leave both the Issue and INDEX unchanged.
Idempotency applies only while the Work is Active; a request targeting a
Completed Work is rejected even if equivalent content exists there.

### Transactional Publication

Publishing the Issue file and regenerating `.scratch/INDEX.md` form one
repository transaction. The operation SHALL stage writes and retain sufficient
pre-operation state to provide these observable guarantees:

- success exposes both the new Issue and an INDEX that lists it;
- any returned failure exposes neither the new Issue nor an INDEX change;
- an interruption is recoverable by retrying the same logical input, without a
  duplicate Issue or duplicate number;
- temporary transaction artifacts are not treated as Runtime assets.

The operation SHALL restore the previous INDEX and remove the new Issue if
INDEX publication fails. Exclusive serialization SHALL remain in effect until
publication or rollback completes. The locking mechanism is
implementation-defined.

### Result

A successful Generate Issue result SHALL expose the allocated number, Issue
name, repository-relative path, and whether this invocation created the Issue.
The result SHALL be identical across invocation mechanisms apart from
presentation format.

------

# Synchronize

## Purpose

Synchronize repository state after lifecycle events.

Examples include:

- updating generated references;
- refreshing managed regions;
- synchronizing derived documentation.

Synchronization maintains consistency.

Synchronization does not create engineering knowledge.

------

# Validate

## Purpose

Verify repository correctness.

Validation includes:

- identity integrity;
- relationship integrity;
- lifecycle consistency;
- repository consistency;
- generated artifact verification.

Validation behavior is defined by:

DOS-4002.

------

# Complete

## Purpose

Perform deterministic Runtime completion by orchestrating two lifecycle stages.

Complete Operations SHALL NOT modify synchronized Knowledge.

## Outcome Input

The terminal `outcome` of a Work is an engineering judgement; it SHALL be supplied by the Agent or caller as an input parameter to the Complete Operation. The Documentation Engine does not decide the outcome; it validates and records it. Normative `outcome` values are `succeeded`, `cancelled`, `superseded`, `failed` (see DOS-3002).

An active Work SHALL NOT carry an `outcome` in its PRD front matter before the Complete Operation is invoked; the `outcome` is recorded by the Complete stage (DOS-3004).

## Preflight

Before performing any state change, the Complete Operation SHALL verify, as a single atomic preflight:

- a legal `outcome` value has been supplied;
- Core Runtime Assets exist (PRD.md, at least one `issues/NN-<slug>.md`, HANDOFF.md);
- when `outcome = succeeded`, no Issue remains `open`, `in-progress`, or `blocked` (such Issues MUST be resolved to `done`, `cancelled`, or `superseded`); a non-`succeeded` Work MAY leave non-terminal Issues with the outcome truthfully recorded.

If preflight fails, the Complete Operation SHALL fail without modifying repository state, and the Work SHALL remain in `active/`.

## Lifecycle Stages

Complete Operation orchestrates the following stages:

### Complete Stage

Record the supplied terminal `outcome` in PRD front matter, clean Ephemeral Runtime Content, and atomically move the Work directory from `active/<workstream-slug>/` to `completed/<workstream-slug>/`. Recording the `outcome`, cleaning Ephemeral Content, and moving the directory form a single transaction: if the directory movement fails, the `outcome` write and Ephemeral cleanup SHALL be rolled back so that the Work remains valid under `active/` (no `outcome` persists on an active Work).

This transition marks the Work's entry into its terminal state and releases Ownership.

### Cleanup Stage

Regenerate `.scratch/INDEX.md` to reflect the updated Runtime state.

This stage is idempotent and may be retried independently if it fails.

### Resuming after a Cleanup Failure

When a Complete Operation is re-invoked on a Work that already resides in `completed/`, the operation SHALL resume from the Cleanup stage only; it SHALL NOT re-execute the Complete stage directory movement. This gives the "retry only Cleanup" recovery path a standard, deterministic behavioral contract.

------

# Migrate

## Purpose

Safely evolve repository structure.

Migration includes:

- repository profile upgrades;
- structural transformations;
- compatibility updates.

Migration SHALL preserve repository semantics.

Migration behavior is defined by:

DOS-4004.

------

# Inspect

## Purpose

Evaluate repository state.

Inspection includes:

- repository health;
- documentation metrics;
- consistency reports;
- lifecycle summaries.

Inspection does not modify repository contents.

------

# Operational Characteristics

Every Documentation Operation SHALL satisfy the following properties.

## Explicit

Operations execute intentionally.

Hidden background modifications are discouraged.

------

## Observable

Operations should provide deterministic results.

Failures should be explicit.

------

## Repeatable

Operations should safely support repeated execution.

------

## Repository Safe

Operations should never leave the repository in a partially inconsistent state.

Partial execution should be recoverable.

------

# Documentation Engine

Documentation Operations are executed by the Documentation Engine.

The Documentation Engine is responsible for:

- locating repository artifacts;
- executing operations;
- enforcing specifications;
- preserving repository consistency.

The Documentation Engine does not replace engineering reasoning.

It implements deterministic repository behavior.

------

# Agent Interaction

AI agents should invoke Documentation Operations whenever deterministic maintenance is required.

Typical examples include:

- allocating identifiers;
- validating repository state;
- updating generated indexes;
- completing Runtime.

Agents should not manually reproduce deterministic operations already implemented by the Documentation Engine.

------

# Operation Ordering

Some operations possess ordering constraints.

Typical example:

```text
Synchronize

↓

Validate

↓

Complete
```

Individual specifications define additional ordering requirements.

Documentation Engines SHALL preserve required ordering.

------

# Failure Handling

Documentation Operations SHALL fail explicitly.

Failure should include sufficient information for:

- diagnosis;
- recovery;
- retry.

Failed operations SHALL avoid leaving the repository in an inconsistent state.

------

# Compliance

A compliant Documentation Engine SHALL provide operations capable of:

- generating deterministic artifacts;
- synchronizing repository state;
- validating repository consistency;
- completing Runtime;
- migrating repository structure;
- inspecting repository health.

Implementations may introduce additional operations provided they remain consistent with Documentation OS.

------

# Non-Goals

This specification intentionally does not define:

- CLI syntax;
- API design;
- implementation language;
- execution environment;
- locking mechanism.

These concerns belong to individual Documentation Engine implementations.

------

# References

- DOS-0003 — Core Principles
- DOS-2004 — Runtime Mapping
- DOS-3004 — Work Close Pipeline
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration
- DOS-4005 — Documentation Testing
- DOS-5004 — Documentation Engine
- DOS-5005 — CLI

------

# Summary

Documentation Operations define the deterministic capabilities required to maintain a Documentation OS repository.

Operations maintain repository structure.

They do not create repository knowledge.

This separation between deterministic maintenance and engineering reasoning forms one of the core architectural principles of Documentation OS.
