package store

import (
	"context"
	"sync"
)

// Memory is an in-memory Store backed by a map protected with a RWMutex.
type Memory struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewMemory returns a new in-memory store.
func NewMemory() *Memory {
	return &Memory{data: make(map[string][]byte)}
}

func (m *Memory) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	// Return a copy to avoid data races on the slice.
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (m *Memory) Set(_ context.Context, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	buf := make([]byte, len(value))
	copy(buf, value)
	m.data[key] = buf
	return nil
}
