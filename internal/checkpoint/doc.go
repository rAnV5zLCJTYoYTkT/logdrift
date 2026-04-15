// Package checkpoint provides durable offset tracking for log file tailing.
//
// logdrift uses checkpoints to remember how far it has read into a log file
// so that it can resume from the correct position after a restart, avoiding
// both duplicate processing and missed lines.
//
// # Store
//
// Store persists a single State value (file path + byte offset) as a JSON
// file. Writes are atomic: the new state is written to a temporary file and
// then renamed over the target path, which is safe on all POSIX systems.
//
// # Manager
//
// Manager wraps a Store and provides the higher-level Resume / Commit / Reset
// interface consumed by the watcher. It also handles the case where the
// monitored log file has been replaced (different path stored in the
// checkpoint) by returning offset 0 so processing starts from the top.
package checkpoint
