// Package profiler provides per-key event-rate profiling for logdrift.
//
// Each key (typically a service name or log source) maintains an exponential
// moving average of its observed event rate. When a new sample exceeds the
// historical mean by a configurable threshold multiplier the Observe method
// returns true, signalling that the caller should raise an anomaly alert.
//
// Stale profiles are evicted automatically after a configurable TTL to
// prevent unbounded memory growth in long-running deployments.
package profiler
