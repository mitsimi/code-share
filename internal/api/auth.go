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
		log.Error("failed to decode request body",
			zap.Error(err),
		)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate access token
	accessTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, false)
	if err != nil {
		log.Error("failed to generate access token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, true)
	if err != nil {
		log.Error("failed to generate refresh token",
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

	if err := h.storage.CreateSession(user.ID, sessionToken, refreshTokenResp.Token, refreshTokenResp.ExpiresAt); err != nil {
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
		Expires:  time.Unix(refreshTokenResp.ExpiresAt, 0),
	})

	response := models.AuthResponse{
		Token:        accessTokenResp.Token,
		RefreshToken: refreshTokenResp.Token,
		User:         models.FromDBUser(user),
		ExpiresAt:    accessTokenResp.ExpiresAt,
	}

	log.Info("user signed up successfully",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
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

	// Generate access token
	accessTokenResp, err := auth.GenerateToken(userID, h.secretKey, false)
	if err != nil {
		log.Error("failed to generate access token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshTokenResp, err := auth.GenerateToken(userID, h.secretKey, true)
	if err != nil {
		log.Error("failed to generate refresh token",
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

	if err := h.storage.CreateSession(userID, sessionToken, refreshTokenResp.Token, refreshTokenResp.ExpiresAt); err != nil {
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
		Expires:  time.Unix(refreshTokenResp.ExpiresAt, 0),
	})

	response := models.AuthResponse{
		Token:        accessTokenResp.Token,
		RefreshToken: refreshTokenResp.Token,
		User:         models.FromDBUser(user),
		ExpiresAt:    accessTokenResp.ExpiresAt,
	}

	log.Info("user logged in successfully",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
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
		log.Error("failed to delete session",
			zap.Error(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FromDBUser(user))
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	// Get refresh token from request body
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body",
			zap.Error(err),
		)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Try to get session cookie first
	cookie, err := r.Cookie("session")
	if err == nil {
		// Session cookie exists, try to use it
		session, err := h.storage.GetSession(cookie.Value)
		if err == nil && session.ExpiresAt > time.Now().Unix() {
			// check if the provided refresh tokens matches the sessions one
			if req.RefreshToken != session.RefreshToken {
				log.Error("refresh token mismatch")
				http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
				return
			}

			// Valid session found, use it to refresh tokens
			user, err := h.storage.GetUser(session.UserID)
			if err != nil {
				log.Error("failed to get user from session",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Generate new tokens
			accessTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, false)
			if err != nil {
				log.Error("failed to generate access token",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			refreshTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, true)
			if err != nil {
				log.Error("failed to generate refresh token",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create new session
			sessionToken, err := auth.GenerateRandomToken()
			if err != nil {
				log.Error("failed to generate session token",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Delete old session and create new one
			if err := h.storage.DeleteSession(cookie.Value); err != nil {
				log.Error("failed to delete old session",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := h.storage.CreateSession(user.ID, sessionToken, refreshTokenResp.Token, refreshTokenResp.ExpiresAt); err != nil {
				log.Error("failed to create new session",
					zap.Error(err),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set new session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    sessionToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Unix(refreshTokenResp.ExpiresAt, 0),
			})

			response := models.AuthResponse{
				Token:        accessTokenResp.Token,
				RefreshToken: refreshTokenResp.Token,
				User:         models.FromDBUser(user),
				ExpiresAt:    accessTokenResp.ExpiresAt,
			}

			log.Info("token refreshed successfully via session",
				zap.String("username", user.Username),
				zap.String("email", user.Email),
			)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		// If we get here, either session not found or expired, fall through to JWT refresh
	}

	// Fallback: Try to refresh using JWT token
	claims, err := auth.ValidateToken(req.RefreshToken, h.secretKey)
	if err != nil {
		log.Error("invalid refresh token",
			zap.Error(err),
		)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if !claims.IsRefresh() {
		log.Error("token is not a refresh token")
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	// Get user from JWT claims
	user, err := h.storage.GetUser(claims.UserID)
	if err != nil {
		log.Error("failed to get user from JWT",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate new tokens
	accessTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, false)
	if err != nil {
		log.Error("failed to generate access token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokenResp, err := auth.GenerateToken(user.ID, h.secretKey, true)
	if err != nil {
		log.Error("failed to generate refresh token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create new session
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		log.Error("failed to generate session token",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.storage.CreateSession(user.ID, sessionToken, refreshTokenResp.Token, refreshTokenResp.ExpiresAt); err != nil {
		log.Error("failed to create new session",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set new session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(refreshTokenResp.ExpiresAt, 0),
	})

	response := models.AuthResponse{
		Token:        accessTokenResp.Token,
		RefreshToken: refreshTokenResp.Token,
		User:         models.FromDBUser(user),
		ExpiresAt:    accessTokenResp.ExpiresAt,
	}

	log.Info("token refreshed successfully via JWT",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
