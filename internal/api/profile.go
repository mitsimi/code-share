package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mitsimi.dev/codeShare/internal/models"
)

// GetCurrentUser returns the current user's information
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	// Get session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Error("failed to get session cookie",
			zap.Error(err),
		)
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get session
	session, err := h.storage.GetSession(cookie.Value)
	if err != nil {
		log.Error("failed to get session",
			zap.Error(err),
		)
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Check if session is expired
	if session.ExpiresAt < time.Now().Unix() {
		log.Error("session expired",
			zap.Int64("expires_at", session.ExpiresAt),
		)
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	// Get user
	user, err := h.storage.GetUserByID(session.UserID)
	if err != nil {
		log.Error("failed to get user",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("retrieved current user",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FromDBUser(user))
}

// UpdatePassword handles updating a user's password
func (h *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify current password
	user, err := h.storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash and update new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	if err := h.storage.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateAvatar handles updating a user's avatar URL
func (h *AuthHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)

	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AvatarURL == "" {
		http.Error(w, "Avatar URL cannot be empty", http.StatusBadRequest)
		return
	}

	h.logger.Info("updating user avatar",
		zap.String("user_id", userID),
		zap.String("avatar_url", req.AvatarURL),
	)
	user, err := h.storage.UpdateUserAvatar(userID, req.AvatarURL)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	// For now, return the request data as the response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"avatar": user.Avatar,
	})
}

// UpdateProfile handles updating a user's profile information
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.storage.UpdateUser(userID, req.Username, req.Email)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// For now, return the request data as the response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username": req.Username,
		"email":    req.Email,
	})
}
