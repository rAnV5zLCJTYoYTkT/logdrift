package labeler

import (
	"testing"
)

func TestNew_NoRulesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Contains: "error", Labels: map[string]string{"sev": "high"}}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NoLabelsReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "level", Contains: "error", Labels: map[string]string{}}})
	if err == nil {
		t.Fatal("expected error for empty labels map")
	}
}

func TestApply_MatchingRule(t *testing.T) {
	l, err := New([]Rule{
		{Field: "level", Contains: "error", Labels: map[string]string{"severity": "high"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := l.Apply(map[string]string{"level": "ERROR"})
	if out["severity"] != "high" {
		t.Errorf("expected severity=high, got %q", out["severity"])
	}
}

func TestApply_NoMatchReturnsEmpty(t *testing.T) {
	l, _ := New([]Rule{
		{Field: "level", Contains: "error", Labels: map[string]string{"severity": "high"}},
	})
	out := l.Apply(map[string]string{"level": "info"})
	if len(out) != 0 {
		t.Errorf("expected no labels, got %v", out)
	}
}

func TestApply_MultipleRulesMerged(t *testing.T) {
	l, _ := New([]Rule{
		{Field: "level", Contains: "error", Labels: map[string]string{"severity": "high"}},
		{Field: "service", Contains: "auth", Labels: map[string]string{"team": "identity"}},
	})
	out := l.Apply(map[string]string{"level": "error", "service": "auth-api"})
	if out["severity"] != "high" {
		t.Errorf("expected severity=high")
	}
	if out["team"] != "identity" {
		t.Errorf("expected team=identity")
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	l, _ := New([]Rule{
		{Field: "region", Contains: "us", Labels: map[string]string{"geo": "americas"}},
	})
	out := l.Apply(map[string]string{"level": "info"})
	if len(out) != 0 {
		t.Errorf("expected no labels for missing field, got %v", out)
	}
}
