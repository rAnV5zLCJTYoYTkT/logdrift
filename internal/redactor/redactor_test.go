package redactor_test

import (
	"strings"
	"testing"

	"github.com/user/logdrift/internal/redactor"
)

func TestScrub_Email(t *testing.T) {
	r, err := redactor.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out := r.Scrub("user logged in as alice@example.com today")
	if strings.Contains(out, "alice@example.com") {
		t.Errorf("email not redacted: %q", out)
	}
	if !strings.Contains(out, "[EMAIL]") {
		t.Errorf("expected [EMAIL] placeholder, got: %q", out)
	}
}

func TestScrub_IPAddress(t *testing.T) {
	r, _ := redactor.New()
	out := r.Scrub("request from 192.168.1.42 blocked")
	if strings.Contains(out, "192.168.1.42") {
		t.Errorf("IP not redacted: %q", out)
	}
	if !strings.Contains(out, "[IP]") {
		t.Errorf("expected [IP] placeholder, got: %q", out)
	}
}

func TestScrub_PasswordField(t *testing.T) {
	r, _ := redactor.New()
	out := r.Scrub("auth failed password=s3cr3t!")
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("password value not redacted: %q", out)
	}
}

func TestScrub_Hash(t *testing.T) {
	r, _ := redactor.New()
	hash := "d41d8cd98f00b204e9800998ecf8427e"
	out := r.Scrub("token=" + hash)
	if strings.Contains(out, hash) {
		t.Errorf("hash not redacted: %q", out)
	}
}

func TestScrub_CustomPattern(t *testing.T) {
	r, err := redactor.New(redactor.PatternReplacement{
		Pattern:     `order-\d+`,
		Replacement: "[ORDER]",
	})
	if err != nil {
		t.Fatalf("New with custom pattern: %v", err)
	}
	out := r.Scrub("processed order-98765 successfully")
	if strings.Contains(out, "order-98765") {
		t.Errorf("custom pattern not redacted: %q", out)
	}
	if !strings.Contains(out, "[ORDER]") {
		t.Errorf("expected [ORDER] placeholder, got: %q", out)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := redactor.New(redactor.PatternReplacement{Pattern: `[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestScrub_NoSensitiveData(t *testing.T) {
	r, _ := redactor.New()
	original := "application started successfully"
	out := r.Scrub(original)
	if out != original {
		t.Errorf("clean message mutated: got %q, want %q", out, original)
	}
}
