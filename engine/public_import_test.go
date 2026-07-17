// Package engine_test verifies the Engine through its published Go import path.
package engine_test

import (
	"testing"

	engine "github.com/L1UUUU/Documentation-OS-Specification/engine"
)

// TestPublishedImportPathExposesEngine verifies consumers can import the public Engine package.
func TestPublishedImportPathExposesEngine(t *testing.T) {
	if got := engine.DefaultProfile().Name; got != engine.ProfileName {
		t.Fatalf("DefaultProfile().Name = %q, want %q", got, engine.ProfileName)
	}
}
