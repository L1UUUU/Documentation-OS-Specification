# Appendix C — Reference Repository

**Status:** Informative
**Version:** 1.0
**Category:** Appendix

------

# Purpose

This appendix provides a reference repository illustrating how Documentation OS concepts may be represented within a Single Repository Profile.

Unlike the normative specifications, this appendix is not a required repository layout.

Its purpose is to demonstrate how the concepts defined throughout Documentation OS can be organized into a practical, navigable repository.

Implementations may vary in detail while preserving the same conceptual model.

------

# Relationship to the Specification

The reference repository combines concepts defined across multiple specifications into a single example.

| Specification | Repository Representation        |
| ------------- | -------------------------------- |
| DOS-2001      | Repository Profile               |
| DOS-2002      | Repository Layout                |
| DOS-2003      | Knowledge Mapping                |
| DOS-2004      | Runtime Mapping                  |
| DOS-5001      | Agent Entry Document (AGENTS.md + CLAUDE.md)           |
| DOS-5004      | Documentation Engine integration |

This appendix is illustrative.

The normative requirements remain defined by the corresponding specifications.

------

# High-Level Repository Layout

A Documentation OS repository may resemble the following structure.

```text
repository/

├── AGENTS.md
├── CLAUDE.md
│
├── docs/
│   ├── AGENTS.md
│   ├── CLAUDE.md
│   │
│   ├── architecture/
│   ├── adr/
│   ├── standards/
│   ├── inbox/
│   │
│   ├── indexes/
│   ├── templates/
│   └── specification/
│
├── .scratch/
│   ├── AGENTS.md
│   ├── CLAUDE.md
│   ├── INDEX.md
│   │
│   ├── active/
│   │   ├── <workstream-slug>/
│   │   └── <workstream-slug>/
│   └── completed/
│       ├── <workstream-slug>/
│       └── <workstream-slug>/
│
├── src/
├── tests/
│
└── .dos/
```

This layout demonstrates one possible implementation of the Single Repository Profile.

------

# Repository Layers

The repository consists of several logical layers.

```text
Repository

├── Source Code

├── Managed Information
│   ├── Knowledge
│   └── Runtime

├── Staging Information

├── Repository Guidance

└── Documentation Engine
```

These layers are conceptual.

Their physical locations are determined by the Repository Profile.

------

# Knowledge Domain

The Knowledge Domain represents long-term repository understanding.

Example layout:

```text
docs/

├── architecture/

├── adr/

├── standards/

├── inbox/   (Staging Information, not Knowledge)

├── indexes/

└── templates/
```

The three directories architecture/, adr/, and standards/ hold Knowledge.

The inbox/ directory holds Staging Information rather than Knowledge.

Knowledge evolves continuously.

Knowledge survives individual implementations.

------

# Runtime Domain

Runtime contains active engineering work.

Example layout:

```text
.scratch/

├── active/
│   ├── <workstream-slug>/
│   └── <workstream-slug>/

├── completed/
│   ├── <workstream-slug>/
│   └── <workstream-slug>/

├── INDEX.md

├── AGENTS.md

└── CLAUDE.md
```

Each Work remains isolated.

Work state is expressed by directory location: active/ vs completed/.

Completed Runtime is moved from active/ to completed/, preserving its Core Runtime Assets (PRD, Issues, Handoff).

------

# Example Work

A Runtime Work may resemble the following structure.

```text
active/<workstream-slug>/

├── PRD.md

├── issues/
│   ├── 01-<slug>.md
│   └── 02-<slug>.md

└── HANDOFF.md
```

PRD.md holds Work requirements and scoped context.

issues/ holds execution issues tracking implementation progress.

HANDOFF.md holds cross-session execution context.

Work state is expressed by directory location (active/ vs completed/) rather than status metadata.

Core Runtime Assets (PRD, Issues, Handoff) are preserved upon Work completion; only temporary Runtime content is disposable.

No work.yaml is required; PRD.md front matter may hold Work-level relationships.

------

# Agent Entry Points

Documentation OS recommends hierarchical guidance.

Example:

```text
Repository

│

├── AGENTS.md
├── CLAUDE.md (mirror)

│

├── docs/

│     ├── AGENTS.md
│     └── CLAUDE.md (mirror)

│

└── .scratch/

      ├── AGENTS.md
      └── CLAUDE.md (mirror)
```

Each Agent Entry file narrows repository context.

Lower-level guidance supplements higher-level guidance.

------

# Documentation Engine

Documentation Engine artifacts should remain separated from repository Knowledge.

Example:

```text
.dos/

├── metadata/

├── cache/

├── generated/

├── indexes/

└── state/
```

The exact structure is implementation-defined.

Knowledge should never depend upon internal engine files.

------

# Knowledge Relationships

Knowledge Categories remain independent while referencing one another.

Example:

```text
Architecture

│

├── references ADR

│

└── references Standards
```

Relationships improve navigation.

They do not imply ownership.

------

# Runtime Relationships

Runtime references Knowledge.

Knowledge does not depend upon Runtime.

```text
Work

↓

PRD

↓

Issues

↓

Implementation

↓

Knowledge Synchronization

↓

Knowledge
```

This one-way dependency prevents Runtime from becoming long-term repository storage.

------

# Documentation Lifecycle Example

The following illustrates a typical Documentation OS workflow.

```text
Inbox

↓

Accepted Work

↓

Runtime (active/)

↓

Implementation

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Validation

↓

Complete (completed/)

↓

Cleanup

↓

Closed
```

The repository evolves without losing implementation knowledge.

------

# Documentation Operations Example

Typical Documentation Engine interaction:

```text
Agent

↓

Documentation Engine

↓

Validation

↓

Repository
```

Or:

```text
Agent

↓

Documentation Engine

↓

Synchronize

↓

Knowledge
```

The Documentation Engine performs deterministic maintenance.

The agent performs engineering reasoning.

------

# Example Navigation

An implementation agent might navigate the repository as follows.

```text
Agent Entry Document

↓

docs/AGENTS.md

↓

Architecture

↓

Referenced ADR

↓

Standards

↓

Related Runtime

↓

Source Code
```

This follows the Reading Strategy defined by Documentation OS.

------

# Example Repository Evolution

Over time, repository evolution resembles the following.

```text
Knowledge

↓

Work Created

↓

Runtime (active/)

↓

Implementation

↓

Knowledge Updated

↓

Runtime Moved to completed/

↓

Knowledge
```

Knowledge continually improves.

Active Runtime continually leaves the active execution context. Core Runtime Assets remain preserved under completed/.

------

# Implementation Notes

Repository implementations are encouraged to:

- preserve the separation of Knowledge and Runtime;
- expose clear Agent entry points;
- isolate Documentation Engine artifacts;
- preserve Core Runtime Assets (PRD, Issues, Handoff) upon Work completion;
- maintain stable identities for globally identified Artifacts;
- minimize coupling between documentation categories.

Additional directories may exist, but normative Single Repository paths and filenames SHALL remain unchanged.

------

# Relationship to Appendix A

Appendix A explains how the specifications depend upon one another.

This appendix demonstrates how those specifications may appear together inside a practical repository.

------

# Relationship to Appendix B

Appendix B explains how readers should navigate the Documentation OS Specification.

This appendix explains how AI agents and developers may navigate a Documentation OS repository.

------

# Summary

This appendix presents a reference implementation of a Documentation OS repository.

It demonstrates how Knowledge, Runtime, Agent guidance, and the Documentation Engine can coexist within a Single Repository Profile while preserving the conceptual separation defined throughout the Documentation OS Specification.

The repository illustrated here is intentionally representative rather than prescriptive.

Its purpose is to help implementers translate the Documentation OS specifications into a practical repository architecture while remaining free to choose implementation details appropriate to their own projects.