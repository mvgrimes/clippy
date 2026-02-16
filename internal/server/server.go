package server

import (
	"net/http"

	"github.com/mvgrimes/clippy/internal/store"
)

// New creates an http.Handler with all routes registered.
func New(s store.Store) http.Handler {
	h := &Handlers{Store: s}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", h.handleIndex)
	mux.HandleFunc("GET /logo.svg", h.handleLogo)

	mux.HandleFunc("GET /@", h.handleGet)
	mux.HandleFunc("GET /@/{key...}", h.handleGet)

	mux.HandleFunc("POST /@", h.handleSet)
	mux.HandleFunc("PUT /@", h.handleSet)
	mux.HandleFunc("POST /@/{key...}", h.handleSet)
	mux.HandleFunc("PUT /@/{key...}", h.handleSet)

	return mux
}
