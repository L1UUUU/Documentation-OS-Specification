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
