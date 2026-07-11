# DOS-3001 — Document Lifecycle

**Status:** Stable
**Version:** 1.0
**Category:** Lifecycle

------

# Abstract

This specification defines the lifecycle of documentation artifacts within Documentation OS.

Every managed documentation artifact evolves through a well-defined lifecycle.

The lifecycle establishes when documents are created, updated, validated, archived, and retired.

Lifecycle management ensures that documentation evolves together with the repository rather than becoming disconnected from implementation.

------

# Purpose

The purpose of the Document Lifecycle is to provide a consistent lifecycle model for all managed documentation.

A document should never appear, change, or disappear without a corresponding lifecycle event.

Lifecycle management improves:

- consistency;
- traceability;
- automation;
- maintainability.

------

# Scope

This specification applies to all managed documentation artifacts, including:

- Architecture
- ADR
- Standards
- Inbox items
- Runtime artifacts

Repository-generated artifacts participate only in the portions of the lifecycle applicable to generated content.

------

# Lifecycle Principles

The Document Lifecycle follows the principles established by:

- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles

In particular:

- Documentation is lifecycle-driven.
- Knowledge evolves continuously.
- Runtime is temporary.
- Documentation Operations maintain consistency.

------

# Lifecycle States

Every managed document progresses through one or more of the following conceptual states.

```text
Created

↓

Active

↓

Updated

↓

Validated

↓

Archived

↓

Retired
```

Not every document necessarily reaches every state.

However, lifecycle transitions shall remain explicit.

------

# State Definitions

## Created

The document has been introduced into the repository.

Identity has been allocated.

Ownership has been established.

The document is not yet considered authoritative.

------

## Active

The document actively participates in repository knowledge or execution.

Active documents may evolve continuously.

Most repository documents spend the majority of their lifetime in this state.

------

## Updated

The document has been modified as the result of a lifecycle event.

Examples include:

- implementation;
- architectural change;
- engineering decision;
- standards update.

Updating a document does not imply correctness.

Validation must follow.

------

## Validated

Repository validation has confirmed that the document satisfies structural requirements.

Validation verifies:

- identifiers;
- references;
- managed regions;
- generated content;
- repository consistency.

Validation does not certify engineering correctness.

------

## Archived

The document is no longer active but remains intentionally preserved.

Examples include:

- completed Runtime artifacts;
- historical execution records.

Knowledge documents are archived only in exceptional circumstances.

Runtime documents commonly enter this state.

------

## Retired

The document has permanently left the active documentation system.

Retired documents are no longer maintained.

Their identities remain reserved.

Retirement differs from deletion.

Retirement preserves repository history.

------

# Lifecycle Events

Lifecycle transitions occur only through explicit events.

Typical events include:

| Event       | Result    |
| ----------- | --------- |
| Create      | Created   |
| Synchronize | Updated   |
| Validate    | Validated |
| Archive     | Archived  |
| Retire      | Retired   |

Documentation shall not transition implicitly.

------

# Lifecycle Ownership

Every document possesses lifecycle ownership.

Ownership determines:

- who may modify the document;
- which events affect the document;
- when synchronization occurs;
- when validation is required.

Ownership is defined by documentation category rather than individual contributor.

------

# Knowledge Lifecycle

Knowledge documents typically follow:

```text
Create

↓

Active

↓

Update

↓

Validate

↓

Active
```

Knowledge normally remains active throughout repository lifetime.

Knowledge is rarely archived.

------

# Runtime Lifecycle

Runtime documents typically follow:

```text
Create

↓

Active

↓

Update

↓

Validate

↓

Archive

↓

Retire
```

Runtime is expected to complete its lifecycle.

Persistent Runtime indicates unfinished work.

------

# Generated Documents

Generated documentation participates in the lifecycle as managed artifacts.

Generation is considered an Update event.

Generated documents should always be reproducible.

Manual modification outside managed regions is prohibited.

------

# Lifecycle and Documentation Operations

Documentation Operations may trigger lifecycle transitions.

Typical examples include:

- identifier allocation;
- validation;
- generated content update;
- archive operations.

Operations maintain lifecycle consistency.

They do not determine engineering intent.

------

# Lifecycle Invariants

The following invariants shall always remain true.

## LI-1

Every managed document has exactly one current lifecycle state.

------

## LI-2

Lifecycle transitions are explicit.

------

## LI-3

Validation follows modification.

------

## LI-4

Archived Runtime shall not become Active again.

If additional work is required, a new Runtime artifact should be created.

------

## LI-5

Retired identities shall never be reused.

------

# Compliance

A compliant Documentation OS implementation SHALL ensure:

- every managed document participates in a lifecycle;
- lifecycle transitions are explicit;
- lifecycle ownership is defined;
- validation follows updates;
- archived Runtime remains immutable.

------

# Non-Goals

This specification intentionally does not define:

- Runtime execution workflow;
- Knowledge synchronization rules;
- Work lifecycle;
- Documentation Operations.

These concerns are specified separately.

------

# References

- DOS-0002 — Design Philosophy
- DOS-0003 — Core Principles
- DOS-1002 — Knowledge Model
- DOS-1003 — Runtime Model
- DOS-3002 — Runtime Lifecycle
- DOS-4002 — Validation

------

# Summary

The Document Lifecycle defines how managed documentation evolves throughout its lifetime.

Every document progresses through explicit lifecycle states governed by ownership, validation, and Documentation Operations.

Lifecycle management ensures that documentation remains consistent, traceable, and synchronized with repository evolution.