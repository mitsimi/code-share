package api

import (
	"encoding/json"
	"net/http"
	"slices"

	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// SnippetHandler handles snippet-related HTTP requests
type SnippetHandler struct {
	storage storage.Storage
	logger  *zap.Logger
}

// NewSnippetHandler creates a new snippet handler
func NewSnippetHandler(storage storage.Storage) *SnippetHandler {
	return &SnippetHandler{
		storage: storage,
		logger:  logger.Log,
	}
}

// GetSnippets returns all snippets
func (h *SnippetHandler) GetSnippets(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	snippets, err := h.storage.GetSnippets()
	if err != nil {
		log.Error("failed to get snippets",
			zap.Error(err),
		)
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

	log.Debug("retrieved snippets",
		zap.Int("count", len(snippets)),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSnippet returns a specific snippet
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
	)

	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		log.Error("failed to get snippet",
			zap.Error(err),
		)
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

	log.Debug("retrieved snippet",
		zap.String("title", snippet.Title),
		zap.String("author", snippet.Author),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSnippet creates a new snippet
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req models.SnippetCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := GetUserID(r)
	if userID == "" {
		log.Error("no user ID in context")
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	snippet := models.Snippet{
		Title:   req.Title,
		Content: req.Content,
		Author:  userID,
		Likes:   0,
		IsLiked: false,
	}

	log.Debug("creating new snippet",
		zap.String("title", snippet.Title),
		zap.String("content", snippet.Content),
		zap.String("author", snippet.Author),
	)

	id, err := h.storage.CreateSnippet(snippet)
	if err != nil {
		log.Error("failed to create snippet",
			zap.Error(err),
			zap.String("title", snippet.Title),
			zap.String("userId", snippet.Author),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s, _ := h.storage.GetSnippet(id)
	log.Info("created new snippet",
		zap.String("id", id),
		zap.String("title", s.Title),
		zap.String("author", s.Author),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// UpdateSnippet updates an existing snippet
func (h *SnippetHandler) UpdateSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
	)

	var req models.SnippetCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		log.Error("failed to get snippet for update",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	snippet.Title = req.Title
	snippet.Content = req.Content

	if err := h.storage.UpdateSnippet(snippet); err != nil {
		log.Error("failed to update snippet",
			zap.Error(err),
			zap.String("title", snippet.Title),
			zap.String("author", snippet.Author),
		)
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

	log.Info("updated snippet",
		zap.String("title", snippet.Title),
		zap.String("author", snippet.Author),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSnippet deletes a snippet
func (h *SnippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
	)

	if err := h.storage.DeleteSnippet(id); err != nil {
		log.Error("failed to delete snippet",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Info("deleted snippet")
	w.WriteHeader(http.StatusNoContent)
}

// ToggleLikeSnippet toggles the like status of a snippet
func (h *SnippetHandler) ToggleLikeSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
	)

	// Parse the action from query parameters
	action := r.URL.Query().Get("action")
	if action != "like" && action != "unlike" {
		log.Error("invalid action",
			zap.String("action", action),
		)
		http.Error(w, "Invalid action. Must be 'like' or 'unlike'", http.StatusBadRequest)
		return
	}

	if err := h.storage.ToggleLikeSnippet(id, action == "like"); err != nil {
		log.Error("failed to toggle like",
			zap.Error(err),
			zap.String("action", action),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get the updated snippet
	snippet, err := h.storage.GetSnippet(id)
	if err != nil {
		log.Error("failed to get updated snippet",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("toggled snippet like",
		zap.String("action", action),
		zap.Int("likes", int(snippet.Likes)),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}
