# DOS-0003 — Core Principles

**Status:** Draft
**Version:** 1.0
**Category:** Foundation

------

# Abstract

This specification defines the immutable principles that govern Documentation OS.

Unlike **DOS-0002 – Design Philosophy**, which explains the reasoning behind the architecture, this specification defines the normative rules that every compliant implementation MUST satisfy.

These principles form the contractual foundation of Documentation OS.

Implementations MAY extend the system, but SHALL NOT violate these principles.

------

# Purpose

The purpose of this specification is to establish a stable set of engineering rules that remain true regardless of:

- repository profile;
- implementation language;
- documentation tooling;
- runtime environment;
- AI model;
- future extensions.

Every Documentation OS implementation is evaluated against these principles.

------

# Principle 1 — Repository Is the Only Source of Truth

The repository SHALL be the only persistent source of project knowledge.

Documentation OS SHALL NOT require external persistent state to understand the repository.

Examples of prohibited persistent state include:

- internal databases;
- synchronization databases;
- hidden metadata stores;
- duplicated knowledge repositories.

Derived artifacts MAY exist if they are reproducible from repository contents.

------

# Principle 2 — Runtime and Knowledge Are Independent Domains

Documentation OS separates information into two independent domains:

- Runtime
- Knowledge

Runtime supports execution.

Knowledge preserves understanding.

Neither domain replaces the other.

Knowledge SHALL NOT depend upon Runtime.

Runtime MAY reference Knowledge.

This directional dependency guarantees that long-term project understanding remains available after Runtime has been completed.

------

# Principle 3 — Runtime Is Finite

Runtime exists only while work is active.

Temporary Runtime content SHALL eventually be completed, synchronized into Knowledge, or discarded.

Core Runtime Assets (PRD, Issues, Handoff) SHALL be preserved upon Work completion.

Runtime SHALL NOT become permanent project storage.

A growing Runtime area indicates lifecycle failure.

------

# Principle 4 — Knowledge Represents Current Truth

Knowledge SHALL describe the current project.

Knowledge SHALL NOT become a historical archive.

Historical reasoning belongs to ADR.

Execution history belongs to Runtime.

Architecture describes the present.

Standards describe expected future behavior.

------

# Principle 5 — Knowledge Synchronization Is Mandatory

Implementation frequently produces new knowledge.

Whenever implementation changes project understanding, that knowledge SHALL be synchronized before Runtime is completed.

The minimum sequence is:

```text
Implementation

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Validation

↓
Complete
```

Skipping Knowledge Synchronization violates Documentation OS.

------

# Principle 6 — Documentation Operations Maintain Structure

Documentation Operations exist to maintain repository consistency.

Operations MAY:

- allocate identifiers;
- generate indexes;
- maintain references;
- synchronize generated files;
- validate repository structure;
- perform deterministic migration.

Documentation Operations SHALL NOT:

- invent architecture;
- create engineering standards;
- generate design rationale;
- replace engineering judgement.

Knowledge creation remains the responsibility of humans or AI agents.

------

# Principle 7 — Deterministic Information Shall Be Automated

Any information that can be deterministically derived SHALL be maintained automatically.

Examples include:

- indexes;
- reverse references;
- generated summaries;
- compatibility files;
- managed metadata.

Manual synchronization of deterministic information SHOULD be eliminated.

------

# Principle 8 — Authored Knowledge Is Authoritative

Authored knowledge SHALL remain authoritative.

Knowledge may be authored by humans or AI agents; once accepted through review it is authoritative Authored Knowledge regardless of its origin.

Generated content SHALL NOT overwrite authored content.

Documentation OS distinguishes three categories:

- Authored Knowledge
- Generated Content
- Managed Regions

Only Generated Content and Managed Regions MAY be regenerated automatically.

------

# Principle 9 — Generated Content Must Be Reproducible

Generated content SHALL satisfy deterministic generation.

Identical repository state SHALL always produce identical generated output.

Generation SHALL NOT depend on:

- timestamps;
- execution order;
- machine identity;
- operating system;
- user identity.

unless explicitly required by another specification.

------

# Principle 10 — Every Knowledge Artifact Has Ownership

Every persistent knowledge artifact SHALL possess explicit ownership.

Ownership defines:

- who updates it;
- when it changes;
- how it evolves;
- how it is validated.

Ownership SHALL be defined by lifecycle rather than by individual developer.

------

# Principle 11 — Documentation Evolves Through Lifecycle

Documentation SHALL evolve through explicit lifecycle transitions.

Documentation SHALL NOT be modified arbitrarily.

Meaningful documentation changes correspond to meaningful engineering events.

Lifecycle events define:

- creation;
- synchronization;
- validation;
- archival;
- retirement.

------

# Principle 12 — Validation Protects Correctness

Validation ensures repository correctness.

Validation SHALL verify deterministic properties.

Examples include:

- identifier uniqueness;
- reference integrity;
- generated artifact consistency;
- lifecycle consistency;
- managed region integrity.

Validation SHALL NOT perform semantic engineering review.

Semantic review belongs to humans or AI reasoning systems.

------

# Principle 13 — Health Detects Knowledge Decay

Health and Validation serve different purposes.

Validation answers:

> Is the repository structurally correct?

Health answers:

> Is project knowledge gradually becoming unhealthy?

Health SHALL report observations rather than enforce correctness.

Health SHOULD assist maintainers in preventing long-term documentation decay.

------

# Principle 14 — Explicit Is Better Than Implicit

Documentation OS favors explicit representation.

Examples include:

- explicit lifecycle;
- explicit ownership;
- explicit identifiers;
- explicit relationships;
- explicit operations.

Implicit conventions SHOULD be minimized.

Explicit systems are easier for humans and AI to understand.

------

# Principle 15 — References Improve Navigation

References exist to improve discoverability.

References SHALL NOT limit repository exploration.

Documentation OS recommends navigation paths.

It does not restrict them.

AI agents MAY explore beyond documented references whenever additional understanding is required.

------

# Principle 16 — One Knowledge System

Documentation OS maintains one shared knowledge system.

Separate documentation for:

- AI;
- humans;
- tooling;

SHOULD NOT exist.

Different consumers may present knowledge differently.

The underlying knowledge SHALL remain singular.

------

# Compliance

An implementation claiming Documentation OS compliance SHALL satisfy every principle defined in this specification.

Extensions MAY introduce additional capabilities.

Extensions SHALL NOT invalidate any existing principle.

Breaking a Core Principle constitutes a breaking change to Documentation OS itself.

------

# References

- DOS-0001 — Documentation OS
- DOS-0002 — Design Philosophy
- DOS-0004 — Terminology
- DOS-1001 — Information Model

------

# Summary

The Core Principles define the immutable engineering contract of Documentation OS.

All subsequent specifications, repository profiles, documentation engines, and tooling derive their normative behavior from these principles.

Any future evolution of Documentation OS SHALL preserve these principles unless introduced as a new major version of the specification.