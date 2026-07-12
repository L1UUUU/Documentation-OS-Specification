# DOS-1003 — Runtime Model

**Status:** Draft
**Version:** 1.0
**Category:** Model

------

# Abstract

This specification defines the Runtime Model of Documentation OS.

Runtime represents the temporary execution domain of a software project. It exists to support planning, clarification, implementation, collaboration, and delivery.

Unlike Knowledge, Runtime is transient by design. Active Runtime artifacts exist only while Work remains active, and are eventually completed, discarded, or transformed into persistent Knowledge.

The Runtime Model defines the semantics of Runtime independently of any repository layout or implementation profile.

------

# Purpose

The purpose of the Runtime Model is to establish a consistent semantic model for all temporary project information.

Runtime provides the execution context required to evolve a repository while preventing temporary artifacts from polluting long-term project knowledge.

Runtime is execution-oriented.

Knowledge is understanding-oriented.

------

# Design Goals

The Runtime Model is designed to satisfy the following objectives.

## RT-1 Support Active Work

Runtime exists to enable implementation.

It should provide sufficient context for humans and AI agents to execute work efficiently.

------

## RT-2 Preserve Temporary Context

Implementation often requires temporary information that should not become permanent project knowledge.

Runtime provides an explicit home for these artifacts.

------

## RT-3 Enable Structured Execution

Runtime should organize execution into well-defined units.

Planning, clarification, implementation, and closure should follow explicit lifecycle transitions.

------

## RT-4 Prevent Knowledge Pollution

Temporary execution artifacts should never become permanent repository knowledge simply because they were useful during implementation.

Runtime exists specifically to isolate temporary information.

------

## RT-5 Produce Better Knowledge

Runtime is valuable because it creates new understanding.

Successful Runtime should improve the project's Knowledge through Knowledge Synchronization.

------

# Definition

Runtime is the temporary information domain responsible for supporting active engineering work.

Runtime contains execution context rather than persistent understanding.

Typical Runtime information includes:

- active planning;
- requirement clarification;
- implementation decomposition;
- execution artifacts;
- temporary analysis.

Active Runtime exists only while Work remains active.

------

# Runtime Characteristics

Every Runtime artifact possesses the following characteristics.

------

## Temporary

Runtime has limited lifetime as an active execution context.

Every Runtime artifact SHALL eventually leave the active execution context.

After leaving the active context, a Runtime artifact transitions into one of:

- Completed Runtime — core assets preserved (immutable business content);
- discarded — removed because it carries no lasting value;
- transformed — promoted into persistent Knowledge through Knowledge Synchronization.

Runtime is never intended to become permanent Managed Information.

Completed Runtime core assets (PRD, Issues, Handoff) are preserved as immutable historical records, but they are no longer active Runtime and do not participate in future execution.

The INDEX.md at `.scratch/INDEX.md` is generated content that may be refreshed by the Documentation Engine.

------

## Execution-Oriented

Runtime supports execution.

It provides context necessary to complete implementation work.

Unlike Knowledge, Runtime does not attempt to describe the repository indefinitely.

------

## Mutable

Runtime changes frequently.

Planning evolves.

Issues change.

Clarifications are refined.

Execution status progresses.

Frequent modification is expected.

------

## Disposable

Active Runtime artifacts may be discarded once their purpose has been fulfilled.

Runtime distinguishes between Core Runtime Assets and Ephemeral Runtime Content:

- Core Runtime Assets (PRD, Issues, Handoff) are preserved upon Work completion and retained in the completed/ directory.
- Ephemeral Runtime Content may be discarded once its purpose is fulfilled.

Completion preserves Core Runtime Assets as immutable history; discarding removes Ephemeral content entirely.

Both outcomes end active Runtime status.

Preservation of active Runtime is not the primary objective; Knowledge extraction is.

------

## Knowledge-Producing

Runtime is expected to generate new project understanding.

That understanding is transferred into the Knowledge domain through Knowledge Synchronization.

------

# Runtime Structure

Documentation OS intentionally defines Runtime as an abstract domain.

The internal organization of Runtime is determined by Repository Profiles.

However, every Runtime implementation should contain equivalent conceptual responsibilities.

Typical Runtime responsibilities include:

- incoming work;
- active work;
- implementation planning;
- execution artifacts;
- closure.

Within each Work, the Core Runtime Assets are:

- PRD.md — the canonical entry document and Work definition;
- issues/ — individual execution items (NN-<slug>.md format);
- HANDOFF.md — execution context for cross-agent/session/phase handoff.

These Core Runtime Assets are preserved upon Work completion.

Repository Profiles determine how these responsibilities are represented.

------

# Work

Work is the primary execution unit within Runtime.

A Work represents one bounded engineering activity.

A Work may include:

- one clarified objective;
- one implementation plan;
- one or more implementation tasks;
- temporary execution artifacts.

Work begins when an accepted effort enters Runtime and its active/<workstream-slug>/ workspace is created.

Work reaches its Completed terminal state when the Complete stage moves its workspace to `completed/`; the subsequent Cleanup stage finalizes repository state and may be retried independently.

------

# Runtime Lifecycle

Runtime participates in a lifecycle distinct from Knowledge.

The conceptual lifecycle is:

```text
Idea

↓

Clarification

↓

Planning

↓

Implementation

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Validation

↓

Completed
```

Repository Profiles may implement this lifecycle differently while preserving identical semantics.

------

# Runtime and Knowledge

Runtime and Knowledge serve complementary purposes.

| Runtime            | Knowledge           |
| ------------------ | ------------------- |
| Temporary          | Persistent          |
| Execution          | Understanding       |
| Mutable            | Stable              |
| Work-Oriented      | Repository-Oriented |
| Produces Knowledge | Preserves Knowledge |

Neither domain replaces the other.

Successful Runtime continuously improves Knowledge.

------

# Runtime Dependencies

Runtime may depend upon:

- Knowledge;
- repository source code;
- implementation artifacts.

Knowledge SHALL NOT depend upon Runtime.

This directional dependency prevents long-term understanding from becoming coupled to temporary execution artifacts.

------

# Runtime Ownership

Runtime ownership belongs to active Work.

Ownership includes responsibility for:

- maintaining execution context;
- updating implementation progress;
- initiating Knowledge Impact Analysis;
- completing the Work Close Pipeline.

Ownership terminates when Runtime is completed.

Knowledge ownership then continues independently.

------

# Runtime Completion

Runtime completion does not occur when implementation finishes.

Runtime completion occurs when:

1. implementation is complete;
2. Knowledge Impact Analysis has been performed;
3. Knowledge Synchronization has finished;
4. repository validation succeeds;
5. The Work directory has moved from active/<workstream-slug>/ to completed/<workstream-slug>/.

The Work is considered Completed upon completion of the Complete stage (directory movement). INDEX.md regeneration belongs to the subsequent Cleanup stage, which is idempotent and may be retried independently; it is not a precondition for Runtime Completion.

Completion therefore represents successful transition rather than successful coding.

------

# Repository Independence

The Runtime Model intentionally avoids specifying:

- repository directories;
- file names;
- document templates;
- implementation tooling.

These concerns belong to Repository Profiles.

For example, the Single Repository Profile may implement Runtime using dedicated repository locations.

Another profile may choose a completely different organization while preserving the Runtime semantics defined here.

------

# Compliance

A Documentation OS implementation SHALL satisfy the following requirements.

- Runtime SHALL remain separate from Knowledge.
- Runtime SHALL support explicit Work.
- Runtime SHALL participate in lifecycle transitions.
- Runtime SHALL eventually leave the active execution context.
- Runtime SHALL produce Knowledge through Knowledge Synchronization.
- Active Runtime SHALL NOT become permanent repository storage.
- Completed Runtime core assets SHALL be preserved.
- Knowledge SHALL NOT depend upon Runtime.

------

# Non-Goals

This specification intentionally does not define:

- repository layout;
- directory names;
- planning document formats;
- issue formats;
- execution tooling;
- project management methodology.

These concerns are specified by Repository Profiles and implementation specifications.

------

# References

- DOS-0001 — Documentation OS
- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles
- DOS-1001 — Information Model
- DOS-1002 — Knowledge Model
- DOS-2004 — Runtime Mapping
- DOS-3002 — Runtime Lifecycle
- DOS-3004 — Work Close Pipeline

------

# Summary

The Runtime Model defines the temporary execution domain of Documentation OS.

Runtime exists to support engineering work rather than preserve project understanding.

It is:

- temporary;
- execution-oriented;
- mutable;
- disposable;
- knowledge-producing.

Successful Runtime concludes by transferring newly created understanding into the Knowledge domain before leaving the repository's active execution context.