package baseline

import (
	"errors"
	"math"
	"sync"
)

// RollingStats maintains a rolling window of log event counts
// and computes mean and standard deviation for anomaly detection.
type RollingStats struct {
	mu      sync.Mutex
	window  []float64
	capacity int
	sum     float64
	sumSq   float64
}

// NewRollingStats creates a new RollingStats with the given window capacity.
func NewRollingStats(capacity int) (*RollingStats, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}
	return &RollingStats{
		window:   make([]float64, 0, capacity),
		capacity: capacity,
	}, nil
}

// Add inserts a new value into the rolling window.
func (r *RollingStats) Add(value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.window) == r.capacity {
		old := r.window[0]
		r.window = r.window[1:]
		r.sum -= old
		r.sumSq -= old * old
	}

	r.window = append(r.window, value)
	r.sum += value
	r.sumSq += value * value
}

// Mean returns the mean of values in the current window.
func (r *RollingStats) Mean() (float64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n := float64(len(r.window))
	if n == 0 {
		return 0, errors.New("no data in window")
	}
	return r.sum / n, nil
}

// StdDev returns the population standard deviation of values in the current window.
func (r *RollingStats) StdDev() (float64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n := float64(len(r.window))
	if n == 0 {
		return 0, errors.New("no data in window")
	}

	mean := r.sum / n
	variance := (r.sumSq / n) - (mean * mean)
	if variance < 0 {
		variance = 0
	}
	return math.Sqrt(variance), nil
}

// IsAnomaly returns true if value deviates more than threshold standard deviations from the mean.
func (r *RollingStats) IsAnomaly(value, threshold float64) (bool, error) {
	mean, err := r.Mean()
	if err != nil {
		return false, err
	}
	stddev, err := r.StdDev()
	if err != nil {
		return false, err
	}
	if stddev == 0 {
		return value != mean, nil
	}
	zScore := math.Abs(value-mean) / stddev
	return zScore > threshold, nil
}

// Len returns the current number of values in the window.
func (r *RollingStats) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.window)
}
