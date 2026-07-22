package engine

import "context"

// validationConformanceHooks is an internal seam. The only non-nil adapter is
// compiled by conformance_validation_fault.go under documentation_conformance.
// Normal Engine builds expose no constructor or configuration path for it.
type validationConformanceHooks interface {
	AfterSynchronize()
	BeforeValidate(context.Context) error
}
