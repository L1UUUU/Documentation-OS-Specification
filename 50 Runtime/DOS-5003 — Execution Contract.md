# DOS-5003 — Execution Contract

**Status:** Draft
**Version:** 1.0
**Category:** Runtime

------

# Abstract

This specification defines the Execution Contract between AI agents and a Documentation OS repository.

The Execution Contract establishes the responsibilities, boundaries, and interaction rules that every agent follows while operating within a Documentation OS environment.

Unlike the Reading Strategy, which defines how an agent acquires repository understanding, the Execution Contract defines how an agent behaves after execution begins.

Its objective is to ensure that multiple agents, documentation tooling, and repository lifecycle remain consistent.

------

# Purpose

The purpose of the Execution Contract is to standardize repository interaction.

Every Documentation OS-compliant agent should behave predictably regardless of:

- AI model;
- orchestration framework;
- implementation language;
- execution environment.

The contract minimizes inconsistent repository modifications and ensures deterministic cooperation between agents and the Documentation Engine.

------

# Scope

This specification applies to every AI agent capable of modifying a Documentation OS repository.

Examples include:

- implementation agents;
- planning agents;
- documentation agents;
- review agents;
- orchestration agents.

Read-only analysis agents should follow this specification where applicable.

------

# Design Principles

The Execution Contract follows the following principles.

## EC-1 Repository Before Agent

The repository defines the rules.

Agents adapt to the repository.

Repositories never adapt to individual agents.

------

## EC-2 Respect Documentation Ownership

Agents SHALL respect Ownership defined by Documentation OS.

Agents should modify only artifacts appropriate to the current lifecycle stage.

------

## EC-3 Deterministic Operations

Agents should delegate deterministic maintenance tasks to the Documentation Engine.

Agents should avoid manually reproducing deterministic operations.

------

## EC-4 Explicit Lifecycle

Repository state SHALL progress through explicit lifecycle transitions.

Agents SHALL NOT bypass mandatory lifecycle stages.

------

## EC-5 Knowledge Preservation

Agents SHALL preserve existing repository Knowledge.

Implementation SHALL improve repository understanding rather than replace it.

------

# Execution Preconditions

Before modifying a repository, an agent SHALL:

1. complete Agent Entry;
2. establish sufficient repository understanding;
3. determine affected Knowledge;
4. identify the current Runtime context.

Execution should not begin before these preconditions have been satisfied.

------

# Agent Responsibilities

During execution an agent is responsible for:

- understanding repository context;
- performing engineering reasoning;
- maintaining Runtime artifacts;
- participating in Knowledge Impact Analysis;
- initiating Knowledge Synchronization when required.

Before declaring Work Complete, an agent SHALL ensure that the workstream's Core Runtime Assets are present and finalized:
- `PRD.md` (the workstream specification);
- `issues/` (the Work's issue definitions; at least one `NN-<slug>.md` SHALL exist before Complete);
- `HANDOFF.md` (cross-agent and cross-session transfer documentation).

Agents are responsible for engineering decisions.

They are not responsible for deterministic repository maintenance.

------

# Documentation Engine Responsibilities

The Documentation Engine is responsible for deterministic operations.

Typical responsibilities include:

- identifier allocation;
- generated artifact updates;
- Validation;
- Runtime completion (moving `active/<workstream-slug>/` → `completed/<workstream-slug>/`) and `.scratch/INDEX.md` regeneration;
- repository Migration.

Agents should invoke these operations rather than implementing them manually.

------

# Repository Modification Rules

Repository modifications should follow the following order.

```text
Understand

↓

Reason

↓

Modify

↓

Synchronize

↓

Validate
```

Skipping intermediate stages is discouraged.

Repository modifications should remain traceable.

------

# Knowledge Modification

Knowledge should only be modified when repository understanding changes.

Agents should avoid:

- cosmetic documentation updates unrelated to implementation;
- speculative architectural changes;
- undocumented engineering decisions.

Knowledge modifications should be supported by completed engineering work or explicit design activities.

------

# Runtime Modification

Agents may freely update active Runtime artifacts while Work remains active.

Typical Runtime modifications include:

- updating `PRD.md` (workstream specification and evolution);
- adding or updating `issues/` (issue definitions and resolution);
- maintaining `HANDOFF.md` (cross-agent and cross-session transfer documentation).

Runtime should remain temporary.

------

# Documentation Operations

Whenever deterministic maintenance is required, agents SHALL invoke Documentation Operations.

Typical examples include:

- Validate;
- Complete;
- Generate;
- Synchronize.

Manual execution of deterministic repository maintenance is discouraged.

------

# Concurrent Execution

Documentation OS permits multiple agents to operate within the same repository.

Agents should cooperate through repository state rather than hidden communication.

Shared repository artifacts should remain the primary coordination mechanism.

Repository-specific concurrency mechanisms are implementation-defined.

------

# Error Handling

When execution fails:

- Runtime should remain recoverable;
- repository consistency should be preserved;
- failures should remain explicit;
- lifecycle state should accurately reflect incomplete execution.

Agents should avoid leaving partially synchronized repository Knowledge.

------

# Completion

Execution completes when:

- implementation objectives have been achieved;
- Knowledge Impact Analysis has completed;
- Knowledge Synchronization has completed;
- Validation succeeds;
- the Complete stage has moved the Work to `completed/` (the Work Close Pipeline's Cleanup stage may follow independently);
- `PRD.md`, `issues/` (containing at least one Issue), and `HANDOFF.md` are present and finalized.

Implementation completion alone does not satisfy the Execution Contract.

------

# Execution Invariants

The following invariants SHALL always remain true.

## ECI-1

Repository Knowledge SHALL remain authoritative.

------

## ECI-2

Runtime SHALL remain temporary.

------

## ECI-3

Documentation Operations SHALL perform deterministic maintenance.

------

## ECI-4

Lifecycle stages SHALL NOT be bypassed.

------

## ECI-5

Agents SHALL preserve repository consistency throughout execution.

------

# Compliance

A Documentation OS-compliant agent SHALL:

- complete Agent Entry before execution;
- follow the Reading Strategy;
- respect documentation Ownership;
- invoke Documentation Operations when deterministic maintenance is required;
- complete the Work Close Pipeline before considering Work finished.

------

# Non-Goals

This specification intentionally does not define:

- prompting strategies;
- reasoning techniques;
- implementation algorithms;
- software engineering methodology;
- project management practices.

These concerns remain implementation-specific.

------

# References

- DOS-3002 — Runtime Lifecycle
- DOS-3004 — Work Close Pipeline
- DOS-3005 — Ownership
- DOS-4001 — Documentation Operations
- DOS-5001 — Agent Entry
- DOS-5002 — Reading Strategy
- DOS-5004 — Documentation Engine

------

# Summary

The Execution Contract defines the behavioral rules that every Documentation OS-compliant agent follows.

It establishes a clear separation of responsibilities between:

- AI agents, responsible for engineering reasoning;
- the Documentation Engine, responsible for deterministic repository maintenance.

By standardizing execution behavior, the Execution Contract enables multiple agents, tooling, and repository lifecycle operations to cooperate consistently while preserving the integrity of the Documentation OS repository.