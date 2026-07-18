package engine

import "errors"

// ErrorCode is a stable machine-readable classification for public operation errors.
type ErrorCode string

// LifecycleStage identifies the stable public stage at which a lifecycle
// operation failed. Values are suitable for logs, metrics, and control flow.
type LifecycleStage string

const (
	LifecycleStageBegin       LifecycleStage = "begin"
	LifecycleStageSynchronize LifecycleStage = "synchronize"
	LifecycleStageValidate    LifecycleStage = "validate"
	LifecycleStageComplete    LifecycleStage = "complete"
	LifecycleStageCleanup     LifecycleStage = "cleanup"
)

// LifecycleError preserves an operation error while attaching its stable
// lifecycle stage. Error text and errors.Is/errors.As behavior remain intact.
type LifecycleError struct {
	stage LifecycleStage
	err   error
}

func (e *LifecycleError) Error() string { return e.err.Error() }
func (e *LifecycleError) Unwrap() error { return e.err }

// FailureStage returns the machine-readable lifecycle stage.
func (e *LifecycleError) FailureStage() LifecycleStage { return e.stage }

// FailureStageOf returns the first lifecycle stage recorded in an error chain.
func FailureStageOf(err error) (LifecycleStage, bool) {
	var target *LifecycleError
	if !errors.As(err, &target) || target.FailureStage() == "" {
		return "", false
	}
	return target.FailureStage(), true
}

// withLifecycleStage adds a default stage without replacing a more precise
// stage already present in the error chain (for example Cleanup within Complete).
func withLifecycleStage(stage LifecycleStage, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := FailureStageOf(err); ok {
		return err
	}
	return &LifecycleError{stage: stage, err: err}
}

const (
	// ErrorCodeInvalidInput means a caller supplied an unsupported value.
	ErrorCodeInvalidInput ErrorCode = "invalid-input"
	// ErrorCodePreflight means a lifecycle operation was rejected before mutation.
	ErrorCodePreflight ErrorCode = "preflight-failed"
	// ErrorCodeInvalidRepository means managed repository state is malformed.
	ErrorCodeInvalidRepository ErrorCode = "invalid-repository"
	// ErrorCodeConflict means the request contradicts persisted state.
	ErrorCodeConflict ErrorCode = "conflict"
	// ErrorCodeInternal is the stable fallback for unclassified implementation or I/O failures.
	ErrorCodeInternal ErrorCode = "internal"
)

// ErrorCodeOf maps an error chain to a stable public classification.
func ErrorCodeOf(err error) ErrorCode {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return ErrorCodeInvalidInput
	case errors.Is(err, ErrConflict), errors.Is(err, ErrIdentityConflict):
		return ErrorCodeConflict
	case errors.Is(err, ErrPreflight):
		return ErrorCodePreflight
	case errors.Is(err, ErrInvalidRepository):
		return ErrorCodeInvalidRepository
	default:
		return ErrorCodeInternal
	}
}
