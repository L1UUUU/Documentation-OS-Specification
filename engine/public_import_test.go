// Package engine_test verifies the Engine through its published Go import path.
package engine_test

import (
	"context"
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
	if engine.EngineVersion != "0.1.0-rc.8" || engine.CLIVersion != "0.1.0-rc.8" {
		t.Fatalf("public implementation versions = engine %q, CLI %q, want rc.8 candidate", engine.EngineVersion, engine.CLIVersion)
	}
}

func TestPublishedLifecycleStageValuesAreStable(t *testing.T) {
	stages := []struct {
		stage engine.LifecycleStage
		want  string
	}{
		{engine.LifecycleStageBegin, "begin"},
		{engine.LifecycleStageIssue, "issue"},
		{engine.LifecycleStageSynchronize, "synchronize"},
		{engine.LifecycleStageValidate, "validate"},
		{engine.LifecycleStageComplete, "complete"},
		{engine.LifecycleStageCleanup, "cleanup"},
	}
	for _, test := range stages {
		if got := string(test.stage); got != test.want {
			t.Errorf("LifecycleStage = %q, want %q", got, test.want)
		}
	}
}

// TestPublishedCreateIssueContract verifies external modules can use the public Issue API.
func TestPublishedCreateIssueContract(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if _, err := instance.BeginWork(engine.BeginInput{Slug: "public-issue", Title: "Public issue"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	result, err := instance.CreateIssue(engine.CreateIssueInput{
		WorkSlug: "public-issue",
		Slug:     "external-consumer",
		Title:    "External consumer",
		Status:   "open",
		Body:     "Exercise the published API.\n",
	})
	if err != nil {
		t.Fatalf("CreateIssue() error = %v", err)
	}
	if result.Number != 1 || result.Name != "01-external-consumer.md" || !result.Created {
		t.Fatalf("CreateIssue() result = %+v", result)
	}
	retried, err := instance.CreateIssueContext(context.Background(), engine.CreateIssueInput{
		WorkSlug: "public-issue",
		Slug:     "external-consumer",
		Title:    "External consumer",
		Status:   "open",
		Body:     "Exercise the published API.\n",
	})
	if err != nil || retried.Number != 1 || retried.Created {
		t.Fatalf("CreateIssueContext(retry) result = %+v, %v", retried, err)
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

	_, err = instance.CreateIssue(engine.CreateIssueInput{})
	if !errors.Is(err, engine.ErrInvalidInput) {
		t.Fatalf("CreateIssue() error = %v, want ErrInvalidInput", err)
	}
	if code := engine.ErrorCodeOf(err); code != engine.ErrorCodeInvalidInput {
		t.Fatalf("CreateIssue() code = %q, want %q", code, engine.ErrorCodeInvalidInput)
	}
	if stage, ok := engine.FailureStageOf(err); !ok || stage != engine.LifecycleStageIssue {
		t.Fatalf("CreateIssue() stage = %q, %v, want %q, true", stage, ok, engine.LifecycleStageIssue)
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
	report, err := instance.ValidateContext(context.Background())
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

// TestPublishedValidateWrapperContract keeps the original context-free entry
// point available while ValidateContext gives cancellable consumers the same
// report contract.
func TestPublishedValidateWrapperContract(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	wrapper, err := instance.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	contextual, err := instance.ValidateContext(context.Background())
	if err != nil {
		t.Fatalf("ValidateContext() error = %v", err)
	}
	if wrapper.Status != contextual.Status || len(wrapper.Issues) != len(contextual.Issues) {
		t.Fatalf("Validate wrapper = %+v, ValidateContext = %+v", wrapper, contextual)
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
