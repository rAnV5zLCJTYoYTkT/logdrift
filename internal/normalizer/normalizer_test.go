package normalizer_test

import (
	"testing"

	"github.com/user/logdrift/internal/normalizer"
)

func TestNormalize_UUID(t *testing.T) {
	n := normalizer.New()
	got := n.Normalize("request id=550e8400-e29b-41d4-a716-446655440000 failed")
	want := "request id=<uuid> failed"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalize_IP(t *testing.T) {
	n := normalizer.New()
	got := n.Normalize("connection from 192.168.1.42 rejected")
	want := "connection from <ip> rejected"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalize_Numbers(t *testing.T) {
	n := normalizer.New()
	got := n.Normalize("retried 3 times after 500ms")
	want := "retried <N> times after <N>ms"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalize_Hex(t *testing.T) {
	n := normalizer.New()
	got := n.Normalize("checksum mismatch: deadbeef12345678")
	want := "checksum mismatch: <hex>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalize_QuotedStrings_Disabled(t *testing.T) {
	n := normalizer.New()
	got := n.Normalize(`key "some-dynamic-value" not found`)
	if got != `key "some-dynamic-value" not found` {
		t.Fatalf("unexpected mutation: %q", got)
	}
}

func TestNormalize_QuotedStrings_Enabled(t *testing.T) {
	n := normalizer.New(normalizer.WithQuotedStrings())
	got := n.Normalize(`key "some-dynamic-value" not found`)
	want := `key "<str>" not found`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalize_StableAcrossCalls(t *testing.T) {
	n := normalizer.New()
	a := n.Normalize("user 123 logged in from 10.0.0.1")
	b := n.Normalize("user 456 logged in from 10.0.0.2")
	if a != b {
		t.Fatalf("expected same canonical form, got %q vs %q", a, b)
	}
}
