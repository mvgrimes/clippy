# clip

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

# Named paste
curl -d 'my content' http://localhost:8080/@/mykey
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

| Flag     | Default     | Description       |
|----------|-------------|-------------------|
| `--host` | `0.0.0.0`  | Host to bind to   |
| `--port` | `8080`      | Port to listen on |

## Client flags

| Flag       | Default                  | Description |
|------------|--------------------------|-------------|
| `--server` | `http://localhost:8080`  | Server URL  |

## Development

```bash
go test ./...
go vet ./...
```
