# DOS-5004 — Documentation Engine

**Status:** Draft
**Version:** 1.0
**Category:** Runtime

------

# Abstract

This specification defines the Documentation Engine of Documentation OS.

The Documentation Engine is the deterministic execution component responsible for maintaining repository consistency.

Unlike AI agents, which perform engineering reasoning, the Documentation Engine performs specification-defined repository operations.

The Documentation Engine is the operational core of Documentation OS.

It transforms the Documentation OS specifications into executable repository behavior.

------

# Purpose

The purpose of the Documentation Engine is to centralize all deterministic repository maintenance.

Without a Documentation Engine:

- different AI agents implement repository operations differently;
- repository conventions gradually diverge;
- lifecycle behavior becomes inconsistent;
- documentation maintenance becomes difficult to automate.

The Documentation Engine provides one authoritative implementation of Documentation OS behavior.

------

# Scope

The Documentation Engine is responsible for implementing Documentation Operations defined by Documentation OS.

Typical responsibilities include:

- Documentation Operations;
- Validation;
- Health evaluation;
- Migration;
- Runtime lifecycle support;
- repository metadata management.

The Documentation Engine does not perform engineering reasoning.

------

# Design Principles

The Documentation Engine follows the following principles.

## DE-1 Deterministic

Given identical repository state and identical operation parameters, the Documentation Engine shall always produce identical results.

------

## DE-2 Stateless

The Documentation Engine should derive behavior entirely from repository contents.

Persistent external state should not be required.

Repository state is the source of truth.

------

## DE-3 Specification-Driven

The Documentation Engine implements Documentation OS specifications.

Implementation behavior shall never contradict specification behavior.

------

## DE-4 Tool Independent

The Documentation Engine is independent of:

- AI models;
- orchestration frameworks;
- editors;
- IDEs;
- MCP clients.

All integrations interact through the same operational contract.

------

## DE-5 Repository Local

All operations execute against the current repository.

The Documentation Engine shall not depend upon cloud services unless explicitly configured by an implementation.

------

# Responsibilities

The Documentation Engine is responsible for deterministic repository maintenance.

Examples include:

- allocating identifiers;
- validating repository consistency;
- maintaining generated artifacts;
- synchronizing managed metadata;
- moving active Runtime workstreams to completed (Complete Work: `active/<workstream-slug>/` → `completed/<workstream-slug>/`);
- generating and refreshing `.scratch/INDEX.md`;
- evaluating repository Health;
- executing Migration.

Engineering decisions remain outside the scope of the Documentation Engine.

------

# Architecture

Conceptually, the Documentation Engine occupies the following position.

```text
AI Agent

↓

Documentation Engine

↓

Repository
```

The agent decides **what** should happen.

The Documentation Engine determines **how deterministic repository operations are executed**.

------

# Operation Execution

Every Documentation Operation follows the same conceptual execution model.

```text
Receive Request

↓

Load Repository

↓

Execute Specification

↓

Validate Result

↓

Return Outcome
```

The internal implementation remains implementation-defined.

The observable behavior shall remain specification-compliant.

------

# Repository Discovery

Before executing any operation, the Documentation Engine shall determine:

- Documentation OS version;
- active Repository Profile;
- repository entry points;
- managed documentation locations.

Repository discovery should be automatic.

------

# Repository Metadata

The Documentation Engine may maintain repository metadata required for deterministic operation.

Examples include:

- allocated identifiers;
- generated indexes (including `.scratch/INDEX.md`);
- managed metadata;
- compatibility information.

Repository metadata shall remain reproducible whenever possible, with `.scratch/INDEX.md` maintained as a stateless derivative of the active and completed workstreams.

------

# Interaction with AI Agents

AI agents interact with the Documentation Engine through Documentation Operations.

Typical interaction model:

```text
AI Agent

↓

Validate()

↓

Documentation Engine

↓

Validation Report
```

The Documentation Engine exposes capabilities.

The AI agent decides when to invoke them.

------

# Interaction with Repository Profiles

Repository Profiles define repository representation.

The Documentation Engine interprets repository contents according to the active Repository Profile.

Adding a new Repository Profile should not require changes to Documentation OS semantics.

------

# Error Handling

The Documentation Engine shall provide explicit error reporting.

Typical error categories include:

- invalid repository structure;
- unsupported Repository Profile;
- validation failure;
- migration failure;
- operation failure.

Errors should identify:

- operation;
- affected artifact;
- failure reason;
- recommended recovery.

------

# Extensibility

Implementations may introduce additional Documentation Operations.

Additional operations shall satisfy the Documentation Engine principles and remain compatible with Documentation OS.

Extensions shall not redefine existing specification behavior.

------

# Implementation Independence

Documentation OS intentionally does not prescribe implementation language.

Valid implementations may include:

- Go
- Rust
- Java
- Python
- other languages

All implementations should expose equivalent observable behavior.

------

# Relationship to CLI

The Documentation Engine is independent of any command-line interface.

CLI implementations invoke Documentation Operations.

The CLI is a presentation layer.

The Documentation Engine is the execution layer.

------

# Relationship to MCP

The Documentation Engine is independent of Model Context Protocol (MCP).

An MCP server may expose Documentation Operations as tools.

Documentation Engine behavior remains unchanged regardless of invocation mechanism.

------

# Compliance

A Documentation OS-compliant Documentation Engine SHALL:

- implement Documentation Operations;
- preserve repository consistency;
- remain deterministic;
- remain specification-driven;
- remain independent of AI reasoning;
- support Repository Profiles.

------

# Non-Goals

This specification intentionally does not define:

- command-line syntax;
- REST APIs;
- RPC protocols;
- MCP schemas;
- implementation architecture;
- storage engine implementation.

These concerns belong to specific implementations.

------

# References

- DOS-4001 — Documentation Operations
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration
- DOS-5003 — Execution Contract
- DOS-5005 — CLI

------

# Summary

The Documentation Engine is the deterministic execution core of Documentation OS.

It translates Documentation OS specifications into executable repository operations while remaining completely independent of engineering reasoning.

By separating deterministic repository maintenance from AI decision-making, the Documentation Engine enables consistent behavior across different AI agents, editors, orchestration systems, and execution environments.