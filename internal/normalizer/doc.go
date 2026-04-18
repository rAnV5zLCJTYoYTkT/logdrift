// Package normalizer provides a Normalizer that replaces dynamic tokens
// (UUIDs, IP addresses, numbers, hex strings, quoted literals) in log
// messages with stable placeholders.
//
// Normalizing messages before fingerprinting or deduplication greatly
// reduces cardinality and improves anomaly-detection accuracy.
//
// Basic usage:
//
//	n := normalizer.New(normalizer.WithQuotedStrings())
//	canonical := n.Normalize(rawMessage)
package normalizer
