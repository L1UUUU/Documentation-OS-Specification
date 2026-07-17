// This file verifies the public Knowledge synchronization declaration contract.
package engine

import (
	"os"
	"testing"
)

// TestSynchronizeAcceptsKnowledgeChangeDeclaration verifies callers can report changed Knowledge.
func TestSynchronizeAcceptsKnowledgeChangeDeclaration(t *testing.T) {
	instance := newTestEngine(t)

	result, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactChanged})
	if err != nil {
		t.Fatalf("Synchronize() error = %v", err)
	}
	if result.NoKnowledgeChange {
		t.Fatal("Synchronize() reported no Knowledge change for a changed declaration")
	}
}

// TestSynchronizeRejectsInvalidKnowledgeImpact verifies declarations are validated before synchronization.
func TestSynchronizeRejectsInvalidKnowledgeImpact(t *testing.T) {
	instance := newTestEngine(t)
	if err := os.Remove(instance.path(instance.Profile.IndexPath)); err != nil {
		t.Fatalf("remove INDEX before test: %v", err)
	}

	if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: "unknown"}); err == nil {
		t.Fatal("Synchronize() should reject an unknown Knowledge impact declaration")
	}
	if _, err := os.Stat(instance.path(instance.Profile.IndexPath)); !os.IsNotExist(err) {
		t.Fatalf("invalid declaration changed INDEX, stat error = %v", err)
	}
}

// TestSynchronizePreservesNoChangeCompatibility verifies explicit and legacy no-change calls.
func TestSynchronizePreservesNoChangeCompatibility(t *testing.T) {
	instance := newTestEngine(t)

	explicit, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange})
	if err != nil {
		t.Fatalf("Synchronize(no-change) error = %v", err)
	}
	if !explicit.NoKnowledgeChange {
		t.Fatal("Synchronize(no-change) did not report the explicit no-change result")
	}

	legacy, err := instance.Synchronize()
	if err != nil {
		t.Fatalf("legacy Synchronize() error = %v", err)
	}
	if !legacy.NoKnowledgeChange {
		t.Fatal("legacy Synchronize() no longer preserves the no-change result")
	}

	if _, err := instance.Synchronize(
		SyncInput{KnowledgeImpact: KnowledgeImpactChanged},
		SyncInput{KnowledgeImpact: KnowledgeImpactNoChange},
	); err == nil {
		t.Fatal("Synchronize() should reject multiple Knowledge impact declarations")
	}
}
