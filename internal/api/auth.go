package api

import (
	"encoding/json"
	"net/http"
	"time"

	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	storage   storage.Storage
	logger    *zap.Logger
	secretKey string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(storage storage.Storage, secretKey string) *AuthHandler {
	return &AuthHandler{
		storage:   storage,
		logger:    logger.Log,
		secretKey: secretKey,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate email
	if !isValidEmail(req.Email) {
		log.Error("invalid email format", zap.String("email", req.Email))
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
		log.Error("invalid password", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.storage.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Error("failed to create user",
			zap.Error(err),
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	// Create tokens and session
	response, sessionToken, err := h.createTokensAndSession(user.ID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Info("user signed up successfully",
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Authenticate user
	userID, err := h.storage.Login(req.Email, req.Password)
	if err != nil {
		log.Error("failed to login",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create tokens and session
	response, sessionToken, err := h.createTokensAndSession(userID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Info("user logged in successfully",
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

	// Delete session
	if err := h.storage.DeleteSession(cookie.Value); err != nil {
		log.Warn("failed to delete session from storage",
			zap.Error(err),
			zap.String("session_token", cookie.Value[:8]+"..."), // Only log first 8 chars for security
		)
		// Continue anyway because if the deletion failed which means that there is no session
		// We delete the cookie so its not dangling around
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	log.Info("user logged out successfully")
	w.WriteHeader(http.StatusNoContent)
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	// Get refresh token from request body
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var userID string
	var logMessage string

	// Try session-based refresh first
	if cookie, err := r.Cookie("session"); err == nil {
		if session, err := h.storage.GetSession(cookie.Value); err == nil && session.ExpiresAt > time.Now().Unix() {
			// Validate refresh token matches session
			if req.RefreshToken != session.RefreshToken {
				log.Error("refresh token mismatch")
				http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
				return
			}

			// Delete old session before creating new one
			if err := h.storage.DeleteSession(cookie.Value); err != nil {
				log.Error("failed to delete old session", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			userID = session.UserID
			logMessage = "token refreshed successfully via session"
		}
	}

	// Fallback to JWT-based refresh
	if userID == "" {
		claims, err := auth.ValidateToken(req.RefreshToken, h.secretKey)
		if err != nil {
			log.Error("invalid refresh token", zap.Error(err))
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		if !claims.IsRefresh {
			log.Error("token is not a refresh token")
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		userID = claims.UserID
		logMessage = "token refreshed successfully via JWT"
	}

	// Create new tokens and session
	response, sessionToken, err := h.createTokensAndSession(userID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set new session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Info(logMessage,
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProfileRequest represents the request body for updating a user's profile
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdatePasswordRequest represents the request body for updating a user's password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// UpdateAvatarRequest represents the request body for updating a user's avatar
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatarUrl"`
}
