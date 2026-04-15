package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logdrift/internal/alert"
	"github.com/yourorg/logdrift/internal/config"
	"github.com/yourorg/logdrift/internal/pipeline"
	"github.com/yourorg/logdrift/internal/report"
	"github.com/yourorg/logdrift/internal/watcher"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "logdrift: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(configPath())
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	w, err := watcher.New(cfg.LogFile, cfg.PollInterval)
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}

	notifier := alert.NewNotifier(os.Stderr)
	recorder := report.NewRecorder(os.Stdout)

	p, err := pipeline.New(cfg.WindowSize, cfg.Threshold, cfg.Cooldown, notifier, recorder)
	if err != nil {
		return fmt.Errorf("creating pipeline: %w", err)
	}

	lines := w.Run(ctx)
	p.Run(ctx, lines)

	return recorder.Flush()
}

func configPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "logdrift.yaml"
}
