# DOS-2003 — Knowledge Mapping

**Status:** Stable
**Version:** 1.0
**Category:** Repository Profile

------

# Abstract

This specification defines how the conceptual Knowledge Model is mapped into concrete Knowledge Categories within the Single Repository Profile.

The Knowledge Model defines **what Knowledge is**.

This specification defines **how Knowledge is organized**.

Repository organization exists to improve maintainability, discoverability, ownership, and lifecycle management.

Knowledge Categories represent conceptual responsibilities rather than filesystem structures.

------

# Purpose

The purpose of Knowledge Mapping is to decompose persistent project knowledge into independent categories with clear responsibilities.

Each category answers one primary question.

Categories should complement rather than duplicate one another.

Repository directories are merely one implementation of these categories.

------

# Design Principles

Knowledge Mapping follows the following principles.

## KM-1 Single Responsibility

Every Knowledge Category should answer one primary question.

A category should not attempt to become a general-purpose documentation area.

------

## KM-2 Minimal Overlap

Knowledge Categories should minimize duplicated information.

If two categories continuously duplicate content, their responsibilities are incorrectly defined.

------

## KM-3 Stable Knowledge

Knowledge Categories should evolve slowly.

They should represent long-term understanding rather than temporary implementation details.

------

## KM-4 Explicit Ownership

Each Knowledge Category has explicit lifecycle ownership.

Knowledge synchronization updates categories according to their responsibilities.

------

# Knowledge Categories

Documentation OS v1 defines four normative Knowledge Categories.

```text
Knowledge

├── Architecture

├── ADR

├── Standards

└── Inbox
```

Future specifications may introduce additional categories.

------

# Architecture

## Purpose

Architecture describes the current structure of the project.

Architecture answers:

> How is the project organized today?

Typical contents include:

- system structure;
- module responsibilities;
- interfaces;
- boundaries;
- dependency relationships;
- invariants.

Architecture intentionally avoids:

- implementation history;
- temporary planning;
- future speculation.

Architecture should always describe the current system.

------

# Architecture Decision Records

## Purpose

Architecture Decision Records explain significant engineering decisions.

ADR answers:

> Why does the system look like this?

Each ADR should document:

- context;
- decision;
- consequences;
- alternatives considered.

Architecture explains the system.

ADR explains its evolution.

These responsibilities shall remain separate.

------

# Standards

## Purpose

Standards define engineering conventions.

Standards answer:

> How should future work be performed?

Examples include:

- coding conventions;
- documentation conventions;
- testing requirements;
- architectural rules;
- repository conventions.

Standards are normative.

They guide future implementation rather than describing existing implementation.

------

# Inbox

## Purpose

Inbox stores unresolved repository knowledge.

Inbox answers:

> What has been observed but not yet analysed?

Inbox exists before Runtime.

Items in Inbox have not yet become implementation work.

Examples include:

- architectural concerns;
- documentation inconsistencies;
- technical debt observations;
- improvement ideas;
- repository issues.

Inbox is intentionally lightweight.

It does not contain implementation plans.

------

# Category Relationships

Knowledge Categories collaborate through explicit responsibilities.

Typical relationships include:

```text
Inbox

↓

ADR

↓

Architecture

↓

Standards
```

Example:

An observation enters Inbox.

↓

Analysis determines that an architectural decision is required.

↓

An ADR records the decision.

↓

Architecture is updated.

↓

Standards are modified if future engineering behavior changes.

This represents one possible Knowledge evolution path.

------

# Category Ownership

Each Knowledge Category possesses independent ownership.

| Category     | Primary Owner                      |
| ------------ | ---------------------------------- |
| Architecture | Current system structure           |
| ADR          | Design rationale                   |
| Standards    | Engineering conventions            |
| Inbox        | Unresolved repository observations |

Ownership determines synchronization responsibilities.

------

# Knowledge Synchronization

Knowledge Synchronization updates categories selectively.

Examples:

Implementation changes module boundaries.

↓

Architecture updated.

Implementation changes engineering conventions.

↓

Standards updated.

Implementation introduces major design decision.

↓

ADR created.

Implementation reveals future work.

↓

Inbox updated.

Not every implementation affects every category.

Knowledge Impact Analysis determines which categories require synchronization.

------

# Category Independence

Knowledge Categories remain independent.

Architecture should not become:

- design history;
- engineering guide;
- issue tracker.

Standards should not duplicate Architecture.

Inbox should not become Runtime.

ADR should not become Architecture.

Clear boundaries improve long-term maintainability.

------

# Repository Mapping

The Single Repository Profile maps Knowledge Categories into repository directories.

A typical mapping is:

```text
docs/

├── architecture/

├── adr/

├── standards/

└── inbox/
```

This mapping belongs to the Single Repository Profile.

The categories themselves belong to the Knowledge Model.

Future Repository Profiles may choose different directory structures while preserving identical semantics.

------

# Compliance

A compliant Single Repository implementation SHALL satisfy the following requirements.

- Every Knowledge artifact belongs to one Knowledge Category.
- Categories preserve distinct responsibilities.
- Inbox contains observations rather than implementation work.
- Architecture describes current structure.
- ADR records significant decisions.
- Standards define engineering conventions.
- Repository layout preserves category semantics.

------

# Non-Goals

This specification intentionally does not define:

- Runtime organization;
- PRD structure;
- Issue management;
- lifecycle operations;
- Documentation Operations.

These concerns are specified elsewhere.

------

# References

- DOS-1002 — Knowledge Model
- DOS-2001 — Single Repository Profile
- DOS-2002 — Repository Layout
- DOS-2004 — Runtime Mapping
- DOS-3003 — Knowledge Impact Analysis

------

# Summary

Knowledge Mapping organizes persistent project knowledge into four independent Knowledge Categories:

- Architecture
- ADR
- Standards
- Inbox

Each category has a single, well-defined responsibility.

Together they form the persistent Knowledge domain of the Single Repository Profile while remaining independent of Runtime and implementation details.