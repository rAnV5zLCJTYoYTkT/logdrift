package formatter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logdrift/logdrift/internal/formatter"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func sampleEntry() formatter.Entry {
	return formatter.Entry{
		Timestamp: fixedTime,
		Level:     "ERROR",
		Message:   "latency spike detected",
		ZScore:    3.75,
		Mean:      120.5,
		StdDev:    15.2,
	}
}

func TestParseFormat_Known(t *testing.T) {
	cases := []struct {
		input string
		want  formatter.Format
	}{
		{"text", formatter.Text},
		{"TEXT", formatter.Text},
		{"", formatter.Text},
		{"json", formatter.JSON},
		{"JSON", formatter.JSON},
	}
	for _, tc := range cases {
		got, err := formatter.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Unknown(t *testing.T) {
	_, err := formatter.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention the bad value, got: %v", err)
	}
}

func TestRender_TextFormat(t *testing.T) {
	f := formatter.New(formatter.Text)
	out := f.Render(sampleEntry())

	for _, want := range []string{"ERROR", "z=3.75", "mean=120.50", "stddev=15.20", "latency spike detected", "2024-06-01T12:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q; got: %s", want, out)
		}
	}
}

func TestRender_JSONFormat(t *testing.T) {
	f := formatter.New(formatter.JSON)
	out := f.Render(sampleEntry())

	for _, want := range []string{`"level":"ERROR"`, `"z_score":3.75`, `"message":"latency spike detected"`} {
		if !strings.Contains(out, want) {
			t.Errorf("json output missing %q; got: %s", want, out)
		}
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected valid JSON object, got: %s", out)
	}
}
