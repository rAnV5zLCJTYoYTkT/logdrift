package splitter_test

import (
	"testing"

	"github.com/user/logdrift/internal/splitter"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := splitter.New("", []string{"a"})
	if err != splitter.ErrEmptyDelimiter {
		t.Fatalf("expected ErrEmptyDelimiter, got %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := splitter.New("|", nil)
	if err != splitter.ErrNoFields {
		t.Fatalf("expected ErrNoFields, got %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := splitter.New("|", []string{"level", "msg"})
	if err != nil || s == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSplit_AllFields(t *testing.T) {
	s, _ := splitter.New("|", []string{"level", "msg", "service"})
	result := s.Split("ERROR | something broke | auth")
	if result["level"] != "ERROR" {
		t.Errorf("level: got %q", result["level"])
	}
	if result["msg"] != "something broke" {
		t.Errorf("msg: got %q", result["msg"])
	}
	if result["service"] != "auth" {
		t.Errorf("service: got %q", result["service"])
	}
}

func TestSplit_MissingColumns(t *testing.T) {
	s, _ := splitter.New("|", []string{"level", "msg", "service"})
	result := s.Split("WARN")
	if result["msg"] != "" {
		t.Errorf("expected empty msg, got %q", result["msg"])
	}
	if result["service"] != "" {
		t.Errorf("expected empty service, got %q", result["service"])
	}
}

func TestSplit_ExtraColumnsIgnored(t *testing.T) {
	s, _ := splitter.New(",", []string{"a", "b"})
	result := s.Split("x,y,z,w")
	if result["b"] != "y,z,w" {
		// SplitN with n=len(fields) keeps remainder in last field
		t.Errorf("unexpected b: %q", result["b"])
	}
}

func TestSplit_TrimsWhitespace(t *testing.T) {
	s, _ := splitter.New("|", []string{"level", "msg"})
	result := s.Split("  INFO  |  hello world  ")
	if result["level"] != "INFO" {
		t.Errorf("level: got %q", result["level"])
	}
	if result["msg"] != "hello world" {
		t.Errorf("msg: got %q", result["msg"])
	}
}
