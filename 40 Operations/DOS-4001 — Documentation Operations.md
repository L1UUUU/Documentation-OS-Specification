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

## Lifecycle Stages

Complete Operation orchestrates the following stages:

### Complete Stage

Move the Work directory from `active/<workstream-slug>/` to `completed/<workstream-slug>/`.

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
- concurrency model.

These concerns belong to individual Documentation Engine implementations.

------

# References

- DOS-0003 — Core Principles
- DOS-3004 — Work Close Pipeline
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration
- DOS-4005 — Documentation Testing
- DOS-5004 — Documentation Engine

------

# Summary

Documentation Operations define the deterministic capabilities required to maintain a Documentation OS repository.

Operations maintain repository structure.

They do not create repository knowledge.

This separation between deterministic maintenance and engineering reasoning forms one of the core architectural principles of Documentation OS.