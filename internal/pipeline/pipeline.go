// Package pipeline wires together the watcher, parser, filter, baseline,
// throttle, formatter, and alert components into a single processing loop.
package pipeline

import (
	"context"
	"fmt"
	"io"

	"github.com/yourorg/logdrift/internal/alert"
	"github.com/yourorg/logdrift/internal/baseline"
	"github.com/yourorg/logdrift/internal/metrics"
	"github.com/yourorg/logdrift/internal/parser"
)

// Notifier is the sink for anomaly alerts.
type Notifier interface {
	Notify(a alert.Alert)
}

// Pipeline processes log lines from a channel and emits alerts.
type Pipeline struct {
	window   int
	stdDev   float64
	notifier Notifier
	registry *metrics.Registry
}

// New creates a Pipeline. window must be >= 2.
func New(window int, threshold float64, n Notifier) (*Pipeline, error) {
	if window < 2 {
		return nil, fmt.Errorf("pipeline: window must be >= 2, got %d", window)
	}
	return &Pipeline{
		window:   window,
		stdDev:   threshold,
		notifier: n,
		registry: metrics.NewRegistry(),
	}, nil
}

// Registry returns the internal metrics registry for inspection.
func (p *Pipeline) Registry() *metrics.Registry { return p.registry }

// Run reads lines from src until it is closed or ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context, src <-chan string, w io.Writer) error {
	stats, err := baseline.NewRollingStats(p.window)
	if err != nil {
		return err
	}
	parsed := p.registry.Counter("lines.parsed")
	skipped := p.registry.Counter("lines.skipped")
	anomalous := p.registry.Counter("lines.anomalous")

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-src:
			if !ok {
				return nil
			}
			entry, err := parser.Parse(line)
			if err != nil {
				skipped.Inc()
				continue
			}
			parsed.Inc()
			if entry.Latency == 0 {
				continue
			}
			v := entry.Latency
			if stats.IsAnomaly(v, p.stdDev) {
				anomalous.Inc()
				a := alert.NewAlert(entry.Level, fmt.Sprintf(
					"latency %.2fms exceeds %.1f sigma baseline", v, p.stdDev,
				))
				p.notifier.Notify(a)
			}
			stats.Add(v)
		}
	}
}
