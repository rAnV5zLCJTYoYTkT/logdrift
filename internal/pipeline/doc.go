// Package pipeline wires together the parser, rolling baseline, and alert
// notifier into a single streaming processing unit.
//
// Usage:
//
//	notifier := alert.NewNotifier(os.Stderr)
//	p, err := pipeline.New(pipeline.Config{
//		WindowSize: 100,
//		Threshold:  2.5,
//		Notifier:   notifier,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := p.Run(os.Stdin); err != nil {
//		log.Fatal(err)
//	}
//
// Each log line is parsed for a latency value. Lines without a latency field
// are silently skipped. Once the rolling window has enough data, any line
// whose latency deviates from the mean by more than Threshold standard
// deviations triggers an alert via the Notifier.
package pipeline
