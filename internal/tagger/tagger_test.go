package tagger_test

import (
	"testing"

	"github.com/user/logdrift/internal/parser"
	"github.com/user/logdrift/internal/tagger"
)

func line(msg string) parser.LogLine {
	return parser.LogLine{Message: msg}
}

func TestNew_EmptyTagReturnsError(t *testing.T) {
	_, err := tagger.New([]tagger.Rule{{Tag: "", Keywords: []string{"foo"}}})
	if err != tagger.ErrEmptyTag {
		t.Fatalf("expected ErrEmptyTag, got %v", err)
	}
}

func TestNew_NoKeywordsReturnsError(t *testing.T) {
	_, err := tagger.New([]tagger.Rule{{Tag: "db", Keywords: nil}})
	if err != tagger.ErrNoKeywords {
		t.Fatalf("expected ErrNoKeywords, got %v", err)
	}
}

func TestTag_MatchesSingleRule(t *testing.T) {
	tr, _ := tagger.New([]tagger.Rule{{Tag: "db", Keywords: []string{"sql", "query"}}})
	tags := tr.Tag(line("executing sql statement"))
	if len(tags) != 1 || tags[0] != "db" {
		t.Fatalf("expected [db], got %v", tags)
	}
}

func TestTag_CaseInsensitive(t *testing.T) {
	tr, _ := tagger.New([]tagger.Rule{{Tag: "auth", Keywords: []string{"LOGIN"}}})
	tags := tr.Tag(line("user login attempt"))
	if len(tags) != 1 || tags[0] != "auth" {
		t.Fatalf("expected [auth], got %v", tags)
	}
}

func TestTag_MultipleRulesMatch(t *testing.T) {
	tr, _ := tagger.New([]tagger.Rule{
		{Tag: "db", Keywords: []string{"query"}},
		{Tag: "slow", Keywords: []string{"timeout"}},
	})
	tags := tr.Tag(line("query timeout exceeded"))
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
}

func TestTag_NoMatch(t *testing.T) {
	tr, _ := tagger.New([]tagger.Rule{{Tag: "db", Keywords: []string{"sql"}}})
	tags := tr.Tag(line("all systems nominal"))
	if len(tags) != 0 {
		t.Fatalf("expected no tags, got %v", tags)
	}
}

func TestTag_NoDuplicateTags(t *testing.T) {
	tr, _ := tagger.New([]tagger.Rule{{Tag: "db", Keywords: []string{"sql", "query"}}})
	tags := tr.Tag(line("sql query executed"))
	if len(tags) != 1 {
		t.Fatalf("expected exactly 1 tag, got %v", tags)
	}
}
