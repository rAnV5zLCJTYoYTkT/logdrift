// Package compactor provides a sliding-window log compaction layer that
// collapses repeated log entries sharing the same fingerprint into a single
// Entry with an occurrence count.
//
// Entries are evicted once they have not been seen for longer than the
// configured TTL, after which the next occurrence is treated as new.
//
// Typical usage:
//
//	c, err := compactor.New(30 * time.Second)
//	if err != nil { ... }
//	entry, isNew := c.Add(fingerprint, message)
//	if isNew {
//	    // forward to downstream pipeline
//	}
package compactor
