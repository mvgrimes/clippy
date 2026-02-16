package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mvgrimes/clippy/internal/store"
)

func setup() (*httptest.Server, store.Store) {
	s := store.NewMemory()
	handler := New(s)
	ts := httptest.NewServer(handler)
	return ts, s
}

func TestIndex(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html, got %s", ct)
	}
}

func TestSetAndGetDefault(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	// POST default paste
	resp, err := http.Post(ts.URL+"/@", "text/plain", strings.NewReader("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("POST expected 200, got %d", resp.StatusCode)
	}

	// GET default paste
	resp, err = http.Get(ts.URL + "/@")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("GET expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", string(body))
	}
}

func TestSetAndGetNamed(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/@/mykey", "text/plain", strings.NewReader("named paste"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	resp, err = http.Get(ts.URL + "/@/mykey")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "named paste" {
		t.Fatalf("expected %q, got %q", "named paste", string(body))
	}
}

func TestGetNotFound(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/@/nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestEmptyBody(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/@", "text/plain", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPutMethod(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/@/putkey", strings.NewReader("via put"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("PUT expected 200, got %d", resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/@/putkey")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "via put" {
		t.Fatalf("expected %q, got %q", "via put", string(body))
	}
}
