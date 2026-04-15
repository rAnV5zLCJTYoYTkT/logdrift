// Package config handles loading, parsing, and validating logdrift's YAML
// configuration file.
//
// # File format
//
// Configuration is expressed in YAML. A minimal valid file looks like:
//
//	watch:
//	  file: /var/log/app.log
//	baseline:
//	  window_size: 60
//	  threshold: 3.0
//
// # Defaults
//
// Fields that are omitted receive sensible defaults:
//   - watch.poll_interval  → 500ms
//   - alert.cooldown       → 10s
//   - report.output_path   → "-" (stdout)
package config
