// Package backoff implements an exponential back-off strategy with optional
// per-key jitter, suitable for retrying transient failures such as alert
// delivery or remote sink writes inside logdrift pipelines.
//
// Usage:
//
//	b, err := backoff.New(100*time.Millisecond, 30*time.Second, 2.0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//		if err := send(); err != nil {
//			time.Sleep(b.Next("sink"))
//			continue
//		}
//		b.Reset("sink")
//		break
//	}
package backoff
