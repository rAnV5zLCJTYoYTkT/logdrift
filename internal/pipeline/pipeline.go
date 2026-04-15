// Package pipeline wires together the log-processing stages: watcher →
// parser → filter → dedup → baseline → alert → report.
package pipeline

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/user/logdrift/internal/alert"
	"github.com/user/logdrift/internal/baseline"
	"github.com/user/logdrift/internal/dedup"
	"github.com/user/logdrift/internal/filter"
	"github.com/user/logdrift/internal/parser"
	"github.com/user/logdrift/internal/report"
	"github.com/user/logdrift/internal/watcher"
)

// Options holds the configuration for a Pipeline.
type Options struct {
	LogFile    string
	WindowSize int
	Threshold  float64
	MinLevel   string
	Pattern    string
	DedupTTL   time.Duration
	PollInterval time.Duration
	Notifier   alert.Notifier
	ReportOut  io.Writer
}

// Pipeline coordinates all processing stages for a single log source.
type Pipeline struct {
	opts     Options
	stats    *baseline.RollingStats
	filter   *filter.Filter
	dedup    *dedup.Filter
	notifier alert.Notifier
	recorder *report.Recorder
}

// New validates options and constructs a ready-to-run Pipeline.
func New(opts Options) (*Pipeline, error) {
	if opts.WindowSize < 2 {
		return nil, fmt.Errorf("pipeline: window size must be at least 2, got %d", opts.WindowSize)
	}

	rs, err := baseline.NewRollingStats(opts.WindowSize)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}

	lvl, err := filter.ParseLevel(opts.MinLevel)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}

	f, err := filter.New(lvl, opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}

	notifier := opts.Notifier
	if notifier == nil {
		notifier = alert.NewNotifier(nil)
	}

	recorder := report.NewRecorder(opts.ReportOut)

	return &Pipeline{
		opts:     opts,
		stats:    rs,
		filter:   f,
		dedup:    dedup.New(opts.DedupTTL),
		notifier: notifier,
		recorder: recorder,
	}, nil
}

// Run starts watching the log file and processing lines until ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context) error {
	w, err := watcher.New(p.opts.LogFile, p.opts.PollInterval)
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	lines, err := w.Run(ctx)
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	for raw := range lines {
		entry, parseErr := parser.Parse(raw)
		if parseErr != nil {
			continue
		}

		if !p.filter.Allow(entry) {
			continue
		}

		p.recorder.Record(entry)

		if entry.Latency > 0 {
			if p.stats.IsAnomaly(entry.Latency, p.opts.Threshold) {
				if !p.dedup.IsDuplicate(entry.Message) {
					a := alert.NewAlert(entry)
					p.notifier.Notify(a)
					p.recorder.RecordAnomaly(entry)
				}
			}
			p.stats.Add(entry.Latency)
		}
	}

	return p.recorder.Flush()
}
