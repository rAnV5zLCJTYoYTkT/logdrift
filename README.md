# logdrift

A lightweight CLI tool to detect and alert on anomalous log patterns using rolling statistical baselines.

---

## Installation

```bash
go install github.com/yourusername/logdrift@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logdrift.git
cd logdrift
go build -o logdrift .
```

---

## Usage

Pipe logs directly into `logdrift` or point it at a file:

```bash
# Stream from a file
logdrift --input /var/log/app.log

# Pipe from journalctl
journalctl -f | logdrift --threshold 2.5

# Set a custom rolling window (in seconds)
logdrift --input app.log --window 60 --threshold 3.0
```

### Flags

| Flag | Default | Description |
|-------------|---------|--------------------------------------|
| `--input` | stdin | Log file path or use stdin |
| `--window` | `30` | Rolling baseline window in seconds |
| `--threshold` | `2.0` | Std deviation threshold for alerts |
| `--format` | `text` | Output format: `text` or `json` |

### Example Output

```
[ALERT] 2024-11-03T14:22:01Z — anomaly detected: error rate 4.8σ above baseline
[INFO]  2024-11-03T14:22:05Z — baseline stable (window: 30s, samples: 142)
```

---

## How It Works

`logdrift` maintains a rolling statistical baseline of log event rates and pattern frequencies. When incoming log activity deviates beyond a configurable standard deviation threshold, it emits an alert — making it easy to catch error spikes, traffic anomalies, or unexpected silence in your logs.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

MIT © 2024 yourusername