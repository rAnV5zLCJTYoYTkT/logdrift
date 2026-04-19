package tracer

import "errors"

var (
	errEmptyPattern = errors.New("tracer: pattern must not be empty")
	errEmptyGroup   = errors.New("tracer: group must not be empty")
)
