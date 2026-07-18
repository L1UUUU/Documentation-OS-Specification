// Package main tests the thin CLI delegation surface.
package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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

// TestRunIssueCreateJSONAndIdempotent verifies the CLI delegates Issue creation
// and exposes the Engine's deterministic result to automation.
func TestRunIssueCreateJSONAndIdempotent(t *testing.T) {
	root := initializeCLIRepository(t)
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := instance.BeginWork(engine.BeginInput{Slug: "cli-issue", Title: "CLI Issue"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	bodyPath := filepath.Join(t.TempDir(), "body.md")
	if err := os.WriteFile(bodyPath, []byte("Implement the CLI seam.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	args := []string{"--root", root, "--json", "issue", "create", "cli-issue", "wire-command", "--title", "Wire command", "--status=open", "--body-file", bodyPath}

	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("run(create) code = %d, stderr = %s", code, stderr.String())
	}
	var created engine.CreateIssueResult
	if err := json.Unmarshal(stdout.Bytes(), &created); err != nil {
		t.Fatalf("decode create JSON: %v; output = %s", err, stdout.String())
	}
	if created != (engine.CreateIssueResult{Number: 1, Name: "01-wire-command.md", Path: ".scratch/active/cli-issue/issues/01-wire-command.md", Created: true}) {
		t.Fatalf("create result = %+v", created)
	}
	persisted, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(created.Path)))
	if err != nil {
		t.Fatalf("ReadFile(created Issue) error = %v", err)
	}
	if !strings.Contains(string(persisted), "Implement the CLI seam.\n") {
		t.Fatalf("created Issue = %s", persisted)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("run(retry) code = %d, stderr = %s", code, stderr.String())
	}
	var retried engine.CreateIssueResult
	if err := json.Unmarshal(stdout.Bytes(), &retried); err != nil {
		t.Fatalf("decode retry JSON: %v; output = %s", err, stdout.String())
	}
	created.Created = false
	if retried != created {
		t.Fatalf("retry result = %+v, want %+v", retried, created)
	}
}

// TestRunIssueCreateReportsConflict verifies contradictory retries retain the
// Engine's stable conflict classification.
func TestRunIssueCreateReportsConflict(t *testing.T) {
	root := initializeCLIRepository(t)
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := instance.BeginWork(engine.BeginInput{Slug: "conflict", Title: "Conflict"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	bodyPath := writeCLIBody(t, "Same body.\n")
	first := []string{"--root", root, "issue", "create", "conflict", "same", "--title", "Original", "--status", "open", "--body-file", bodyPath}
	var stdout, stderr bytes.Buffer
	if code := run(first, &stdout, &stderr); code != 0 {
		t.Fatalf("run(first) code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Created Issue 01-same.md") {
		t.Fatalf("human output = %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	conflicting := []string{"--root", root, "--json", "issue", "create", "conflict", "same", "--title", "Different", "--status", "open", "--body-file", bodyPath}
	if code := run(conflicting, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), `"code":"conflict"`) {
		t.Fatalf("run(conflict) code = %d, stdout = %s, stderr = %s", code, stdout.String(), stderr.String())
	}
}

// TestRunIssueCreateRejectsInvalidInvocation verifies CLI-local parsing and
// body-file failures are invalid usage rather than Engine failures.
func TestRunIssueCreateRejectsInvalidInvocation(t *testing.T) {
	root := initializeCLIRepository(t)
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing subcommand", args: []string{"issue"}, want: "issue requires `create`"},
		{name: "missing title", args: []string{"issue", "create", "work", "task", "--status", "open", "--body-file", "body.md"}, want: "--title"},
		{name: "unknown option", args: []string{"issue", "create", "work", "task", "--title", "Task", "--status", "open", "--body-file", "body.md", "--unknown"}, want: "unknown issue create option"},
		{name: "missing body file", args: []string{"issue", "create", "work", "task", "--title", "Task", "--status", "open", "--body-file", filepath.Join(root, "missing.md")}, want: "read --body-file"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := append([]string{"--root", root, "--json"}, test.args...)
			if code := run(args, &stdout, &stderr); code != 2 || !strings.Contains(stderr.String(), `"code":"invalid-input"`) || !strings.Contains(stderr.String(), test.want) {
				t.Fatalf("run() code = %d, stdout = %s, stderr = %s", code, stdout.String(), stderr.String())
			}
		})
	}
}

// TestRunIssueCreateReportsEngineErrors verifies valid CLI input reaches the
// Engine and retains its machine-readable error classification.
func TestRunIssueCreateReportsEngineErrors(t *testing.T) {
	root := initializeCLIRepository(t)
	bodyPath := writeCLIBody(t, "Engine-owned validation.\n")
	tests := []struct {
		name     string
		workSlug string
		status   string
		code     string
	}{
		{name: "missing work", workSlug: "missing", status: "open", code: "work-not-found"},
		{name: "invalid status", workSlug: "valid-work", status: "unknown", code: "invalid-input"},
	}
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := instance.BeginWork(engine.BeginInput{Slug: "valid-work", Title: "Valid Work"}); err != nil {
		t.Fatalf("BeginWork() error = %v", err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := []string{"--root", root, "--json", "issue", "create", test.workSlug, "task", "--title", "Task", "--status", test.status, "--body-file", bodyPath}
			if exit := run(args, &stdout, &stderr); exit != 1 || !strings.Contains(stderr.String(), `"code":"`+test.code+`"`) {
				t.Fatalf("run() code = %d, stdout = %s, stderr = %s", exit, stdout.String(), stderr.String())
			}
		})
	}
}

// TestCLIUsageDocumentsIssueCreate verifies the new command is discoverable.
func TestCLIUsageDocumentsIssueCreate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "issue create <work-slug> <issue-slug> --title TITLE --status STATUS --body-file PATH") {
		t.Fatalf("CLI usage = %s", stdout.String())
	}
}

func initializeCLIRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	instance, err := engine.New(root)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := instance.Initialize(); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	return root
}

func writeCLIBody(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "body.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
