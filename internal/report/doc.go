// Package report provides a lightweight summary recorder for logdrift pipeline
// runs. It tracks the number of log lines processed, lines skipped due to parse
// errors, and alerts emitted during anomaly detection.
//
// Usage:
//
//	rec := report.NewRecorder(os.Stdout)
//
//	// during pipeline execution:
//	rec.IncProcessed()
//	rec.IncSkipped()
//	rec.RecordAlert("GET /health")
//
//	// at the end of the run:
//	rec.Flush() // prints the summary table
package report
