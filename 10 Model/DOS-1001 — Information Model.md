# DOS-1001 — Information Model

**Status:** Draft
**Version:** 1.0
**Category:** Model

------

# Abstract

This specification defines the conceptual information model of Documentation OS.

The Information Model establishes the highest level abstraction for all repository information.

It defines how Documentation OS classifies information independently of any repository structure, filesystem layout, implementation language, or tooling.

Repository Profiles implement this model.

They do not redefine it.

------

# Purpose

The purpose of the Information Model is to establish a stable conceptual foundation for project information.

Traditional documentation systems organize documents by directory.

Documentation OS organizes information by responsibility.

Directories are implementation details.

Information domains are architectural concepts.

This distinction enables Documentation OS to support multiple repository profiles while preserving identical semantics.

------

# Design Goals

The Information Model is designed to satisfy the following goals.

## Repository Independence

The model MUST remain independent of filesystem layout.

The following concepts are intentionally excluded:

- directories;
- filenames;
- Markdown templates;
- repository conventions.

These concerns belong to Repository Profiles.

------

## Stable Semantics

The meaning of project information should remain constant regardless of implementation.

Changing repository layout MUST never change the meaning of project knowledge.

------

## Clear Responsibility

Every information artifact should possess a single primary responsibility.

Responsibilities should not overlap.

Ambiguous documents should be decomposed into separate information artifacts.

------

## Lifecycle Awareness

Information should participate in explicit lifecycle transitions.

Lifecycle behavior is part of the information model rather than an implementation detail.

------

# Top-Level Domains

Documentation OS organizes all repository information into three top-level categories: Managed Information, Staging Information, and Repository Guidance.

```text
Repository Information

├── Managed Information
│   ├── Knowledge
│   └── Runtime
├── Staging Information
│   └── Inbox
└── Repository Guidance
    └── Agent Entry Document
```

Managed Information is divided into two mutually exclusive domains: Knowledge and Runtime.

Every managed information artifact SHALL belong to exactly one of these two domains.

Staging Information and Repository Guidance are not Managed Information and belong to neither domain.

Staging Information and Repository Guidance are defined separately below.

------

# Knowledge Domain

The Knowledge domain represents persistent understanding of the project.

Knowledge answers questions such as:

- What exists?
- Why does it exist?
- How should it evolve?
- Which rules govern future work?

Knowledge evolves slowly.

Knowledge survives multiple implementation cycles.

Knowledge is authoritative.

Knowledge is defined further in:

DOS-1002 — Knowledge Model.

------

# Runtime Domain

The Runtime domain represents temporary execution information.

Runtime answers questions such as:

- What is currently being implemented?
- Which work is active?
- Which plans exist?
- Which implementation artifacts are temporary?

Runtime evolves continuously.

Active Runtime exists only while Work remains active.

Completed Runtime core assets (PRD, Issues, Handoff) are preserved as immutable historical records.

Runtime is defined further in:

DOS-1003 — Runtime Model.

------

# Repository Guidance

Repository Guidance consists of control artifacts that govern how agents and Documentation Engines enter, read, and operate on the repository.

Repository Guidance is not Managed Information.

It belongs to neither the Knowledge domain nor the Runtime domain.

The Agent Entry Document (AGENTS.md + CLAUDE.md) is the normative Repository Guidance document.

It is defined further in:

DOS-5001 — Agent Entry.

------

# Staging Information

Staging Information holds unresolved observations that have not yet been promoted into Knowledge or discarded.

Staging Information is not Managed Information.

It belongs to neither the Knowledge domain nor the Runtime domain.

Staging items are intentionally lightweight.

They do not possess stable identity, do not participate in lifecycle transitions, and do not carry ownership obligations.

Each staging item is expected to be promoted into a Knowledge Category or discarded.

The Inbox is the normative Staging Information staging area.

It is defined further in:

DOS-2003 — Knowledge Mapping.

------

# Domain Independence

Knowledge and Runtime are conceptually independent.

Knowledge SHALL remain understandable without Runtime.

Runtime MAY reference Knowledge.

Knowledge SHALL NOT depend upon Runtime.

The dependency direction is therefore:

```text
Runtime

↓

Knowledge
```

This relationship guarantees that long-term understanding remains available even after Runtime has been completed.

------

# Information Lifecycle

Repository information continuously flows through the Documentation OS lifecycle.

The conceptual flow is:

```text
Idea

↓

Runtime

↓

Implementation

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Knowledge
```

This flow describes the evolution of information.

It does not prescribe implementation details.

Lifecycle specifications define the operational behavior separately.

------

# Information Identity

Every managed information artifact SHALL possess a stable identity.

Identity is independent of:

- filename;
- repository location;
- implementation profile.

Identity specifications are defined in:

DOS-1004 — Identity Model.

------

# Information Relationships

Information artifacts may possess semantic relationships.

Examples include:

- references;
- dependencies;
- rationale;
- implementation guidance.

Relationships improve navigation.

Relationships do not imply ownership.

Relationship semantics are defined in:

DOS-1005 — Relationship Model.

------

# Information Ownership

Every managed information artifact SHALL possess lifecycle ownership.

Ownership determines:

- who maintains it;
- when it changes;
- how it evolves;
- which lifecycle transitions affect it.

Ownership belongs to the artifact rather than the storage location.

------

# Repository Profiles

Repository Profiles provide concrete mappings from this abstract model into repository structures.

Examples include:

- Single Repository Profile

Future profiles may include:

- Workspace Profile
- Multi Repository Profile

Profiles SHALL preserve the semantics defined by this Information Model.

Profiles SHALL NOT redefine Knowledge or Runtime.

------

# Compliance

A Documentation OS implementation SHALL satisfy the following requirements.

- Every managed information artifact belongs to exactly one domain (Knowledge or Runtime).
- The Knowledge and Runtime domains remain conceptually independent.
- Repository Guidance documents are not Managed Information.
- Staging Information items are not Managed Information and carry no identity, lifecycle, or ownership obligations.
- Repository layout SHALL NOT define semantics.
- Information identity remains stable across profiles.
- Profiles preserve model semantics.

------

# Non-Goals

This specification intentionally does not define:

- repository directories;
- filenames;
- Markdown templates;
- generated files;
- documentation operations;
- implementation tooling.

These concerns are addressed by later specifications.

------

# References

- DOS-0001 — Documentation OS
- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-1004 — Identity Model
- DOS-1005 — Relationship Model

------

# Summary

The Information Model defines the conceptual structure of all repository information.

It organizes repository information into three top-level categories:

- Managed Information — the Knowledge and Runtime domains
- Staging Information — the Inbox
- Repository Guidance — the Agent Entry Document

Every subsequent Documentation OS specification builds upon this abstraction.

Repository Profiles implement the model.

They do not modify it.