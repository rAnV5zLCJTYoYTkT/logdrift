package router_test

import (
	"testing"

	"github.com/user/logdrift/internal/router"
	"github.com/user/logdrift/internal/severity"
)

func TestNew_NoRulesReturnsError(t *testing.T) {
	_, err := router.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptySinkReturnsError(t *testing.T) {
	_, err := router.New([]router.Rule{{Sink: "", MinLevel: severity.Info}})
	if err == nil {
		t.Fatal("expected error for empty sink name")
	}
}

func TestRoute_MatchesByLevel(t *testing.T) {
	r, err := router.New([]router.Rule{
		{Sink: "alerts", MinLevel: severity.Error},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sink, ok := r.Route(router.Entry{Level: severity.Error, Message: "boom"})
	if !ok || sink != "alerts" {
		t.Fatalf("expected alerts/true, got %q/%v", sink, ok)
	}
}

func TestRoute_BelowMinLevelNoMatch(t *testing.T) {
	r, _ := router.New([]router.Rule{
		{Sink: "alerts", MinLevel: severity.Error},
	})
	_, ok := r.Route(router.Entry{Level: severity.Info, Message: "ok"})
	if ok {
		t.Fatal("expected no match for info level")
	}
}

func TestRoute_MatchesByContains(t *testing.T) {
	r, _ := router.New([]router.Rule{
		{Sink: "audit", MinLevel: severity.Info, Contains: "login"},
	})
	sink, ok := r.Route(router.Entry{Level: severity.Info, Message: "user login succeeded"})
	if !ok || sink != "audit" {
		t.Fatalf("expected audit/true, got %q/%v", sink, ok)
	}
}

func TestRoute_ContainsMismatchNoMatch(t *testing.T) {
	r, _ := router.New([]router.Rule{
		{Sink: "audit", MinLevel: severity.Info, Contains: "login"},
	})
	_, ok := r.Route(router.Entry{Level: severity.Info, Message: "disk full"})
	if ok {
		t.Fatal("expected no match when contains not found")
	}
}

func TestRoute_FirstRuleWins(t *testing.T) {
	r, _ := router.New([]router.Rule{
		{Sink: "first", MinLevel: severity.Warn},
		{Sink: "second", MinLevel: severity.Info},
	})
	sink, ok := r.Route(router.Entry{Level: severity.Error, Message: "oops"})
	if !ok || sink != "first" {
		t.Fatalf("expected first/true, got %q/%v", sink, ok)
	}
}
