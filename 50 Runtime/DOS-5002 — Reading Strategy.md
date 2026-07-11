# DOS-5002 — Reading Strategy

**Status:** Stable
**Version:** 1.0
**Category:** Runtime

------

# Abstract

This specification defines the Reading Strategy for AI agents operating within a Documentation OS repository.

The Reading Strategy establishes a deterministic approach for discovering, selecting, and consuming repository knowledge.

Its purpose is to maximize understanding while minimizing unnecessary repository traversal.

The Reading Strategy recommends efficient navigation.

It does not restrict repository exploration.

------

# Purpose

The purpose of the Reading Strategy is to enable AI agents to efficiently establish sufficient repository understanding before performing engineering work.

Without a structured reading strategy, agents often:

- read too little and miss critical context;
- read too much and waste context window;
- repeatedly read identical information;
- ignore authoritative documentation.

Documentation OS provides a standardized strategy to balance completeness and efficiency.

------

# Scope

The Reading Strategy applies to every AI agent interacting with a Documentation OS repository.

Examples include:

- implementation agents;
- documentation agents;
- review agents;
- planning agents;
- orchestration agents.

The strategy is independent of the underlying AI model.

------

# Design Principles

The Reading Strategy follows the following principles.

## RS-1 Guidance Before Exploration

Agents should follow repository guidance before performing repository-wide discovery.

Repository documentation should reduce unnecessary search.

------

## RS-2 Progressive Understanding

Repository understanding should grow incrementally.

Agents should begin with high-level documentation before reading implementation details.

------

## RS-3 Task-Oriented Reading

Agents should read only the documentation required to perform the current task.

Understanding should be sufficient rather than exhaustive.

------

## RS-4 Knowledge Before Runtime

Persistent Knowledge should be understood before Runtime artifacts.

Runtime depends upon Knowledge.

Knowledge does not depend upon Runtime.

------

## RS-5 Relationships Guide Navigation

Document Relationships define preferred navigation paths.

Agents may continue exploring beyond documented relationships when necessary.

Relationships improve efficiency.

They do not restrict exploration.

------

# Reading Hierarchy

Documentation OS recommends the following reading hierarchy.

```text
AGENTS.md

↓

Repository Guidance

↓

Relevant Knowledge

↓

Related Knowledge

↓

Runtime

↓

Implementation
```

Each level refines repository understanding.

------

# Stage 1 — Repository Guidance

The first stage establishes repository context.

Typical reading includes:

- AGENTS.md;
- repository conventions;
- Repository Profile;
- Documentation OS guidance.

This stage answers:

> How is this repository organized?

------

# Stage 2 — Relevant Knowledge

Agents identify the Knowledge directly related to the task.

Typical categories include:

- Architecture;
- ADR;
- Standards;
- Inbox.

Agents should avoid reading unrelated Knowledge.

------

# Stage 3 — Related Knowledge

After understanding primary Knowledge, agents follow documented Relationships.

Examples include:

Architecture

↓

Referenced ADR

↓

Related Standards

Reading expands naturally according to repository semantics.

------

# Stage 4 — Runtime

Only after Knowledge has been established should Runtime be consulted.

Typical Runtime artifacts include:

- clarified requirements;
- implementation plans;
- execution tasks;
- temporary notes.

Runtime provides execution context.

It does not replace Knowledge.

------

# Stage 5 — Implementation

Implementation details should normally be consulted after repository understanding has been established.

Examples include:

- source code;
- tests;
- configuration;
- implementation artifacts.

Implementation should confirm repository understanding rather than define it.

------

# Reading Decisions

Agents should continuously evaluate whether additional reading is necessary.

Typical questions include:

- Is current Knowledge sufficient?
- Are documented Relationships available?
- Does Runtime reference additional Knowledge?
- Has implementation revealed new context?

Reading should stop when sufficient understanding has been achieved.

------

# Relationship Traversal

Relationships should guide repository exploration.

Example:

```text
Architecture

↓

ADR

↓

Standards

↓

Related Architecture
```

Agents may stop traversal whenever additional documents no longer improve task understanding.

------

# Repository Search

Repository-wide search is permitted.

However, search should supplement documented navigation.

Search should not become the primary navigation mechanism.

Repository guidance and documented Relationships remain authoritative.

------

# Context Management

Agents should avoid repeatedly reading identical documentation.

Previously established repository understanding should be reused whenever possible.

Documentation Operations may assist by identifying relevant documentation.

Context management strategy remains implementation-defined.

------

# Documentation Operations

Documentation Operations may assist the Reading Strategy by:

- locating related artifacts;
- generating navigation indexes;
- resolving stable identities;
- validating references.

Documentation Operations improve navigation.

They do not replace repository understanding.

------

# Reading Completion

Reading is considered complete when:

- repository guidance has been understood;
- relevant Knowledge has been reviewed;
- required Runtime has been consulted;
- sufficient implementation context has been established.

Complete understanding of the entire repository is not required.

------

# Compliance

A Documentation OS-compliant agent SHALL:

- begin with repository guidance;
- prioritize Knowledge over Runtime;
- use documented Relationships where available;
- treat repository search as supplementary;
- establish sufficient understanding before implementation.

------

# Non-Goals

This specification intentionally does not define:

- prompt engineering;
- reasoning algorithms;
- retrieval implementation;
- context window management;
- ranking algorithms.

These concerns remain implementation-specific.

------

# References

- DOS-1002 — Knowledge Model
- DOS-1005 — Relationship Model
- DOS-2003 — Knowledge Mapping
- DOS-5001 — Agent Entry
- DOS-5003 — Execution Contract
- DOS-5004 — Documentation Engine

------

# Summary

The Reading Strategy defines how AI agents progressively build repository understanding.

Agents begin with repository guidance, follow documented Knowledge and Relationships, consult Runtime only when necessary, and finally access implementation details.

This strategy enables efficient, consistent, and deterministic repository navigation while preserving the flexibility to explore beyond documented guidance when required.