# Documentation OS Specification

**Version:** 1.0
**Status:** Draft
**Specification Series:** DOS (Documentation Operating System)

------

# 1. Introduction

Documentation OS (DOS) is a specification for building documentation systems that evolve together with software.

Unlike traditional documentation approaches that treat documentation as static files organized by directories, Documentation OS treats documentation as a living system composed of:

- persistent Knowledge;
- temporary Runtime;
- deterministic Operations;
- explicit Lifecycles;
- standardized Repository Profiles.

The objective of Documentation OS is to make documentation behave more like an operating system than a collection of documents.

------

# 2. Why Documentation OS Exists

Modern software projects increasingly rely on AI-assisted development.

However, existing documentation systems suffer from common problems:

- documentation gradually becomes outdated;
- implementation knowledge disappears after development;
- architecture decisions are difficult to trace;
- repository organization differs across projects;
- AI agents repeatedly rediscover project knowledge.

Documentation OS addresses these problems by introducing a standardized documentation architecture.

Documentation is no longer treated as passive reference material.

Instead, it becomes an active component of software development.

------

# 3. Design Goals

Documentation OS is designed around the following goals.

## Persistent Knowledge

Project understanding should survive individual implementations.

Knowledge should remain useful throughout the lifetime of the repository.

------

## Temporary Runtime

Execution artifacts should remain temporary.

Implementation context should not permanently pollute repository knowledge.

------

## Deterministic Repository Behavior

Repository maintenance should be deterministic.

AI agents should perform engineering reasoning.

Documentation Engines should perform deterministic repository operations.

------

## Repository Independence

Conceptual models should remain independent of repository layout.

Different Repository Profiles may organize repositories differently while preserving identical semantics.

------

## AI-Native Collaboration

Documentation OS is designed for collaboration between:

- humans;
- AI agents;
- automation;
- Documentation Engines.

All participants operate using the same repository model.

------

# 4. Core Concepts

Documentation OS is built upon five foundational concepts.

## Knowledge

Persistent understanding of the repository.

Knowledge answers questions such as:

- What exists?
- Why does it exist?
- How should future work proceed?

Knowledge survives implementation.

------

## Runtime

Temporary execution context.

Runtime exists only while Work is active.

Runtime leaves the active execution context after implementation, transferring new understanding into Knowledge. Core Runtime Assets (PRD, Issues, Handoff) are preserved upon Work completion; only temporary Runtime content is disposable.

------

## Work

A Work represents one bounded engineering activity.

Every Work progresses through an explicit lifecycle.

Work concludes only after repository Knowledge has been synchronized.

------

## Documentation Operations

Documentation Operations perform deterministic repository maintenance.

Examples include:

- Validation;
- Synchronization;
- Migration;
- Health evaluation.

Operations never replace engineering reasoning.

------

## Documentation Engine

The Documentation Engine executes Documentation Operations.

It serves as the deterministic execution core of Documentation OS.

AI agents decide **what** should happen.

The Documentation Engine determines **how deterministic repository operations are performed**.

------

# 5. Architectural Overview

Documentation OS separates repository concerns into independent layers.

```text
Foundation

↓

Model

↓

Repository Profile

↓

Lifecycle

↓

Operations

↓

Runtime
```

Each layer builds upon the previous one.

Higher layers SHALL NOT redefine lower-layer concepts; lower layers SHALL NOT contradict higher-layer normative requirements.

Reference is a cross-cutting support section rather than a layer.

It provides terminology, conformance, and versioning referenced by all layers.

Reference does not participate in the vertical dependency chain.

------

# 6. Specification Structure

The Documentation OS Specification consists of the following chapters.

## 00 — Foundation

Defines:

- Documentation OS
- Design Philosophy
- Core Principles
- Terminology

------

## 10 — Model

Defines the implementation-independent conceptual model.

Includes:

- Information Model
- Knowledge Model
- Runtime Model
- Identity Model
- Relationship Model

------

## 20 — Repository Profile

Defines how conceptual models are represented inside repositories.

Includes:

- Single Repository Profile
- Repository Layout
- Knowledge Mapping
- Runtime Mapping
- Repository Conventions

------

## 30 — Lifecycle

Defines how documentation evolves.

Includes:

- Document Lifecycle
- Runtime Lifecycle
- Knowledge Impact Analysis
- Work Close Pipeline
- Ownership

------

## 40 — Operations

Defines deterministic repository operations.

Includes:

- Documentation Operations
- Validation
- Health
- Migration
- Documentation Testing

------

## 50 — Runtime

Defines interaction between AI agents and Documentation OS repositories.

Includes:

- Agent Entry
- Reading Strategy
- Execution Contract
- Documentation Engine
- CLI

------

## 60 — Reference

Provides common normative language and compliance information.

Includes:

- Normative Language
- Conformance
- Versioning
- Change Log

------

## Appendix

Provides supplementary material.

Examples include:

- dependency graphs;
- reading guides;
- repository examples.

Appendices are informative rather than normative.

------

# 7. Reading Guide

Readers are encouraged to follow the specification in the following order.

For general understanding:

```text
README

↓

Foundation

↓

Model
```

For repository implementation:

```text
Profile

↓

Lifecycle

↓

Operations
```

For Documentation Engine implementation:

```text
Operations

↓

Runtime

↓

Reference
```

Reading specifications out of order is permitted.

However, lower-numbered specifications define concepts used by later chapters.

------

# 8. Repository Profiles

Documentation OS separates conceptual models from repository representation.

Repository Profiles define how repositories organize documentation.

Version 1.0 defines one normative profile:

- Single Repository Profile

Future versions may define additional profiles such as:

- Workspace Profile;
- Monorepo Profile;
- Multi-Repository Profile.

All Repository Profiles preserve the same conceptual semantics.

------

# 9. Documentation Engine

The Documentation Engine is the reference execution model of Documentation OS.

Its responsibilities include:

- Documentation Operations;
- Validation;
- Migration;
- Health evaluation;
- repository maintenance.

The Documentation Engine intentionally excludes engineering reasoning.

This separation enables different AI systems to cooperate while producing identical repository behavior.

------

# 10. Conformance

A Documentation OS implementation is considered compliant when it satisfies the normative requirements defined throughout this specification.

Conformance applies independently to:

- repositories;
- Repository Profiles;
- Documentation Engines;
- AI agents;
- CLI implementations.

Detailed conformance requirements are defined in:

**DOS-6003 — Conformance**

------

# 11. Versioning

Documentation OS evolves through explicit specification versions.

The specification version is independent from:

- Repository Profile version;
- Documentation Engine version;
- repository version.

Versioning rules are defined in:

**DOS-6004 — Versioning**

------

# 12. Normative Language

The keywords:

- SHALL
- MUST
- SHOULD
- MAY
- RECOMMENDED

are interpreted according to the definitions provided in:

**DOS-6002 — Normative Language**

Unless explicitly stated otherwise, all normative statements within this specification use those definitions.

------

# 13. Future Evolution

Documentation OS is designed to evolve incrementally.

Future versions may introduce:

- additional Repository Profiles;
- new Documentation Operations;
- expanded Runtime capabilities;
- richer Documentation Engine integrations;
- new appendices and implementation guidance.

Future specifications shall preserve the conceptual models established by Documentation OS unless explicitly superseded by a later specification.

------

# 14. Relationship to Implementations

This specification intentionally avoids prescribing implementation details.

Valid implementations may include:

- standalone Documentation Engines;
- CLI applications;
- MCP servers;
- IDE extensions;
- Git hooks;
- CI/CD integrations.

All implementations should expose behavior consistent with this specification.

------

# 15. Summary

Documentation OS defines a complete operating model for repository documentation.

Rather than treating documentation as a collection of static files, Documentation OS organizes repository knowledge into:

- persistent Knowledge;
- temporary Runtime;
- explicit Lifecycles;
- deterministic Operations;
- standardized Repository Profiles.

Together, these concepts provide a foundation for AI-native software development in which documentation continuously evolves with the software it describes.

This README serves as the entry point to the Documentation OS Specification.

Subsequent specifications define each concept in detail while preserving the layered architecture described in this document.