package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/clip/internal/store"
)

// Handlers holds HTTP handler methods and their dependencies.
type Handlers struct {
	Store store.Store
}

func (h *Handlers) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>clip</title></head>
<body>
<h1>clip — simple pastebin</h1>
<pre>
# Store a paste (default)
curl -d 'hello world' http://`+r.Host+`/@

# Retrieve the default paste
curl http://`+r.Host+`/@

# Store a named paste
curl -d 'my content' http://`+r.Host+`/@/mykey

# Retrieve a named paste
curl http://`+r.Host+`/@/mykey
</pre>
</body>
</html>`)
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
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}
