// Package engine tests the public Documentation OS engine operations.
package engine

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestBeginWorkCreatesCallerDefinedCoreAssets verifies the public Begin Work contract.
func TestBeginWorkCreatesCallerDefinedCoreAssets(t *testing.T) {
	instance := newTestEngine(t)

	result, err := instance.BeginWork(BeginInput{
		Slug:  "write-lifecycle",
		Title: "Write lifecycle",
		Issues: []BeginIssue{
			{Slug: "implement-adapter", Title: "Implement adapter", Status: "open"},
			{Slug: "verify-recovery", Title: "Verify recovery", Status: "in-progress"},
		},
	})
	if err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	if result.Slug != "write-lifecycle" || result.Path == "" || result.IndexPath != ".scratch/INDEX.md" {
		t.Fatalf("BeginWork() result = %+v", result)
	}

	prd := readText(t, filepath.Join(result.Path, "PRD.md"))
	if !strings.Contains(prd, "# Write lifecycle") {
		t.Fatalf("PRD does not contain caller title:\n%s", prd)
	}
	firstIssue := readText(t, filepath.Join(result.Path, "issues", "01-implement-adapter.md"))
	if !strings.Contains(firstIssue, "status: open") || !strings.Contains(firstIssue, `title: "Implement adapter"`) {
		t.Fatalf("first Issue does not contain caller fields:\n%s", firstIssue)
	}
	secondIssue := readText(t, filepath.Join(result.Path, "issues", "02-verify-recovery.md"))
	if !strings.Contains(secondIssue, "status: in-progress") || !strings.Contains(secondIssue, `title: "Verify recovery"`) {
		t.Fatalf("second Issue does not contain caller fields:\n%s", secondIssue)
	}
	assertFileExists(t, filepath.Join(result.Path, "HANDOFF.md"))
	if index := readText(t, filepath.Join(instance.Root, ".scratch", "INDEX.md")); !strings.Contains(index, "### write-lifecycle") || !strings.Contains(index, "01-implement-adapter.md [open] Implement adapter") {
		t.Fatalf("INDEX does not contain begun Work:\n%s", index)
	}
	if report, validateErr := instance.Validate(); validateErr != nil || !report.Passed() {
		t.Fatalf("begun Work should validate, error = %v, report = %s", validateErr, report)
	}
}

// TestBeginWorkSameInputRetryIsIdempotent verifies retries return persisted Work without writes.
func TestBeginWorkSameInputRetryIsIdempotent(t *testing.T) {
	instance := newTestEngine(t)
	input := BeginInput{
		Slug:  "retry-begin",
		Title: "Retry begin",
		Issues: []BeginIssue{
			{Slug: "implement", Title: "Implement", Status: "open"},
		},
	}
	first, err := instance.BeginWork(input)
	if err != nil {
		t.Fatalf("first BeginWork() error = %v", err)
	}
	paths := []string{
		filepath.Join(first.Path, "PRD.md"),
		filepath.Join(first.Path, "HANDOFF.md"),
		filepath.Join(first.Path, "issues", "01-implement.md"),
		filepath.Join(instance.Root, ".scratch", "INDEX.md"),
	}
	type fileState struct {
		content string
		modTime time.Time
	}
	before := map[string]fileState{}
	for _, path := range paths {
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Fatalf("stat %s: %v", path, statErr)
		}
		before[path] = fileState{content: readText(t, path), modTime: info.ModTime()}
	}
	time.Sleep(20 * time.Millisecond)

	retry, err := instance.BeginWork(input)
	if err != nil {
		t.Fatalf("retry BeginWork() error = %v", err)
	}
	if retry != first {
		t.Fatalf("retry BeginWork() result = %+v, want persisted %+v", retry, first)
	}
	for _, path := range paths {
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Fatalf("stat retried %s: %v", path, statErr)
		}
		if got := readText(t, path); got != before[path].content {
			t.Errorf("retry changed %s content", path)
		}
		if !info.ModTime().Equal(before[path].modTime) {
			t.Errorf("retry rewrote %s", path)
		}
	}
}

// TestBeginWorkDifferentInputRetryConflicts verifies persisted Begin input is authoritative.
func TestBeginWorkDifferentInputRetryConflicts(t *testing.T) {
	instance := newTestEngine(t)
	input := BeginInput{Slug: "conflicting-begin", Title: "Original title"}
	if _, err := instance.BeginWork(input); err != nil {
		t.Fatalf("first BeginWork() error = %v", err)
	}
	before := readText(t, filepath.Join(instance.Root, ".scratch", "active", input.Slug, "PRD.md"))

	_, err := instance.BeginWork(BeginInput{Slug: input.Slug, Title: "Different title"})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("conflicting BeginWork() error = %v, want ErrConflict", err)
	}
	if code := ErrorCodeOf(err); code != ErrorCodeConflict {
		t.Fatalf("conflicting BeginWork() code = %q, want %q", code, ErrorCodeConflict)
	}
	if after := readText(t, filepath.Join(instance.Root, ".scratch", "active", input.Slug, "PRD.md")); after != before {
		t.Fatal("conflicting BeginWork() changed persisted Work")
	}
}

// TestBeginWorkRollsBackInitialIssueFailure verifies the creation transaction spans initial Issues.
func TestBeginWorkRollsBackInitialIssueFailure(t *testing.T) {
	instance := newTestEngine(t)
	indexPath := filepath.Join(instance.Root, ".scratch", "INDEX.md")
	indexBefore := readText(t, indexPath)
	originalWrite := instance.writeFileAtomic
	instance.writeFileAtomic = func(path string, data []byte, mode os.FileMode) error {
		if filepath.Base(path) == "02-fail.md" {
			return errors.New("injected initial Issue failure")
		}
		return originalWrite(path, data, mode)
	}

	_, err := instance.BeginWork(BeginInput{
		Slug:  "rollback-begin",
		Title: "Rollback begin",
		Issues: []BeginIssue{
			{Slug: "created", Title: "Created before failure", Status: "open"},
			{Slug: "fail", Title: "Injected failure", Status: "open"},
		},
	})
	if err == nil {
		t.Fatal("BeginWork() should report the injected initial Issue failure")
	}
	assertFileNotExists(t, filepath.Join(instance.Root, ".scratch", "active", "rollback-begin"))
	if indexAfter := readText(t, indexPath); indexAfter != indexBefore {
		t.Fatal("failed BeginWork() did not restore INDEX")
	}
}

// TestBeginWorkRollsBackIndexFailure verifies the final generated artifact is in the transaction.
func TestBeginWorkRollsBackIndexFailure(t *testing.T) {
	instance := newTestEngine(t)
	indexPath := filepath.Join(instance.Root, ".scratch", "INDEX.md")
	indexBefore := readText(t, indexPath)
	originalWrite := instance.writeFileAtomic
	instance.writeFileAtomic = func(path string, data []byte, mode os.FileMode) error {
		if path == indexPath {
			return errors.New("injected INDEX failure")
		}
		return originalWrite(path, data, mode)
	}

	_, err := instance.BeginWork(BeginInput{Slug: "rollback-index", Title: "Rollback index"})
	if err == nil {
		t.Fatal("BeginWork() should report the injected INDEX failure")
	}
	assertFileNotExists(t, filepath.Join(instance.Root, ".scratch", "active", "rollback-index"))
	if indexAfter := readText(t, indexPath); indexAfter != indexBefore {
		t.Fatal("failed BeginWork() did not restore INDEX after generated artifact failure")
	}
}

// TestBeginWorkConcurrentSameInputIsIdempotent verifies concurrent callers publish one complete Work.
func TestBeginWorkConcurrentSameInputIsIdempotent(t *testing.T) {
	instance := newTestEngine(t)
	input := BeginInput{
		Slug:  "concurrent-same-begin",
		Title: "Concurrent same begin",
		Issues: []BeginIssue{
			{Slug: "implement", Title: "Implement", Status: "open"},
		},
	}
	const callers = 8
	start := make(chan struct{})
	results := make(chan WorkResult, callers)
	errorsFound := make(chan error, callers)
	var workers sync.WaitGroup
	for index := 0; index < callers; index++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			result, err := instance.BeginWork(input)
			results <- result
			errorsFound <- err
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	close(errorsFound)

	for err := range errorsFound {
		if err != nil {
			t.Errorf("concurrent same-input BeginWork() error = %v", err)
		}
	}
	var expected WorkResult
	for result := range results {
		if expected == (WorkResult{}) {
			expected = result
			continue
		}
		if result != expected {
			t.Errorf("concurrent BeginWork() result = %+v, want %+v", result, expected)
		}
	}
	assertFileExists(t, filepath.Join(instance.Root, ".scratch", "active", input.Slug, "PRD.md"))
	assertFileExists(t, filepath.Join(instance.Root, ".scratch", "active", input.Slug, "HANDOFF.md"))
	assertFileExists(t, filepath.Join(instance.Root, ".scratch", "active", input.Slug, "issues", "01-implement.md"))
	if report, err := instance.Validate(); err != nil || !report.Passed() {
		t.Fatalf("concurrent BeginWork left invalid repository, error = %v, report = %s", err, report)
	}
}

// TestBeginWorkConcurrentDifferentInputConflictsWithoutOverwriting verifies the winner remains authoritative.
func TestBeginWorkConcurrentDifferentInputConflictsWithoutOverwriting(t *testing.T) {
	instance := newTestEngine(t)
	inputs := []BeginInput{
		{Slug: "concurrent-different-begin", Title: "First title", Issues: []BeginIssue{{Slug: "first", Title: "First", Status: "open"}}},
		{Slug: "concurrent-different-begin", Title: "Second title", Issues: []BeginIssue{{Slug: "second", Title: "Second", Status: "open"}}},
	}
	type attempt struct {
		input  BeginInput
		result WorkResult
		err    error
	}
	start := make(chan struct{})
	attempts := make(chan attempt, len(inputs))
	var workers sync.WaitGroup
	for _, input := range inputs {
		input := input
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			result, err := instance.BeginWork(input)
			attempts <- attempt{input: input, result: result, err: err}
		}()
	}
	close(start)
	workers.Wait()
	close(attempts)

	var winner *attempt
	conflicts := 0
	for current := range attempts {
		if current.err == nil {
			copy := current
			winner = &copy
			continue
		}
		if !errors.Is(current.err, ErrConflict) || ErrorCodeOf(current.err) != ErrorCodeConflict {
			t.Errorf("concurrent different-input BeginWork() error = %v, want stable conflict", current.err)
			continue
		}
		conflicts++
	}
	if winner == nil || conflicts != 1 {
		t.Fatalf("concurrent attempts winner = %+v, conflicts = %d", winner, conflicts)
	}
	prd := readText(t, filepath.Join(winner.result.Path, "PRD.md"))
	if !strings.Contains(prd, "# "+winner.input.Title) {
		t.Fatalf("persisted Work does not match winning input:\n%s", prd)
	}
	assertFileExists(t, filepath.Join(winner.result.Path, "issues", "01-"+winner.input.Issues[0].Slug+".md"))
	if report, err := instance.Validate(); err != nil || !report.Passed() {
		t.Fatalf("different-input race left invalid repository, error = %v, report = %s", err, report)
	}
}

// TestBeginWorkIgnoresCrashedStaging verifies a partial staging tree is never exposed as active Work.
func TestBeginWorkIgnoresCrashedStaging(t *testing.T) {
	instance := newTestEngine(t)
	orphan := filepath.Join(instance.Root, ".scratch", ".begin-staging", "crashed-call")
	writeText(t, filepath.Join(orphan, "PRD.md"), "partial staging content\n")
	if report, err := instance.Inspect(); err != nil || report.ActiveWorks != 0 {
		t.Fatalf("Inspect() exposed crashed staging, error = %v, report = %+v", err, report)
	}
	assertFileNotExists(t, filepath.Join(instance.Root, ".scratch", "active", "crashed-call"))

	result, err := instance.BeginWork(BeginInput{Slug: "after-crash", Title: "After crash"})
	if err != nil {
		t.Fatalf("BeginWork() after crashed staging error = %v", err)
	}
	assertFileExists(t, filepath.Join(result.Path, "PRD.md"))
	assertFileExists(t, filepath.Join(result.Path, "HANDOFF.md"))
	assertFileExists(t, orphan)
}

// TestBeginWorkRetryRepairsPostPublishIndexCrash verifies recovery after rename but before INDEX update.
func TestBeginWorkRetryRepairsPostPublishIndexCrash(t *testing.T) {
	instance := newTestEngine(t)
	input := BeginInput{Slug: "post-publish-crash", Title: "Post publish crash"}
	if _, err := instance.BeginWork(input); err != nil {
		t.Fatalf("first BeginWork() error = %v", err)
	}
	indexPath := filepath.Join(instance.Root, ".scratch", "INDEX.md")
	writeText(t, indexPath, "stale index from simulated crash\n")

	result, err := instance.BeginWork(input)
	if err != nil {
		t.Fatalf("recovery BeginWork() error = %v", err)
	}
	if result.Slug != input.Slug {
		t.Fatalf("recovery BeginWork() result = %+v", result)
	}
	if index := readText(t, indexPath); !strings.Contains(index, "### post-publish-crash") {
		t.Fatalf("recovery BeginWork() did not repair INDEX:\n%s", index)
	}
}

// TestGenerateWorkCreatesCoreAssetsAndDeterministicIndex verifies Work creation and INDEX generation.
func TestGenerateWorkCreatesCoreAssetsAndDeterministicIndex(t *testing.T) {
	engine := newTestEngine(t)

	if _, err := engine.GenerateWork("add-engine"); err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}

	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "active", "add-engine", "PRD.md"))
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "active", "add-engine", "issues"))
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "active", "add-engine", "HANDOFF.md"))

	first, err := os.ReadFile(filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	if err != nil {
		t.Fatalf("read generated INDEX: %v", err)
	}
	if !strings.Contains(string(first), "## Active Works") || !strings.Contains(string(first), "add-engine") {
		t.Fatalf("generated INDEX does not list the active Work:\n%s", first)
	}

	if _, err := engine.GenerateIndex(); err != nil {
		t.Fatalf("GenerateIndex() error = %v", err)
	}
	second, err := os.ReadFile(filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	if err != nil {
		t.Fatalf("read regenerated INDEX: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("INDEX is not reproducible")
	}
}

// TestInspectReturnsIssueSummaries verifies callers can inspect Issue identity, title, and DOS status.
func TestInspectReturnsIssueSummaries(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("inspect-issues")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-query-engine.md"), "---\nstatus: in-progress\ntitle: Query the engine\n---\n")

	report, err := engine.Inspect()
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if len(report.Works) != 1 {
		t.Fatalf("Inspect() returned %d Works, want 1", len(report.Works))
	}
	want := IssueSummary{
		Name:   "01-query-engine.md",
		Path:   ".scratch/active/inspect-issues/issues/01-query-engine.md",
		Title:  "Query the engine",
		Status: "in-progress",
	}
	if len(report.Works[0].Issues) != 1 || report.Works[0].Issues[0] != want {
		t.Fatalf("Inspect() Issues = %+v, want [%+v]", report.Works[0].Issues, want)
	}
}

// TestGenerateWorkDuplicateLeavesRepositoryUnchanged verifies Generate rollback on preflight failure.
func TestGenerateWorkDuplicateLeavesRepositoryUnchanged(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.GenerateWork("duplicate-work"); err != nil {
		t.Fatalf("first GenerateWork() error = %v", err)
	}
	before := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	if _, err := engine.GenerateWork("duplicate-work"); err == nil {
		t.Fatal("duplicate GenerateWork() should fail")
	}
	if after := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); after != before {
		t.Fatal("duplicate GenerateWork() changed INDEX")
	}
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "active", "duplicate-work", "PRD.md"))
}

// TestInitializeRepairsMissingCanonicalMirror verifies existing guidance is preserved.
func TestInitializeRepairsMissingCanonicalMirror(t *testing.T) {
	root := t.TempDir()
	writeText(t, filepath.Join(root, "CLAUDE.md"), "custom entry\n")
	engine, err := New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if got := readText(t, filepath.Join(root, "AGENTS.md")); got != "custom entry\n" {
		t.Fatalf("AGENTS.md = %q, want existing CLAUDE.md content", got)
	}
}

// TestValidateAllowsActiveWorkWithEmptyIssues verifies the active Work exception.
func TestValidateAllowsActiveWorkWithEmptyIssues(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("empty-issues")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	prdBefore := readText(t, filepath.Join(work.Path, "PRD.md"))
	indexBefore := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))

	report, err := engine.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !report.Passed() {
		t.Fatalf("active Work with empty issues should pass, got:\n%s", report)
	}
	if prdAfter := readText(t, filepath.Join(work.Path, "PRD.md")); prdAfter != prdBefore {
		t.Fatal("Validation changed the active PRD")
	}
	if indexAfter := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); indexAfter != indexBefore {
		t.Fatal("Validation changed INDEX")
	}
}

// TestCompleteRejectsSucceededWorkWithoutIssue verifies Complete preflight.
func TestCompleteRejectsSucceededWorkWithoutIssue(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.GenerateWork("missing-issue"); err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}

	_, err := engine.Complete("missing-issue", OutcomeSucceeded)
	if err == nil {
		t.Fatal("Complete() should reject a Work without an Issue")
	}
	if !errors.Is(err, ErrPreflight) {
		t.Fatalf("Complete() error = %v, want ErrPreflight", err)
	}
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "active", "missing-issue"))
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "completed", "missing-issue"))

	prd, readErr := os.ReadFile(filepath.Join(engine.Root, ".scratch", "active", "missing-issue", "PRD.md"))
	if readErr != nil {
		t.Fatalf("read PRD after rejected Complete: %v", readErr)
	}
	if strings.Contains(string(prd), "outcome:") {
		t.Fatalf("rejected Complete wrote an outcome:\n%s", prd)
	}
}

// TestCompletePreservesCoreAssetsAndSurfacesOutcome verifies successful completion.
func TestCompletePreservesCoreAssetsAndSurfacesOutcome(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("finish-engine")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-implement-engine.md"), "---\nstatus: done\ntitle: Implement engine\n---\n\nDone.\n")

	result, err := engine.Complete("finish-engine", OutcomeSucceeded)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if !result.Completed || !result.CleanupCompleted {
		t.Fatalf("Complete() result = %+v, want completed and cleaned up", result)
	}
	completed := filepath.Join(engine.Root, ".scratch", "completed", "finish-engine")
	assertFileExists(t, filepath.Join(completed, "PRD.md"))
	assertFileExists(t, filepath.Join(completed, "issues", "01-implement-engine.md"))
	assertFileExists(t, filepath.Join(completed, "HANDOFF.md"))
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "active", "finish-engine"))

	prd := readText(t, filepath.Join(completed, "PRD.md"))
	if !strings.Contains(prd, "outcome: succeeded") {
		t.Fatalf("completed PRD does not record outcome:\n%s", prd)
	}
	index := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	if !strings.Contains(index, "Outcome: succeeded") {
		t.Fatalf("INDEX does not surface outcome:\n%s", index)
	}
	report, validateErr := engine.Validate()
	if validateErr != nil || !report.Passed() {
		t.Fatalf("completed Work should validate, error = %v, report = %s", validateErr, report)
	}
}

// TestCompleteCancelledAllowsNonTerminalIssue verifies non-succeeded completion semantics.
func TestCompleteCancelledAllowsNonTerminalIssue(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("cancel-engine")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-partial.md"), "---\nstatus: in-progress\n---\n")

	if _, err := engine.Complete("cancel-engine", OutcomeCancelled); err != nil {
		t.Fatalf("cancelled Complete() error = %v", err)
	}
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "completed", "cancel-engine", "issues", "01-partial.md"))
}

// TestCompletePreflightLeavesWorkUnchanged verifies Complete preflight is non-mutating.
func TestCompletePreflightLeavesWorkUnchanged(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("move-failure")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-finish.md"), "---\nstatus: done\n---\n")
	writeText(t, filepath.Join(work.Path, "temporary.md"), "temporary context\n")
	if err := os.Mkdir(filepath.Join(work.Path, "empty-temporary"), 0o755); err != nil {
		t.Fatalf("create temporary directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(engine.Root, ".scratch", "completed", "move-failure"), []byte("occupied"), 0o644); err != nil {
		t.Fatalf("create occupied completed path: %v", err)
	}

	if _, err := engine.Complete("move-failure", OutcomeSucceeded); err == nil {
		t.Fatal("Complete() should fail when completed destination is occupied by a file")
	}
	assertFileExists(t, filepath.Join(work.Path, "temporary.md"))
	assertFileExists(t, filepath.Join(work.Path, "empty-temporary"))
	if strings.Contains(readText(t, filepath.Join(work.Path, "PRD.md")), "outcome:") {
		t.Fatal("failed Complete() left outcome on active PRD")
	}
}

// TestCompletedWorkResumesCleanup verifies idempotent retry after Complete.
func TestCompletedWorkResumesCleanup(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("retry-cleanup")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-finish.md"), "---\nstatus: done\n---\n")
	if _, err := engine.Complete("retry-cleanup", OutcomeSucceeded); err != nil {
		t.Fatalf("initial Complete() error = %v", err)
	}
	if err := os.Remove(filepath.Join(engine.Root, ".scratch", "INDEX.md")); err != nil {
		t.Fatalf("remove INDEX: %v", err)
	}

	result, err := engine.Complete("retry-cleanup", OutcomeSucceeded)
	if err != nil {
		t.Fatalf("cleanup retry error = %v", err)
	}
	if !result.RetriedCleanup || !result.CleanupCompleted {
		t.Fatalf("cleanup retry result = %+v", result)
	}
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "completed", "retry-cleanup", "PRD.md"))
	assertFileExists(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
}

// TestCompletedWorkRejectsConflictingOutcome verifies retries honor the persisted outcome.
func TestCompletedWorkRejectsConflictingOutcome(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("retry-conflict")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-finish.md"), "---\nstatus: done\n---\n")
	if _, err := engine.Complete("retry-conflict", OutcomeSucceeded); err != nil {
		t.Fatalf("initial Complete() error = %v", err)
	}

	result, err := engine.Complete("retry-conflict", OutcomeFailed)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("conflicting Complete() error = %v, want ErrConflict", err)
	}
	if ErrorCodeOf(err) != ErrorCodeConflict {
		t.Fatalf("conflicting Complete() code = %q, want %q", ErrorCodeOf(err), ErrorCodeConflict)
	}
	if result.Outcome != OutcomeSucceeded {
		t.Fatalf("conflicting Complete() result outcome = %q, want persisted %q", result.Outcome, OutcomeSucceeded)
	}
}

// TestVersionReportsCompatibilityMatrix verifies all independent version dimensions.
func TestVersionReportsCompatibilityMatrix(t *testing.T) {
	compatibility := Compatibility()
	if compatibility.RepositoryVersion != "unknown" {
		t.Fatalf("repository-independent compatibility = %+v", compatibility)
	}
	info := newTestEngine(t).Version()
	if info.SpecificationVersion != SpecificationVersion || info.SpecificationStatus != SpecificationStatus || info.SpecificationRevision != SpecificationRevision {
		t.Fatalf("specification compatibility = %+v", info)
	}
	if info.RepositoryProfile != ProfileName || info.RepositoryProfileVersion != RepositoryProfileVersion {
		t.Fatalf("profile compatibility = %+v", info)
	}
	if info.EngineVersion != EngineVersion || info.CLIVersion != CLIVersion {
		t.Fatalf("implementation compatibility = %+v", info)
	}
	if info.SpecificationRevision != "13" || info.EngineVersion != "0.1.0-rc.5" || info.CLIVersion != "0.1.0-rc.5" {
		t.Fatalf("rc.5 release candidate compatibility = %+v", info)
	}
}

// TestValidateRejectsActiveOutcome verifies the observable lifecycle invariant.
func TestValidateRejectsActiveOutcome(t *testing.T) {
	engine := newTestEngine(t)
	work, err := engine.GenerateWork("invalid-active-outcome")
	if err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}
	prdPath := filepath.Join(work.Path, "PRD.md")
	updated, setErr := setFrontMatterField([]byte(readText(t, prdPath)), "outcome", OutcomeSucceeded)
	if setErr != nil {
		t.Fatalf("set active outcome fixture: %v", setErr)
	}
	writeText(t, prdPath, string(updated))

	report, err := engine.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if report.Passed() {
		t.Fatal("Validate() should reject an active Work declaring outcome")
	}
	if !strings.Contains(report.String(), "active Work SHALL NOT declare an outcome") {
		t.Fatalf("validation report does not explain the failure:\n%s", report)
	}
}

// TestValidateRejectsBrokenRelationship verifies target resolution diagnostics.
func TestValidateRejectsBrokenRelationship(t *testing.T) {
	engine := newTestEngine(t)
	writeText(t, filepath.Join(engine.Root, "docs", "architecture", "0001-system.md"), "---\nid: ARCH-0001\nstatus: active\nrelationships:\n  - type: references\n    target: ADR-0001\n---\n\nSystem.\n")

	report, err := engine.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if report.Passed() || !strings.Contains(report.String(), "ADR-0001 does not resolve") {
		t.Fatalf("broken relationship was not reported:\n%s", report)
	}
}

// TestHealthAndMigrationExposeDeterministicOperations verifies advisory and migration surfaces.
func TestHealthAndMigrationExposeDeterministicOperations(t *testing.T) {
	engine := newTestEngine(t)
	health, err := engine.Health()
	if err != nil {
		t.Fatalf("Health() error = %v", err)
	}
	if len(health.Categories) != 6 {
		t.Fatalf("Health() categories = %d, want 6", len(health.Categories))
	}
	migration, err := engine.Migrate(SpecificationVersion)
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if !migration.Success || migration.After != SpecificationVersion {
		t.Fatalf("Migrate() result = %+v", migration)
	}
}

// TestAllocateKnowledgeIdentifierIntegratesDraft verifies draft allocation and reference regeneration.
func TestAllocateKnowledgeIdentifierIntegratesDraft(t *testing.T) {
	engine := newTestEngine(t)
	draft := filepath.Join(engine.Root, "docs", "adr", "ADR-DRAFT-engine.md")
	writeText(t, draft, "---\nstatus: created\ntitle: Engine decision\nrelationships: []\n---\n\nDecision.\n")
	architecture := filepath.Join(engine.Root, "docs", "architecture", "0001-system.md")
	writeText(t, architecture, "---\nid: ARCH-0001\nstatus: active\nrelationships:\n  - type: references\n    target: ADR-DRAFT-engine\n---\n\nSystem.\n")

	result, err := engine.AllocateKnowledgeIdentifier("adr", "docs/adr/ADR-DRAFT-engine.md")
	if err != nil {
		t.Fatalf("AllocateKnowledgeIdentifier() error = %v", err)
	}
	if result.Identifier != "ADR-0001" {
		t.Fatalf("allocated identifier = %q, want ADR-0001", result.Identifier)
	}
	assertFileNotExists(t, draft)
	finalPath := filepath.Join(engine.Root, "docs", "adr", "0001-engine.md")
	assertFileExists(t, finalPath)
	if !strings.Contains(readText(t, finalPath), "id: ADR-0001") {
		t.Fatalf("final ADR has no allocated id:\n%s", readText(t, finalPath))
	}
	if strings.Contains(readText(t, architecture), "ADR-DRAFT-engine") || !strings.Contains(readText(t, architecture), "ADR-0001") {
		t.Fatalf("managed reference was not regenerated:\n%s", readText(t, architecture))
	}
}

// newTestEngine creates an initialized Single Repository fixture.
func newTestEngine(t *testing.T) *Engine {
	t.Helper()
	root := t.TempDir()
	engine, err := New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	return engine
}

// writeText writes UTF-8 fixture content.
func writeText(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// readText reads UTF-8 fixture content.
func readText(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

// assertFileExists asserts that a fixture path exists.
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

// assertFileNotExists asserts that a fixture path is absent.
func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s not to exist, stat error = %v", path, err)
	}
}
