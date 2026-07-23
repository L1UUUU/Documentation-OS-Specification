//go:build documentation_conformance

package engine_test

import (
	"context"
	"testing"

	engine "github.com/L1UUUU/Documentation-OS-Specification/engine"
)

// TestPublishedConformanceFaultContract proves the tagged surface is usable
// through the public module boundary rather than only by in-package tests.
func TestPublishedConformanceFaultContract(t *testing.T) {
	plan, err := engine.NewConformanceFaultPlan(engine.ConformanceFaultSpec{
		Point: engine.ConformanceFaultAfterOutcomePersisted, WorkSlug: "public-fault",
		FirstAttempt: 1, Count: 1,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := engine.NewConformanceWithFaultPlan(t.TempDir(), plan)
	if err != nil {
		t.Fatal(err)
	}
	if instance == nil {
		t.Fatal("NewConformanceWithFaultPlan returned nil Engine")
	}
	if err := plan.Trigger(context.Background(), engine.ConformanceFaultAfterWorkMoved, "other"); err != nil {
		t.Fatalf("non-target Trigger: %v", err)
	}
}
