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
