package scorer

import "fmt"

// Threshold decides whether a score crosses a configured alert boundary.
type Threshold struct {
	warn     float64
	critical float64
}

// Level describes the result of a threshold evaluation.
type Level int

const (
	LevelOK       Level = iota
	LevelWarn
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelWarn:
		return "warn"
	case LevelCritical:
		return "critical"
	default:
		return "ok"
	}
}

// NewThreshold returns a Threshold. warn and critical must be in (0,1] and
// warn must be strictly less than critical.
func NewThreshold(warn, critical float64) (*Threshold, error) {
	if warn <= 0 || warn > 1 {
		return nil, fmt.Errorf("scorer: warn threshold %f out of range (0,1]", warn)
	}
	if critical <= 0 || critical > 1 {
		return nil, fmt.Errorf("scorer: critical threshold %f out of range (0,1]", critical)
	}
	if warn >= critical {
		return nil, fmt.Errorf("scorer: warn (%f) must be less than critical (%f)", warn, critical)
	}
	return &Threshold{warn: warn, critical: critical}, nil
}

// Evaluate returns the alert level for the given score.
func (t *Threshold) Evaluate(score float64) Level {
	switch {
	case score >= t.critical:
		return LevelCritical
	case score >= t.warn:
		return LevelWarn
	default:
		return LevelOK
	}
}
