package emitter

import "errors"

// Sentinel errors exposed for callers that need to distinguish failure modes.
var (
	// ErrNoSinks is returned by Emit when no sinks have been registered.
	ErrNoSinks = errors.New("emitter: no sinks registered")

	// ErrEmptySinkName is returned by Register when an empty name is supplied.
	ErrEmptySinkName = errors.New("emitter: sink name must not be empty")

	// ErrNilSink is returned by Register when a nil Sink is supplied.
	ErrNilSink = errors.New("emitter: sink must not be nil")
)
