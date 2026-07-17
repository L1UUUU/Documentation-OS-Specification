// Package main contains CLI-level Documentation OS conformance tests.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIConformanceFullLifecycle verifies operation integration through the CLI boundary.
func TestCLIConformanceFullLifecycle(t *testing.T) {
	root := t.TempDir()
	runCLI(t, root, "init")
	runCLI(t, root, "generate", "work", "cli-lifecycle")
	issuePath := filepath.Join(root, ".scratch", "active", "cli-lifecycle", "issues", "01-implement.md")
	issue := []byte("---\nstatus: done\ntitle: Implement lifecycle\n---\n")
	if err := os.WriteFile(issuePath, issue, 0o644); err != nil {
		t.Fatalf("write Issue fixture: %v", err)
	}

	syncOutput := runCLI(t, root, "--json", "sync", "--knowledge-impact", "no-change")
	if !strings.Contains(syncOutput, `"no_knowledge_change": true`) {
		t.Fatalf("sync JSON output = %s", syncOutput)
	}
	validateOutput := runCLI(t, root, "--json", "validate")
	if !strings.Contains(validateOutput, `"status": "passed"`) {
		t.Fatalf("validate JSON output = %s", validateOutput)
	}
	completeOutput := runCLI(t, root, "--json", "complete", "cli-lifecycle", "--outcome", "succeeded")
	for _, fragment := range []string{`"completed": true`, `"cleanup_completed": true`, `"outcome": "succeeded"`} {
		if !strings.Contains(completeOutput, fragment) {
			t.Errorf("complete JSON does not contain %s: %s", fragment, completeOutput)
		}
	}
	inspectOutput := runCLI(t, root, "--json", "inspect")
	if !strings.Contains(inspectOutput, `"active_works": 0`) || !strings.Contains(inspectOutput, `"completed_works": 1`) {
		t.Fatalf("inspect JSON output = %s", inspectOutput)
	}
	healthOutput := runCLI(t, root, "--json", "health")
	if !strings.Contains(healthOutput, `"score":`) || !strings.Contains(healthOutput, `"level":`) {
		t.Fatalf("health JSON output = %s", healthOutput)
	}
}

// runCLI executes one CLI command and returns stdout after asserting success.
func runCLI(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := append([]string{"--root", root}, args...)
	var stdout, stderr bytes.Buffer
	if code := run(command, &stdout, &stderr); code != 0 {
		t.Fatalf("run(%v) code = %d, stderr = %s", args, code, stderr.String())
	}
	return stdout.String()
}
