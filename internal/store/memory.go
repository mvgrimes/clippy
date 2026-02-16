package store

import (
	"container/list"
	"context"
	"fmt"
	"sync"
)

// Memory is an in-memory Store backed by a map protected with a RWMutex.
type Memory struct {
	mu          sync.RWMutex
	data        map[string][]byte
	order       *list.List            // insertion order for eviction (front = oldest)
	elements    map[string]*list.Element // key → list element
	totalSize   int64
	maxPerClip  int64 // 0 means unlimited
	maxTotal    int64 // 0 means unlimited
}

// NewMemory returns a new in-memory store.
func NewMemory(opts ...MemoryOption) *Memory {
	m := &Memory{
		data:     make(map[string][]byte),
		order:    list.New(),
		elements: make(map[string]*list.Element),
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// MemoryOption configures a Memory store.
type MemoryOption func(*Memory)

// WithMaxPerClip sets the maximum size of a single clip in bytes.
func WithMaxPerClip(n int64) MemoryOption {
	return func(m *Memory) { m.maxPerClip = n }
}

// WithMaxTotal sets the maximum total size of all clips in bytes.
func WithMaxTotal(n int64) MemoryOption {
	return func(m *Memory) { m.maxTotal = n }
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

	size := int64(len(value))

	if m.maxPerClip > 0 && size > m.maxPerClip {
		return fmt.Errorf("%w: %d bytes exceeds limit of %d", ErrTooLarge, size, m.maxPerClip)
	}

	// Account for replacing an existing key.
	oldSize := int64(0)
	if old, ok := m.data[key]; ok {
		oldSize = int64(len(old))
	}

	newTotal := m.totalSize - oldSize + size

	// Evict oldest clips until we fit within maxTotal.
	if m.maxTotal > 0 {
		for newTotal > m.maxTotal && m.order.Len() > 0 {
			front := m.order.Front()
			evictKey := front.Value.(string)
			// Don't evict the key we're about to write.
			if evictKey == key {
				m.order.MoveToBack(front)
				if m.order.Front().Value.(string) == key {
					// Only our key remains; can't evict further.
					break
				}
				continue
			}
			newTotal -= int64(len(m.data[evictKey]))
			delete(m.data, evictKey)
			delete(m.elements, evictKey)
			m.order.Remove(front)
		}
	}

	buf := make([]byte, size)
	copy(buf, value)
	m.data[key] = buf

	// Update ordering: move to back if exists, otherwise push back.
	if el, ok := m.elements[key]; ok {
		m.order.MoveToBack(el)
	} else {
		el = m.order.PushBack(key)
		m.elements[key] = el
	}

	m.totalSize = newTotal
	return nil
}
