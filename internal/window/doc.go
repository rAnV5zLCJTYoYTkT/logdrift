// Package window implements a sliding time-window counter used by the
// logdrift pipeline to track event rates over a rolling period.
//
// A Slider divides the configured duration into discrete time slots
// (buckets). Each call to Add places the delta into the current slot,
// creating a new one when the previous slot has aged past its width.
// Stale buckets are lazily evicted so that Total and Buckets always
// reflect only the events that fall within the active window.
//
// Typical usage:
//
//	s, err := window.New(60, time.Minute) // 60 one-second buckets
//	if err != nil { ... }
//	s.Add(1)          // record an event
//	total := s.Total() // events in the last minute
package window
