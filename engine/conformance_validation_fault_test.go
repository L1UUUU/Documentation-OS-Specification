//go:build documentation_conformance

package engine

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
)

func TestConformanceValidationFaultFailsOnceAndSamePipelineRecovers(t *testing.T) {
	root := t.TempDir()
	events := make([]ValidationFaultEvent, 0, 3)
	injector, err := NewValidationFaultPlan(ValidationFaultSpec{
		WorkSlug: "validation-recovery", FirstAttempt: 1, Count: 1,
	}, func(event ValidationFaultEvent) { events = append(events, event) })
	if err != nil {
		t.Fatalf("NewValidationFaultPlan: %v", err)
	}
	instance, err := NewConformance(root, "validation-recovery", injector)
	if err != nil {
		t.Fatalf("NewConformance: %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	work, err := instance.BeginWork(BeginInput{
		Slug: "validation-recovery", Title: "Validation recovery",
		Issues: []BeginIssue{{Slug: "implementation", Title: "Implementation", Status: "open"}},
	})
	if err != nil {
		t.Fatalf("BeginWork: %v", err)
	}
	writeText(t, filepath.Join(work.Path, "issues", "01-implementation.md"), "---\nstatus: done\ntitle: Implementation\n---\n")
	before, err := instance.Inspect()
	if err != nil {
		t.Fatalf("Inspect before: %v", err)
	}

	if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
		t.Fatalf("Synchronize first attempt: %v", err)
	}
	_, err = instance.ValidateContext(context.Background())
	if err == nil {
		t.Fatal("ValidateContext first attempt error = nil, want injected failure")
	}
	if stage, ok := FailureStageOf(err); !ok || stage != LifecycleStageValidate {
		t.Fatalf("FailureStageOf(first) = %q, %v", stage, ok)
	}
	if code := ErrorCodeOf(err); code != ErrorCodeInvalidRepository {
		t.Fatalf("ErrorCodeOf(first) = %q", code)
	}
	afterFailure, err := instance.Inspect()
	if err != nil {
		t.Fatalf("Inspect after failure: %v", err)
	}
	if afterFailure.ActiveWorks != 1 || afterFailure.CompletedWorks != 0 {
		t.Fatalf("after failure active=%d completed=%d", afterFailure.ActiveWorks, afterFailure.CompletedWorks)
	}

	if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
		t.Fatalf("Synchronize retry: %v", err)
	}
	report, err := instance.ValidateContext(context.Background())
	if err != nil || !report.Passed() {
		t.Fatalf("ValidateContext retry error=%v report=%s", err, report)
	}
	if _, err := instance.Complete("validation-recovery", OutcomeSucceeded); err != nil {
		t.Fatalf("Complete retry: %v", err)
	}
	after, err := instance.Inspect()
	if err != nil {
		t.Fatalf("Inspect after retry: %v", err)
	}
	if after.ActiveWorks != 0 || after.CompletedWorks != 1 || after.KnowledgeDocuments != before.KnowledgeDocuments {
		t.Fatalf("after retry = %+v, before knowledge=%d", after, before.KnowledgeDocuments)
	}
	if len(after.Works) != 1 || len(after.Works[0].Issues) != 1 {
		t.Fatalf("retry duplicated Work or Issue: %+v", after.Works)
	}
	assertValidationFaultKinds(t, events, []ValidationFaultEventKind{
		ValidationFaultActivated, ValidationFaultTriggered, ValidationFaultExhausted,
	})
}

func TestConformanceValidationFaultDefaultOffTargetingCancellationAndConcurrency(t *testing.T) {
	t.Run("default off", func(t *testing.T) {
		instance := newTestEngine(t)
		if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
			t.Fatal(err)
		}
		if report, err := instance.ValidateContext(context.Background()); err != nil || !report.Passed() {
			t.Fatalf("default ValidateContext error=%v report=%s", err, report)
		}
	})

	t.Run("non target", func(t *testing.T) {
		injector, err := NewValidationFaultPlan(ValidationFaultSpec{WorkSlug: "target", FirstAttempt: 1, Count: 1}, nil)
		if err != nil {
			t.Fatal(err)
		}
		instance, err := NewConformance(t.TempDir(), "other", injector)
		if err != nil {
			t.Fatal(err)
		}
		if err := instance.Initialize(); err != nil {
			t.Fatal(err)
		}
		if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
			t.Fatal(err)
		}
		if report, err := instance.ValidateContext(context.Background()); err != nil || !report.Passed() {
			t.Fatalf("non-target error=%v report=%s", err, report)
		}
	})

	t.Run("cancel does not consume", func(t *testing.T) {
		injector, err := NewValidationFaultPlan(ValidationFaultSpec{WorkSlug: "target", FirstAttempt: 1, Count: 1}, nil)
		if err != nil {
			t.Fatal(err)
		}
		instance, err := NewConformance(t.TempDir(), "target", injector)
		if err != nil {
			t.Fatal(err)
		}
		if err := instance.Initialize(); err != nil {
			t.Fatal(err)
		}
		if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := instance.ValidateContext(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled error=%v", err)
		}
		if _, err := instance.ValidateContext(context.Background()); err == nil {
			t.Fatal("first live validation did not consume the fault")
		}
	})

	t.Run("first attempt and count", func(t *testing.T) {
		injector, err := NewValidationFaultPlan(ValidationFaultSpec{WorkSlug: "target", FirstAttempt: 2, Count: 2}, nil)
		if err != nil {
			t.Fatal(err)
		}
		for attempt, wantFailure := range []bool{false, true, true, false} {
			instance, err := NewConformance(t.TempDir(), "target", injector)
			if err != nil {
				t.Fatal(err)
			}
			if err := instance.Initialize(); err != nil {
				t.Fatal(err)
			}
			if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
				t.Fatal(err)
			}
			_, err = instance.ValidateContext(context.Background())
			if (err != nil) != wantFailure {
				t.Fatalf("attempt %d error=%v, wantFailure=%v", attempt+1, err, wantFailure)
			}
		}
		snapshot := injector.Snapshot()
		if snapshot.Attempts != 4 || snapshot.Remaining != 0 || !snapshot.Exhausted {
			t.Fatalf("snapshot=%+v", snapshot)
		}
	})

	t.Run("concurrent once", func(t *testing.T) {
		injector, err := NewValidationFaultPlan(ValidationFaultSpec{WorkSlug: "target", FirstAttempt: 1, Count: 1}, nil)
		if err != nil {
			t.Fatal(err)
		}
		const callers = 8
		var wg sync.WaitGroup
		errs := make(chan error, callers)
		for range callers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				instance, err := NewConformance(t.TempDir(), "target", injector)
				if err != nil {
					errs <- err
					return
				}
				if err := instance.Initialize(); err != nil {
					errs <- err
					return
				}
				if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange}); err != nil {
					errs <- err
					return
				}
				_, err = instance.ValidateContext(context.Background())
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)
		failed := 0
		for err := range errs {
			if err != nil {
				failed++
			}
		}
		if failed != 1 {
			t.Fatalf("concurrent failures=%d, want 1", failed)
		}
	})
}

func assertValidationFaultKinds(t *testing.T, got []ValidationFaultEvent, want []ValidationFaultEventKind) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("events=%+v, want kinds=%+v", got, want)
	}
	for index := range want {
		if got[index].Kind != want[index] {
			t.Fatalf("events[%d]=%+v, want kind=%s", index, got[index], want[index])
		}
	}
}
