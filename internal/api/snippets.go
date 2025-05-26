package api

import (
	"codeShare/internal/models"
	"codeShare/internal/storage"
	"encoding/json"
	"net/http"
	"slices"

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

	// Add isLiked field to each snippet
	response := make([]struct {
		models.Snippet
		IsLiked bool `json:"is_liked"`
	}, len(snippets))

	for i, snippet := range snippets {
		response[i] = struct {
			models.Snippet
			IsLiked bool `json:"is_liked"`
		}{
			Snippet: snippet,
			IsLiked: snippet.IsLiked,
		}
	}

	slices.SortFunc(response, func(a, b struct {
		models.Snippet
		IsLiked bool `json:"is_liked"`
	}) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	json.NewEncoder(w).Encode(response)
}

// GetSnippet returns a specific snippet
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add isLiked field to the response
	response := struct {
		models.Snippet
		IsLiked bool `json:"is_liked"`
	}{
		Snippet: snippet,
		IsLiked: snippet.IsLiked,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateSnippet creates a new snippet
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var req models.SnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snippet := models.Snippet{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
		Likes:   0,
		IsLiked: false,
	}

	id, err := h.storage.CreateSnippet(snippet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s, _ := h.storage.GetSnippet(id)
	json.NewEncoder(w).Encode(s)
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

	response := struct {
		models.Snippet
		IsLiked bool `json:"is_liked"`
	}{
		Snippet: snippet,
		IsLiked: snippet.IsLiked,
	}

	json.NewEncoder(w).Encode(response)
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

	// Parse the action from query parameters
	action := r.URL.Query().Get("action")
	if action != "like" && action != "unlike" {
		http.Error(w, "Invalid action. Must be 'like' or 'unlike'", http.StatusBadRequest)
		return
	}

	if err := h.storage.ToggleLikeSnippet(id, action == "like"); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get the updated snippet
	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snippet)
}
