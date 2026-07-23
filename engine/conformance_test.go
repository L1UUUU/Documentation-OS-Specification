// Package engine contains executable conformance scenarios from DOS-4005.
package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestConformanceCompleteRollsBackWhenMoveFails verifies DOS-4005 scenario 11.
func TestConformanceCompleteRollsBackWhenMoveFails(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("rollback-move")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-finish.md"), "---\nstatus: done\ntitle: Finish\n---\n")
	writeText(t, filepath.Join(work.Path, "notes.md"), "temporary context\n")

	originalRename := engine.renameFile
	engine.renameFile = func(oldPath, newPath string) error {
		if oldPath == work.Path {
			return errors.New("injected Work move failure")
		}
		return originalRename(oldPath, newPath)
	}

	if _, err := engine.Complete("rollback-move", OutcomeSucceeded); err == nil {
		t.Fatal("Complete() should report the injected directory movement failure")
	}
	activePath := filepath.Join(engine.Root, ".scratch", "active", "rollback-move")
	assertFileExists(t, filepath.Join(activePath, "PRD.md"))
	assertFileExists(t, filepath.Join(activePath, "notes.md"))
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "completed", "rollback-move"))
	if prd := readText(t, filepath.Join(activePath, "PRD.md")); strings.Contains(prd, "outcome:") {
		t.Fatalf("failed Complete() left outcome on the active Work:\n%s", prd)
	}
}

// TestConformanceGenerateWorkRollsBackMidCreationFailure verifies DOS-4005 scenario 4.
func TestConformanceGenerateWorkRollsBackMidCreationFailure(t *testing.T) {
	engine := newTestEngine(t)
	indexPath := filepath.Join(engine.Root, ".scratch", "INDEX.md")
	indexBefore := readText(t, indexPath)
	originalWrite := engine.writeFileAtomic
	engine.writeFileAtomic = func(path string, data []byte, mode os.FileMode) error {
		if filepath.Base(path) == "HANDOFF.md" {
			return errors.New("injected HANDOFF write failure")
		}
		return originalWrite(path, data, mode)
	}

	if _, err := engine.GenerateWork("rollback-generate"); err == nil || !strings.Contains(err.Error(), "rolled back") {
		t.Fatalf("GenerateWork() error = %v, want explicit rollback diagnostic", err)
	}
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "active", "rollback-generate"))
	if indexAfter := readText(t, indexPath); indexAfter != indexBefore {
		t.Fatal("failed GenerateWork() did not restore INDEX")
	}
}

// TestConformanceCompleteCleanupFailureIsRetriable verifies DOS-4005 scenario 5.
func TestConformanceCompleteCleanupFailureIsRetriable(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("retry-failed-cleanup")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-finish.md"), "---\nstatus: done\ntitle: Finish\n---\n")

	indexPath := filepath.Join(engine.Root, ".scratch", "INDEX.md")
	if err := os.Remove(indexPath); err != nil {
		t.Fatalf("remove INDEX fixture: %v", err)
	}
	if err := os.Mkdir(indexPath, 0o755); err != nil {
		t.Fatalf("block INDEX regeneration: %v", err)
	}

	result, err := engine.Complete("retry-failed-cleanup", OutcomeSucceeded)
	if err == nil {
		t.Fatal("Complete() should report Cleanup failure")
	}
	if stage, ok := FailureStageOf(err); !ok || stage != LifecycleStageCleanup {
		t.Fatalf("Complete() failure stage = %q, %v, want %q, true", stage, ok, LifecycleStageCleanup)
	}
	if !result.Completed || result.CleanupCompleted {
		t.Fatalf("Complete() result = %+v, want Completed with Cleanup pending", result)
	}
	completedPath := filepath.Join(engine.Root, ".scratch", "completed", "retry-failed-cleanup")
	assertFileExists(t, filepath.Join(completedPath, "PRD.md"))
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "active", "retry-failed-cleanup"))
	if prd := readText(t, filepath.Join(completedPath, "PRD.md")); !strings.Contains(prd, "outcome: succeeded") {
		t.Fatalf("Completed Work does not retain outcome:\n%s", prd)
	}

	if err := os.Remove(indexPath); err != nil {
		t.Fatalf("remove Cleanup blocker: %v", err)
	}
	retry, err := engine.Complete("retry-failed-cleanup", OutcomeSucceeded)
	if err != nil {
		t.Fatalf("Complete() Cleanup retry error = %v", err)
	}
	if !retry.Completed || !retry.CleanupCompleted || !retry.RetriedCleanup {
		t.Fatalf("Complete() retry result = %+v", retry)
	}
	assertFileExists(t, indexPath)
}

// TestConformanceConcurrentADRAllocationIsAtomic verifies DOS-4005 scenario 13.
func TestConformanceConcurrentADRAllocationIsAtomic(t *testing.T) {
	engine := newTestEngine(t)
	const draftCount = 8
	for index := 1; index <= draftCount; index++ {
		path := filepath.Join(engine.Root, "docs", "adr", fmt.Sprintf("ADR-DRAFT-decision-%d.md", index))
		writeText(t, path, fmt.Sprintf("---\nstatus: created\ntitle: Decision %d\nrelationships: []\n---\n", index))
	}

	type allocation struct {
		result AllocationResult
		err    error
	}
	start := make(chan struct{})
	allocations := make(chan allocation, draftCount)
	var workers sync.WaitGroup
	for index := 1; index <= draftCount; index++ {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()
			<-start
			result, err := engine.AllocateKnowledgeIdentifier("adr", fmt.Sprintf("docs/adr/ADR-DRAFT-decision-%d.md", index))
			allocations <- allocation{result: result, err: err}
		}(index)
	}
	close(start)
	workers.Wait()
	close(allocations)

	identifiers := make(map[string]bool, draftCount)
	for allocation := range allocations {
		if allocation.err != nil {
			t.Errorf("concurrent AllocateKnowledgeIdentifier() error = %v", allocation.err)
			continue
		}
		if identifiers[allocation.result.Identifier] {
			t.Errorf("identifier %s was allocated more than once", allocation.result.Identifier)
		}
		identifiers[allocation.result.Identifier] = true
	}
	if len(identifiers) != draftCount {
		t.Fatalf("allocated %d distinct identifiers, want %d", len(identifiers), draftCount)
	}
	report, err := engine.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !report.Passed() {
		t.Fatalf("concurrent allocation left an invalid repository:\n%s", report)
	}
}

// TestConformanceSucceededWorkRejectsNonTerminalIssue verifies DOS-4005 scenario 9.
func TestConformanceSucceededWorkRejectsNonTerminalIssue(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("unfinished-success")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-unfinished.md"), "---\nstatus: in-progress\ntitle: Unfinished\n---\n")
	indexBefore := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))

	if _, err := engine.Complete("unfinished-success", OutcomeSucceeded); !errors.Is(err, ErrPreflight) {
		t.Fatalf("Complete() error = %v, want ErrPreflight", err)
	}
	activePath := filepath.Join(engine.Root, ".scratch", "active", "unfinished-success")
	assertFileExists(t, filepath.Join(activePath, "PRD.md"))
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "completed", "unfinished-success"))
	if prd := readText(t, filepath.Join(activePath, "PRD.md")); strings.Contains(prd, "outcome:") {
		t.Fatalf("rejected Complete() wrote outcome:\n%s", prd)
	}
	if indexAfter := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); indexAfter != indexBefore {
		t.Fatal("rejected Complete() changed INDEX")
	}
}

// TestConformanceFailedCompleteCanRetryIdempotently verifies DOS-4005 scenario 6.
func TestConformanceFailedCompleteCanRetryIdempotently(t *testing.T) {
	retried := newTestEngine(t)
	retriedWork, err := retried.GenerateWork("retry-pipeline")
	if err != nil {
		t.Fatalf("GenerateWork() retried fixture error = %v", err)
	}
	retriedIssue := filepath.Join(retriedWork.Path, "issues", "01-finish.md")
	writeText(t, retriedIssue, "---\nstatus: in-progress\ntitle: Finish\n---\n")
	if _, err := retried.Complete("retry-pipeline", OutcomeSucceeded); !errors.Is(err, ErrPreflight) {
		t.Fatalf("first Complete() error = %v, want ErrPreflight", err)
	}
	writeText(t, retriedIssue, "---\nstatus: done\ntitle: Finish\n---\n")
	if _, err := retried.Complete("retry-pipeline", OutcomeSucceeded); err != nil {
		t.Fatalf("retried Complete() error = %v", err)
	}

	control := newTestEngine(t)
	controlWork, err := control.GenerateWork("retry-pipeline")
	if err != nil {
		t.Fatalf("GenerateWork() control fixture error = %v", err)
	}
	writeText(t, filepath.Join(controlWork.Path, "issues", "01-finish.md"), "---\nstatus: done\ntitle: Finish\n---\n")
	if _, err := control.Complete("retry-pipeline", OutcomeSucceeded); err != nil {
		t.Fatalf("control Complete() error = %v", err)
	}

	paths := []string{
		".scratch/INDEX.md",
		".scratch/completed/retry-pipeline/PRD.md",
		".scratch/completed/retry-pipeline/HANDOFF.md",
		".scratch/completed/retry-pipeline/issues/01-finish.md",
	}
	for _, relative := range paths {
		retriedContent := readText(t, retried.path(relative))
		controlContent := readText(t, control.path(relative))
		if retriedContent != controlContent {
			t.Errorf("retried %s differs from uninterrupted execution", relative)
		}
	}
}

// TestConformanceCancelledWorkCompletesStandardPipeline verifies DOS-4005 scenario 10.
func TestConformanceCancelledWorkCompletesStandardPipeline(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("cancelled-pipeline")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-withdrawn.md"), "---\nstatus: in-progress\ntitle: Withdrawn work\n---\n")
	if _, err := engine.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
		t.Fatalf("Synchronize() error = %v", err)
	}
	if report, err := engine.Validate(); err != nil || !report.Passed() {
		t.Fatalf("pre-Complete Validate() error = %v, report = %s", err, report)
	}

	result, err := engine.Complete("cancelled-pipeline", OutcomeCancelled)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if !result.Completed || !result.CleanupCompleted || result.Outcome != OutcomeCancelled {
		t.Fatalf("Complete() result = %+v", result)
	}
	completedPath := filepath.Join(engine.Root, ".scratch", "completed", "cancelled-pipeline")
	assertFileExists(t, filepath.Join(completedPath, "PRD.md"))
	assertFileExists(t, filepath.Join(completedPath, "issues", "01-withdrawn.md"))
	assertFileExists(t, filepath.Join(completedPath, "HANDOFF.md"))
	if prd := readText(t, filepath.Join(completedPath, "PRD.md")); !strings.Contains(prd, "outcome: cancelled") {
		t.Fatalf("Completed PRD does not record cancelled outcome:\n%s", prd)
	}
	if index := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); !strings.Contains(index, "Outcome: cancelled") {
		t.Fatalf("INDEX does not surface cancelled outcome:\n%s", index)
	}
	if report, err := engine.Validate(); err != nil || !report.Passed() {
		t.Fatalf("post-Complete Validate() error = %v, report = %s", err, report)
	}
}

// TestConformanceNonSuccessOutcomesCompleteConsistently verifies every
// exceptional terminal outcome permits unfinished Issues and remains
// idempotent when Cleanup is retried with the persisted outcome.
func TestConformanceNonSuccessOutcomesCompleteConsistently(t *testing.T) {
	for _, outcome := range []string{OutcomeCancelled, OutcomeSuperseded, OutcomeFailed} {
		t.Run(outcome, func(t *testing.T) {
			instance := newTestEngine(t)
			slug := outcome + "-pipeline"
			work, err := instance.GenerateWork(slug)
			if err != nil {
				t.Fatalf("GenerateWork() error = %v", err)
			}
			writeText(t, filepath.Join(work.Path, "issues", "01-unfinished.md"), "---\nstatus: in-progress\ntitle: Unfinished\n---\n")

			result, err := instance.Complete(slug, outcome)
			if err != nil {
				t.Fatalf("Complete() error = %v", err)
			}
			if result.Outcome != outcome || !result.Completed || !result.CleanupCompleted {
				t.Fatalf("Complete() result = %+v, want terminal outcome %q", result, outcome)
			}
			retry, err := instance.Complete(slug, outcome)
			if err != nil {
				t.Fatalf("Complete() retry error = %v", err)
			}
			if retry.Outcome != outcome || !retry.RetriedCleanup || !retry.CleanupCompleted {
				t.Fatalf("Complete() retry = %+v, want idempotent Cleanup for %q", retry, outcome)
			}
		})
	}
}

// TestConformanceIndexSurfacesCompletedOutcome verifies DOS-4005 scenario 12.
func TestConformanceIndexSurfacesCompletedOutcome(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("superseded-index")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-replaced.md"), "---\nstatus: blocked\ntitle: Replaced approach\n---\n")
	if _, err := engine.Complete("superseded-index", OutcomeSuperseded); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if _, err := engine.GenerateIndex(); err != nil {
		t.Fatalf("GenerateIndex() error = %v", err)
	}

	index := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	expectedFragments := []string{
		"### superseded-index",
		".scratch/completed/superseded-index/PRD.md",
		".scratch/completed/superseded-index/HANDOFF.md",
		"Outcome: superseded",
		"01-replaced.md [blocked] Replaced approach",
	}
	for _, fragment := range expectedFragments {
		if !strings.Contains(index, fragment) {
			t.Errorf("INDEX does not contain %q:\n%s", fragment, index)
		}
	}
}

// TestConformanceADRAllocationRollsBackWhenValidationGateFails verifies DOS-4005 scenario 14.
func TestConformanceADRAllocationRollsBackWhenValidationGateFails(t *testing.T) {
	engine := newTestEngine(t)
	draftPath := filepath.Join(engine.Root, "docs", "adr", "ADR-DRAFT-gated.md")
	draftContent := "---\nstatus: created\ntitle: Gated decision\nrelationships: []\n---\n"
	writeText(t, draftPath, draftContent)
	if report, err := engine.Validate(); err != nil || !report.Passed() {
		t.Fatalf("draft identity should be exempt before integration, error = %v, report = %s", err, report)
	}
	referencePath := filepath.Join(engine.Root, "docs", "architecture", "0001-system.md")
	referenceContent := "---\nid: ARCH-0001\nstatus: active\nrelationships:\n  - type: references\n    target: ADR-DRAFT-gated\n---\n"
	writeText(t, referencePath, referenceContent)
	writeText(t, filepath.Join(engine.Root, "docs", "standards", "0001-broken.md"), "---\nid: STD-0001\nstatus: active\nrelationships:\n  - type: references\n    target: ADR-9999\n---\n")

	if _, err := engine.AllocateKnowledgeIdentifier("adr", "docs/adr/ADR-DRAFT-gated.md"); err == nil {
		t.Fatal("AllocateKnowledgeIdentifier() should fail its Relationship Validation gate")
	}
	assertFileExists(t, draftPath)
	assertFileNotExists(t, filepath.Join(engine.Root, "docs", "adr", "0001-gated.md"))
	if draft := readText(t, draftPath); draft != draftContent {
		t.Fatalf("allocation rollback changed the draft:\n%s", draft)
	}
	if reference := readText(t, referencePath); reference != referenceContent {
		t.Fatalf("allocation rollback did not restore managed references:\n%s", reference)
	}
}

// TestConformanceMigrationPreservesKnowledgeSemantics verifies the migration testing category.
func TestConformanceMigrationPreservesKnowledgeSemantics(t *testing.T) {
	engine := newTestEngine(t)
	adrPath := filepath.Join(engine.Root, "docs", "adr", "0001-runtime.md")
	adrContent := "---\nid: ADR-0001\nstatus: active\nrelationships: []\n---\n\nRuntime decision.\n"
	writeText(t, adrPath, adrContent)
	architecturePath := filepath.Join(engine.Root, "docs", "architecture", "0001-system.md")
	architectureContent := "---\nid: ARCH-0001\nstatus: active\nrelationships:\n  - type: references\n    target: ADR-0001\n---\n\nSystem.\n"
	writeText(t, architecturePath, architectureContent)
	if _, err := engine.GenerateWork("migration-work"); err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"), "stale index\n")

	report, err := engine.Migrate(SpecificationVersion)
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if !report.Success || report.After != SpecificationVersion {
		t.Fatalf("Migrate() report = %+v", report)
	}
	if content := readText(t, adrPath); content != adrContent {
		t.Fatalf("migration changed ADR semantics:\n%s", content)
	}
	if content := readText(t, architecturePath); content != architectureContent {
		t.Fatalf("migration changed relationship semantics:\n%s", content)
	}
	if index := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); !strings.Contains(index, "### migration-work") {
		t.Fatalf("migration did not regenerate INDEX:\n%s", index)
	}
	validation, err := engine.Validate()
	if err != nil {
		t.Fatalf("post-migration Validate() error = %v", err)
	}
	if !validation.Passed() {
		t.Fatalf("post-migration repository is invalid:\n%s", validation)
	}
}

// TestConformanceSingleRepositoryProfileInitialization verifies the repository profile category.
func TestConformanceSingleRepositoryProfileInitialization(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.GenerateWork("profile-fixture"); err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	directories := []string{
		"docs/architecture",
		"docs/adr",
		"docs/standards",
		"docs/inbox",
		".scratch/active",
		".scratch/completed",
		".scratch/active/profile-fixture/issues",
	}
	for _, relative := range directories {
		info, err := os.Stat(engine.path(relative))
		if err != nil || !info.IsDir() {
			t.Errorf("required directory %s is missing or invalid: %v", relative, err)
		}
	}
	files := []string{
		"AGENTS.md",
		"CLAUDE.md",
		".scratch/AGENTS.md",
		".scratch/CLAUDE.md",
		".scratch/INDEX.md",
		".scratch/active/profile-fixture/PRD.md",
		".scratch/active/profile-fixture/HANDOFF.md",
	}
	for _, relative := range files {
		info, err := os.Stat(engine.path(relative))
		if err != nil || !info.Mode().IsRegular() {
			t.Errorf("required file %s is missing or invalid: %v", relative, err)
		}
	}
	if report, err := engine.Validate(); err != nil || !report.Passed() {
		t.Fatalf("initialized profile Validate() error = %v, report = %s", err, report)
	}
}

// Ensure the injected seam retains the production signature.
var _ func(string, string) error = os.Rename
