// Package engine_test verifies the Engine through its published Go import path.
package engine_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	engine "github.com/L1UUUU/Documentation-OS-Specification/engine"
)

// TestPublishedImportPathExposesEngine verifies consumers can import the public Engine package.
func TestPublishedImportPathExposesEngine(t *testing.T) {
	if got := engine.DefaultProfile().Name; got != engine.ProfileName {
		t.Fatalf("DefaultProfile().Name = %q, want %q", got, engine.ProfileName)
	}
}

// TestPublishedLifecycleFailureStageContract verifies consumers can classify
// failures without parsing messages or defining adapter-local operation names.
func TestPublishedLifecycleFailureStageContract(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	_, err = instance.BeginWork(engine.BeginInput{Slug: "invalid title", Title: "Invalid"})
	if !errors.Is(err, engine.ErrInvalidInput) {
		t.Fatalf("BeginWork() error = %v, want ErrInvalidInput", err)
	}
	if stage, ok := engine.FailureStageOf(err); !ok || stage != engine.LifecycleStageBegin {
		t.Fatalf("BeginWork() stage = %q, %v, want %q, true", stage, ok, engine.LifecycleStageBegin)
	}
	var lifecycleErr *engine.LifecycleError
	if !errors.As(err, &lifecycleErr) {
		t.Fatalf("BeginWork() error %T is not a LifecycleError", err)
	}

	_, err = instance.Synchronize(engine.SyncInput{})
	if stage, ok := engine.FailureStageOf(err); !ok || stage != engine.LifecycleStageSynchronize {
		t.Fatalf("Synchronize() stage = %q, %v, want %q, true", stage, ok, engine.LifecycleStageSynchronize)
	}
	_, err = instance.Complete("missing", "unexpected")
	if stage, ok := engine.FailureStageOf(err); !ok || stage != engine.LifecycleStageComplete {
		t.Fatalf("Complete() stage = %q, %v, want %q, true", stage, ok, engine.LifecycleStageComplete)
	}

	if err := os.Remove(filepath.Join(root, "AGENTS.md")); err != nil {
		t.Fatalf("remove required file: %v", err)
	}
	report, err := instance.Validate()
	if err != nil {
		t.Fatalf("Validate() operational error = %v", err)
	}
	failure := report.Failure()
	if !errors.Is(failure, engine.ErrInvalidRepository) {
		t.Fatalf("Validation failure = %v, want ErrInvalidRepository", failure)
	}
	if stage, ok := engine.FailureStageOf(failure); !ok || stage != engine.LifecycleStageValidate {
		t.Fatalf("Validation stage = %q, %v, want %q, true", stage, ok, engine.LifecycleStageValidate)
	}
}

// TestPublishedMissingRepositoryRootHasStableClassification verifies a caller
// can distinguish an unmaterialized root from an unclassified I/O failure.
func TestPublishedMissingRepositoryRootHasStableClassification(t *testing.T) {
	_, err := engine.New(filepath.Join(t.TempDir(), "not-materialized"))
	if code := engine.ErrorCodeOf(err); code != engine.ErrorCodeInvalidRepository {
		t.Fatalf("New() code = %q, want %q (error = %v)", code, engine.ErrorCodeInvalidRepository, err)
	}
	if _, ok := engine.FailureStageOf(err); ok {
		t.Fatalf("repository construction error unexpectedly has a lifecycle stage: %v", err)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("New() error no longer preserves the filesystem cause: %v", err)
	}
}
