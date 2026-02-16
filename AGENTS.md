# Agent Guidelines for clippy

## Project Overview
`clippy` is a simple pastebin service with a CLI client, written in Go 1.24 with Cobra.

## Architecture
- `cmd/` — Cobra CLI commands (root, server, paste, get)
- `internal/store/` — Pluggable storage backend interface + in-memory implementation
- `internal/server/` — HTTP handlers and route registration

## Key Patterns
- The `Store` interface (`internal/store/store.go`) uses `context.Context` on all methods
- `store.ErrNotFound` sentinel error distinguishes 404 from 500 in handlers
- Default paste uses empty string key `""`
- HTTP handlers read `r.Body` directly (no JSON parsing)
- Go 1.22+ method-based `ServeMux` patterns (e.g., `GET /@/{key...}`)

## Adding a Storage Backend
1. Create `internal/store/<name>.go` implementing `store.Store`
2. Add constructor and any config needed
3. Add `--store` flag to `cmd/server.go` and wire up the new backend

## Testing
```bash
go test ./...
```

## Conventions
- Keep handlers thin — business logic belongs in the store layer
- Always return `store.ErrNotFound` for missing keys (not nil data)
- Copy byte slices in the memory store to prevent data races
