package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/api/dto"
	"mitsimi.dev/codeShare/internal/services"
	ws "mitsimi.dev/codeShare/internal/websocket"

	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SnippetHandler handles snippet-related HTTP requests
type SnippetHandler struct {
	snippets    repository.SnippetRepository
	likes       repository.LikeRepository
	bookmarks   repository.BookmarkRepository
	viewTracker *services.ViewTracker
	wsHub       *ws.Hub
	logger      *zap.Logger
}

// NewSnippetHandler creates a new snippet handler
func NewSnippetHandler(
	snippets repository.SnippetRepository,
	likes repository.LikeRepository,
	bookmarks repository.BookmarkRepository,
	viewTracker *services.ViewTracker,
	wsHub *ws.Hub,
) *SnippetHandler {
	return &SnippetHandler{
		snippets:    snippets,
		likes:       likes,
		bookmarks:   bookmarks,
		viewTracker: viewTracker,
		wsHub:       wsHub,
		logger:      logger.Log,
	}
}

// GetSnippets returns all snippets
func (h *SnippetHandler) GetSnippets(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(zap.String("request_id", requestID), zap.String("user_id", userID))

	snippets, err := h.snippets.GetAll(r.Context(), userID)
	if err != nil {
		log.Error("failed to get snippets",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.SnippetResponse, len(snippets))
	for i, snippet := range snippets {
		responses[i] = dto.ToSnippetResponse(snippet)
	}

	slices.SortFunc(responses, func(a, b dto.SnippetResponse) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	log.Debug("retrieved snippets",
		zap.Int("count", len(responses)),
	)

	api.WriteSuccess(w, http.StatusOK, "Snippets retrieved successfully", responses)
}

// GetSnippet returns a specific snippet
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
		zap.String("user_id", userID),
	)

	snippet, err := h.snippets.GetByID(r.Context(), id, userID)
	if err != nil {
		log.Error("failed to get snippet",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Track view asynchronously to avoid blocking the response
	go func() {
		// Create a separate context with timeout for view tracking
		// Don't use r.Context() as it gets cancelled when the response is sent
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.viewTracker.TrackView(ctx, r, id, userID); err != nil {
			log.Error("failed to track view",
				zap.Error(err),
				zap.String("snippet_id", id),
				zap.String("user_id", userID),
			)
		} else {
			// Broadcast view count update after successfully tracking the view
			if h.wsHub != nil {
				// Get the updated snippet to get the new view count
				if updatedSnippet, err := h.snippets.GetByID(ctx, id, userID); err == nil {
					h.wsHub.BroadcastSnippetStatsUpdate(id, &updatedSnippet.Views, &updatedSnippet.Likes)
				}
			}
		}
	}()

	response := dto.ToSnippetResponse(snippet)

	log.Debug("retrieved snippet",
		zap.String("id", response.ID),
		zap.String("title", response.Title),
		zap.String("author", response.Author.ID),
	)

	api.WriteSuccess(w, http.StatusOK, "Snippet retrieved successfully", response)
}

// CreateSnippet creates a new snippet
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(zap.String("request_id", requestID), zap.String("user_id", userID))

	var req dto.CreateSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Title == "" || req.Content == "" {
		log.Error("invalid snippet data",
			zap.String("title", req.Title),
			zap.String("content", req.Content),
		)
		api.WriteError(w, http.StatusBadRequest, "Title and content cannot be empty")
		return
	}

	if req.Language == "" {
		req.Language = "plaintext" // Default language if not provided
	}

	if userID == "" {
		log.Error("no user ID in context")
		api.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Convert to domain model
	domainSnippet := dto.ToDomainSnippet(req, userID)
	domainSnippet.ID = uuid.New().String()
	domainSnippet.CreatedAt = time.Now()
	domainSnippet.UpdatedAt = time.Now()

	log.Debug("creating new snippet",
		zap.String("title", domainSnippet.Title),
		zap.String("content", domainSnippet.Content),
		zap.String("author", userID),
	)

	if err := h.snippets.Create(r.Context(), domainSnippet); err != nil {
		log.Error("failed to create snippet",
			zap.Error(err),
			zap.String("title", domainSnippet.Title),
			zap.String("userId", domainSnippet.Author.ID),
		)
		api.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s, _ := h.snippets.GetByID(r.Context(), domainSnippet.ID, userID)
	response := dto.ToSnippetResponse(s)

	log.Debug("created new snippet",
		zap.String("id", s.ID),
		zap.String("title", response.Title),
		zap.String("author", response.Author.ID),
	)

	w.Header().Set("Location", "/snippets/"+s.ID)
	api.WriteSuccess(w, http.StatusCreated, "Snippet created successfully", response)
}

// UpdateSnippet updates an existing snippet
func (h *SnippetHandler) UpdateSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
		zap.String("user_id", userID),
	)

	var req dto.UpdateSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	snippet, err := h.snippets.GetByID(r.Context(), id, userID)
	if err != nil {
		log.Error("failed to get snippet for update",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if snippet.Author.ID != userID {
		log.Error("unauthorized update attempt",
			zap.String("snippet_author", snippet.Author.ID),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusForbidden, "Only the author can update this snippet")
		return
	}

	// Update domain model
	dto.UpdateDomainSnippet(snippet, req)

	if err := h.snippets.Update(r.Context(), snippet); err != nil {
		log.Error("failed to update snippet",
			zap.Error(err),
			zap.String("title", snippet.Title),
			zap.String("author", snippet.Author.ID),
		)
		api.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert to response DTO
	response := dto.ToSnippetResponse(snippet)
	log.Debug("updated snippet",
		zap.String("title", response.Title),
		zap.String("author", response.Author.ID),
	)

	// Broadcast content update to both snippet detail and list subscribers
	if h.wsHub != nil {
		h.wsHub.BroadcastSnippetContentUpdate(
			snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Language,
		)
	}

	api.WriteSuccess(w, http.StatusOK, "Snippet updated successfully", response)
}

// DeleteSnippet deletes a snippet
func (h *SnippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", snippetID),
		zap.String("user_id", userID),
	)

	snippet, err := h.snippets.GetByID(r.Context(), snippetID, userID)
	if err != nil {
		log.Error("failed to get snippet",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if snippet.Author.ID != userID {
		log.Error("unauthorized deletion attempt",
			zap.String("snippet_author", snippet.Author.ID),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusForbidden, "Only the author can delete this snippet")
		return
	}

	if err := h.snippets.Delete(r.Context(), snippetID); err != nil {
		log.Error("failed to delete snippet",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	log.Debug("deleted snippet")
	api.WriteSuccess(w, http.StatusOK, "Snippet deleted successfully", nil)
}

// ToggleLikeSnippet toggles the like status of a snippet
func (h *SnippetHandler) ToggleLikeSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
		zap.String("user_id", userID),
	)

	// Parse the action from query parameters
	action := r.URL.Query().Get("action")
	if action == "" {
		action = "like"
	}

	if action != "like" && action != "unlike" {
		log.Error("invalid action",
			zap.String("action", action),
		)
		api.WriteError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	if err := h.likes.ToggleLike(r.Context(), userID, id, action == "like"); err != nil {
		log.Error("failed to toggle like",
			zap.Error(err),
			zap.String("action", action),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Get the updated snippet
	snippet, err := h.snippets.GetByID(r.Context(), id, userID)
	if err != nil {
		log.Error("failed to get updated snippet",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.wsHub != nil {
		h.wsHub.BroadcastUserAction(userID, ws.UserActionData{
			Action:    action,
			SnippetID: id,
			Value:     snippet.IsLiked,
		})
		h.wsHub.BroadcastSnippetStatsUpdate(id, &snippet.Views, &snippet.Likes)
	}

	log.Info("toggled snippet like",
		zap.String("action", action),
	)
	api.WriteSuccess(w, http.StatusOK, "Snippet like toggled successfully", dto.ToSnippetResponse(snippet))
}

// ToggleSaveSnippet toggles the save status of a snippet
func (h *SnippetHandler) ToggleSaveSnippet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("snippet_id", id),
		zap.String("user_id", userID),
	)

	// Parse the action from query parameters
	action := r.URL.Query().Get("action")
	if action == "" {
		action = "save"
	}

	if action != "save" && action != "unsave" {
		log.Error("invalid action",
			zap.String("action", action),
		)
		api.WriteError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	if err := h.bookmarks.ToggleSave(r.Context(), userID, id, action == "save"); err != nil {
		log.Error("failed to toggle save",
			zap.Error(err),
			zap.String("action", action),
		)
		api.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Get the updated snippet
	snippet, err := h.snippets.GetByID(r.Context(), id, userID)
	if err != nil {
		log.Error("failed to get updated snippet",
			zap.Error(err),
		)
		api.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.wsHub != nil {
		h.wsHub.BroadcastUserAction(userID, ws.UserActionData{
			Action:    action,
			SnippetID: id,
			Value:     snippet.IsSaved,
		})
		h.wsHub.BroadcastSnippetStatsUpdate(id, &snippet.Views, &snippet.Likes)
	}

	log.Debug("toggled snippet save",
		zap.String("action", action),
	)
	api.WriteSuccess(w, http.StatusOK, "Snippet save toggled successfully", dto.ToSnippetResponse(snippet))
}
