// Package batcher provides a size- and time-bounded batching primitive for
// logdrift pipelines.
//
// A Batcher accumulates string items (log lines) and emits them as a slice
// to a caller-supplied callback whenever either of two conditions is met:
//
//  1. The batch reaches its configured capacity (size-triggered flush).
//  2. The configured timeout elapses since the last flush (time-triggered flush).
//
// Both conditions are safe for concurrent use. The flush callback is invoked
// in its own goroutine so that Add and Flush never block waiting for
// downstream processing to complete.
package batcher
