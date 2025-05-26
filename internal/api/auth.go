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

// SignupRequest represents the data needed to create a new user
type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the data needed to login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt int64       `json:"expires_at"`
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create user
	userID, err := h.storage.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Error("failed to create user",
			zap.Error(err),
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the created user
	user, err := h.storage.GetUser(userID)
	if err != nil {
		log.Error("failed to get created user",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token, err := auth.GenerateToken(userID, h.secretKey)
	if err != nil {
		log.Error("failed to generate token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		log.Error("failed to generate session token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.storage.CreateSession(userID, sessionToken, expiresAt); err != nil {
		log.Error("failed to create session",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(expiresAt, 0),
	})

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	log.Info("user signed up successfully",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	// Get user details
	user, err := h.storage.GetUser(userID)
	if err != nil {
		log.Error("failed to get user",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token, err := auth.GenerateToken(userID, h.secretKey)
	if err != nil {
		log.Error("failed to generate token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		log.Error("failed to generate session token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.storage.CreateSession(userID, sessionToken, expiresAt); err != nil {
		log.Error("failed to create session",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(expiresAt, 0),
	})

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	log.Info("user logged in successfully",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)
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
		log.Error("failed to delete session",
			zap.Error(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
	user, err := h.storage.GetUser(session.UserID)
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
	json.NewEncoder(w).Encode(user)
}
