# DOS-2005 — Repository Conventions

**Status:** Draft
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

When multiple valid approaches exist, one SHOULD be selected and applied consistently.

------

## RC-2 Explicit Structure

Repository conventions SHOULD be explicit rather than implicit.

Repository users SHOULD NOT rely on undocumented conventions.

------

## RC-3 Deterministic Organization

Repository organization SHOULD produce identical structure under identical circumstances.

Documentation Operations SHOULD NOT rely on heuristic behavior.

------

## RC-4 Human and Machine Readability

Repository conventions SHOULD remain understandable by both humans and automated tooling.

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

Managed documentation categories requiring stable ordering use deterministic numbering.

Examples include:

- ADRs
- Architecture documents
- Standards

Number allocation SHALL be monotonic.

Allocated numbers SHALL NOT be reused.

## Allocation Timing and Concurrent Work

A Knowledge artifact (for example an ADR) MAY be authored as a draft while its Work is active. To avoid numbering collisions between parallel Works, a draft artifact SHALL NOT pre-allocate a final identifier; it uses a non-final local placeholder (for example `ADR-DRAFT`, or a Work-scoped working name) until integration.

Final identifier allocation SHALL be performed by the Documentation Engine at integration time (PR preparation or the equivalent merge-to-trunk step), against the integration target repository state. Allocation SHALL be:

- atomic — a single deterministic step assigns the next available number to one artifact at a time;
- monotonic — a newly allocated number is greater than all previously allocated numbers in its category;
- non-reusable — a number withdrawn after its artifact is retired SHALL NOT be reissued.

If two integrations collide (concurrent branches each produced an artifact in the same category), the later integration SHALL receive a newly allocated number; the Documentation Engine SHALL regenerate all managed references to the renamed artifact so that long-lived identity references remain valid.

------

# Identifier Format

The Single Repository Profile assigns every identity-managed Knowledge artifact a stable identifier of the form:

```text
PREFIX-NNNN
```

`PREFIX` identifies the Knowledge Category, and `NNNN` is a zero-padded four-digit monotonic number.

The normative prefixes are:

| Prefix | Category        | Example   |
| ------ | --------------- | --------- |
| ARCH   | Architecture    | ARCH-0003 |
| ADR    | ADR             | ADR-0007  |
| STD    | Standards       | STD-0012  |

Inbox observations do not require a stable identifier; they may use lightweight local names until promoted into a Category.

Runtime Work uses a workstream slug rather than a global identifier.

------

# Identifier Resolution

The Single Repository Profile resolves an identifier to a repository path by prefix:

| Identifier | Path                          |
| ---------- | ----------------------------- |
| ARCH-0003  | docs/architecture/0003-*.md   |
| ADR-0007   | docs/adr/0007-*.md            |
| STD-0012   | docs/standards/0012-*.md      |

For documents, the numeric segment forms the leading segment of the filename.

The remaining filename segment is human-readable and non-normative.

Runtime Work uses a workstream slug for resolution:

| Workstream slug          | Path                                    |
| ------------------------ | --------------------------------------- |
| `<workstream-slug>`      | `.scratch/active/<workstream-slug>/`    |
| `<workstream-slug>`      | `.scratch/completed/<workstream-slug>/` |

Long-lived Knowledge references SHALL target identifiers rather than filenames. Runtime references SHALL use workstream slugs or Work-scoped paths (see DOS-2004); Runtime Works are not addressable by identifiers.

------

# Work Outcome

A Runtime Work records its terminal result in the PRD's front matter `outcome` field when it reaches `Completed`. Normative `outcome` values are:

- `succeeded`
- `cancelled`
- `superseded`
- `failed`

An active Work SHALL NOT declare an `outcome`. The field is set by the Complete stage (DOS-3004) and surfaced by INDEX.md for Completed Works (DOS-2004).

------

# Relationship Representation

Relationships for Knowledge documents are declared in a YAML front matter block at the top of each Knowledge document.

A minimal Knowledge document example:

```yaml
---
id: ADR-0007
status: active
relationships:
  - type: affects
    target: ARCH-0002
---
```

`id` is the artifact identifier.

`status` is the Knowledge artifact lifecycle state.

Valid status values for Knowledge front matter are:

- `created`
- `active`
- `updated`
- `validated`
- `archived`
- `retired`

`relationships` is a list of `{ type, target }` entries, where `target` references another identifier.

Valid relationship types are:

- `references`
- `depends-on`
- `implements`
- `produces`
- `affects`
- `supersedes`

These relationship types form the normative vocabulary for expressing relationships between managed artifacts.

Relationship type semantics are defined by DOS-1005 — Relationship Model.

Work-level relationships are declared in the PRD.md front matter at the top of the PRD document.

Issues and HANDOFF documents do not require front matter unless they need to declare specific relationships.

Reverse references and navigation indexes MAY be generated from these declarations by Documentation Operations.

------

# Repository Entry Points

The repository exposes the following standard entry points.

| Entry      | Purpose                            |
| ---------- | ---------------------------------- |
| AGENTS.md  | Agent entry (canonical source)     |
| CLAUDE.md  | Agent entry (mirror of AGENTS.md)  |
| docs/      | Persistent Knowledge               |
| .scratch/  | Runtime                            |

Additional repository content may exist.

These entry points remain stable throughout repository evolution.

------

# Local Guidance

Subdirectories may contain local guidance documents.

Local guidance refines repository behavior within a limited scope.

Local guidance SHALL NOT contradict repository-level guidance.

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

The `CLAUDE.md` mirror of `AGENTS.md` is one such compatibility mechanism: it lets Claude-compatible tooling reuse the same Agent Entry Document without duplicating its content.

Examples include:

- `CLAUDE.md` as a content-equivalent mirror of `AGENTS.md`;
- generated navigation artifacts;
- repository metadata.

Compatibility mechanisms should always be derivable from repository contents.

The recommended mechanism to keep `CLAUDE.md` and `AGENTS.md` equivalent is a symbolic link (`CLAUDE.md` → `AGENTS.md`); generators, hooks, or Validation rules MAY be used instead.

Because symbolic links are not reliably preserved across all platforms and version-control clients, equivalence SHALL NOT depend solely on link behavior.

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

- Repository entry points SHALL remain stable.
- Naming SHALL follow documented conventions.
- Managed numbering SHALL remain deterministic.
- Managed Knowledge artifacts SHALL declare an identifier in the format defined above.
- Relationships for Knowledge documents SHALL be declared through front matter targeting identifiers.
- Work-level relationships SHALL be declared in the PRD.md front matter or referenced via workstream slug.
- Local guidance SHALL NOT contradict repository guidance.
- Generated content SHALL remain distinguishable from human-authored content.

------

# Non-Goals

This specification intentionally does not define:

- full document templates;
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