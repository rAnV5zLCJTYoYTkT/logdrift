// Package pipeline wires together the logdrift processing stages: watch →
// parse → filter → enrich → redact → baseline → alert → report.
package pipeline

import (
	"context"
	"fmt"
	"io"

	"github.com/user/logdrift/internal/alert"
	"github.com/user/logdrift/internal/baseline"
	"github.com/user/logdrift/internal/filter"
	"github.com/user/logdrift/internal/parser"
	"github.com/user/logdrift/internal/redactor"
	"github.com/user/logdrift/internal/report"
	"github.com/user/logdrift/internal/watcher"
)

// Options configure a Pipeline.
type Options struct {
	WindowSize  int
	Threshold   float64
	MinLevel    string
	Pattern     string
	RedactPII   bool
	ExtraRules  []redactor.PatternReplacement
	Notifier    alert.Notifier
	Recorder    *report.Recorder
	PollInterval int // seconds; 0 → watcher default
}

// Pipeline coordinates the full log-processing workflow.
type Pipeline struct {
	opts     Options
	redactor *redactor.Redactor
	stats    *baseline.RollingStats
	filter   *filter.Filter
	notifier alert.Notifier
	recorder *report.Recorder
}

// New validates options and constructs a ready-to-run Pipeline.
func New(opts Options) (*Pipeline, error) {
	if opts.WindowSize < 2 {
		return nil, fmt.Errorf("pipeline: window_size must be >= 2, got %d", opts.WindowSize)
	}
	stats, err := baseline.NewRollingStats(opts.WindowSize)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}
	var red *redactor.Redactor
	if opts.RedactPII {
		red, err = redactor.New(opts.ExtraRules...)
		if err != nil {
			return nil, fmt.Errorf("pipeline: redactor: %w", err)
		}
	}
	lvl, err := filter.ParseLevel(opts.MinLevel)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}
	f := filter.New(lvl, opts.Pattern)
	n := opts.Notifier
	if n == nil {
		n = alert.NewNotifier(io.Discard)
	}
	r := opts.Recorder
	if r == nil {
		r = report.NewRecorder(io.Discard)
	}
	return &Pipeline{
		opts:     opts,
		redactor: red,
		stats:    stats,
		filter:   f,
		notifier: n,
		recorder: r,
	}, nil
}

// Run tails the file at path, processing lines until ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context, path string) error {
	w, err := watcher.New(path)
	if err != nil {
		return fmt.Errorf("pipeline: watcher: %w", err)
	}
	lines, err := w.Run(ctx)
	if err != nil {
		return fmt.Errorf("pipeline: watcher.Run: %w", err)
	}
	for raw := range lines {
		entry, ok := parser.Parse(raw)
		if !ok {
			continue
		}
		if p.redactor != nil {
			entry.Message = p.redactor.Scrub(entry.Message)
		}
		if !p.filter.Allow(entry) {
			continue
		}
		p.recorder.Record(entry)
		if entry.Latency > 0 {
			if p.stats.IsAnomaly(entry.Latency, p.opts.Threshold) {
				a := alert.NewAlert(entry, p.opts.Threshold)
				p.notifier.Notify(a)
				p.recorder.RecordAnomaly(entry)
			}
			p.stats.Add(entry.Latency)
		}
	}
	return p.recorder.Flush()
}
