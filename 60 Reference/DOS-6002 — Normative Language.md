# DOS-6002 — Normative Language

**Status:** Draft
**Version:** 1.0
**Category:** Reference

------

# Abstract

This specification defines the normative language and glossary used throughout Documentation OS.

Its purpose is twofold:

- define the meaning of normative requirement keywords;
- define common specification vocabulary used consistently across Documentation OS.

This specification provides the interpretation rules for every normative statement within the Documentation OS Specification.

------

# Purpose

Documentation OS contains both:

- descriptive statements;
- normative requirements.

This specification distinguishes between them.

Normative keywords establish implementation requirements.

Descriptive language explains concepts.

Implementations claiming Documentation OS compliance SHALL interpret normative keywords according to this specification.

------

# Normative Language

Documentation OS adopts the following requirement keywords.

Unless explicitly stated otherwise, these keywords SHALL be interpreted exactly as defined below.

------

# SHALL

**SHALL** indicates an absolute requirement.

A compliant implementation is required to satisfy every SHALL statement.

Failure to satisfy a SHALL statement results in non-conformance.

Example:

> Every managed Artifact SHALL belong to exactly one domain.

------

# MUST

**MUST** is equivalent to SHALL.

Documentation OS treats SHALL and MUST identically.

Implementations may use either keyword for readability.

------

# SHALL NOT

**SHALL NOT** indicates an absolute prohibition.

A compliant implementation does not perform the prohibited behavior.

Example:

> Runtime SHALL NOT become permanent repository storage.

------

# MUST NOT

**MUST NOT** is equivalent to SHALL NOT.

------

# SHOULD

**SHOULD** indicates a strong recommendation.

A compliant implementation is expected to follow the recommendation unless there is a well-understood reason not to.

Deviation should be deliberate rather than accidental.

------

# SHOULD NOT

**SHOULD NOT** indicates that an action is generally discouraged.

Exceptional situations may justify deviation.

------

# MAY

**MAY** indicates optional behavior.

Implementations are free to support or omit optional behavior provided all mandatory requirements remain satisfied.

------

# RECOMMENDED

**RECOMMENDED** is equivalent to SHOULD.

------

# OPTIONAL

**OPTIONAL** is equivalent to MAY.

------

# Informative Statements

Statements that do not use normative keywords are informative.

Informative statements:

- explain concepts;
- provide examples;
- describe rationale;
- improve readability.

Informative statements do not define conformance requirements.

------

# Specification Vocabulary

The following vocabulary is used consistently throughout Documentation OS.

------

## Specification

Defined normatively in DOS-0004 — Terminology.

DOS-0004 is the single source of truth for core terminology; this section references rather than redefines it.

------

## Implementation

Defined normatively in DOS-0004 — Terminology.

------

## Conformance

The degree to which an implementation satisfies Documentation OS normative requirements.

Conformance is evaluated according to:

DOS-6003 — Conformance.

------

## Repository State

The complete observable state of a Documentation OS repository at a particular point in time.

Repository State includes:

- documentation;
- metadata;
- Runtime;
- generated artifacts.

Repository State excludes external systems.

------

## Deterministic

A process is **deterministic** when identical inputs always produce identical observable outputs.

Determinism is one of the core principles of Documentation OS.

------

## Managed Artifact

Defined as "Artifact" in DOS-0004 — Terminology.

------

## Generated Artifact

Defined as "Generated Content" in DOS-0004 — Terminology.

------

## Managed Region

Defined normatively in DOS-0004 — Terminology.

------

## Active Repository

A repository currently participating in the Documentation OS lifecycle.

An Active Repository may contain both Knowledge and Runtime.

------

## Completed Runtime

Runtime that has completed the Work Close Pipeline and no longer participates in active execution.

Completed Runtime core assets (PRD, Issues, Handoff) are preserved.

They remain immutable; only the generated INDEX.md may be regenerated.

------

## Issue Status

The normative lifecycle state of a Runtime Issue, with values open, in-progress, done, blocked (defined in DOS-2004).

------

## Compatible

An implementation is **Documentation OS Compatible** when it satisfies the applicable conformance requirements defined by this specification.

Compatibility is determined by behavior rather than implementation language.

------

# Interpretation Rules

When interpreting Documentation OS specifications:

1. Normative keywords override descriptive wording.
2. Terminology definitions override informal language.
3. Higher-level specifications SHALL NOT redefine concepts defined at lower levels; lower-level specifications SHALL NOT contradict higher-level normative requirements.

   These two clauses govern different concerns. Redefinition governs where a concept is authoritatively defined: lower layers own their concepts and higher layers build upon them without redefining them. Contradiction governs compliance priority: lower-layer specifications must obey the normative requirements that higher layers establish upon those concepts.
4. Repository Profiles SHALL preserve conceptual semantics.
5. Implementations SHALL preserve externally observable behavior even if internal implementations differ.

------

# Relationship to Terminology

This specification defines:

- normative language;
- interpretation rules.

Conceptual definitions are specified separately in:

DOS-0004 — Terminology.

------

# Relationship to Conformance

Normative keywords defined in this specification determine the mandatory requirements evaluated by:

DOS-6003 — Conformance.

------

# References

- README
- DOS-0004 — Terminology
- DOS-6003 — Conformance

------

# Summary

This specification establishes the normative language of Documentation OS.

By standardizing requirement keywords such as SHALL, MUST, SHOULD, and MAY, together with common specification vocabulary, Documentation OS ensures that all specifications are interpreted consistently by humans, AI agents, Documentation Engines, and Repository Profile implementations.

Every normative requirement in Documentation OS derives its meaning from the definitions contained in this specification.