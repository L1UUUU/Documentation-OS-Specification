# DOS-6003 — Conformance

**Status:** Draft
**Version:** 1.0
**Category:** Reference

------

# Abstract

This specification defines the conformance requirements for Documentation OS.

Conformance establishes the criteria by which repositories, Repository Profiles, Documentation Engines, AI agents, and related tooling may claim compatibility with Documentation OS.

Documentation OS defines behavioral conformance rather than implementation conformance.

Implementations are free to differ internally provided their externally observable behavior remains consistent with this specification.

------

# Purpose

The purpose of Conformance is to establish a common definition of compatibility.

Conformance enables:

- interoperable Documentation Engines;
- portable Repository Profiles;
- predictable AI agent behavior;
- repository portability;
- specification evolution.

Without conformance requirements, implementations cannot reliably cooperate.

------

# Scope

Conformance applies independently to:

- Documentation OS repositories;
- Repository Profiles;
- Documentation Engines;
- AI agents;
- CLI implementations.

Each implementation category possesses its own conformance requirements.

Compliance in one category does not imply compliance in another.

------

# Conformance Principles

Conformance follows the following principles.

## CF-1 Behavioral Compatibility

Documentation OS evaluates observable behavior.

Internal implementation details remain implementation-defined.

------

## CF-2 Specification Compliance

Conformance is determined exclusively by normative requirements defined within the Documentation OS Specification.

Undocumented implementation behavior SHALL NOT affect conformance.

------

## CF-3 Independent Components

Documentation OS components are evaluated independently.

Examples include:

- Repository
- Repository Profile
- Documentation Engine
- Agent
- CLI

Each component may achieve conformance separately.

------

## CF-4 Version Awareness

Conformance is always evaluated against a specific Documentation OS specification version.

Implementations SHALL declare the specification version they support.

------

# Repository Conformance

A Documentation OS Repository SHALL satisfy the following requirements.

The repository SHALL:

- expose a valid Repository Profile;
- maintain Knowledge and Runtime separation;
- preserve Knowledge artifact identities (identity-managed artifacts: ARCH, ADR, STD);
- preserve Core Runtime Assets (PRD, Issues, Handoff) upon Work completion;
- participate in documented lifecycle behavior;
- support Documentation Operations.

Repository structure SHALL conform to the active Repository Profile.

------

# Repository Profile Conformance

A Repository Profile SHALL:

- preserve Documentation OS conceptual models;
- define deterministic repository representation;
- maintain semantic equivalence;
- support Documentation Operations;
- remain compatible with Documentation OS lifecycle semantics.

Repository Profiles SHALL NOT redefine Knowledge or Runtime.

------

# Documentation Engine Conformance

A Documentation Engine SHALL:

- implement Documentation Operations;
- support Validation;
- support Health evaluation;
- support Migration;
- support Documentation Testing (provide the conformance tests required by DOS-4005);
- preserve repository consistency;
- execute deterministic operations.

Documentation Engines SHALL implement observable behavior consistent with Documentation OS.

------

# Agent Conformance

A Documentation OS-compliant AI agent SHALL:

- perform Agent Entry;
- follow the Reading Strategy;
- respect Ownership;
- follow lifecycle behavior;
- invoke Documentation Operations for deterministic maintenance.

Agents remain responsible for engineering reasoning.

Documentation Engines remain responsible for deterministic repository behavior.

------

# CLI Conformance

A Documentation OS-compliant CLI SHALL:

- expose Documentation Engine capabilities;
- remain implementation-independent;
- avoid embedding repository logic;
- support deterministic execution.

The CLI is evaluated independently from the Documentation Engine.

------

# Conformance Levels

Documentation OS defines three conformance levels.

## Level 1 — Repository Compatible

Supports:

- Repository Profile;
- documentation organization;
- lifecycle-compatible repository structure.

No Documentation Engine required.

------

## Level 2 — Engine Compatible

Supports:

- Repository compatibility;
- Documentation Engine;
- Documentation Operations;
- Validation;
- Migration;
- Health;
- Documentation Testing.

Suitable for automated documentation management.

------

## Level 3 — Ecosystem Compatible

Supports:

- Repository compatibility;
- Documentation Engine compatibility;
- AI Agent compatibility;
- CLI compatibility;
- complete Documentation OS lifecycle.

Represents full Documentation OS implementation.

------

# Conformance Matrix

| Component                | L1   | L2   | L3   |
| ------------------------ | ---- | ---- | ---- |
| Repository               | ✓    | ✓    | ✓    |
| Repository Profile       | ✓    | ✓    | ✓    |
| Documentation Engine     |      | ✓    | ✓    |
| Documentation Operations |      | ✓    | ✓    |
| Documentation Testing    |      | ✓    | ✓    |
| Validation               |      | ✓    | ✓    |
| Health                   |      | ✓    | ✓    |
| Migration                |      | ✓    | ✓    |
| AI Agent                 |      |      | ✓    |
| CLI                      |      |      | ✓    |

------

# Conformance Claims

Implementations claiming Documentation OS compatibility should declare:

- Documentation OS specification version;
- supported Repository Profile;
- implemented conformance level.

Example:

```text
Documentation OS Specification: 1.0

Repository Profile: Single Repository

Conformance Level: Level 3
```

------

# Conformance Checklist

A Documentation OS implementation should satisfy the following checklist.

## Repository

- Repository Profile implemented.
- Knowledge and Runtime separated.
- Managed documentation organized correctly.
- Lifecycle-compatible repository.

------

## Documentation Engine

- Documentation Operations implemented.
- Validation implemented.
- Health implemented.
- Migration implemented.
- Documentation Testing provided (DOS-4005).
- Deterministic behavior verified.

------

## Agent

- Agent Entry implemented.
- Reading Strategy implemented.
- Execution Contract followed.
- Documentation Operations delegated correctly.

------

## CLI

- Documentation Engine exposed.
- Repository logic delegated.
- Automation supported.
- Deterministic behavior preserved.

------

# Verification

Conformance should be verified using Documentation Testing.

Documentation Testing evaluates implementation behavior against the normative requirements defined by Documentation OS.

Testing procedures are defined in:

DOS-4005 — Documentation Testing.

------

# Non-Conformance

An implementation SHALL NOT claim Documentation OS compatibility if it knowingly violates mandatory (SHALL or MUST) requirements defined by the applicable specification version.

Implementations may extend Documentation OS provided that:

- existing normative behavior is preserved;
- extensions do not redefine existing semantics;
- extensions remain clearly distinguishable from the core specification.

------

# References

- README
- DOS-4005 — Documentation Testing
- DOS-5004 — Documentation Engine
- DOS-5005 — CLI
- DOS-6002 — Normative Language
- DOS-6004 — Versioning

------

# Summary

This specification defines how Documentation OS compatibility is evaluated.

Conformance is based on observable behavior rather than implementation details.

By defining independent conformance requirements for repositories, Repository Profiles, Documentation Engines, AI agents, and CLI implementations, Documentation OS enables a consistent ecosystem in which independently developed components can interoperate while remaining faithful to the same specification.