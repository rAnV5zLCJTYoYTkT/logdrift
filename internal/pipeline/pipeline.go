// Package pipeline wires together the watcher, parser, baseline, sampler,
// alert notifier, and report recorder into a single cohesive run loop.
package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/logdrift/internal/alert"
	"github.com/yourorg/logdrift/internal/baseline"
	"github.com/yourorg/logdrift/internal/parser"
	"github.com/yourorg/logdrift/internal/report"
	"github.com/yourorg/logdrift/internal/sampler"
)

// Pipeline orchestrates log ingestion, anomaly detection, and alerting.
type Pipeline struct {
	window   int
	thresh   float64
	cooldown time.Duration
	notifier *alert.Notifier
	recorder *report.Recorder
}

// New constructs a Pipeline. window must be >= 2.
func New(window int, thresh float64, cooldown time.Duration, n *alert.Notifier, rec *report.Recorder) (*Pipeline, error) {
	if window < 2 {
		return nil, fmt.Errorf("pipeline: window must be >= 2, got %d", window)
	}
	return &Pipeline{
		window:   window,
		thresh:   thresh,
		cooldown: cooldown,
		notifier: n,
		recorder: rec,
	}, nil
}

// Run consumes lines from the channel until it is closed or ctx is cancelled.
// Each line is parsed; valid lines with latency are evaluated against a rolling
// baseline. Anomalies trigger an alert (subject to sampler cooldown).
func (p *Pipeline) Run(ctx context.Context, lines <-chan string) {
	stats := make(map[string]*baseline.RollingStats)
	samp := sampler.New(p.cooldown)

	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-lines:
			if !ok {
				return
			}
			entry, err := parser.Parse(line)
			if err != nil {
				p.recorder.IncSkipped()
				continue
			}
			p.recorder.IncProcessed()
			if entry.Latency == 0 {
				continue
			}
			key := entry.Method + " " + entry.Path
			rs, exists := stats[key]
			if !exists {
				rs, _ = baseline.NewRollingStats(p.window)
				stats[key] = rs
			}
			rs.Add(entry.Latency)
			if rs.IsAnomaly(entry.Latency, p.thresh) && samp.Allow(key) {
				a := alert.NewAlert(key, entry.Latency, rs)
				p.notifier.Notify(a)
				p.recorder.RecordAlert(key)
			}
		}
	}
}
