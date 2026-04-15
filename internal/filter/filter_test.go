package filter_test

import (
	"regexp"
	"testing"

	"github.com/user/logdrift/internal/filter"
	"github.com/user/logdrift/internal/parser"
)

func line(level, msg string) parser.LogLine {
	return parser.LogLine{Level: level, Message: msg}
}

func TestParseLevel_Known(t *testing.T) {
	cases := []struct {
		input string
		want  filter.Level
	}{
		{"DEBUG", filter.LevelDebug},
		{"INFO", filter.LevelInfo},
		{"WARN", filter.LevelWarn},
		{"ERROR", filter.LevelError},
	}
	for _, tc := range cases {
		got, ok := filter.ParseLevel(tc.input)
		if !ok || got != tc.want {
			t.Errorf("ParseLevel(%q) = %v, %v; want %v, true", tc.input, got, ok, tc.want)
		}
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	_, ok := filter.ParseLevel("TRACE")
	if ok {
		t.Error("expected ok=false for unknown level")
	}
}

func TestAllow_MinLevel(t *testing.T) {
	f := filter.New(filter.Options{MinLevel: filter.LevelWarn})

	if f.Allow(line("DEBUG", "verbose")) {
		t.Error("DEBUG should be blocked by WARN min level")
	}
	if f.Allow(line("INFO", "info msg")) {
		t.Error("INFO should be blocked by WARN min level")
	}
	if !f.Allow(line("WARN", "warning")) {
		t.Error("WARN should pass WARN min level")
	}
	if !f.Allow(line("ERROR", "boom")) {
		t.Error("ERROR should pass WARN min level")
	}
}

func TestAllow_Pattern(t *testing.T) {
	f := filter.New(filter.Options{
		MinLevel: filter.LevelDebug,
		Pattern:  regexp.MustCompile(`timeout`),
	})

	if !f.Allow(line("ERROR", "connection timeout occurred")) {
		t.Error("message containing 'timeout' should be allowed")
	}
	if f.Allow(line("ERROR", "disk full")) {
		t.Error("message without 'timeout' should be blocked")
	}
}

func TestAllow_UnknownLevel(t *testing.T) {
	f := filter.New(filter.Options{MinLevel: filter.LevelDebug})
	if f.Allow(line("TRACE", "anythingtt.Error("unknown log level should be blocked")
	}
}

func TestAllow_NoPattern_AllowsAllLevels(t *testing.T) {
	f := filter.New(filter.Options{MinLevel: filter.LevelDebug})
	for _, lvl := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if !f.Allow(line(lvl, "some message")) {
			t.Errorf("expected level %s to be allowed", lvl)
		}
	}
}
