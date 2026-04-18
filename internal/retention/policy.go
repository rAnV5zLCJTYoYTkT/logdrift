package retention

import (
	"errors"
	"time"
)

// ParseTTL converts a human-readable duration string into a Policy.
// Accepted examples: "24h", "7d" (days are converted to hours), "30m".
func ParseTTL(s string) (Policy, error) {
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return Policy{}, errors.New("retention: invalid day duration: " + s)
		}
		return Policy{TTL: days * 24}, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return Policy{}, errors.New("retention: invalid duration: " + s)
	}
	if d <= 0 {
		return Policy{}, errors.New("retention: duration must be positive")
	}
	return Policy{TTL: d}, nil
}
