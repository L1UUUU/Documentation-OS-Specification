# DOS-3003 — Knowledge Impact Analysis

**Status:** Draft
**Version:** 1.0
**Category:** Lifecycle

------

# Abstract

This specification defines Knowledge Impact Analysis (KIA), the lifecycle stage responsible for determining how implementation or execution activity — whether it ended successfully, was abandoned, or was superseded — affects the repository's persistent Knowledge.

Knowledge Impact Analysis bridges the Runtime and Knowledge domains.

Its purpose is not to modify documentation directly, but to determine **whether** Knowledge should change, **which** Knowledge Categories are affected, and **what** synchronization actions are required.

Knowledge Impact Analysis is mandatory for every completed Work.

------

# Purpose

A Work's implementation or execution activity frequently changes repository understanding.

However, not every Work requires documentation updates.

Knowledge Impact Analysis exists to answer one fundamental question:

> **What has the repository learned from this Work?**

Only after this question has been answered can Knowledge Synchronization occur.

------

# Scope

Knowledge Impact Analysis applies to every Work whose implementation or execution activity has ended (achieved, abandoned, or superseded), including a Work being terminated before its objectives were reached.

It evaluates the impact of the Work's implementation or execution activity on the Knowledge domain.

It does **not** perform documentation updates.

It does **not** perform repository validation.

It produces a structured impact assessment that becomes the input to Knowledge Synchronization.

------

# Lifecycle Position

Knowledge Impact Analysis occupies the following position within the Runtime Lifecycle.

```text
Implementation  (MAY be partial or absent for a non-succeeded outcome)

↓

Knowledge Impact Analysis

↓

Knowledge Synchronization

↓

Validation
```

Implementation MAY be partial or absent for a Work terminated before its objectives were reached; Knowledge Impact Analysis then evaluates whatever activity did occur (frequently no impact). Knowledge Impact Analysis SHALL occur exactly once for every completed Work.

------

# Objectives

Knowledge Impact Analysis has four primary objectives.

## KIA-1 Detect Knowledge Changes

Determine whether the Work's implementation or execution activity has changed repository understanding.

------

## KIA-2 Identify Affected Categories

Identify which Knowledge Categories require synchronization.

Possible categories include:

- Architecture
- ADR
- Standards

A Work's implementation or execution activity may also surface unresolved concerns that do not fit any category; these are recorded as Inbox staging items rather than Knowledge.

------

## KIA-3 Prevent Documentation Drift

Ensure that repository Knowledge remains synchronized with the Work's implementation or execution activity.

------

## KIA-4 Minimize Unnecessary Updates

Documentation should only change when repository understanding changes.

Cosmetic implementation or execution changes should not trigger unnecessary documentation updates.

------

# Inputs

Knowledge Impact Analysis evaluates information from multiple sources.

Typical inputs include:

- implementation or execution activity that has ended (successful, abandoned, or superseded);
- Runtime artifacts;
- implementation notes;
- repository Knowledge;
- existing ADRs;
- repository Standards.

Repository Profiles may define additional inputs.

------

# Outputs

Knowledge Impact Analysis produces a structured impact assessment.

The assessment addresses the three Knowledge Categories below, and may additionally identify unresolved concerns to record as Inbox staging items.

The assessment should answer the following questions.

## Architecture

Has repository structure changed?

If yes:

- which architectural areas changed;
- which documents require synchronization.

------

## ADR

Was a significant engineering decision made?

If yes:

- should an ADR be created?
- should an existing ADR be updated?

------

## Standards

Did the Work's implementation or execution activity introduce or modify engineering conventions?

If yes:

- which Standards require synchronization?

------

## Inbox (Staging Output)

Did the Work's implementation or execution activity discover unresolved repository concerns?

If yes:

- should new Inbox observations be created?

------

# Decision Matrix

Knowledge Impact Analysis determines required synchronization actions.

Typical outcomes include:

| Observation                         | Action                                                                                                            |
| ----------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| No repository understanding changed | Execute Knowledge Synchronization as a no-change operation; record that no Knowledge edits were required         |
| Architecture changed                | Update Architecture                                                                                               |
| Significant design decision         | Create or update ADR                                                                                              |
| Engineering convention changed      | Update Standards                                                                                                  |
| New unresolved concern discovered   | Record Inbox staging item                                                                                         |

Multiple outcomes may occur within a single Work.

The Knowledge Synchronization stage is mandatory for every completed Work and SHALL NOT be skipped (DOS-3004). A no-impact KIA result therefore yields a no-change synchronization run — the stage is executed and records an explicit no-change result — rather than a skipped stage.

An explicit no-change result SHALL be exposed as deterministic Documentation Operation output, so that a no-change Work is observable rather than silent. Repository Profiles or implementations MAY additionally persist it in HANDOFF.md or another defined Runtime record.

------

# Responsibilities

Knowledge Impact Analysis is responsible for:

- evaluating repository impact;
- identifying affected Knowledge Categories;
- preparing synchronization decisions.

Knowledge Impact Analysis is **not** responsible for:

- editing documentation;
- repository validation;
- complete management.

These responsibilities belong to later lifecycle stages.

------

# Category Evaluation

Every Knowledge Category should be evaluated independently.

Typical evaluation sequence:

```text
Architecture

↓

ADR

↓

Standards
```

Unresolved concerns that do not fit any category are recorded as Inbox staging items after category evaluation.

Repository implementations may use different evaluation order provided identical results are produced.

------

# Relationship Awareness

Knowledge Impact Analysis should consider documented Relationships.

For example:

- updating Architecture may affect Standards;
- creating an ADR may require Architecture references;
- modifying Standards may affect Architecture guidance.

Relationship analysis improves synchronization completeness.

Relationship semantics are defined separately.

------

# Invariants

The following invariants SHALL always remain true.

## KI-1

Knowledge Impact Analysis occurs after a Work's implementation or execution activity has ended (successfully, by abandonment, or by supersession).

------

## KI-2

Knowledge Impact Analysis occurs before Knowledge Synchronization.

------

## KI-3

Knowledge Impact Analysis never modifies Knowledge directly.

------

## KI-4

Every affected Knowledge Category SHALL be explicitly identified.

------

## KI-5

No Runtime artifact SHALL be completed before Knowledge Impact Analysis has completed.

------

# Documentation Operations

Documentation Operations may assist Knowledge Impact Analysis by:

- locating related documentation;
- identifying affected artifacts;
- validating references;
- generating impact reports.

Documentation Operations SHALL NOT determine engineering impact autonomously.

Engineering judgement remains the responsibility of humans or AI agents.

------

# Failure Handling

Knowledge Impact Analysis may conclude that:

- no Knowledge changes are required;
- synchronization cannot proceed;
- additional engineering clarification is required.

In such cases, the Work SHALL NOT bypass the remaining lifecycle.

Failure should remain explicit.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every completed Work performs Knowledge Impact Analysis;
- affected Knowledge Categories are explicitly identified;
- Knowledge is not modified during analysis;
- Knowledge Synchronization receives a complete impact assessment.

------

# Non-Goals

This specification intentionally does not define:

- documentation editing;
- synchronization algorithms;
- repository validation;
- complete behavior;
- engineering review methodology.

These concerns are specified separately.

------

# References

- DOS-1002 — Knowledge Model
- DOS-1005 — Relationship Model
- DOS-2003 — Knowledge Mapping
- DOS-3002 — Runtime Lifecycle
- DOS-3004 — Work Close Pipeline
- DOS-4001 — Documentation Operations

------

# Summary

Knowledge Impact Analysis is the bridge between Runtime and Knowledge.

It evaluates what a Work's implementation or execution activity has taught the repository, determines which Knowledge Categories require synchronization (or records an explicit no-change result), and prepares the repository for the Knowledge Synchronization stage.

Knowledge Impact Analysis never modifies documentation directly.

It determines **what should change**, allowing subsequent lifecycle stages to update the repository in a controlled and deterministic manner.