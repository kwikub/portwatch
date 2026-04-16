# portwatch

A lightweight CLI daemon that monitors and logs port activity changes on a host.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Monitor specific ports and log to a file:

```bash
portwatch start --ports 80,443,8080 --interval 5s --log /var/log/portwatch.log
```

Run a one-time snapshot of active ports:

```bash
portwatch scan
```

### Example Output

```
2024/01/15 10:23:01 [OPEN]   port 8080 (tcp) — process: nginx (pid 1234)
2024/01/15 10:23:11 [CLOSED] port 3000 (tcp) — process: node (pid 5678)
2024/01/15 10:23:21 [OPEN]   port 5432 (tcp) — process: postgres (pid 91011)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `10s` | Polling interval |
| `--ports` | all | Comma-separated list of ports to watch |
| `--log` | stdout | Path to log file |
| `--format` | `text` | Output format (`text` or `json`) |

## License

MIT © 2024 [Your Name](https://github.com/yourusername)