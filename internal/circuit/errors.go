package circuit

import "fmt"

// StateError is returned when an operation is invalid for the current state.
type StateError struct {
	Current State
	Op      string
}

func (e *StateError) Error() string {
	return fmt.Sprintf("circuit: operation %q not allowed in state %v", e.Op, e.Current)
}

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}
