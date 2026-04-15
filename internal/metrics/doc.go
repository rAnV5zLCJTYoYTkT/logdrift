// Package metrics provides lightweight, concurrency-safe counters and a
// named registry for tracking runtime statistics within logdrift.
//
// # Counter
//
// A Counter is a monotonically increasing uint64 that can be incremented
// atomically from multiple goroutines without external locking.
//
// # Registry
//
// A Registry groups named counters and exposes a point-in-time Snapshot
// suitable for inclusion in structured reports or debug output.
//
// Example:
//
//	reg := metrics.NewRegistry()
//	reg.Counter("lines.parsed").Inc()
//	reg.Counter("lines.anomalous").Inc()
//	fmt.Println(reg.Snapshot())
package metrics
