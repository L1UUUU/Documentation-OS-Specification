// Package main tests the thin CLI delegation surface.
package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/L1UUUU/Documentation-OS-Specification/engine"
)

// TestRunVersionJSONReturnsCompatibilityMatrix verifies machine-readable version negotiation.
func TestRunVersionJSONReturnsCompatibilityMatrix(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"--json", "version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	var got engine.VersionInfo
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("decode version JSON: %v; output = %s", err, stdout.String())
	}
	if got.SpecificationStatus != engine.SpecificationStatus || got.SpecificationRevision != engine.SpecificationRevision || got.RepositoryProfileVersion != engine.RepositoryProfileVersion {
		t.Fatalf("version matrix = %+v", got)
	}
	if got.EngineVersion != engine.EngineVersion || got.CLIVersion != engine.CLIVersion {
		t.Fatalf("version matrix = %+v", got)
	}
}

// TestRunValidateJSON verifies command-local JSON output and engine delegation.
func TestRunValidateJSON(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "validate", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "passed"`) {
		t.Fatalf("JSON validation output = %s", stdout.String())
	}
}

// TestRunCompleteRequiresCallerOutcome verifies CLI usage validation.
func TestRunCompleteRequiresCallerOutcome(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"complete", "missing", "--root", root}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "--outcome") {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
}

// TestRunSyncRequiresKnowledgeImpact verifies automation cannot omit the caller declaration.
func TestRunSyncRequiresKnowledgeImpact(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "sync"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "--knowledge-impact") {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
}

// TestRunSyncMissingKnowledgeImpactHasStableJSONCode verifies CLI usage classification.
func TestRunSyncMissingKnowledgeImpactHasStableJSONCode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--json", "sync"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), `"code":"invalid-input"`) {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
}

// TestRunSyncReportsChangedKnowledgeImpact verifies JSON exposes the caller declaration.
func TestRunSyncReportsChangedKnowledgeImpact(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "--json", "sync", "--knowledge-impact", "changed"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"no_knowledge_change": false`) {
		t.Fatalf("sync JSON output = %s", stdout.String())
	}
}

// TestRunSyncAcceptsEqualsStyleNoChange verifies compatibility with existing option syntax.
func TestRunSyncAcceptsEqualsStyleNoChange(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "--json", "sync", "--knowledge-impact=no-change"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"no_knowledge_change": true`) {
		t.Fatalf("sync JSON output = %s", stdout.String())
	}
}

// TestRunSyncRejectsMissingKnowledgeImpactValue verifies malformed options are usage errors.
func TestRunSyncRejectsMissingKnowledgeImpactValue(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "sync", "--knowledge-impact"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "requires changed or no-change") {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
}

// TestRunSyncHumanOutputDistinguishesChangedKnowledge verifies interactive output is truthful.
func TestRunSyncHumanOutputDistinguishesChangedKnowledge(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "sync", "--knowledge-impact=changed"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Knowledge edits declared") || strings.Contains(stdout.String(), "no Knowledge edits required") {
		t.Fatalf("sync human output = %s", stdout.String())
	}
}

// TestRunSyncRejectsInvalidKnowledgeImpact verifies Engine validation reaches CLI callers.
func TestRunSyncRejectsInvalidKnowledgeImpact(t *testing.T) {
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"--root", root, "--json", "sync", "--knowledge-impact=unknown"}, &stdout, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), `"code":"invalid-input"`) || !strings.Contains(stderr.String(), `invalid Knowledge impact`) {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
}

// TestCLIUsageDocumentsSyncKnowledgeImpact verifies the required declaration is discoverable.
func TestCLIUsageDocumentsSyncKnowledgeImpact(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sync --knowledge-impact <changed|no-change>") {
		t.Fatalf("CLI usage = %s", stdout.String())
	}
}
