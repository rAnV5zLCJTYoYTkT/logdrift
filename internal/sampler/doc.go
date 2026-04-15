// Package sampler implements a thread-safe, key-based rate limiter used by
// logdrift to suppress duplicate anomaly alerts within a configurable cooldown
// window.
//
// # Overview
//
// When the pipeline detects an anomalous log line it may emit many alerts in
// quick succession for the same source (e.g. the same HTTP path or error
// message). The Sampler prevents alert fatigue by ensuring that at most one
// alert per key is forwarded to the Notifier within the cooldown period.
//
// # Usage
//
//	s := sampler.New(30 * time.Second)
//
//	if s.Allow(alertKey) {
//		notifier.Notify(alert)
//	}
//
// Call Evict periodically (e.g. via a time.Ticker) in long-running processes
// to reclaim memory when key cardinality is high.
package sampler
