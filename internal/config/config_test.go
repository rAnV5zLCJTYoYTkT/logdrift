package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cfg-*.yaml")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
watch:
  file: /var/log/app.log
  poll_interval: 1s
baseline:
  window_size: 60
  threshold: 3.0
alert:
  cooldown: 30s
report:
  output_path: /tmp/report.json
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Watch.File != "/var/log/app.log" {
		t.Errorf("watch.file = %q, want /var/log/app.log", cfg.Watch.File)
	}
	if cfg.Baseline.WindowSize != 60 {
		t.Errorf("window_size = %d, want 60", cfg.Baseline.WindowSize)
	}
	if cfg.Baseline.Threshold != 3.0 {
		t.Errorf("threshold = %f, want 3.0", cfg.Baseline.Threshold)
	}
}

func TestLoad_AppliesDefaults(t *testing.T) {
	path := writeTemp(t, `
watch:
  file: /tmp/test.log
baseline:
  window_size: 10
  threshold: 2.5
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Watch.PollInterval != 500*time.Millisecond {
		t.Errorf("poll_interval = %v, want 500ms", cfg.Watch.PollInterval)
	}
	if cfg.Alert.Cooldown != 10*time.Second {
		t.Errorf("cooldown = %v, want 10s", cfg.Alert.Cooldown)
	}
	if cfg.Report.OutputPath != "-" {
		t.Errorf("output_path = %q, want -", cfg.Report.OutputPath)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/cfg.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidWindowSize(t *testing.T) {
	path := writeTemp(t, `
watch:
  file: /tmp/test.log
baseline:
  window_size: 1
  threshold: 2.0
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for window_size < 2")
	}
}

func TestLoad_MissingWatchFile(t *testing.T) {
	path := writeTemp(t, `
baseline:
  window_size: 10
  threshold: 2.0
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing watch.file")
	}
}
