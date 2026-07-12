# DOS-5005 — CLI

**Status:** Draft
**Version:** 1.0
**Category:** Runtime

------

# Abstract

This specification defines the Command-Line Interface (CLI) of Documentation OS.

The CLI provides a standardized user interface for invoking Documentation Engine capabilities.

Unlike the Documentation Engine, which defines execution behavior, the CLI defines how humans, AI agents, scripts, and automation interact with those capabilities.

The CLI is an interface layer.

It does not implement Documentation OS semantics.

------

# Purpose

The purpose of the CLI is to expose Documentation Engine functionality through a consistent, scriptable interface.

The CLI should enable:

- manual repository maintenance;
- automation;
- CI integration;
- AI agent invocation;
- editor integration.

The CLI should remain lightweight.

All repository logic belongs to the Documentation Engine.

------

# Scope

The CLI exposes Documentation Operations implemented by the Documentation Engine.

Typical operations include:

- Validation;
- Health evaluation;
- Migration;
- Runtime completion;
- repository inspection.

The CLI intentionally avoids embedding engineering reasoning.

------

# Design Principles

The CLI follows the following principles.

## CLI-1 Thin Interface

The CLI is a presentation layer.

Business logic belongs exclusively to the Documentation Engine.

------

## CLI-2 Deterministic

Executing identical commands against identical repository state shall produce identical results.

------

## CLI-3 Script Friendly

CLI output should support:

- shell scripting;
- automation;
- continuous integration;
- AI orchestration.

Human-readable output should remain available.

------

## CLI-4 Discoverable

Users should be able to discover available commands through the CLI itself.

Documentation should not be required for basic exploration.

------

## CLI-5 Stable

Public commands should remain stable across Documentation OS releases whenever possible.

Breaking changes should occur only through documented version transitions.

------

# Architecture

The CLI occupies the following position within Documentation OS.

```text
User / AI Agent

↓

CLI

↓

Documentation Engine

↓

Repository
```

The CLI translates user intent into Documentation Operations.

The Documentation Engine performs the work.

------

# Command Categories

Documentation OS defines the following normative command categories.

```text
CLI

├── Inspect

├── Validate

├── Health

├── Synchronize

├── Complete

├── Migrate

└── Repository
```

Individual implementations may introduce additional commands.

------

# Inspect Commands

Inspect commands expose repository information without modifying repository state.

Typical capabilities include:

- Documentation OS version;
- active Repository Profile;
- repository summary;
- Runtime overview (including active and completed workstreams, based on `.scratch/INDEX.md`);
- Knowledge overview.

Inspection shall remain read-only.

------

# Validation Commands

Validation commands execute repository Validation.

Results should include:

- overall status;
- validation failures;
- warnings;
- affected artifacts.

Validation commands shall not modify repository state.

------

# Health Commands

Health commands evaluate repository sustainability.

Typical output includes:

- Health score;
- category summaries;
- documentation debt;
- recommended improvements.

Health remains advisory.

------

# Synchronization Commands

Synchronization commands invoke Documentation Operations responsible for maintaining repository consistency.

Typical capabilities include:

- generated artifact synchronization;
- managed metadata synchronization;
- documentation synchronization.

Synchronization should remain deterministic.

------

# Complete Commands

Complete commands execute Runtime completion.

Typical behavior includes:

- moving `active/<workstream-slug>/` → `completed/<workstream-slug>/`;
- regenerating `.scratch/INDEX.md`;
- updating repository metadata.

Complete commands shall preserve Core Runtime Assets (PRD.md, issues/, HANDOFF.md) and shall not delete them.

------

# Migration Commands

Migration commands perform repository evolution.

Typical capabilities include:

- Repository Profile migration;
- Documentation OS upgrades;
- repository transformation.

Migration commands should conclude with Validation.

------

# Repository Commands

Repository commands expose repository-level capabilities.

Examples include:

- repository initialization;
- repository diagnostics;
- Documentation OS version information;
- repository configuration inspection.

Repository-specific behavior remains implementation-defined.

------

# Output Formats

The CLI should support multiple output formats.

Typical formats include:

- human-readable;
- machine-readable;
- structured formats suitable for automation.

Output representation remains implementation-defined.

Output semantics remain standardized.

------

# Exit Status

CLI commands should return deterministic exit status.

Typical outcomes include:

| Exit Status   | Meaning                                            |
| ------------- | -------------------------------------------------- |
| Success       | Operation completed successfully                   |
| Warning       | Operation completed with non-blocking observations |
| Failure       | Operation failed                                   |
| Invalid Usage | Command invocation error                           |

Specific numeric values remain implementation-defined.

------

# Automation

The CLI should support automated execution.

Typical integration scenarios include:

- continuous integration;
- Git hooks;
- scheduled maintenance;
- AI orchestration;
- editor extensions.

Automation should invoke the same Documentation Operations available to interactive users.

------

# Relationship to Documentation Engine

The CLI shall not duplicate Documentation Engine behavior.

The relationship is:

```text
CLI

↓

Documentation Engine

↓

Repository
```

The Documentation Engine remains the single source of operational behavior.

------

# Relationship to AI Agents

AI agents may invoke the CLI directly.

Alternatively, agents may invoke Documentation Operations through:

- MCP;
- native APIs;
- embedded integrations.

Invocation mechanism shall not affect Documentation Engine behavior.

------

# Compliance

A Documentation OS-compliant CLI SHALL:

- expose Documentation Engine capabilities;
- remain a thin interface layer;
- avoid embedding repository logic;
- support deterministic execution;
- support automation.

------

# Non-Goals

This specification intentionally does not define:

- command names;
- command-line syntax;
- option flags;
- configuration files;
- shell completion;
- implementation language.

These concerns belong to individual CLI implementations.

------

# References

- DOS-4001 — Documentation Operations
- DOS-4002 — Validation
- DOS-4003 — Health
- DOS-4004 — Migration
- DOS-5004 — Documentation Engine

------

# Summary

The Documentation OS CLI provides a standardized interface for interacting with the Documentation Engine.

It exposes Documentation Operations through a stable, scriptable interface while remaining intentionally thin.

By separating interface from execution, the CLI ensures that humans, AI agents, automation, and external tooling all interact with the same deterministic Documentation Engine behavior, preserving consistency across every execution environment.