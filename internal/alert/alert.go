package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert holds information about a detected anomaly.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Value     float64
	Mean      float64
	StdDev    float64
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf(
		"[%s] %s | value=%.2f mean=%.2f stddev=%.2f | %s",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Value,
		a.Mean,
		a.StdDev,
		a.Message,
	)
}

// Notifier dispatches alerts to a configured writer.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stderr is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stderr
	}
	return &Notifier{out: w}
}

// Notify formats and writes the alert to the configured writer.
func (n *Notifier) Notify(a Alert) error {
	_, err := fmt.Fprintln(n.out, a.String())
	return err
}

// NewAlert constructs an Alert with the current timestamp.
func NewAlert(level Level, value, mean, stddev float64, message string) Alert {
	return Alert{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Value:     value,
		Mean:      mean,
		StdDev:    stddev,
	}
}
