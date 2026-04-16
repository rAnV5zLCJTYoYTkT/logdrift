package tagger

import "errors"

// ErrEmptyTag is returned when a rule has a blank tag name.
var ErrEmptyTag = errors.New("tagger: rule has empty tag")

// ErrNoKeywords is returned when a rule has no keywords.
var ErrNoKeywords = errors.New("tagger: rule has no keywords")
