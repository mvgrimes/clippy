package server

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/mvgrimes/clippy/internal/store"
)

//go:embed logo.svg
var logoSVG []byte

//go:embed index.html
var indexHTML string

var indexTmpl = template.Must(template.New("index").Parse(indexHTML))

// Handlers holds HTTP handler methods and their dependencies.
type Handlers struct {
	Store store.Store
}

func (h *Handlers) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	indexTmpl.Execute(w, struct{ Host string }{Host: r.Host})
}

func (h *Handlers) handleLogo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(logoSVG)
}

func (h *Handlers) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	data, err := h.Store.Get(r.Context(), key)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}

func (h *Handlers) handleSet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	if err := h.Store.Set(r.Context(), key, body); err != nil {
		if errors.Is(err, store.ErrTooLarge) {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}
