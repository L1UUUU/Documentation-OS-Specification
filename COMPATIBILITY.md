# Compatibility

This repository is the permanent Go module for the Documentation OS Engine:

```text
github.com/L1UUUU/Documentation-OS-Specification
```

The Engine remains in this specification repository until evidence from at
least two independent consumers justifies a repository split. Published
consumers should import:

```go
import "github.com/L1UUUU/Documentation-OS-Specification/engine"
```

## v0.1.0-rc.4 matrix

| Dimension | Supported value |
| --- | --- |
| Specification | 1.0 Draft, revision 13 |
| Repository Profile | Single Repository Profile 1.0 |
| Engine | 0.1.0-rc.4 |
| CLI | 0.1.0-rc.4 |
| Target conformance | Level 2 preview |
| Go | 1.22 or newer |
| CI platforms | Linux, Windows, macOS |

The `engine/v0.1.0-rc.4` tag is published. Its tag workflow passed on Linux,
Windows, and macOS, including external resolution without a `replace`
directive. rc.4 is intended for Kanban consumer validation; it does not claim
complete Level 3 conformance and is not the stable `v0.1.0` release.

Consumers should call `Engine.Version()` or `dos --json version` and compare
all independent dimensions. A matching specification number without a
matching status and revision is not a compatibility guarantee.

## Published v0.1.0-rc.4 contract

- `BeginWork(BeginInput)` atomically creates caller-defined core Work assets.
  An identical active retry returns the persisted result without writing;
  different input for the same slug returns `conflict`.
- Every synchronization requires an explicit `KnowledgeImpact` value.
- Completed Work retries must repeat the persisted terminal outcome.
- Use `errors.Is` with exported sentinel errors for control flow.
- Use `engine.ErrorCodeOf` or CLI JSON `code` values at integration boundaries.
- Do not parse human-readable error messages.
- `CreateIssue(CreateIssueInput)` atomically appends caller-authored Issue
  files to Active Work, allocates monotonically increasing two-digit numbers,
  is idempotent for identical retries, rejects conflicting slug reuse, and
  reconciles the derived INDEX.
- `LifecycleStage`, `LifecycleError`, and
  `FailureStageOf` provide stable lifecycle stage values (`begin`, `issue`,
  `synchronize`, `validate`, `complete`, or `cleanup`). `CreateIssue` reports
  all preflight, conflict, lock, write, and INDEX failures as the `issue` stage
  while preserving its existing `errors.Is` and `ErrorCodeOf` classification.
- A failed `ValidationReport` remains data from `Validate`; workflow gates can
  call `ValidationReport.Failure()` to enter the same error contract.
- The CLI exposes the same Issue contract through
  `dos issue create <work-slug> <issue-slug> --title TITLE --status STATUS --body-file PATH`.

Consumers pinned to rc.3 cannot use `CreateIssue`, lifecycle-stage APIs, or
`ValidationReport.Failure`; consumers can use those APIs by resolving rc.4.

Repository construction is not a lifecycle stage. `New` classifies a missing
or non-directory root as `invalid-repository` while retaining the underlying
filesystem cause. Whether an expected pre-worktree path should be opened at all
remains consumer policy, not Engine lifecycle policy.

Breaking corrections may be made between this RC and `v0.1.0`. After the
stable release, public compatibility is governed by semantic versioning while
the module remains at major version zero.
