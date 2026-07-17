package engine_test

import (
	"fmt"
	"os"

	"github.com/L1UUUU/Documentation-OS-Specification/engine"
)

// ExampleEngine demonstrates the minimum explicit synchronization contract.
func ExampleEngine() {
	root, err := os.MkdirTemp("", "documentation-os-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(root)

	instance, err := engine.New(root)
	if err != nil {
		panic(err)
	}
	if err := instance.Initialize(); err != nil {
		panic(err)
	}
	result, err := instance.Synchronize(engine.SyncInput{
		KnowledgeImpact: engine.KnowledgeImpactNoChange,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.KnowledgeImpact)
	// Output: no-change
}
