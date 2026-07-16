// Package main tests the thin CLI delegation surface.
package main

import (
	"bytes"
	"strings"
	"testing"

	"documentation-os/engine"
)

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
