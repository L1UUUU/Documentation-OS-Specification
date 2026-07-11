# DOS-6001 — Terminology

**Status:** Stable
**Version:** 1.0
**Category:** Reference

------

# Abstract

This specification defines the normative terminology used throughout Documentation OS.

Its purpose is to ensure that every specification uses consistent language.

Documentation OS distinguishes carefully between concepts that are often treated interchangeably in traditional documentation systems.

Every term defined here possesses one canonical meaning.

Later specifications shall use these terms consistently.

------

# Purpose

The Terminology specification establishes a shared vocabulary for:

- specification authors;
- Documentation Engine implementations;
- Repository Profiles;
- AI agents;
- repository maintainers.

Definitions in this specification are normative.

If another specification appears to conflict with these definitions, this specification takes precedence.

------

# Terminology Principles

Documentation OS terminology follows the following principles.

## TM-1 One Concept, One Term

Every important concept should have exactly one preferred term.

Multiple synonyms are discouraged.

------

## TM-2 Stable Meaning

Terminology should remain stable across Documentation OS versions whenever practical.

Changing terminology is considered a breaking specification change.

------

## TM-3 Layer Independence

Terminology describes conceptual models rather than repository implementations.

Terms should remain valid regardless of Repository Profile.

------

# Repository

A **Repository** is the complete documentation and source code environment managed under a single Documentation OS implementation.

A Repository contains:

- source code;
- Knowledge;
- Runtime;
- Documentation Operations;
- repository metadata.

Repository organization is defined by a Repository Profile.

------

# Knowledge

**Knowledge** is the persistent understanding of the repository.

Knowledge describes:

- current architecture;
- engineering decisions;
- repository standards;
- unresolved repository observations.

Knowledge survives implementation.

Knowledge belongs to the Knowledge domain.

------

# Runtime

**Runtime** is the temporary execution domain of the repository.

Runtime supports active engineering work.

Runtime exists only while Work remains active.

Runtime concludes through the Work Close Pipeline.

------

# Work

A **Work** is the primary execution unit of Documentation OS.

A Work represents one bounded engineering activity.

Every Work possesses:

- identity;
- lifecycle;
- ownership;
- Runtime artifacts.

A Work concludes only after repository Knowledge has been synchronized.

------

# Artifact

An **Artifact** is any managed documentation object participating in Documentation OS.

Examples include:

- Architecture documents;
- ADRs;
- Standards;
- Inbox items;
- Runtime artifacts.

Artifacts possess identity and lifecycle.

------

# Identity

**Identity** is the stable identifier assigned to a managed Artifact.

Identity remains unchanged throughout the Artifact lifecycle.

Identity is independent of:

- filename;
- directory;
- Repository Profile.

------

# Relationship

A **Relationship** is a semantic connection between two Artifacts.

Relationships describe meaning.

They do not describe repository layout.

Relationships are expressed using stable identities.

------

# Lifecycle

A **Lifecycle** defines the sequence of states through which an Artifact or Work evolves.

Lifecycle governs:

- creation;
- modification;
- validation;
- archival;
- retirement.

Lifecycle progression is explicit.

------

# Knowledge Synchronization

**Knowledge Synchronization** is the process of transferring repository understanding from Runtime into Knowledge.

Knowledge Synchronization updates persistent repository documentation after implementation has completed.

Synchronization follows Knowledge Impact Analysis.

------

# Knowledge Impact Analysis

**Knowledge Impact Analysis (KIA)** is the evaluation performed after implementation to determine which Knowledge Categories require synchronization.

KIA determines:

- what changed;
- what Knowledge is affected;
- what synchronization actions are required.

KIA does not modify documentation.

------

# Documentation Operation

A **Documentation Operation** is a deterministic repository maintenance operation executed by the Documentation Engine.

Examples include:

- Validation;
- Migration;
- Health evaluation;
- Synchronization;
- Archive.

Documentation Operations never perform engineering reasoning.

------

# Documentation Engine

The **Documentation Engine** is the deterministic execution component implementing Documentation Operations.

The Documentation Engine maintains repository consistency.

It does not make engineering decisions.

------

# Repository Profile

A **Repository Profile** defines how Documentation OS concepts are represented inside a repository.

Repository Profiles define representation.

They do not redefine semantics.

Examples include:

- Single Repository Profile.

------

# Agent

An **Agent** is an autonomous or semi-autonomous system capable of interacting with a Documentation OS repository.

Agents perform:

- repository understanding;
- engineering reasoning;
- implementation.

Agents delegate deterministic maintenance to the Documentation Engine.

------

# AGENTS.md

**AGENTS.md** is the canonical repository entry document for AI agents.

It defines:

- repository overview;
- documentation architecture;
- execution guidance;
- Documentation Operations;
- reading strategy.

AGENTS.md is guidance rather than project Knowledge.

------

# Inbox

**Inbox** is the Knowledge Category responsible for recording unresolved repository observations.

Inbox items represent repository understanding that has not yet become implementation work.

Inbox belongs to the Knowledge domain.

It is not part of Runtime.

------

# Architecture

**Architecture** describes the current structure of the repository.

Architecture answers:

> What exists?

Architecture intentionally excludes historical rationale.

------

# Architecture Decision Record (ADR)

An **Architecture Decision Record (ADR)** records significant engineering decisions.

An ADR answers:

> Why does the current Architecture exist?

ADR complements Architecture.

It does not replace it.

------

# Standards

**Standards** define engineering conventions governing future repository evolution.

Standards answer:

> How should future work proceed?

Standards are normative.

------

# Validation

**Validation** is the deterministic verification of repository correctness.

Validation evaluates repository structure.

Validation does not evaluate engineering quality.

------

# Health

**Health** is the long-term assessment of repository sustainability.

Health evaluates trends rather than correctness.

Health is advisory.

------

# Migration

**Migration** is the deterministic transformation of repository representation while preserving repository semantics.

Migration changes structure.

Migration does not change meaning.

------

# CLI

The **CLI** is the command-line interface exposing Documentation Operations.

The CLI is an interface layer.

The Documentation Engine performs the underlying operations.

------

# Documentation OS

**Documentation OS** is the complete specification defining:

- conceptual models;
- Repository Profiles;
- lifecycle behavior;
- Documentation Operations;
- AI interaction.

Documentation OS standardizes documentation systems for AI-native software development.

------

# References

- README
- DOS-1001 — Information Model
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-4001 — Documentation Operations
- DOS-5004 — Documentation Engine

------

# Summary

This specification establishes the canonical terminology of Documentation OS.

All subsequent specifications, implementations, Repository Profiles, Documentation Engines, and AI agents should use these definitions consistently to preserve clarity, interoperability, and long-term specification stability.