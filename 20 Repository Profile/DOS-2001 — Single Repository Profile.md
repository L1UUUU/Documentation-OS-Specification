# DOS-2001 — Single Repository Profile

**Status:** Stable
**Version:** 1.0
**Category:** Repository Profile

------

# Abstract

This specification defines the normative Repository Profile for a single software repository.

A Repository Profile maps the abstract models defined by Documentation OS into concrete repository structures.

This specification is the reference implementation profile for Documentation OS v1.

Future profiles (such as Workspace Profile or Multi-Repository Profile) SHALL preserve the conceptual models while providing alternative repository organizations.

------

# Purpose

The purpose of the Single Repository Profile is to provide a deterministic repository organization that:

- is understandable by humans;
- is efficient for AI agents;
- minimizes duplicated knowledge;
- supports deterministic documentation operations;
- supports the complete Documentation OS lifecycle.

This specification defines repository semantics rather than engineering workflow.

------

# Relationship to the Model Layer

The previous specifications define concepts.

This specification defines one implementation.

For example:

```text
Knowledge Model
        │
        ▼
docs/
```

Runtime Model

↓

.scratch/

Information Model

↓

Repository

The conceptual models remain unchanged.

Only their repository representation is defined here.

------

# Repository Structure

A compliant Single Repository SHALL contain the following documentation structure.

```text
repository/

├── AGENTS.md
├── CLAUDE.md
│
├── docs/
│
├── .scratch/
│
└── ...
```

Additional project directories MAY exist.

Documentation OS defines only the documentation-related structure.

------

# Top-Level Components

The repository consists of three documentation components.

## AGENTS.md

Defines the entry contract for AI agents.

Responsibilities include:

- repository overview;
- documentation entry points;
- required Documentation Operations;
- execution contract;
- reading strategy.

AGENTS.md is not Managed Information.

It is Repository Guidance for agents.

A content-equivalent mirror named `CLAUDE.md` SHALL accompany `AGENTS.md` at the repository root.

Both files are valid entry points.

See DOS-5001 — Agent Entry for the synchronization requirement.

------

## docs/

The permanent Knowledge domain.

Contains:

- Architecture
- ADR
- Standards
- Inbox

Knowledge stored here persists across implementation activities.

------

## .scratch/

The Runtime domain.

Contains temporary execution artifacts.

Typical examples include:

- PRDs
- generated Issues
- execution notes
- temporary planning artifacts

The contents of `.scratch/` exist only while Work remains active.

------

# Domain Mapping

The Single Repository Profile maps the Repository Information model into concrete locations.

Managed Information domains:

| Domain    | Repository Location |
| --------- | ------------------- |
| Knowledge | docs/               |
| Runtime   | .scratch/           |

Repository Guidance:

| Concept        | Repository Location              |
| -------------- | -------------------------------- |
| Agent Contract | AGENTS.md (+ CLAUDE.md mirror)   |

This mapping is specific to the Single Repository Profile.

Other profiles may use different layouts.

------

# Knowledge Mapping

Knowledge Categories are represented beneath:

```text
docs/
```

This specification intentionally does not define the internal layout.

Knowledge Mapping is specified separately.

------

# Runtime Mapping

Runtime artifacts are represented beneath:

```text
.scratch/
```

This specification intentionally does not define internal Runtime organization.

Runtime Mapping is specified separately.

------

# Repository Responsibilities

The repository shall satisfy the following responsibilities.

## Repository stores Knowledge.

Knowledge survives Runtime.

------

## Repository stores Runtime.

Runtime supports implementation.

------

## Repository exposes an Agent Entry.

Agents should be able to initialize repository understanding from the Agent Entry Document (`AGENTS.md`, mirrored as `CLAUDE.md`).

------

## Repository supports deterministic Documentation Operations.

Operations shall function without requiring external persistent state.

------

# Repository Independence

This specification intentionally avoids defining:

- Architecture layout;
- ADR numbering;
- Runtime organization;
- Inbox organization;
- generated files.

These concerns are delegated to dedicated Repository Profile specifications.

------

# Compliance

A compliant Single Repository implementation SHALL:

- expose one Knowledge domain;
- expose one Runtime domain;
- provide an Agent Entry;
- preserve the conceptual models defined by Documentation OS.

Additional repository content MAY exist provided these guarantees remain valid.

------

# Future Profiles

Future Documentation OS profiles may include:

- Workspace Profile
- Monorepo Profile
- Multi-Repository Profile
- Cloud Documentation Profile

These profiles SHALL preserve the semantics defined by the Model layer.

Only repository representation may differ.

------

# References

- DOS-1001 — Information Model
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-2002 — Repository Layout
- DOS-2003 — Knowledge Mapping
- DOS-2004 — Runtime Mapping

------

# Summary

The Single Repository Profile is the reference implementation profile for Documentation OS.

It maps:

- Knowledge → `docs/`
- Runtime → `.scratch/`
- Agent Entry → `AGENTS.md` (+ `CLAUDE.md` mirror)

while preserving the implementation-independent semantics defined by the Documentation OS Model layer.