package engine

import "errors"

// ErrorCode is a stable machine-readable classification for public operation errors.
type ErrorCode string

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
