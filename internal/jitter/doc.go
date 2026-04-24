// Package jitter provides randomised delay helpers for use in retry loops
// and backoff strategies.
//
// # Overview
//
// When multiple goroutines or processes fail simultaneously and retry at the
// same fixed interval they can create a thundering-herd effect that overwhelms
// the downstream system.  Adding a small random component — jitter — spreads
// the retries over time and smooths out the load spike.
//
// # Usage
//
//	j, err := jitter.New(10*time.Millisecond, 90*time.Millisecond)
//	if err != nil {
//		log.Fatal(err)
//	}
//	time.Sleep(j.Next()) // sleeps between 10 ms and 100 ms
package jitter
