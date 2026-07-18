# Unreleased lifecycle hardening

- Added `LifecycleStage`, `LifecycleError`, and `FailureStageOf` so consumers
  can distinguish Begin, Synchronize, Validate, Complete, and post-commit
  Cleanup failures without parsing text or maintaining adapter-local names.
- Added `ValidationReport.Failure()` for workflow gates while preserving the
  existing `Validate` report-as-data API.
- Missing and non-directory repository roots now use the stable
  `invalid-repository` code and preserve their filesystem cause. Expected
  unmaterialized worktrees remain a consumer policy decision.
- Cleanup failure remains retriable and is now explicitly reported as the
  `cleanup` stage after terminal Work state has been persisted.

# Documentation Engine v0.1.0-rc.3

- Added `BeginWork(BeginInput)` for atomic creation of caller-defined PRD,
  HANDOFF, initial Issues, and INDEX.
- Identical active Begin retries return the persisted result without writes;
  retries with different input return the stable `conflict` error category.
- Work assets are prepared in an exclusive staging directory and published by
  atomic rename under the repository lock, so concurrent calls cannot overwrite
  or remove each other's Work.
- A retry repairs the INDEX after interruption between Work publication and
  INDEX generation; abandoned staging directories are never exposed as Work.

# Documentation Engine v0.1.0-rc.2

This release candidate freezes the first public Engine contract for consumer
validation.

## Compatibility claim

```text
Documentation Engine: 0.1.0-rc.2
Documentation CLI: 0.1.0-rc.2
Supported Specification: 1.0 Draft, revision 12
Supported Profile: Single Repository Profile 1.0
Target conformance: Level 2 preview
```

## Contract highlights

- `Synchronize(SyncInput)` requires an explicit `changed` or `no-change`
  Knowledge impact declaration.
- `KnowledgeImpact` is a typed public enum and is observable in `SyncResult`.
- Cleanup retries read the completed PRD and reject a conflicting outcome.
- `VersionInfo` and `dos --json version` expose the full compatibility matrix.
- `ErrorCodeOf` and CLI JSON errors expose stable machine-readable categories.
- The public Engine module path is fixed as
  `github.com/L1UUUU/Documentation-OS-Specification/engine`. The Engine is a
  nested Go module so Windows consumers do not need to download specification
  documents whose normative filenames contain Unicode punctuation.

## Known limitations

- `rc.1` could not be consumed with `go get` on Windows because the root module
  archive included normative filenames containing an em dash. `rc.2` publishes
  the Engine as a nested module and is the candidate Kanban validates.
- This RC awaits validation by the Kanban consumer before stable `v0.1.0`.
- The release claims Level 2 preview, not complete Level 3 conformance.
- Unclassified filesystem and implementation failures use the `internal`
  error code; finer recovery diagnostics are scheduled for lifecycle hardening.
- Historical migration and production write-lifecycle evidence are outside
  this release candidate.
