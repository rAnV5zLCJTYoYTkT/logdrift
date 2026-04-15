// Package redactor scrubs personally identifiable information (PII) and other
// sensitive tokens from log messages before logdrift processes or forwards
// them.
//
// It ships with built-in patterns for common data types:
//
//   - Email addresses  → [EMAIL]
//   - IPv4 addresses   → [IP]
//   - Credit-card numbers (Visa / Mastercard / Amex) → [CARD]
//   - Key=value credential pairs (password, token, api_key …) → [REDACTED]
//   - Long hex strings / hashes → [HASH]
//
// Additional patterns can be injected at construction time via
// PatternReplacement values passed to New.
package redactor
