# DOS-2005 — Repository Conventions

**Status:** Stable
**Version:** 1.0
**Category:** Repository Profile

------

# Abstract

This specification defines the repository-wide conventions for the Single Repository Profile.

While previous specifications define repository structure and information organization, this specification establishes the common conventions that ensure consistency across the entire repository.

Repository Conventions provide predictable behavior for humans, AI agents, and Documentation Operations.

------

# Purpose

The purpose of Repository Conventions is to ensure that all repositories implementing the Single Repository Profile follow a consistent set of engineering conventions.

Consistent conventions improve:

- readability;
- maintainability;
- discoverability;
- automation;
- interoperability.

Repository Conventions intentionally avoid prescribing engineering methodology.

They define repository behavior rather than development process.

------

# Convention Principles

Repository Conventions follow the following principles.

## RC-1 Consistency Over Preference

Repository-wide consistency is preferred over individual preference.

When multiple valid approaches exist, one should be selected and applied consistently.

------

## RC-2 Explicit Structure

Repository conventions should be explicit rather than implicit.

Repository users should not rely on undocumented conventions.

------

## RC-3 Deterministic Organization

Repository organization should produce identical structure under identical circumstances.

Documentation Operations should not rely on heuristic behavior.

------

## RC-4 Human and Machine Readability

Repository conventions should remain understandable by both humans and automated tooling.

Readable repositories reduce onboarding effort and improve automation reliability.

------

# Naming Conventions

Repository names should be descriptive.

Document names should clearly communicate their purpose.

Names should avoid implementation-specific terminology whenever possible.

Generated artifacts should follow deterministic naming rules.

Repository Profiles may define additional naming conventions for specific document categories.

------

# Numbering

Managed documentation categories requiring stable ordering should use deterministic numbering.

Examples include:

- ADRs
- Architecture documents
- Standards
- Work artifacts

Number allocation shall be monotonic.

Allocated numbers shall never be reused.

Specific numbering formats are defined by individual document specifications.

------

# Repository Entry Points

The repository exposes the following standard entry points.

| Entry     | Purpose              |
| --------- | -------------------- |
| AGENTS.md | Agent entry          |
| docs/     | Persistent Knowledge |
| .scratch/ | Runtime              |

Additional repository content may exist.

These entry points remain stable throughout repository evolution.

------

# Local Guidance

Subdirectories may contain local guidance documents.

Local guidance refines repository behavior within a limited scope.

Local guidance shall not contradict repository-level guidance.

Repository-level guidance always remains authoritative.

------

# Generated Content

Automatically generated content should remain clearly distinguishable from human-authored content.

Generated regions should be explicitly identified.

Implementations should avoid modifying human-authored content outside managed regions.

Generated content specifications are defined separately.

------

# Documentation Operations

Repository Conventions assume deterministic Documentation Operations.

Operations should:

- preserve repository consistency;
- avoid destructive behavior;
- maintain generated artifacts;
- preserve human-authored knowledge.

Repository Conventions do not define individual operations.

Operation behavior is defined by the Documentation Operations specifications.

------

# Compatibility

The Single Repository Profile includes compatibility mechanisms where required.

Examples include:

- compatibility guidance files;
- generated navigation artifacts;
- repository metadata.

Compatibility mechanisms should always be derivable from repository contents.

------

# Repository Evolution

Repositories evolve continuously.

Repository evolution should prioritize:

- stable entry points;
- backward-compatible organization;
- minimal structural disruption.

Large structural changes should occur only when justified by significant architectural improvements.

------

# Compliance

A compliant Single Repository implementation SHALL satisfy the following requirements.

- Repository entry points shall remain stable.
- Naming shall follow documented conventions.
- Managed numbering shall remain deterministic.
- Local guidance shall not contradict repository guidance.
- Generated content shall remain distinguishable from human-authored content.

------

# Non-Goals

This specification intentionally does not define:

- documentation templates;
- Markdown formatting;
- lifecycle behavior;
- documentation operations;
- implementation tooling.

These concerns are specified elsewhere.

------

# References

- DOS-2001 — Single Repository Profile
- DOS-2002 — Repository Layout
- DOS-2003 — Knowledge Mapping
- DOS-2004 — Runtime Mapping
- DOS-4001 — Documentation Operations

------

# Summary

Repository Conventions establish the common behavioral rules shared by all Single Repository implementations.

Together with the Repository Layout, Knowledge Mapping, and Runtime Mapping, these conventions provide a predictable and deterministic repository environment for both humans and AI agents while preserving the conceptual models defined by Documentation OS.