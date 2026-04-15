package pipeline

import (
	"bufio"
	"io"

	"github.com/user/logdrift/internal/alert"
	"github.com/user/logdrift/internal/baseline"
	"github.com/user/logdrift/internal/parser"
)

// Config holds the configuration for the pipeline.
type Config struct {
	WindowSize  int
	Threshold   float64
	Notifier    *alert.Notifier
}

// Pipeline reads log lines, maintains rolling stats, and emits alerts.
type Pipeline struct {
	cfg   Config
	stats *baseline.RollingStats
}

// New creates a new Pipeline with the given configuration.
func New(cfg Config) (*Pipeline, error) {
	stats, err := baseline.NewRollingStats(cfg.WindowSize)
	if err != nil {
		return nil, err
	}
	return &Pipeline{cfg: cfg, stats: stats}, nil
}

// Run reads lines from r until EOF, processing each log line.
// Anomalous lines are sent to the configured Notifier.
func (p *Pipeline) Run(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parser.Parse(line)
		if err != nil {
			// skip unparseable lines
			continue
		}

		if entry.Latency == 0 {
			continue
		}

		latency := float64(entry.Latency.Milliseconds())
		isAnomaly := p.stats.IsAnomaly(latency, p.cfg.Threshold)
		p.stats.Add(latency)

		if isAnomaly {
			a := alert.NewAlert(entry.Level, line)
			p.cfg.Notifier.Notify(a)
		}
	}
	return scanner.Err()
}
