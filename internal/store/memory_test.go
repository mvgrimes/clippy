package store

import (
	"context"
	"errors"
	"testing"
)

func TestMemory_GetNotFound(t *testing.T) {
	s := NewMemory()
	_, err := s.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_SetAndGet(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()

	if err := s.Set(ctx, "k1", []byte("hello")); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("expected %q, got %q", "hello", string(got))
	}
}

func TestMemory_Overwrite(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()

	_ = s.Set(ctx, "k", []byte("first"))
	_ = s.Set(ctx, "k", []byte("second"))

	got, _ := s.Get(ctx, "k")
	if string(got) != "second" {
		t.Fatalf("expected %q, got %q", "second", string(got))
	}
}

func TestMemory_DefaultKey(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()

	_ = s.Set(ctx, "", []byte("default"))
	got, err := s.Get(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "default" {
		t.Fatalf("expected %q, got %q", "default", string(got))
	}
}

func TestMemory_IsolatedKeys(t *testing.T) {
	s := NewMemory()
	ctx := context.Background()

	_ = s.Set(ctx, "a", []byte("alpha"))
	_ = s.Set(ctx, "b", []byte("bravo"))

	got, _ := s.Get(ctx, "a")
	if string(got) != "alpha" {
		t.Fatalf("expected %q, got %q", "alpha", string(got))
	}
	got, _ = s.Get(ctx, "b")
	if string(got) != "bravo" {
		t.Fatalf("expected %q, got %q", "bravo", string(got))
	}
}

func TestMemory_MaxPerClip(t *testing.T) {
	s := NewMemory(WithMaxPerClip(10))
	ctx := context.Background()

	// Within limit: OK
	if err := s.Set(ctx, "k", []byte("short")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Exceeds limit: ErrTooLarge
	err := s.Set(ctx, "k", []byte("this is way too long"))
	if !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}

	// Original value should remain
	got, _ := s.Get(ctx, "k")
	if string(got) != "short" {
		t.Fatalf("expected %q, got %q", "short", string(got))
	}
}

func TestMemory_MaxTotalEviction(t *testing.T) {
	s := NewMemory(WithMaxTotal(15))
	ctx := context.Background()

	_ = s.Set(ctx, "a", []byte("aaaaa"))     // 5 bytes, total=5
	_ = s.Set(ctx, "b", []byte("bbbbb"))     // 5 bytes, total=10
	_ = s.Set(ctx, "c", []byte("ccccc"))     // 5 bytes, total=15
	_ = s.Set(ctx, "d", []byte("ddddd"))     // 5 bytes, would be 20 → evict "a"

	// "a" should be evicted
	_, err := s.Get(ctx, "a")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for evicted key, got %v", err)
	}

	// "b", "c", "d" should exist
	for _, key := range []string{"b", "c", "d"} {
		if _, err := s.Get(ctx, key); err != nil {
			t.Fatalf("key %q should exist, got %v", key, err)
		}
	}
}

func TestMemory_MaxTotalEvictsMultiple(t *testing.T) {
	s := NewMemory(WithMaxTotal(10))
	ctx := context.Background()

	_ = s.Set(ctx, "a", []byte("aaa"))    // 3 bytes
	_ = s.Set(ctx, "b", []byte("bbb"))    // 3 bytes, total=6
	_ = s.Set(ctx, "c", []byte("cccccccc")) // 8 bytes, needs to evict a+b

	// Both "a" and "b" should be evicted
	for _, key := range []string{"a", "b"} {
		_, err := s.Get(ctx, key)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("key %q should be evicted, got %v", key, err)
		}
	}

	got, _ := s.Get(ctx, "c")
	if string(got) != "cccccccc" {
		t.Fatalf("expected %q, got %q", "cccccccc", string(got))
	}
}

func TestMemory_OverwriteAdjustsTotalSize(t *testing.T) {
	s := NewMemory(WithMaxTotal(10))
	ctx := context.Background()

	_ = s.Set(ctx, "a", []byte("aaaaaaaaaa")) // 10 bytes, total=10
	_ = s.Set(ctx, "a", []byte("bb"))          // 2 bytes, total=2 (replaced)

	// Should have room for more
	if err := s.Set(ctx, "b", []byte("cccccccc")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
