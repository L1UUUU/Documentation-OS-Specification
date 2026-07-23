package engine

import "context"

// lifecycleConformanceHooks is an internal seam. The only non-nil adapters are
// compiled under documentation_conformance. Normal Engine builds expose no
// constructor or configuration path for fault injection.
type lifecycleConformanceHooks interface {
	AfterSynchronize()
	BeforeValidate(context.Context) error
	AfterOutcomePersisted(string) error
	AfterWorkMoved(string) error
}
