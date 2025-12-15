package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type noteStore struct {
	mu    sync.RWMutex
	notes map[string]Note
}

func newNoteStore() *noteStore {
	return &noteStore{notes: map[string]Note{}}
}

func (s *noteStore) list() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Note, 0, len(s.notes))
	for _, n := range s.notes {
		out = append(out, n)
	}
	return out
}

func (s *noteStore) get(id string) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n, ok := s.notes[id]
	return n, ok
}

func (s *noteStore) create(title, content string) Note {
	n := Note{
		ID:        newID(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	s.notes[n.ID] = n
	s.mu.Unlock()

	return n
}

func (s *noteStore) update(id, title, content string) (Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, ok := s.notes[id]
	if !ok {
		return Note{}, false
	}
	if title != "" {
		n.Title = title
	}
	if content != "" {
		n.Content = content
	}
	n.UpdatedAt = time.Now().UTC()
	s.notes[id] = n
	return n, true
}

func (s *noteStore) delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[id]; !ok {
		return false
	}
	delete(s.notes, id)
	return true
}

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func main() {
	store := newNoteStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.list())
		case http.MethodPost:
			var req createNoteRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
				return
			}
			n := store.create(strings.TrimSpace(req.Title), req.Content)
			writeJSON(w, http.StatusCreated, n)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/notes/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/notes/")
		if id == "" || strings.Contains(id, "/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			n, ok := store.get(id)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, n)
		case http.MethodPut:
			var req updateNoteRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
				return
			}
			var title, content string
			if req.Title != nil {
				title = strings.TrimSpace(*req.Title)
			}
			if req.Content != nil {
				content = *req.Content
			}
			n, ok := store.update(id, title, content)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, n)
		case http.MethodDelete:
			if ok := store.delete(id); !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	h := corsMiddleware(mux)

	addr := ":8080"
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dev-friendly default for Vite.
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
