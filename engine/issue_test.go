package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

func TestCreateIssueAllocatesFirstAndNextNumberWithoutFillingGaps(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.BeginWork(BeginInput{Slug: "issue-allocation", Title: "Issue allocation"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}

	first, err := engine.CreateIssue(CreateIssueInput{WorkSlug: "issue-allocation", Slug: "first", Title: "First", Status: "open", Body: "Implement the first task.\n"})
	if err != nil {
		t.Fatalf("CreateIssue(first) error = %v", err)
	}
	if first.Number != 1 || first.Name != "01-first.md" || !first.Created || first.Path != ".scratch/active/issue-allocation/issues/01-first.md" {
		t.Fatalf("CreateIssue(first) = %+v", first)
	}

	issuesPath := filepath.Join(engine.Root, ".scratch", "active", "issue-allocation", "issues")
	if err := os.Rename(filepath.Join(issuesPath, "01-first.md"), filepath.Join(issuesPath, "02-first.md")); err != nil {
		t.Fatalf("renumber fixture: %v", err)
	}
	writeText(t, filepath.Join(issuesPath, "01-existing.md"), "---\nstatus: done\ntitle: Existing\n---\nExisting body.\n")

	next, err := engine.CreateIssue(CreateIssueInput{WorkSlug: "issue-allocation", Slug: "third", Title: "Third", Status: "in-progress", Body: "Implement the third task.\n"})
	if err != nil {
		t.Fatalf("CreateIssue(next) error = %v", err)
	}
	if next.Number != 3 || next.Name != "03-third.md" || !next.Created {
		t.Fatalf("CreateIssue(next) = %+v", next)
	}
	index := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	if !strings.Contains(index, "03-third.md [in-progress] Third") {
		t.Fatalf("INDEX does not list new Issue:\n%s", index)
	}
}

func TestCreateIssueIsIdempotentAndRejectsConflictingSlugReuse(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.BeginWork(BeginInput{Slug: "issue-retry", Title: "Issue retry"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	input := CreateIssueInput{WorkSlug: "issue-retry", Slug: "same", Title: "Same", Status: "blocked", Body: "Waiting for a dependency.\n"}
	created, err := engine.CreateIssue(input)
	if err != nil {
		t.Fatalf("CreateIssue() error = %v", err)
	}
	issuePath := filepath.Join(engine.Root, filepath.FromSlash(created.Path))
	beforeIssue := readText(t, issuePath)
	indexPath := filepath.Join(engine.Root, ".scratch", "INDEX.md")
	writeText(t, indexPath, "stale INDEX\n")

	retried, err := engine.CreateIssue(input)
	if err != nil {
		t.Fatalf("CreateIssue(retry) error = %v", err)
	}
	if retried != (CreateIssueResult{Number: 1, Name: "01-same.md", Path: ".scratch/active/issue-retry/issues/01-same.md", Created: false}) {
		t.Fatalf("CreateIssue(retry) = %+v", retried)
	}
	if got := readText(t, issuePath); got != beforeIssue {
		t.Fatal("idempotent retry rewrote Issue")
	}
	beforeIndex := readText(t, indexPath)
	if !strings.Contains(beforeIndex, "01-same.md [blocked] Same") {
		t.Fatalf("idempotent retry did not reconcile stale INDEX:\n%s", beforeIndex)
	}

	for _, change := range []CreateIssueInput{
		{WorkSlug: input.WorkSlug, Slug: input.Slug, Title: "Different", Status: input.Status, Body: input.Body},
		{WorkSlug: input.WorkSlug, Slug: input.Slug, Title: input.Title, Status: "done", Body: input.Body},
		{WorkSlug: input.WorkSlug, Slug: input.Slug, Title: input.Title, Status: input.Status, Body: "Different body.\n"},
	} {
		if _, err := engine.CreateIssue(change); !errors.Is(err, ErrConflict) || ErrorCodeOf(err) != ErrorCodeConflict {
			t.Fatalf("CreateIssue(conflict) error = %v, code = %q, want ErrConflict / %q", err, ErrorCodeOf(err), ErrorCodeConflict)
		} else {
			assertIssueFailureStage(t, err)
		}
		if got := readText(t, issuePath); got != beforeIssue {
			t.Fatal("conflict changed Issue")
		}
		if got := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); got != beforeIndex {
			t.Fatal("conflict changed INDEX")
		}
	}
}

func TestCreateIssueRejectsInvalidInputWithoutMutation(t *testing.T) {
	engine := newTestEngine(t)
	if _, err := engine.BeginWork(BeginInput{Slug: "invalid-issue-input", Title: "Invalid issue input"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	valid := CreateIssueInput{WorkSlug: "invalid-issue-input", Slug: "valid", Title: "Valid", Status: "open", Body: "Valid body.\n"}
	cases := map[string]CreateIssueInput{
		"work slug":       {WorkSlug: "Invalid", Slug: valid.Slug, Title: valid.Title, Status: valid.Status, Body: valid.Body},
		"issue slug":      {WorkSlug: valid.WorkSlug, Slug: "Bad-Slug", Title: valid.Title, Status: valid.Status, Body: valid.Body},
		"status":          {WorkSlug: valid.WorkSlug, Slug: valid.Slug, Title: valid.Title, Status: "ready-for-agent", Body: valid.Body},
		"empty title":     {WorkSlug: valid.WorkSlug, Slug: valid.Slug, Title: " \t", Status: valid.Status, Body: valid.Body},
		"multiline title": {WorkSlug: valid.WorkSlug, Slug: valid.Slug, Title: "two\nlines", Status: valid.Status, Body: valid.Body},
		"empty body":      {WorkSlug: valid.WorkSlug, Slug: valid.Slug, Title: valid.Title, Status: valid.Status, Body: " \r\n\t"},
	}
	indexBefore := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md"))
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := engine.CreateIssue(input); !errors.Is(err, ErrInvalidInput) || ErrorCodeOf(err) != ErrorCodeInvalidInput {
				t.Fatalf("CreateIssue() error = %v, code = %q, want ErrInvalidInput / %q", err, ErrorCodeOf(err), ErrorCodeInvalidInput)
			} else {
				assertIssueFailureStage(t, err)
			}
		})
	}
	entries, err := os.ReadDir(filepath.Join(engine.Root, ".scratch", "active", "invalid-issue-input", "issues"))
	if err != nil || len(entries) != 0 {
		t.Fatalf("issues after invalid inputs = %v, %v", entries, err)
	}
	if got := readText(t, filepath.Join(engine.Root, ".scratch", "INDEX.md")); got != indexBefore {
		t.Fatal("invalid input changed INDEX")
	}
}

func TestCreateIssueClassifiesUnavailableAndCorruptWork(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		engine := newTestEngine(t)
		_, err := engine.CreateIssue(validCreateIssueInput("missing"))
		if !errors.Is(err, ErrWorkNotFound) || ErrorCodeOf(err) != ErrorCodeWorkNotFound {
			t.Fatalf("error = %v, code = %q", err, ErrorCodeOf(err))
		}
		assertIssueFailureStage(t, err)
	})
	t.Run("completed", func(t *testing.T) {
		engine := newTestEngine(t)
		moveWorkToCompleted(t, engine, "finished")
		_, err := engine.CreateIssue(validCreateIssueInput("finished"))
		if !errors.Is(err, ErrCompletedWork) || ErrorCodeOf(err) != ErrorCodeCompletedWork {
			t.Fatalf("error = %v, code = %q", err, ErrorCodeOf(err))
		}
		assertIssueFailureStage(t, err)
	})
	t.Run("active non-directory", func(t *testing.T) {
		engine := newTestEngine(t)
		writeText(t, filepath.Join(engine.Root, ".scratch", "active", "occupied"), "not a directory")
		_, err := engine.CreateIssue(validCreateIssueInput("occupied"))
		if !errors.Is(err, ErrInvalidRepository) {
			t.Fatalf("error = %v, want ErrInvalidRepository", err)
		}
		assertIssueFailureStage(t, err)
	})
	t.Run("completed non-directory", func(t *testing.T) {
		engine := newTestEngine(t)
		writeText(t, filepath.Join(engine.Root, ".scratch", "completed", "occupied"), "not a directory")
		_, err := engine.CreateIssue(validCreateIssueInput("occupied"))
		if !errors.Is(err, ErrInvalidRepository) {
			t.Fatalf("error = %v, want ErrInvalidRepository", err)
		}
		assertIssueFailureStage(t, err)
	})
	t.Run("both active and completed", func(t *testing.T) {
		engine := newTestEngine(t)
		beginIssueWork(t, engine, "duplicated")
		if err := os.Mkdir(filepath.Join(engine.Root, ".scratch", "completed", "duplicated"), 0o755); err != nil {
			t.Fatalf("create duplicate completed Work: %v", err)
		}
		_, err := engine.CreateIssue(validCreateIssueInput("duplicated"))
		if !errors.Is(err, ErrInvalidRepository) {
			t.Fatalf("error = %v, want ErrInvalidRepository", err)
		}
		assertIssueFailureStage(t, err)
	})
	t.Run("corrupt active work", func(t *testing.T) {
		engine := newTestEngine(t)
		beginIssueWork(t, engine, "corrupt")
		if err := os.Remove(filepath.Join(engine.Root, ".scratch", "active", "corrupt", "PRD.md")); err != nil {
			t.Fatalf("remove PRD fixture: %v", err)
		}
		_, err := engine.CreateIssue(validCreateIssueInput("corrupt"))
		if !errors.Is(err, ErrInvalidRepository) {
			t.Fatalf("error = %v, want ErrInvalidRepository", err)
		}
		assertIssueFailureStage(t, err)
	})
}

func TestCreateIssueRollsBackIssueAndIndexWhenIndexWriteFails(t *testing.T) {
	engine := newTestEngine(t)
	beginIssueWork(t, engine, "rollback")
	indexPath := filepath.Join(engine.Root, ".scratch", "INDEX.md")
	indexBefore := readText(t, indexPath)
	originalWrite := engine.writeFileAtomic
	injectedIndexErr := errors.New("injected INDEX failure")
	engine.writeFileAtomic = func(path string, data []byte, mode os.FileMode) error {
		if path == indexPath {
			return injectedIndexErr
		}
		return originalWrite(path, data, mode)
	}
	input := validCreateIssueInput("rollback")
	if _, err := engine.CreateIssue(input); !errors.Is(err, injectedIndexErr) {
		t.Fatalf("CreateIssue() error = %v", err)
	} else {
		assertIssueFailureStage(t, err)
	}
	assertFileNotExists(t, filepath.Join(engine.Root, ".scratch", "active", "rollback", "issues", "01-new-issue.md"))
	if got := readText(t, indexPath); got != indexBefore {
		t.Fatal("failed transaction did not restore INDEX")
	}
	engine.writeFileAtomic = originalWrite
	result, err := engine.CreateIssue(input)
	if err != nil || result.Number != 1 || !result.Created {
		t.Fatalf("CreateIssue(retry) = %+v, %v", result, err)
	}
}

func TestCreateIssueClassifiesWriteAndLockFailuresAtIssueStage(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		engine := newTestEngine(t)
		beginIssueWork(t, engine, "write-failure")
		injectedErr := errors.New("injected Issue write failure")
		engine.writeFileAtomic = func(string, []byte, os.FileMode) error { return injectedErr }

		_, err := engine.CreateIssue(validCreateIssueInput("write-failure"))
		if !errors.Is(err, injectedErr) || ErrorCodeOf(err) != ErrorCodeInternal {
			t.Fatalf("CreateIssue() error = %v, code = %q", err, ErrorCodeOf(err))
		}
		assertIssueFailureStage(t, err)
	})

	t.Run("lock", func(t *testing.T) {
		engine := newTestEngine(t)
		beginIssueWork(t, engine, "lock-failure")
		lockPath := filepath.Join(engine.Root, ".scratch", ".locks", "create-issue.lock")
		if err := os.Mkdir(lockPath, 0o755); err != nil {
			t.Fatalf("occupy lock path: %v", err)
		}

		_, err := engine.CreateIssue(validCreateIssueInput("lock-failure"))
		if err == nil || ErrorCodeOf(err) != ErrorCodeInternal {
			t.Fatalf("CreateIssue() error = %v, code = %q", err, ErrorCodeOf(err))
		}
		assertIssueFailureStage(t, err)
	})
}

func TestCreateIssueConcurrentCallsReceiveUniqueNumbers(t *testing.T) {
	engine := newTestEngine(t)
	beginIssueWork(t, engine, "concurrent")
	const count = 20
	results := make(chan CreateIssueResult, count)
	errs := make(chan error, count)
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			input := CreateIssueInput{WorkSlug: "concurrent", Slug: fmt.Sprintf("issue-%02d", index), Title: fmt.Sprintf("Issue %02d", index), Status: "open", Body: fmt.Sprintf("Body %02d.\n", index)}
			result, err := engine.CreateIssue(input)
			if err != nil {
				errs <- err
				return
			}
			results <- result
		}(i)
	}
	wait.Wait()
	close(results)
	close(errs)
	for err := range errs {
		t.Errorf("CreateIssue() error = %v", err)
	}
	numbers := make([]int, 0, count)
	for result := range results {
		numbers = append(numbers, result.Number)
	}
	if len(numbers) != count {
		t.Fatalf("successful results = %d, want %d", len(numbers), count)
	}
	sort.Ints(numbers)
	for index, number := range numbers {
		if number != index+1 {
			t.Fatalf("numbers = %v", numbers)
		}
	}
}

func TestCreateIssueConcurrentIdenticalCallsConverge(t *testing.T) {
	engine := newTestEngine(t)
	beginIssueWork(t, engine, "concurrent-retry")
	const count = 20
	input := validCreateIssueInput("concurrent-retry")
	results := make(chan CreateIssueResult, count)
	errs := make(chan error, count)
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			result, err := engine.CreateIssue(input)
			if err != nil {
				errs <- err
				return
			}
			results <- result
		}()
	}
	wait.Wait()
	close(results)
	close(errs)
	for err := range errs {
		t.Errorf("CreateIssue() error = %v", err)
	}
	created := 0
	for result := range results {
		if result.Number != 1 || result.Name != "01-new-issue.md" {
			t.Errorf("CreateIssue() result = %+v", result)
		}
		if result.Created {
			created++
		}
	}
	if created != 1 {
		t.Fatalf("created results = %d, want 1", created)
	}
}

func TestCreateIssueRejectsCapacityAndMalformedExistingIssues(t *testing.T) {
	t.Run("capacity", func(t *testing.T) {
		engine := newTestEngine(t)
		beginIssueWork(t, engine, "full")
		writeText(t, filepath.Join(engine.Root, ".scratch", "active", "full", "issues", "99-last.md"), "---\nstatus: done\ntitle: Last\n---\nLast body.\n")
		if _, err := engine.CreateIssue(validCreateIssueInput("full")); !errors.Is(err, ErrConflict) || ErrorCodeOf(err) != ErrorCodeConflict {
			t.Fatalf("error = %v, code = %q, want ErrConflict / %q", err, ErrorCodeOf(err), ErrorCodeConflict)
		} else {
			assertIssueFailureStage(t, err)
		}
	})
	for name, filename := range map[string]string{"malformed name": "1-bad.md", "duplicate number": "01-second.md"} {
		t.Run(name, func(t *testing.T) {
			engine := newTestEngine(t)
			beginIssueWork(t, engine, "malformed")
			issues := filepath.Join(engine.Root, ".scratch", "active", "malformed", "issues")
			writeText(t, filepath.Join(issues, "01-first.md"), "---\nstatus: done\ntitle: First\n---\nFirst body.\n")
			writeText(t, filepath.Join(issues, filename), "---\nstatus: done\ntitle: Existing\n---\nExisting body.\n")
			if _, err := engine.CreateIssue(validCreateIssueInput("malformed")); !errors.Is(err, ErrInvalidRepository) {
				t.Fatalf("error = %v, want ErrInvalidRepository", err)
			} else {
				assertIssueFailureStage(t, err)
			}
		})
	}
}

func validCreateIssueInput(workSlug string) CreateIssueInput {
	return CreateIssueInput{WorkSlug: workSlug, Slug: "new-issue", Title: "New issue", Status: "open", Body: "Implement the new issue.\n"}
}

func beginIssueWork(t *testing.T, engine *Engine, slug string) {
	t.Helper()
	if _, err := engine.BeginWork(BeginInput{Slug: slug, Title: titleFromSlug(slug)}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
}

func moveWorkToCompleted(t *testing.T, engine *Engine, slug string) {
	t.Helper()
	beginIssueWork(t, engine, slug)
	active := filepath.Join(engine.Root, ".scratch", "active", slug)
	completed := filepath.Join(engine.Root, ".scratch", "completed", slug)
	if err := os.Rename(active, completed); err != nil {
		t.Fatalf("move Work to completed: %v", err)
	}
}

func assertIssueFailureStage(t *testing.T, err error) {
	t.Helper()
	if stage, ok := FailureStageOf(err); !ok || stage != LifecycleStageIssue {
		t.Fatalf("FailureStageOf(%v) = %q, %v, want %q, true", err, stage, ok, LifecycleStageIssue)
	}
}
