# CLI Reference

CasPaste provides two command-line tools:

- **caspaste** - The server binary
- **caspaste-cli** - The client for interacting with CasPaste servers

## Server CLI (caspaste)

### Basic Usage

```bash
caspaste [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--help` | Show help | - |
| `--version` | Show version | - |
| `--config PATH` | Config directory | Platform-specific |
| `--data PATH` | Data directory | Platform-specific |
| `--cache PATH` | Cache directory | Platform-specific |
| `--log PATH` | Log directory | Platform-specific |
| `--backup PATH` | Backup directory | Platform-specific |
| `--pid PATH` | PID file path | Platform-specific |
| `--address ADDR` | Listen address | `0.0.0.0` |
| `--port PORT` | Listen port | Auto-detect |
| `--mode MODE` | Application mode | `production` |
| `--status` | Show running status | - |
| `--daemon` | Daemonize (detach) | - |
| `--debug` | Enable debug mode | - |

### Service Management

```bash
caspaste --service {install|uninstall|start|stop|restart|status}
```

| Command | Description |
|---------|-------------|
| `install` | Install as system service |
| `uninstall` | Remove system service |
| `start` | Start the service |
| `stop` | Stop the service |
| `restart` | Restart the service |
| `status` | Show service status |

### Maintenance Operations

```bash
caspaste --maintenance {backup|restore|cleanup|reset-admin}
```

| Command | Description |
|---------|-------------|
| `backup` | Create a backup |
| `restore` | Restore from latest backup |
| `restore FILE` | Restore from specific file |
| `cleanup` | Run cleanup operations |
| `reset-admin` | Reset admin credentials |

### Examples

```bash
# Run with custom directories
caspaste --data /data/caspaste --config /etc/caspaste --port 8080

# Run in debug mode
caspaste --debug

# Check if server is healthy
caspaste --status
echo $?  # 0=healthy, 1=unhealthy, 2=degraded

# Install and start as service
sudo caspaste --service install
sudo caspaste --service start
```

## Client CLI (caspaste-cli)

### Configuration

First-time setup:

```bash
caspaste-cli login
```

This creates a config file at:

- **Linux:** `~/.config/casjay-forks/caspaste/cli.yml`
- **macOS:** `~/Library/Application Support/CasPaste/cli.yml`
- **Windows:** `%LOCALAPPDATA%\CasPaste\cli.yml`

### Create Paste

```bash
# From stdin
echo "Hello World" | caspaste-cli new

# From file
caspaste-cli new -f script.py

# With options
caspaste-cli new -f code.go -s go -t "My Code" --expires 1d
```

| Flag | Description |
|------|-------------|
| `-f, --file FILE` | Upload file |
| `-s, --syntax LANG` | Syntax highlighting |
| `-t, --title TITLE` | Paste title |
| `--expires DURATION` | Expiration time |
| `--burn` | Burn after reading |
| `--password PASS` | Password protection |

### Get Paste

```bash
# View paste
caspaste-cli get abc123

# Download to file
caspaste-cli get abc123 -o output.txt

# Get raw content
caspaste-cli get abc123 --raw
```

### List Pastes

```bash
# List recent pastes
caspaste-cli list

# With pagination
caspaste-cli list --limit 50 --offset 100
```

### Shorten URL

```bash
caspaste-cli shorten https://example.com/very/long/url
```

### Configuration File

```yaml
# ~/.config/casjay-forks/caspaste/cli.yml
server: https://paste.example.com
token: your-api-token
default_syntax: plaintext
default_expires: never
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Connection error |
| 3 | Authentication error |
| 4 | Not found |
