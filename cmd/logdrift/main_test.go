package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath_DefaultsToYaml(t *testing.T) {
	// preserve original args
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"logdrift"}
	got := configPath()
	if got != "logdrift.yaml" {
		t.Errorf("expected logdrift.yaml, got %q", got)
	}
}

func TestConfigPath_UsesFirstArg(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"logdrift", "/etc/logdrift/custom.yaml"}
	got := configPath()
	if got != "/etc/logdrift/custom.yaml" {
		t.Errorf("expected /etc/logdrift/custom.yaml, got %q", got)
	}
}

func TestRun_MissingConfigReturnsError(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	tmp := filepath.Join(t.TempDir(), "nonexistent.yaml")
	os.Args = []string{"logdrift", tmp}

	err := run()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRun_InvalidLogFileReturnsError(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "logdrift.yaml")
	content := []byte("log_file: /nonexistent/path/app.log\nwindow_size: 10\nthreshold: 2.5\n")
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"logdrift", cfgPath}

	err := run()
	if err == nil {
		t.Fatal("expected error for invalid log file, got nil")
	}
}
