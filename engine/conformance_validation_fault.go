//go:build documentation_conformance

package engine

import (
	"context"
	"fmt"
	"sync"
)

// ValidationFaultInjector is the narrow conformance-only interface installed
// at the Engine validation seam. It is absent from normal Engine builds.
type ValidationFaultInjector interface {
	BeforeValidate(context.Context, string) error
}

// ValidationFaultSpec identifies a bounded sequence of matching validations.
type ValidationFaultSpec struct {
	WorkSlug     string
	FirstAttempt uint64
	Count        uint64
}

type ValidationFaultEventKind string

const (
	ValidationFaultActivated    ValidationFaultEventKind = "activated"
	ValidationFaultNotTriggered ValidationFaultEventKind = "not-triggered"
	ValidationFaultTriggered    ValidationFaultEventKind = "triggered"
	ValidationFaultExhausted    ValidationFaultEventKind = "exhausted"
)

// ValidationFaultEvent is deliberately path-free and secret-free so a
// consumer may persist it as conformance audit evidence.
type ValidationFaultEvent struct {
	Kind      ValidationFaultEventKind
	WorkSlug  string
	Attempt   uint64
	Remaining uint64
}

type ValidationFaultObserver func(ValidationFaultEvent)

// ValidationFaultPlan is a concurrency-safe, count-bounded conformance adapter.
type ValidationFaultPlan struct {
	mu               sync.Mutex
	spec             ValidationFaultSpec
	observer         ValidationFaultObserver
	attempts         uint64
	remaining        uint64
	exhaustedEmitted bool
}

// ValidationFaultSnapshot is a stable audit view of a plan.
type ValidationFaultSnapshot struct {
	WorkSlug  string
	Attempts  uint64
	Remaining uint64
	Exhausted bool
}

func NewValidationFaultPlan(spec ValidationFaultSpec, observer ValidationFaultObserver) (*ValidationFaultPlan, error) {
	if err := validateSlug(spec.WorkSlug); err != nil {
		return nil, fmt.Errorf("%w: validation fault Work slug: %v", ErrInvalidInput, err)
	}
	if spec.FirstAttempt == 0 || spec.Count == 0 {
		return nil, fmt.Errorf("%w: validation fault first attempt and count must be positive", ErrInvalidInput)
	}
	plan := &ValidationFaultPlan{spec: spec, observer: observer, remaining: spec.Count}
	plan.emit(ValidationFaultEvent{Kind: ValidationFaultActivated, WorkSlug: spec.WorkSlug, Remaining: spec.Count})
	return plan, nil
}

func (p *ValidationFaultPlan) BeforeValidate(ctx context.Context, workSlug string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	if workSlug != p.spec.WorkSlug {
		event := ValidationFaultEvent{Kind: ValidationFaultNotTriggered, WorkSlug: workSlug, Remaining: p.remaining}
		p.mu.Unlock()
		p.emit(event)
		return nil
	}
	p.attempts++
	attempt := p.attempts
	if attempt < p.spec.FirstAttempt || p.remaining == 0 {
		p.mu.Unlock()
		return nil
	}
	p.remaining--
	remaining := p.remaining
	emitExhausted := remaining == 0 && !p.exhaustedEmitted
	if emitExhausted {
		p.exhaustedEmitted = true
	}
	p.mu.Unlock()
	p.emit(ValidationFaultEvent{Kind: ValidationFaultTriggered, WorkSlug: workSlug, Attempt: attempt, Remaining: remaining})
	if emitExhausted {
		p.emit(ValidationFaultEvent{Kind: ValidationFaultExhausted, WorkSlug: workSlug, Attempt: attempt})
	}
	return withLifecycleStage(LifecycleStageValidate, fmt.Errorf("%w: conformance validation fault triggered", ErrInvalidRepository))
}

func (p *ValidationFaultPlan) Snapshot() ValidationFaultSnapshot {
	p.mu.Lock()
	defer p.mu.Unlock()
	return ValidationFaultSnapshot{
		WorkSlug: p.spec.WorkSlug, Attempts: p.attempts, Remaining: p.remaining,
		Exhausted: p.remaining == 0,
	}
}

func (p *ValidationFaultPlan) emit(event ValidationFaultEvent) {
	if p.observer != nil {
		p.observer(event)
	}
}

type conformanceValidationHooks struct {
	mu       sync.Mutex
	armed    bool
	workSlug string
	injector ValidationFaultInjector
}

func (h *conformanceValidationHooks) AfterSynchronize() {
	h.mu.Lock()
	h.armed = true
	h.mu.Unlock()
}

func (h *conformanceValidationHooks) BeforeValidate(ctx context.Context) error {
	h.mu.Lock()
	armed := h.armed
	if armed {
		h.armed = false
	}
	h.mu.Unlock()
	if !armed {
		return nil
	}
	return h.injector.BeforeValidate(ctx, h.workSlug)
}

// NewConformance constructs the real Engine with a conformance-only adapter.
// The symbol does not exist unless documentation_conformance is selected.
func NewConformance(root, workSlug string, injector ValidationFaultInjector) (*Engine, error) {
	if injector == nil {
		return nil, fmt.Errorf("%w: validation fault injector is required", ErrInvalidInput)
	}
	if err := validateSlug(workSlug); err != nil {
		return nil, fmt.Errorf("%w: conformance Work slug: %v", ErrInvalidInput, err)
	}
	instance, err := New(root)
	if err != nil {
		return nil, err
	}
	instance.conformance = &conformanceValidationHooks{workSlug: workSlug, injector: injector}
	return instance, nil
}

var _ ValidationFaultInjector = (*ValidationFaultPlan)(nil)
