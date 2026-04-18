package classifier_test

import (
	"testing"

	"github.com/user/logdrift/internal/classifier"
)

func TestNew_NoRulesReturnsError(t *testing.T) {
	_, err := classifier.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_RuleMissingKeywordsReturnsError(t *testing.T) {
	_, err := classifier.New([]classifier.Rule{
		{Keywords: nil, Category: classifier.CategoryError},
	})
	if err == nil {
		t.Fatal("expected error for rule with no keywords")
	}
}

func TestNew_RuleMissingCategoryReturnsError(t *testing.T) {
	_, err := classifier.New([]classifier.Rule{
		{Keywords: []string{"error"}, Category: ""},
	})
	if err == nil {
		t.Fatal("expected error for rule with empty category")
	}
}

func TestClassify_MatchesFirstRule(t *testing.T) {
	c, err := classifier.New([]classifier.Rule{
		{Keywords: []string{"panic", "fatal"}, Category: classifier.CategoryCritical},
		{Keywords: []string{"error", "failed"}, Category: classifier.CategoryError},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := c.Classify("fatal signal received"); got != classifier.CategoryCritical {
		t.Errorf("expected critical, got %s", got)
	}
	if got := c.Classify("connection failed"); got != classifier.CategoryError {
		t.Errorf("expected error, got %s", got)
	}
}

func TestClassify_CaseInsensitive(t *testing.T) {
	c, _ := classifier.New([]classifier.Rule{
		{Keywords: []string{"WARNING"}, Category: classifier.CategoryWarning},
	})
	if got := c.Classify("warning: disk space low"); got != classifier.CategoryWarning {
		t.Errorf("expected warning, got %s", got)
	}
}

func TestClassify_NoMatchReturnsNoise(t *testing.T) {
	c, _ := classifier.New([]classifier.Rule{
		{Keywords: []string{"error"}, Category: classifier.CategoryError},
	})
	if got := c.Classify("everything is fine"); got != classifier.CategoryNoise {
		t.Errorf("expected noise, got %s", got)
	}
}

func TestClassify_FirstRuleWins(t *testing.T) {
	c, _ := classifier.New([]classifier.Rule{
		{Keywords: []string{"error"}, Category: classifier.CategoryError},
		{Keywords: []string{"error"}, Category: classifier.CategoryWarning},
	})
	if got := c.Classify("error occurred"); got != classifier.CategoryError {
		t.Errorf("expected error category from first rule, got %s", got)
	}
}
