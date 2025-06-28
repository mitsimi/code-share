package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/api/dto"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
)

type UserHandler struct {
	users     repository.UserRepository
	snippets  repository.SnippetRepository
	likes     repository.LikeRepository
	bookmarks repository.BookmarkRepository
	logger    *zap.Logger
}

func NewUserHandler(users repository.UserRepository,
	snippets repository.SnippetRepository,
	likes repository.LikeRepository,
	bookmarks repository.BookmarkRepository) *UserHandler {
	return &UserHandler{
		users:     users,
		snippets:  snippets,
		likes:     likes,
		bookmarks: bookmarks,
		logger:    logger.Log,
	}
}

// ===== Helper methods for common logic =====

// getUserByID is a helper method that retrieves a user by ID
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, userID string) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("queried user_id", userID),
	)

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		log.Warn("failed to get user",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	log.Info("retrieved user",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)

	api.WriteSuccess(w, http.StatusOK, "User retrieved successfully", dto.ToUserResponse(user))
}

// getUserSnippetsByID is a helper method that retrieves snippets by author ID
func (h *UserHandler) getUserSnippetsByID(w http.ResponseWriter, r *http.Request, authorID string, userID string) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("author_id", authorID),
		zap.String("user_id", userID),
	)

	snippets, err := h.snippets.GetAllByAuthor(r.Context(), authorID, userID)
	if err != nil {
		log.Warn("failed to get user snippets",
			zap.Error(err),
			zap.String("author_id", authorID),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusInternalServerError, "Failed to retrieve snippets")
		return
	}

	log.Info("retrieved user snippets",
		zap.Int("count", len(snippets)),
	)

	responses := make([]dto.SnippetResponse, len(snippets))
	for i, snippet := range snippets {
		responses[i] = dto.ToSnippetResponse(snippet)
	}

	api.WriteSuccess(w, http.StatusOK, "User snippets retrieved successfully", responses)
}

// getUserLikedSnippetsByID is a helper method that retrieves liked snippets by user ID
func (h *UserHandler) getUserLikedSnippetsByID(w http.ResponseWriter, r *http.Request, userID string) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
	)

	snippets, err := h.likes.GetLikedSnippets(r.Context(), userID)
	if err != nil {
		log.Warn("failed to get user liked snippets",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusInternalServerError, "Failed to retrieve liked snippets")
		return
	}

	responses := make([]dto.SnippetResponse, len(snippets))
	for i, snippet := range snippets {
		responses[i] = dto.ToSnippetResponse(snippet)
	}

	log.Info("retrieved user liked snippets",
		zap.Int("count", len(snippets)),
	)

	api.WriteSuccess(w, http.StatusOK, "User liked snippets retrieved successfully", responses)
}

// getUserSavedSnippetsByID is a helper method that retrieves saved snippets by user ID
func (h *UserHandler) getUserSavedSnippetsByID(w http.ResponseWriter, r *http.Request, userID string) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
	)

	snippets, err := h.bookmarks.GetSavedSnippets(r.Context(), userID)
	if err != nil {
		log.Warn("failed to get user saved snippets",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		api.WriteError(w, http.StatusInternalServerError, "Failed to retrieve saved snippets")
		return
	}

	responses := make([]dto.SnippetResponse, len(snippets))
	for i, snippet := range snippets {
		responses[i] = dto.ToSnippetResponse(snippet)
	}

	log.Info("retrieved user saved snippets",
		zap.Int("count", len(snippets)),
	)

	api.WriteSuccess(w, http.StatusOK, "User saved snippets retrieved successfully", responses)
}

// updateUserProfileByID is a helper method that updates user profile by ID
func (h *UserHandler) updateUserProfileByID(w http.ResponseWriter, r *http.Request, userID string) {
	var req dto.UpdateUserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" {
		api.WriteError(w, http.StatusBadRequest, "Username and email cannot be empty")
		return
	}
	h.logger.Debug("updating user profile",
		zap.String("user_id", userID),
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	user := &domain.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}
	updatedUser, err := h.users.Update(r.Context(), user)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "Failed to update user info")
		return
	}

	api.WriteSuccess(w, http.StatusOK, "User info updated successfully", dto.ToUserResponse(updatedUser))
}

// updateUserPasswordByID is a helper method that updates user password by ID
func (h *UserHandler) updateUserPasswordByID(w http.ResponseWriter, r *http.Request, userID string) {
	var req dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify current password
	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Current password is incorrect")
		return
	}

	// Hash and update new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	if err := h.users.UpdatePassword(r.Context(), userID, string(hashedPassword)); err != nil {
		api.WriteError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	api.WriteSuccess(w, http.StatusOK, "Password updated successfully", nil)
}

// updateUserAvatarByID is a helper method that updates user avatar by ID
func (h *UserHandler) updateUserAvatarByID(w http.ResponseWriter, r *http.Request, userID string) {
	var req dto.UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AvatarURL == "" {
		api.WriteError(w, http.StatusBadRequest, "Avatar URL cannot be empty")
		return
	}

	h.logger.Debug("updating user avatar",
		zap.String("user_id", userID),
		zap.String("avatar_url", req.AvatarURL),
	)

	err := h.users.UpdateAvatar(r.Context(), userID, req.AvatarURL)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "Failed to update avatar")
		return
	}

	api.WriteSuccess(w, http.StatusOK, "Avatar updated successfully", map[string]string{
		"avatar": req.AvatarURL,
	})
}

// ===== Refactored route handlers =====

// GetUser returns a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	h.getUserByID(w, r, userID)
}

// GetMe returns the authenticated user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.getUserByID(w, r, userID)
}

// GetUserSnippets returns all snippets created by a user
func (h *UserHandler) GetUserSnippets(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "id")
	userID := api.GetUserID(r)
	h.getUserSnippetsByID(w, r, authorID, userID)
}

// GetMySnippets returns all snippets created by the authenticated user
func (h *UserHandler) GetMySnippets(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.getUserSnippetsByID(w, r, userID, userID)
}

// GetUserLikedSnippets returns all snippets liked by a user
func (h *UserHandler) GetUserLikedSnippets(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	h.getUserLikedSnippetsByID(w, r, userID)
}

// GetMyLikedSnippets returns all snippets liked by the authenticated user
func (h *UserHandler) GetMyLikedSnippets(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.getUserLikedSnippetsByID(w, r, userID)
}

// GetUserSavedSnippets returns all snippets saved by a user
func (h *UserHandler) GetUserSavedSnippets(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	h.getUserSavedSnippetsByID(w, r, userID)
}

// GetMySavedSnippets returns all snippets saved by the authenticated user
func (h *UserHandler) GetMySavedSnippets(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.getUserSavedSnippetsByID(w, r, userID)
}

// UpdateProfile handles updating a user's profile information
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserProfileByID(w, r, userID)
}

// UpdateMyProfile handles updating the authenticated user's profile information
func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserProfileByID(w, r, userID)
}

// UpdatePassword handles updating a user's password
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserPasswordByID(w, r, userID)
}

// UpdateMyPassword handles updating the authenticated user's password
func (h *UserHandler) UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserPasswordByID(w, r, userID)
}

// UpdateAvatar handles updating a user's avatar URL
func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserAvatarByID(w, r, userID)
}

// UpdateMyAvatar handles updating the authenticated user's avatar URL
func (h *UserHandler) UpdateMyAvatar(w http.ResponseWriter, r *http.Request) {
	userID := api.GetUserID(r)
	h.updateUserAvatarByID(w, r, userID)
}
