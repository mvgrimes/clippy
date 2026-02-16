<p align="center">
  <img src="logo.svg" width="120" alt="clippy logo">
</p>

# clippy

A simple pastebin service and CLI client, written in Go.

## Quick Start

### Start the server

```bash
go run . server
```

### Store and retrieve pastes

```bash
# Default paste
curl -d 'hello world' http://localhost:8080/@
curl http://localhost:8080/@
echo "hi" | curl -d @- http://localhost:8080/@

# Named paste
curl -d 'my content' http://localhost:8080/@/mykey
echo "hi" | curl -d @- http://localhost:8080/@/mykey
curl http://localhost:8080/@/mykey
```

### CLI client

```bash
# Store from stdin
echo 'hello' | go run . paste
echo 'named' | go run . paste mykey

# Retrieve
go run . get
go run . get mykey
```

## Server flags

| Flag             | Default    | Description                                    |
|------------------|------------|------------------------------------------------|
| `--host`         | `0.0.0.0`  | Host to bind to                               |
| `--port`         | `8080`     | Port to listen on                              |
| `--max-per-clip` |            | Max size per clip (e.g. `1M`, `512K`)          |
| `--max-clips`    |            | Max total size of all clips (e.g. `100M`, `1G`)|

Size values accept optional units: `B`, `K`/`KB`, `M`/`MB`, `G`/`GB` (case-insensitive). Plain numbers are treated as bytes.

When `--max-per-clip` is exceeded the server responds with `413 Request Entity Too Large`. When `--max-clips` is exceeded the oldest clip is evicted to make room.

## Client flags

| Flag       | Default                  | Description |
|------------|--------------------------|-------------|
| `--server` | `http://localhost:8080`  | Server URL  |

## Environment variables

All flags can be set via environment variables. Flags given on the command line take precedence.

| Variable           | Flag             | Command  |
|--------------------|------------------|----------|
| `CLIP_HOST`        | `--host`         | `server` |
| `CLIP_PORT`        | `--port`         | `server` |
| `CLIP_MAX_PER_CLIP`| `--max-per-clip` | `server` |
| `CLIP_MAX_CLIPS`   | `--max-clips`    | `server` |
| `CLIP_SERVER`      | `--server`       | `paste`, `get` |

## Development

```bash
go test ./...
go vet ./...
```
