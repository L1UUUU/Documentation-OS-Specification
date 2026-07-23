//go:build documentation_conformance

package engine

import (
	"context"
	"fmt"
	"sync"
)

// ConformanceFaultPoint names a durable lifecycle boundary. The closed set is
// deliberately narrow so tagged consumers cannot inject arbitrary callbacks.
type ConformanceFaultPoint string

const (
	ConformanceFaultAfterOutcomePersisted ConformanceFaultPoint = "after-outcome-persisted"
	ConformanceFaultAfterWorkMoved        ConformanceFaultPoint = "after-work-moved"
)

type ConformanceFaultEventKind string

const (
	ConformanceFaultActivated    ConformanceFaultEventKind = "activated"
	ConformanceFaultNotTriggered ConformanceFaultEventKind = "not-triggered"
	ConformanceFaultTriggered    ConformanceFaultEventKind = "triggered"
	ConformanceFaultExhausted    ConformanceFaultEventKind = "exhausted"
)

// ConformanceFaultSpec identifies one bounded, targeted fault sequence.
type ConformanceFaultSpec struct {
	Point        ConformanceFaultPoint
	WorkSlug     string
	FirstAttempt uint64
	Count        uint64
}

// ConformanceFaultEvent contains stable audit dimensions only. It never
// exposes repository paths, injected error text, secrets, or unknown values.
type ConformanceFaultEvent struct {
	Kind      ConformanceFaultEventKind
	Point     ConformanceFaultPoint
	WorkSlug  string
	Attempt   uint64
	Remaining uint64
}

type ConformanceFaultObserver func(ConformanceFaultEvent)

// ConformanceFaultPlan is safe for concurrent calls and consumes faults only
// for a live, matching point and Work slug.
type ConformanceFaultPlan struct {
	mu               sync.Mutex
	spec             ConformanceFaultSpec
	observer         ConformanceFaultObserver
	attempts         uint64
	remaining        uint64
	exhaustedEmitted bool
}

type ConformanceFaultSnapshot struct {
	Point     ConformanceFaultPoint
	WorkSlug  string
	Attempts  uint64
	Remaining uint64
	Exhausted bool
}

func NewConformanceFaultPlan(spec ConformanceFaultSpec, observer ConformanceFaultObserver) (*ConformanceFaultPlan, error) {
	if spec.Point != ConformanceFaultAfterOutcomePersisted && spec.Point != ConformanceFaultAfterWorkMoved {
		return nil, fmt.Errorf("%w: unsupported conformance fault point", ErrInvalidInput)
	}
	if err := validateSlug(spec.WorkSlug); err != nil {
		return nil, fmt.Errorf("%w: conformance fault Work slug: %v", ErrInvalidInput, err)
	}
	if spec.FirstAttempt == 0 || spec.Count == 0 {
		return nil, fmt.Errorf("%w: conformance fault first attempt and count must be positive", ErrInvalidInput)
	}
	plan := &ConformanceFaultPlan{spec: spec, observer: observer, remaining: spec.Count}
	plan.emit(ConformanceFaultEvent{
		Kind: ConformanceFaultActivated, Point: spec.Point, WorkSlug: spec.WorkSlug, Remaining: spec.Count,
	})
	return plan, nil
}

func (p *ConformanceFaultPlan) Trigger(ctx context.Context, point ConformanceFaultPoint, workSlug string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	if point != p.spec.Point || workSlug != p.spec.WorkSlug {
		event := ConformanceFaultEvent{
			Kind: ConformanceFaultNotTriggered, Point: point, WorkSlug: workSlug, Remaining: p.remaining,
		}
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
	p.emit(ConformanceFaultEvent{
		Kind: ConformanceFaultTriggered, Point: point, WorkSlug: workSlug, Attempt: attempt, Remaining: remaining,
	})
	if emitExhausted {
		p.emit(ConformanceFaultEvent{
			Kind: ConformanceFaultExhausted, Point: point, WorkSlug: workSlug, Attempt: attempt,
		})
	}
	return fmt.Errorf("%w: conformance lifecycle fault triggered", ErrInvalidRepository)
}

func (p *ConformanceFaultPlan) Snapshot() ConformanceFaultSnapshot {
	p.mu.Lock()
	defer p.mu.Unlock()
	return ConformanceFaultSnapshot{
		Point: p.spec.Point, WorkSlug: p.spec.WorkSlug, Attempts: p.attempts,
		Remaining: p.remaining, Exhausted: p.remaining == 0,
	}
}

func (p *ConformanceFaultPlan) emit(event ConformanceFaultEvent) {
	if p.observer != nil {
		p.observer(event)
	}
}

type conformanceLifecycleHooks struct {
	plan *ConformanceFaultPlan
}

func (*conformanceLifecycleHooks) AfterSynchronize() {}
func (*conformanceLifecycleHooks) BeforeValidate(context.Context) error {
	return nil
}
func (h *conformanceLifecycleHooks) AfterOutcomePersisted(workSlug string) error {
	return h.plan.Trigger(context.Background(), ConformanceFaultAfterOutcomePersisted, workSlug)
}
func (h *conformanceLifecycleHooks) AfterWorkMoved(workSlug string) error {
	return h.plan.Trigger(context.Background(), ConformanceFaultAfterWorkMoved, workSlug)
}

// NewConformanceWithFaultPlan constructs the real Engine with a tagged,
// lifecycle-bound fault plan. This symbol is absent from normal builds.
func NewConformanceWithFaultPlan(root string, plan *ConformanceFaultPlan) (*Engine, error) {
	if plan == nil {
		return nil, fmt.Errorf("%w: conformance fault plan is required", ErrInvalidInput)
	}
	instance, err := New(root)
	if err != nil {
		return nil, err
	}
	instance.conformance = &conformanceLifecycleHooks{plan: plan}
	return instance, nil
}

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

func (h *conformanceValidationHooks) AfterOutcomePersisted(string) error { return nil }
func (h *conformanceValidationHooks) AfterWorkMoved(string) error        { return nil }

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
