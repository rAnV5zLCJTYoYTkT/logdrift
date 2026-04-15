// Package config provides configuration loading and validation for logdrift.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full logdrift runtime configuration.
type Config struct {
	Watch   WatchConfig   `yaml:"watch"`
	Baseline BaselineConfig `yaml:"baseline"`
	Alert   AlertConfig   `yaml:"alert"`
	Report  ReportConfig  `yaml:"report"`
}

// WatchConfig controls file-watching behaviour.
type WatchConfig struct {
	File         string        `yaml:"file"`
	PollInterval time.Duration `yaml:"poll_interval"`
}

// BaselineConfig controls the rolling-stats window.
type BaselineConfig struct {
	WindowSize int     `yaml:"window_size"`
	Threshold  float64 `yaml:"threshold"`
}

// AlertConfig controls alert sampling / cooldown.
type AlertConfig struct {
	Cooldown time.Duration `yaml:"cooldown"`
}

// ReportConfig controls where the summary report is written.
type ReportConfig struct {
	OutputPath string `yaml:"output_path"`
}

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Watch.File == "" {
		return errors.New("config: watch.file must not be empty")
	}
	if c.Baseline.WindowSize < 2 {
		return errors.New("config: baseline.window_size must be >= 2")
	}
	if c.Baseline.Threshold <= 0 {
		return errors.New("config: baseline.threshold must be > 0")
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Watch.PollInterval == 0 {
		c.Watch.PollInterval = 500 * time.Millisecond
	}
	if c.Alert.Cooldown == 0 {
		c.Alert.Cooldown = 10 * time.Second
	}
	if c.Report.OutputPath == "" {
		c.Report.OutputPath = "-" // stdout
	}
}
