// Package main is the entry point for the logdrift CLI tool.
//
// Usage:
//
//	logdrift [config-file]
//
// If no config file path is provided, logdrift looks for logdrift.yaml in the
// current working directory.
//
// logdrift tails the configured log file, parses each line, and maintains a
// rolling statistical baseline of observed latency values. When a new value
// deviates from the baseline by more than the configured threshold (in standard
// deviations), an alert is emitted to stderr. A summary report is written to
// stdout when the process exits.
//
// Example config (logdrift.yaml):
//
//	log_file: /var/log/app/access.log
//	window_size: 200
//	threshold: 3.0
//	cooldown: 30s
//	poll_interval: 500ms
package main
