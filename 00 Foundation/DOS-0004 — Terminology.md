# DOS-0004 — Terminology

**Status:** Stable
**Version:** 1.0
**Category:** Foundation

------

# Abstract

This specification defines the normative terminology of Documentation OS.

A shared vocabulary is essential for ensuring that specifications, implementations, tooling, and documentation all describe the same concepts consistently.

Every term defined in this specification has normative meaning.

Implementations SHOULD use these terms consistently.

Alternative terminology SHOULD be avoided unless explicitly mapped to the definitions contained herein.

------

# Purpose

Documentation OS separates concepts that are often conflated in traditional documentation systems.

Examples include:

- Knowledge vs Runtime
- Architecture vs ADR
- Model vs Profile
- Specification vs Implementation

Without standardized terminology, different implementations may assign different meanings to the same words.

This specification establishes a common language.

------

# Core Concepts

------

# Documentation OS

## Definition

The complete specification governing project knowledge organization, lifecycle, repository profiles, documentation operations, and runtime interaction.

Documentation OS is a specification.

It is not a software implementation.

------

# Repository

## Definition

The complete version-controlled project.

The repository is the authoritative source of:

- source code;
- documentation;
- runtime artifacts;
- metadata;
- relationships.

Documentation OS assumes exactly one repository within the current scope.

------

# Information

## Definition

Any structured information stored within the repository.

Repository information is organized into Managed Information and Repository Guidance.

Managed Information is divided into two domains:

- Knowledge
- Runtime

Every managed information artifact belongs to exactly one domain.

------

# Knowledge

## Definition

Persistent project information describing the current understanding of the system.

Knowledge survives individual implementation work.

Knowledge is authoritative.

Knowledge includes categories such as:

- Architecture
- ADR
- Standards

------

# Runtime

## Definition

Temporary project information supporting active implementation.

Runtime exists only while work is progressing.

Runtime SHALL eventually:

- become archived;
- become discarded;
- or produce persistent Knowledge.

Examples include:

- PRDs
- Issues
- Scratch artifacts

------

# Work

## Definition

A bounded implementation effort managed within the Runtime domain.

A Work represents a unit of execution.

A Work may contain:

- one PRD;
- multiple Issues;
- implementation artifacts.

A Work concludes through the Work Close Pipeline.

------

# Knowledge Category

## Definition

A logical classification of persistent knowledge.

Categories define responsibility rather than storage.

Examples include:

- Architecture
- ADR
- Standards

Repository Profiles determine how categories are mapped into directories.

------

# Repository Profile

## Definition

A concrete mapping between conceptual models and repository layout.

Profiles define:

- directory structure;
- filenames;
- generated artifacts;
- compatibility files.

Examples:

- Single Repository Profile

Future profiles may include:

- Workspace Profile
- Multi Repository Profile

------

# Model

## Definition

An implementation-independent conceptual abstraction.

Models describe semantics.

Models do not define filesystem layout.

Examples:

- Information Model
- Runtime Model
- Identity Model

------

# Specification

## Definition

A normative document defining required behavior.

Specifications describe what must be true.

Specifications do not prescribe implementation details unless explicitly stated.

------

# Implementation

## Definition

A concrete realization of Documentation OS.

Examples include:

- Documentation Engine
- CLI
- Repository Templates
- Validation Tools

Multiple implementations may satisfy the same specification.

------

# Documentation Operation

## Definition

A deterministic operation that maintains repository consistency.

Operations include:

- Generate
- Validate
- Migrate
- Health
- Test

Documentation Operations maintain structure.

They do not create project knowledge.

------

# Documentation Engine

## Definition

The implementation responsible for executing Documentation Operations.

The Documentation Engine interprets Documentation OS specifications and applies deterministic repository transformations.

------

# Managed Region

## Definition

A document region owned by the Documentation Engine.

Managed Regions may be regenerated automatically.

Human-authored content outside Managed Regions shall remain untouched.

------

# Generated Content

## Definition

Repository content produced deterministically by Documentation Operations.

Generated Content shall always be reproducible.

Examples include:

- indexes;
- summaries;
- compatibility files.

------

# Validation

## Definition

The deterministic verification of repository correctness.

Validation verifies structural correctness.

Validation does not evaluate engineering quality.

------

# Health

## Definition

A quality assessment of project knowledge.

Health evaluates long-term maintainability.

Health complements Validation.

Health does not replace Validation.

------

# Knowledge Impact Analysis

## Definition

The process of determining how implementation affects project knowledge.

Knowledge Impact Analysis identifies:

- affected Knowledge;
- required synchronization;
- documentation updates.

Knowledge Impact Analysis precedes Knowledge Synchronization.

------

# Knowledge Synchronization

## Definition

The process of updating persistent Knowledge after implementation.

Knowledge Synchronization transfers newly created understanding from Runtime into the Knowledge domain.

------

# Work Close Pipeline

## Definition

The lifecycle responsible for completing a Work.

The pipeline includes:

Implementation

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Validation

↓

Archive

↓

Close

A Work is not considered complete until the pipeline finishes successfully.

------

# Identity

## Definition

A stable identifier assigned to an artifact.

Identity remains stable throughout the artifact lifecycle.

Identity is independent of filename.

Identity specifications are defined separately.

------

# Relationship

## Definition

A directed semantic connection between documentation artifacts.

Relationships improve navigation and traceability.

Relationships do not imply execution dependency.

------

# Ownership

## Definition

Responsibility for maintaining a documentation artifact.

Ownership defines:

- lifecycle authority;
- update responsibility;
- validation responsibility.

Ownership belongs to lifecycle rather than individuals.

------

# Domain

## Definition

A classification of Managed Information.

Documentation OS defines two domains:

- Knowledge
- Runtime

Every managed information artifact belongs to exactly one domain.

Domains are conceptual.

They are not filesystem directories.

------

# Artifact

## Definition

Any managed documentation object participating in Documentation OS.

Examples include:

- Architecture documents;
- ADRs;
- Standards;
- Inbox observations;
- Runtime artifacts.

Artifacts possess identity and lifecycle.

------

# Agent

## Definition

An autonomous or semi-autonomous system capable of interacting with a Documentation OS repository.

Agents perform:

- repository understanding;
- engineering reasoning;
- implementation.

Agents delegate deterministic maintenance to the Documentation Engine.

------

# Agent Entry Document

## Definition

The canonical repository entry document for AI agents.

Under the Single Repository Profile, `AGENTS.md` is the canonical source of the Agent Entry Document, and `CLAUDE.md` is its content-equivalent mirror.

Both files are valid entry points.

The Agent Entry Document is Repository Guidance rather than Managed Information.

------

# Lifecycle

## Definition

The sequence of states through which an Artifact or Work evolves.

Lifecycle governs:

- creation;
- modification;
- validation;
- archival;
- retirement.

Lifecycle progression is explicit.

------

# Inbox

## Definition

A staging area defined by the Single Repository Profile, not a normative Knowledge Category.

Inbox holds unresolved repository observations that have not yet been classified into Architecture, ADR, or Standards.

Inbox items are expected to be promoted into a Knowledge Category or discarded.

------

# Usage Rules

Implementations SHALL use terminology consistently with this specification.

Alternative names MAY be used in user interfaces.

However, implementation documentation SHOULD map those names back to the normative terminology defined herein.

For example:

```
Task
```

may appear in a UI.

Internally, Documentation OS defines the normative concept as:

```
Work
```

------

# Future Extensions

Future versions may introduce additional terminology.

Existing definitions SHOULD remain stable.

Changing the meaning of an existing term constitutes a breaking specification change and therefore requires a new major version of Documentation OS.

------

# References

- DOS-0001 — Documentation OS
- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles

------

# Summary

This specification establishes the shared vocabulary of Documentation OS.

All subsequent specifications, implementations, repository profiles, documentation engines, and tooling SHALL interpret these terms according to the definitions provided herein.