// Package watchdog monitors a stream of events for elevated error rates.
//
// A Watchdog maintains a rolling window of bucketed event counts and
// computes the error rate on each call to Record. When the rate exceeds
// the configured threshold an AlertFunc callback is invoked, making it
// straightforward to integrate with the alert and pipeline packages.
//
// Example:
//
//	w, err := watchdog.New(time.Minute, 0.1, func(rate float64) {
//		log.Printf("high error rate: %.2f", rate)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Record(isErr)
package watchdog
