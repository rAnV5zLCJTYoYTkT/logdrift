// Package scorer computes a normalised anomaly score in the range [0, 1]
// for each log entry processed by logdrift.
//
// The score is a weighted combination of three components:
//
//   - Latency deviation: how many standard deviations the observed latency
//     sits from the rolling mean (mapped through a sigmoid).
//   - Severity: a linear mapping from the severity rank (debug → fatal).
//   - Error rate: the fraction of recent lines that were errors (0–1).
//
// Weights default to 0.5 / 0.3 / 0.2 and can be overridden via options.
package scorer
