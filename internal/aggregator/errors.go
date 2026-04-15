package aggregator

import "errors"

// ErrInvalidWindow is returned when a non-positive window duration is supplied.
var ErrInvalidWindow = errors.New("aggregator: window duration must be positive")
