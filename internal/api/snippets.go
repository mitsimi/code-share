package api

import (
	"codeShare/internal/models"
	"codeShare/internal/storage"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SnippetHandler handles snippet-related HTTP requests
type SnippetHandler struct {
	storage storage.Storage
}

// NewSnippetHandler creates a new snippet handler
func NewSnippetHandler(storage storage.Storage) *SnippetHandler {
	return &SnippetHandler{storage: storage}
}

// GetSnippets returns all snippets
func (h *SnippetHandler) GetSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.storage.GetSnippets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snippets)
}

// GetSnippet returns a specific snippet
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(snippet)
}

// CreateSnippet creates a new snippet
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var req models.SnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snippet := models.Snippet{
		Title:    req.Title,
		Content:  req.Content,
		Author: req.Author,
	}

	if err := h.storage.CreateSnippet(snippet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(snippet)
}

// UpdateSnippet updates an existing snippet
func (h *SnippetHandler) UpdateSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.SnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	snippet.Title = req.Title
	snippet.Content = req.Content
	snippet.Author = req.Author

	if err := h.storage.UpdateSnippet(snippet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snippet)
}

// DeleteSnippet deletes a snippet
func (h *SnippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.storage.DeleteSnippet(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleLikeSnippet toggles the like status of a snippet
func (h *SnippetHandler) ToggleLikeSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.storage.ToggleLikeSnippet(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
} 