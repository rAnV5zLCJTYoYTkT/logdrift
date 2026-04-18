// Package cursor provides a thread-safe byte-offset tracker for log file
// tailing. It is used by the watcher and checkpoint packages to record
// how far into a file has been consumed, allowing logdrift to resume
// from the correct position after a restart or rotation.
package cursor
