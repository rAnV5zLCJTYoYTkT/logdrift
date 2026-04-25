package trimmer_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/trimmer"
)

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := trimmer.New(nil)
	if err == nil {
		t.Fatal("expected error for nil fields")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := trimmer.New([]string{"message", ""})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := trimmer.New([]string{"message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_TrimsTargetField(t *testing.T) {
	tr, _ := trimmer.New([]string{"message"})
	in := map[string]string{"message": "  hello world  ", "level": "  info  "}
	out := tr.Apply(in)
	if got := out["message"]; got != "hello world" {
		t.Errorf("message: got %q, want %q", got, "hello world")
	}
	// Non-targeted field must remain unchanged.
	if got := out["level"]; got != "  info  " {
		t.Errorf("level: got %q, want %q", got, "  info  ")
	}
}

func TestApply_CollapseInternalWhitespace(t *testing.T) {
	tr, _ := trimmer.New([]string{"message"}, trimmer.WithCollapse())
	in := map[string]string{"message": "  foo   bar   baz  "}
	out := tr.Apply(in)
	if got := out["message"]; got != "foo bar baz" {
		t.Errorf("got %q, want %q", got, "foo bar baz")
	}
}

func TestApply_NoCollapseByDefault(t *testing.T) {
	tr, _ := trimmer.New([]string{"message"})
	in := map[string]string{"message": "  foo   bar  "}
	out := tr.Apply(in)
	if got := out["message"]; got != "foo   bar" {
		t.Errorf("got %q, want %q", got, "foo   bar")
	}
}

func TestApply_UnknownFieldPassedThrough(t *testing.T) {
	tr, _ := trimmer.New([]string{"message"})
	in := map[string]string{"service": "  auth  "}
	out := tr.Apply(in)
	if got := out["service"]; got != "  auth  " {
		t.Errorf("got %q, want %q", got, "  auth  ")
	}
}

func TestApply_EmptyMapReturnsEmptyMap(t *testing.T) {
	tr, _ := trimmer.New([]string{"message"})
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
