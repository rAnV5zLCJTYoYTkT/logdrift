// Package dispatcher provides a fan-out mechanism that delivers each log entry
// to all registered handlers concurrently.
//
// # Overview
//
// A Dispatcher holds a list of Handler implementations. When Dispatch is
// called, every handler runs in its own goroutine. Errors from individual
// handlers are collected and returned together so that a failure in one
// handler does not prevent others from receiving the entry.
//
// # Usage
//
//	d := dispatcher.New()
//	_ = d.Register(myHandler)
//	errs := d.Dispatch(ctx, logLine)
package dispatcher
