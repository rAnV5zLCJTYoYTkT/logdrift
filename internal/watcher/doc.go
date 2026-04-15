// Package watcher implements a poll-based file tailer for logdrift.
//
// It is intentionally dependency-free, relying only on the standard library
// so that logdrift remains a single statically-linked binary.
//
// # Usage
//
//	w := watcher.New("/var/log/app.log", 0) // 0 → DefaultPollInterval
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	go func() {
//		if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//			log.Println("watcher error:", err)
//		}
//	}()
//
//	for line := range w.Lines {
//		// process line
//	}
//
// The watcher seeks to the end of the file before it begins reading,
// so only lines appended after Run is called are emitted.
package watcher
