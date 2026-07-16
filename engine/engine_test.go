// Package engine tests the public Documentation OS engine operations.
package engine

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	if _, err := engine.GenerateWork("empty-issues"); err != nil {
		t.Fatalf("GenerateWork() error = %v", err)
	}

	report, err := engine.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !report.Passed() {
		t.Fatalf("active Work with empty issues should pass, got:\n%s", report)
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
