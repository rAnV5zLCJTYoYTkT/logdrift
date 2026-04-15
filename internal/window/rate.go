package window

import (
	"time"
)

// Rate computes the average events-per-second observed across all active
// buckets in the Slider. Returns 0 when no buckets are present.
func (s *Slider) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	if len(s.buckets) == 0 {
		return 0
	}
	var total int64
	for _, b := range s.buckets {
		total += b.Count
	}
	oldest := s.buckets[0].Timestamp
	elapsed := time.Since(oldest).Seconds()
	if elapsed <= 0 {
		return float64(total)
	}
	return float64(total) / elapsed
}

// Peak returns the highest single-bucket count within the active window.
// Returns 0 when the window is empty.
func (s *Slider) Peak() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	var peak int64
	for _, b := range s.buckets {
		if b.Count > peak {
			peak = b.Count
		}
	}
	return peak
}
