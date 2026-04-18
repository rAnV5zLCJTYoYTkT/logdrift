// Package mapper translates parser.LogLine values into Entry structs suitable
// for downstream pipeline stages.
//
// A Mapper resolves the severity level, fills in a default service name when
// none is present, and copies all other fields verbatim from the parsed line.
//
// Usage:
//
//	m := mapper.New(mapper.WithDefaultService("gateway"))
//	entry, err := m.Map(line)
package mapper
