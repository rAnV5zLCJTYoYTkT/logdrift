package backoff

import "errors"

// Sentinel errors returned by New.
var (
	// ErrInvalidBase is returned when the base delay is not positive.
	ErrInvalidBase = errors.New("backoff: base delay must be positive")

	// ErrInvalidMax is returned when the max delay is less than the base delay.
	ErrInvalidMax = errors.New("backoff: max delay must be >= base delay")

	// ErrInvalidFactor is returned when the multiplicative factor is less than 1.
	ErrInvalidFactor = errors.New("backoff: factor must be >= 1")
)
