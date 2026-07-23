package main

import (
	"os"
	"strings"
	"testing"
)

// TestReleaseWorkflowExercisesNormalAndConformanceBuilds prevents a release
// tag from going green without compiling and executing the build-tag-only
// conformance surface on every supported platform.
func TestReleaseWorkflowExercisesNormalAndConformanceBuilds(t *testing.T) {
	data, err := os.ReadFile("../../.github/workflows/go.yml")
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}
	workflow := strings.ReplaceAll(string(data), "\r\n", "\n")
	testStart := strings.Index(workflow, "jobs:\n  test:")
	raceStart := strings.Index(workflow, "\n  race:")
	if testStart < 0 || raceStart <= testStart {
		t.Fatalf("release workflow test/race job boundaries not found")
	}
	matrixJob := workflow[testStart:raceStart]
	raceJob := workflow[raceStart:]
	for _, required := range []string{
		"os: [ubuntu-latest, windows-latest, macos-latest]",
		"go -C engine test -count=1 ./...",
		"go -C engine test -tags documentation_conformance -count=1 ./...",
		"go -C engine vet ./...",
		"go -C engine vet -tags documentation_conformance ./...",
		"go build ./cmd/dos",
		"go build -tags documentation_conformance ./cmd/dos",
	} {
		if !strings.Contains(matrixJob, required) {
			t.Errorf("three-platform test job missing %q", required)
		}
	}
	for _, required := range []string{
		"runs-on: ubuntu-latest",
		"go -C engine test -race -count=1 ./...",
		"go -C engine test -tags documentation_conformance -race -count=1 ./...",
	} {
		if !strings.Contains(raceJob, required) {
			t.Errorf("Linux race job missing %q", required)
		}
	}
	for _, required := range []string{
		"$env:REF_TYPE -eq 'tag' -and $env:REF_NAME.StartsWith('engine/v')",
		"go run ./scripts/external-consumer-smoke -version '${{ steps.engine-version.outputs.version }}'",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("tag publication workflow missing %q", required)
		}
	}
}
