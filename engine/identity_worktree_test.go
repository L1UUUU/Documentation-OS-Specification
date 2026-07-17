// This file verifies repository-scoped Identity allocation across Git worktrees.
package engine

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestAllocateKnowledgeIdentifierUsesCommonGitDirScope verifies parallel worktrees receive distinct identities.
func TestAllocateKnowledgeIdentifierUsesCommonGitDirScope(t *testing.T) {
	repositoryRoot := filepath.Join(t.TempDir(), "repository")
	worktreeRoot := filepath.Join(t.TempDir(), "worktree")
	if err := os.MkdirAll(repositoryRoot, 0o755); err != nil {
		t.Fatalf("create repository root: %v", err)
	}
	runGit(t, repositoryRoot, "init")
	runGit(t, repositoryRoot, "config", "user.email", "documentation-os@example.invalid")
	runGit(t, repositoryRoot, "config", "user.name", "Documentation OS Test")
	primary, err := New(repositoryRoot)
	if err != nil {
		t.Fatalf("New(primary) error = %v", err)
	}
	if err := primary.Initialize(); err != nil {
		t.Fatalf("Initialize(primary) error = %v", err)
	}
	runGit(t, repositoryRoot, "add", ".")
	runGit(t, repositoryRoot, "commit", "-m", "initialize Documentation OS")
	runGit(t, repositoryRoot, "worktree", "add", "-b", "parallel-work", worktreeRoot)
	parallel, err := New(worktreeRoot)
	if err != nil {
		t.Fatalf("New(parallel) error = %v", err)
	}

	writeText(t, filepath.Join(primary.Root, "docs", "adr", "ADR-DRAFT-primary.md"), "---\nstatus: created\ntitle: Primary decision\nrelationships: []\n---\n")
	writeText(t, filepath.Join(parallel.Root, "docs", "adr", "ADR-DRAFT-parallel.md"), "---\nstatus: created\ntitle: Parallel decision\nrelationships: []\n---\n")

	type allocation struct {
		result AllocationResult
		err    error
	}
	start := make(chan struct{})
	allocations := make(chan allocation, 2)
	var workers sync.WaitGroup
	for _, request := range []struct {
		engine *Engine
		path   string
	}{{primary, "docs/adr/ADR-DRAFT-primary.md"}, {parallel, "docs/adr/ADR-DRAFT-parallel.md"}} {
		workers.Add(1)
		go func(engine *Engine, path string) {
			defer workers.Done()
			<-start
			result, err := engine.AllocateKnowledgeIdentifier("adr", path)
			allocations <- allocation{result: result, err: err}
		}(request.engine, request.path)
	}
	close(start)
	workers.Wait()
	close(allocations)

	identifiers := map[string]bool{}
	for allocation := range allocations {
		if allocation.err != nil {
			t.Errorf("AllocateKnowledgeIdentifier() error = %v", allocation.err)
			continue
		}
		identifiers[allocation.result.Identifier] = true
	}
	if !identifiers["ADR-0001"] || !identifiers["ADR-0002"] || len(identifiers) != 2 {
		t.Fatalf("allocated identities = %v, want ADR-0001 and ADR-0002", identifiers)
	}
}

// TestAllocateKnowledgeIdentifierReportsCrossWorktreeIdentityConflict verifies duplicate identities remain retryable.
func TestAllocateKnowledgeIdentifierReportsCrossWorktreeIdentityConflict(t *testing.T) {
	repositoryRoot := filepath.Join(t.TempDir(), "repository")
	worktreeRoot := filepath.Join(t.TempDir(), "worktree")
	if err := os.MkdirAll(repositoryRoot, 0o755); err != nil {
		t.Fatalf("create repository root: %v", err)
	}
	runGit(t, repositoryRoot, "init")
	runGit(t, repositoryRoot, "config", "user.email", "documentation-os@example.invalid")
	runGit(t, repositoryRoot, "config", "user.name", "Documentation OS Test")
	primary, err := New(repositoryRoot)
	if err != nil {
		t.Fatalf("New(primary) error = %v", err)
	}
	if err := primary.Initialize(); err != nil {
		t.Fatalf("Initialize(primary) error = %v", err)
	}
	runGit(t, repositoryRoot, "add", ".")
	runGit(t, repositoryRoot, "commit", "-m", "initialize Documentation OS")
	runGit(t, repositoryRoot, "worktree", "add", "-b", "conflicting-work", worktreeRoot)

	writeText(t, filepath.Join(repositoryRoot, "docs", "adr", "0001-primary.md"), "---\nid: ADR-0001\nstatus: active\ntitle: Primary decision\nrelationships: []\n---\n")
	writeText(t, filepath.Join(worktreeRoot, "docs", "adr", "0001-parallel.md"), "---\nid: ADR-0001\nstatus: active\ntitle: Parallel decision\nrelationships: []\n---\n")
	draftPath := filepath.Join(repositoryRoot, "docs", "adr", "ADR-DRAFT-next.md")
	writeText(t, draftPath, "---\nstatus: created\ntitle: Next decision\nrelationships: []\n---\n")

	_, err = primary.AllocateKnowledgeIdentifier("adr", "docs/adr/ADR-DRAFT-next.md")
	if !errors.Is(err, ErrIdentityConflict) {
		t.Fatalf("AllocateKnowledgeIdentifier() error = %v, want ErrIdentityConflict", err)
	}
	for _, expected := range []string{"ADR-0001", "rebase", "reallocate", "retry"} {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("conflict error %q does not contain %q", err, expected)
		}
	}
	assertFileExists(t, draftPath)
	assertFileNotExists(t, filepath.Join(repositoryRoot, "docs", "adr", "0002-next.md"))
}

// runGit runs Git commands for a real repository fixture.
func runGit(t *testing.T, root string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", arguments, err, output)
	}
}
