# Documentation Engine v0.1.0-rc.4 release candidate

This candidate implements Documentation OS 1.0 Draft revision 13. It is not a
published release until the `engine/v0.1.0-rc.4` tag is pushed and its Linux,
Windows, and macOS tag workflow succeeds.

- Added `LifecycleStage`, `LifecycleError`, and `FailureStageOf` so consumers
  can distinguish Begin, Issue, Synchronize, Validate, Complete, and post-commit
  Cleanup failures without parsing text or maintaining adapter-local names.
- `CreateIssue` now attaches the stable `issue` stage to preflight, conflict,
  lock, write, and INDEX failures while preserving existing sentinel
  and error-code classification.
- Added `CreateIssue(CreateIssueInput)` for atomic, monotonically numbered Issue
  creation in Active Work, with identical-retry idempotency, conflicting-slug
  rejection, rollback, and derived INDEX reconciliation.
- Added `dos issue create` with JSON and human-readable output over the same
  Engine contract.
- Added `ValidationReport.Failure()` for workflow gates while preserving the
  existing `Validate` report-as-data API.
- Missing and non-directory repository roots now use the stable
  `invalid-repository` code and preserve their filesystem cause. Expected
  unmaterialized worktrees remain a consumer policy decision.
- Cleanup failure remains retriable and is now explicitly reported as the
  `cleanup` stage after terminal Work state has been persisted.
- `Complete` now recovers when a process stops after the active PRD durably
  records its terminal outcome but before the Work directory moves to
  `completed`. A retry with the same outcome resumes cleanup and movement; a
  retry with a different outcome returns `conflict`.
- Added a deterministic post-outcome interruption checkpoint to the Engine's
  package-level conformance harness. This simulates the non-returning crash
  boundary without exposing a new public API or weakening normal rollback.
- Added an external-consumer smoke command for published Engine versions. It
  creates a fresh module outside this repository, forces `GOWORK=off`, rejects any `replace`
  resolution, and compiles the selected published Engine version. Branch CI
  continues to default to the latest published RC (`v0.1.0-rc.3`) before rc.4
  is published; manual and reusable workflow runs accept an explicit version;
  an `engine/v*` tag run derives the module version from that tag. A local
  Windows run against rc.3 passed.

The smoke test validates what an external consumer can download, not
unpublished branch contents. A candidate version cannot pass this remote,
no-`replace` check before its `engine/v*` tag exists upstream. Therefore tag
validation runs after the tag is pushed; a successful tag workflow is required
before creating the corresponding GitHub release or promoting the version to
the branch default.

## v0.1.0 release readiness audit (2026-07-18)

**Decision: NO-GO.** The verifiable real lifecycle cohort remains **0**. Do not
create or push a stable tag or GitHub release yet.

| Gate | Evidence | Decision |
| --- | --- | --- |
| KB-04 real Demand cohort | Kanban commit `1b49150` adds lifecycle evidence logging, but the available `.kanban/runtime` contains one skill-materialization run and no lifecycle attempt log. No repository `daemon.log` is available. Verifiable real lifecycle cohorts: **0**; verifiable real `succeeded`, `cancelled`, `superseded`, `failed`, or crash-recovery closures: **0**. Unit/integration fixtures are not counted as real Demands. | **Blocking** |
| Complete hard-crash window | A conformance test interrupts immediately after outcome persistence. The retry now resumes the active transaction for the same outcome and rejects a conflicting outcome. Full Engine tests and vet pass on Windows. | Cleared locally |
| Published external module | A fresh external module resolves `engine@v0.1.0-rc.3` with `GOWORK=off`, no `replace`, and passes `go test` on Windows. The three-platform CI matrix defaults to this latest published RC on branches and can select a published version for manual, reusable, or tag validation. This proves the published package only; it does not prove Unreleased branch APIs. | Implementation cleared; require green tag-version CI runs before release |

Stable release requires a recorded KB-04 cohort with the planned sample size,
all four terminal outcomes, and real retry/reconcile evidence. It also requires
green Linux, Windows, and macOS runs of the external-consumer smoke. Synthetic
conformance tests establish deterministic Engine behavior but do not substitute
for production dogfood evidence.

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
