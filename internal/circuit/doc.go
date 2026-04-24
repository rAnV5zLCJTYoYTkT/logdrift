// Package circuit provides a thread-safe circuit-breaker for use in logdrift
// pipeline stages that call external sinks or downstream services.
//
// A Breaker starts in the Closed state and allows all calls through. After a
// configurable number of consecutive failures it transitions to Open, rejecting
// all calls with ErrOpen. Once the cooldown period has elapsed it moves to
// HalfOpen and allows a single probe call; a success resets to Closed while
// another failure reopens the circuit.
//
// Usage:
//
//	br, err := circuit.New(5, 10*time.Second)
//	if err != nil { ... }
//
//	if !br.Allow() {
//	    return circuit.ErrOpen
//	}
//	if err := callDownstream(); err != nil {
//	    br.RecordFailure()
//	    return err
//	}
//	br.RecordSuccess()
package circuit
