// This file verifies the public Knowledge synchronization declaration contract.
package engine

import (
	"errors"
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
	if result.KnowledgeImpact != KnowledgeImpactChanged || result.NoKnowledgeChange {
		t.Fatalf("Synchronize() result = %+v, want explicit changed declaration", result)
	}
}

// TestSynchronizeRejectsInvalidKnowledgeImpact verifies declarations are validated before synchronization.
func TestSynchronizeRejectsInvalidKnowledgeImpact(t *testing.T) {
	instance := newTestEngine(t)
	if err := os.Remove(instance.path(instance.Profile.IndexPath)); err != nil {
		t.Fatalf("remove INDEX before test: %v", err)
	}

	if _, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpact("unknown")}); err == nil {
		t.Fatal("Synchronize() should reject an unknown Knowledge impact declaration")
	} else if !errors.Is(err, ErrInvalidInput) || ErrorCodeOf(err) != ErrorCodeInvalidInput {
		t.Fatalf("Synchronize() error = %v, code = %q", err, ErrorCodeOf(err))
	}
	if _, err := os.Stat(instance.path(instance.Profile.IndexPath)); !os.IsNotExist(err) {
		t.Fatalf("invalid declaration changed INDEX, stat error = %v", err)
	}
}

// TestSynchronizeRequiresExplicitNoChange verifies the explicit no-change result.
func TestSynchronizeRequiresExplicitNoChange(t *testing.T) {
	instance := newTestEngine(t)

	explicit, err := instance.Synchronize(SyncInput{KnowledgeImpact: KnowledgeImpactNoChange})
	if err != nil {
		t.Fatalf("Synchronize(no-change) error = %v", err)
	}
	if explicit.KnowledgeImpact != KnowledgeImpactNoChange || !explicit.NoKnowledgeChange {
		t.Fatalf("Synchronize(no-change) result = %+v, want explicit no-change", explicit)
	}
}
