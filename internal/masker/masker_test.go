package masker_test

import (
	"testing"

	"github.com/ourorg/logdrift/internal/masker"
)

func TestNew_NoFieldsReturnsError(t *testing.T) {
	_, err := masker.New(nil)
	if err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}

func TestNew_EmptyFieldNameReturnsError(t *testing.T) {
	_, err := masker.New([]string{"token", ""})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestMask_ReplacesMatchingField(t *testing.T) {
	m, err := masker.New([]string{"password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Mask("user=alice password=secret123 action=login")
	want := "user=alice password=[MASKED] action=login"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_CaseInsensitiveFieldMatch(t *testing.T) {
	m, _ := masker.New([]string{"token"})
	got := m.Mask("Token=abc123 level=info")
	want := "Token=[MASKED] level=info"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_NoMatchLeavesLineUnchanged(t *testing.T) {
	m, _ := masker.New([]string{"secret"})
	input := "user=bob action=read status=ok"
	got := m.Mask(input)
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestMask_MultipleFieldsMasked(t *testing.T) {
	m, _ := masker.New([]string{"password", "token"})
	got := m.Mask("user=carol password=hunter2 token=xyz789 ok=true")
	want := "user=carol password=[MASKED] token=[MASKED] ok=true"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	m, _ := masker.New([]string{"api_key"}, masker.WithPlaceholder("***"))
	got := m.Mask("service=payments api_key=sk_live_abc")
	want := "service=payments api_key=***"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_TokenWithoutEquals(t *testing.T) {
	m, _ := masker.New([]string{"password"})
	input := "bare-token level=warn"
	got := m.Mask(input)
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}
