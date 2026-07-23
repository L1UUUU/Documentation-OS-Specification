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

## v0.1.0-rc.8 candidate matrix

| Dimension | Supported value |
| --- | --- |
| Specification | 1.0 Draft, revision 13 |
| Repository Profile | Single Repository Profile 1.0 |
| Engine | 0.1.0-rc.8 |
| CLI | 0.1.0-rc.8 |
| Target conformance | Level 2 preview |
| Go | 1.22 or newer |
| CI platforms | Linux, Windows, macOS |

This matrix describes the locally prepared rc.8 candidate. It becomes a
published compatibility claim only after the `engine/v0.1.0-rc.8` tag is
pushed and its three-platform tag workflow succeeds. Until then rc.7 remains
the latest GitHub prerelease. rc.8 is intended for Kanban consumer
validation; it does not claim complete Level 3 conformance and is not the
stable `v0.1.0` release.

The `engine/v0.1.0-rc.6` tag exists, but no GitHub Release was created for it.
Its tag workflow did not exercise the `documentation_conformance` build on all
supported platforms, so rc.6 was superseded by rc.7 rather than promoted.

Consumers should call `Engine.Version()` or `dos --json version` and compare
all independent dimensions. A matching specification number without a
matching status and revision is not a compatibility guarantee.

## v0.1.0-rc.8 candidate contract

- `BeginWork(BeginInput)` atomically creates caller-defined core Work assets.
  An identical active retry returns the persisted result without writing;
  different input for the same slug returns `conflict`.
- Every synchronization requires an explicit `KnowledgeImpact` value.
- Completed Work retries must repeat the persisted terminal outcome.
- An Active Work with a legal persisted outcome is a durable Complete recovery
  marker. Validation accepts it so a restarted production pipeline can repeat
  Synchronize and Validate before retrying Complete; a conflicting outcome is
  still rejected without changing the marker.
- Use `errors.Is` with exported sentinel errors for control flow.
- Use `engine.ErrorCodeOf` or CLI JSON `code` values at integration boundaries.
- Do not parse human-readable error messages.
- `CreateIssue(CreateIssueInput)` atomically appends caller-authored Issue
  files to Active Work, allocates monotonically increasing two-digit numbers,
  is idempotent for identical retries, rejects conflicting slug reuse, and
  reconciles the derived INDEX. Cross-process lock contention returns within
  five seconds rather than waiting indefinitely.
- `CreateIssueContext(context.Context, CreateIssueInput)` exposes the same
  atomic Issue contract with caller-controlled cancellation and deadlines
  while waiting for the repository lock.
- `LifecycleStage`, `LifecycleError`, and
  `FailureStageOf` provide stable lifecycle stage values (`begin`, `issue`,
  `synchronize`, `validate`, `complete`, or `cleanup`). `CreateIssue` reports
  all preflight, conflict, lock, write, and INDEX failures as the `issue` stage
  while preserving its existing `errors.Is` and `ErrorCodeOf` classification.
- A failed `ValidationReport` remains data from `Validate`; workflow gates can
  call `ValidationReport.Failure()` to enter the same error contract.
- `ValidateContext(context.Context)` adds cancellation to validation while
  preserving `Validate()` as the backward-compatible wrapper.
- Validate and Complete fault plans exist only in builds compiled with the
  `documentation_conformance` tag. Complete plans can stop immediately after
  durable outcome persistence or immediately after the Active-to-Completed
  move and before INDEX Cleanup. Plans are target-, attempt-, and count-bounded,
  cancellation-aware before consumption, concurrency-safe, and emit only
  path-free stable audit fields. Normal builds expose no injector constructor
  and accept no fault configuration; consumers must not treat the tagged seam
  as a runtime or production API.
- The CLI exposes the same Issue contract through
  `dos issue create <work-slug> <issue-slug> --title TITLE --status STATUS --body-file PATH`.

Consumers pinned to published rc.7 retain `ValidateContext` and its tagged
validate-only fault adapter, but must wait for published rc.8 to consume the
two Complete recovery points. Tag workflows resolve and smoke-test their exact
tag instead of unpublished source.

Repository construction is not a lifecycle stage. `New` classifies a missing
or non-directory root as `invalid-repository` while retaining the underlying
filesystem cause. Whether an expected pre-worktree path should be opened at all
remains consumer policy, not Engine lifecycle policy.

Breaking corrections may be made between this RC and `v0.1.0`. After the
stable release, public compatibility is governed by semantic versioning while
the module remains at major version zero.
