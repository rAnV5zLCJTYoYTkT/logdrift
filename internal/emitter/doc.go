// Package emitter provides a lightweight, thread-safe event broadcaster for
// the logdrift pipeline.
//
// An Emitter holds a named set of Sink implementations. When Emit is called
// with an Event the emitter fans the event out to every registered sink
// concurrently-safe under a read lock, collecting and joining any errors.
//
// Built-in sinks:
//
//	- WriterSink  – serialises events as JSON lines to any io.Writer.
//	- FuncSink    – adapts an arbitrary function to the Sink interface.
//
// Usage:
//
//	e := emitter.New()
//	_ = e.Register("stderr", emitter.NewWriterSink(os.Stderr))
//	_ = e.Emit(emitter.Event{Message: "high error rate detected", Score: 0.92})
package emitter
