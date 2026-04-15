package severity_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/severity"
)

func TestParse_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  severity.Level
	}{
		{"debug", severity.Debug},
		{"INFO", severity.Info},
		{"warn", severity.Warn},
		{"WARNING", severity.Warn},
		{"ERROR", severity.Error},
		{"fatal", severity.Fatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := severity.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParse_UnknownLevel(t *testing.T) {
	_, err := severity.Parse("VERBOSE")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
}

func TestLevel_String(t *testing.T) {
	if got := severity.Error.String(); got != "ERROR" {
		t.Errorf("Error.String() = %q, want %q", got, "ERROR")
	}
}

func TestAtLeast(t *testing.T) {
	if !severity.AtLeast(severity.Error, severity.Warn) {
		t.Error("expected Error >= Warn")
	}
	if severity.AtLeast(severity.Debug, severity.Info) {
		t.Error("expected Debug < Info")
	}
	if !severity.AtLeast(severity.Fatal, severity.Fatal) {
		t.Error("expected Fatal >= Fatal")
	}
}

func TestRank_Ordering(t *testing.T) {
	levels := []severity.Level{
		severity.Debug,
		severity.Info,
		severity.Warn,
		severity.Error,
		severity.Fatal,
	}
	for i := 1; i < len(levels); i++ {
		if severity.Rank(levels[i]) <= severity.Rank(levels[i-1]) {
			t.Errorf("expected rank(%v) > rank(%v)", levels[i], levels[i-1])
		}
	}
}
