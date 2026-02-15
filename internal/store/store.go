package store

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a paste key does not exist.
var ErrNotFound = errors.New("paste not found")

// Store defines the interface for paste storage backends.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
}
