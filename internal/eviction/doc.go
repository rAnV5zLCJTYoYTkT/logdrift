// Package eviction implements a thread-safe, TTL-based key eviction cache for
// use in logdrift pipelines. It is designed to bound memory consumption when
// tracking high-cardinality keys such as log fingerprints or source addresses
// over a sliding time window.
//
// Usage:
//
//	c, err := eviction.New(5 * time.Minute)
//	if err != nil { ... }
//
//	if c.Track(fingerprint) {
//	    // first time seen (or re-appeared after TTL)
//	}
//
//	// periodically purge expired entries
//	c.Evict()
package eviction
