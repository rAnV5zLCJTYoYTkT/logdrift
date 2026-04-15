// Package report provides summary reporting for logdrift anomaly detection runs.
package report

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Summary holds aggregate statistics collected during a pipeline run.
type Summary struct {
	StartTime     time.Time
	EndTime       time.Time
	LinesProcessed int
	LinesSkipped   int
	AlertsEmitted  int
	AnomalyKeys    []string
}

// Duration returns the elapsed time of the run.
func (s *Summary) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

// Recorder accumulates run statistics and can print a final report.
type Recorder struct {
	summary Summary
	w       io.Writer
}

// NewRecorder creates a Recorder that writes reports to w.
// If w is nil, os.Stdout is used.
func NewRecorder(w io.Writer) *Recorder {
	if w == nil {
		w = os.Stdout
	}
	return &Recorder{w: w, summary: Summary{StartTime: time.Now()}}
}

// IncProcessed increments the count of successfully processed lines.
func (r *Recorder) IncProcessed() { r.summary.LinesProcessed++ }

// IncSkipped increments the count of skipped (unparseable) lines.
func (r *Recorder) IncSkipped() { r.summary.LinesSkipped++ }

// RecordAlert increments alert count and records the anomaly key.
func (r *Recorder) RecordAlert(key string) {
	r.summary.AlertsEmitted++
	r.summary.AnomalyKeys = append(r.summary.AnomalyKeys, key)
}

// Flush finalises the summary timestamp and writes the report to the writer.
func (r *Recorder) Flush() {
	r.summary.EndTime = time.Now()
	s := r.summary
	fmt.Fprintf(r.w, "--- logdrift run summary ---\n")
	fmt.Fprintf(r.w, "duration:   %s\n", s.Duration().Round(time.Millisecond))
	fmt.Fprintf(r.w, "processed:  %d\n", s.LinesProcessed)
	fmt.Fprintf(r.w, "skipped:    %d\n", s.LinesSkipped)
	fmt.Fprintf(r.w, "alerts:     %d\n", s.AlertsEmitted)
	if len(s.AnomalyKeys) > 0 {
		fmt.Fprintf(r.w, "anomalies:\n")
		for _, k := range s.AnomalyKeys {
			fmt.Fprintf(r.w, "  - %s\n", k)
		}
	}
}

// GetSummary returns a copy of the current summary.
func (r *Recorder) GetSummary() Summary { return r.summary }
