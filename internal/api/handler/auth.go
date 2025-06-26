package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/api/dto"
	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	users     repository.UserRepository
	sessions  repository.SessionRepository
	logger    *zap.Logger
	secretKey string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(users repository.UserRepository, sessions repository.SessionRepository, secretKey string) *AuthHandler {
	return &AuthHandler{
		users:     users,
		sessions:  sessions,
		logger:    logger.Log,
		secretKey: secretKey,
	}
}

// Signup handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req dto.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		api.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate email
	if !isValidEmail(req.Email) {
		log.Error("invalid email format", zap.String("email", req.Email))
		api.WriteError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
		log.Error("invalid password", zap.Error(err))
		api.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))
		api.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user, err := h.users.Create(r.Context(), &domain.UserCreation{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if repository.IsAlreadyExists(err) {
			log.Error("username or email already exists",
				zap.String("username", req.Username),
				zap.String("email", req.Email),
			)
			api.WriteError(w, http.StatusBadRequest, "Username or email already in use")
			return
		}
		log.Error("failed to create user",
			zap.Error(err),
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		api.WriteError(w, http.StatusBadRequest, "Failed to create user")
		return
	}

	// Create tokens and session
	response, sessionToken, err := h.createTokensAndSession(r.Context(), user.ID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.Error(err))
		api.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Debug("user signed up successfully",
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	api.WriteSuccess(w, http.StatusCreated, "Registration successful", response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	log := h.logger.With(zap.String("request_id", requestID))

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		api.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Authenticate user
	userID, err := h.authenticateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Error("failed to login",
			zap.Error(err),
			zap.String("username", req.Username),
		)
		api.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create tokens and session
	response, sessionToken, err := h.createTokensAndSession(r.Context(), userID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.Error(err))
		api.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Debug("user logged in successfully",
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	api.WriteSuccess(w, http.StatusOK, "Login successful", response)
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
		api.WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Delete session
	if err := h.sessions.Delete(r.Context(), cookie.Value); err != nil {
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

	log.Debug("user logged out successfully")
	api.WriteSuccess(w, http.StatusOK, "Logout successful", nil)
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	userID := api.GetUserID(r)
	log := h.logger.With(zap.String("request_id", requestID), zap.String("user_id", userID))

	// Get refresh token from request body
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		api.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var logMessage string

	// Try session-based refresh first
	var sessionBasedSuccess bool
	if cookie, err := r.Cookie("session"); err == nil {
		if session, err := h.sessions.GetByToken(r.Context(), cookie.Value); err == nil && session.ExpiresAt > time.Now().Unix() {
			// Validate refresh token matches session
			if req.RefreshToken != session.RefreshToken {
				log.Error("refresh token mismatch")
				api.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
				return
			}

			// Delete old session before creating new one
			if err := h.sessions.Delete(r.Context(), cookie.Value); err != nil {
				log.Error("failed to delete old session", zap.Error(err))
				api.WriteError(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			if session.UserID != userID {
				log.Error("session user ID mismatch", zap.String("session_user_id", session.UserID), zap.String("request_user_id", userID))
				api.WriteError(w, http.StatusUnauthorized, "Invalid session")
				return
			}

			logMessage = "token refreshed successfully via session"
		}
	}

	// Fallback to JWT-based refresh
	if !sessionBasedSuccess {
		claims, err := auth.ValidateToken(req.RefreshToken, h.secretKey)
		if err != nil {
			log.Error("invalid refresh token", zap.Error(err))
			api.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		if !claims.IsRefresh {
			log.Error("token is not a refresh token")
			api.WriteError(w, http.StatusUnauthorized, "Invalid token type")
			return
		}

		if claims.UserID != userID {
			log.Error("token user ID mismatch", zap.String("token_user_id", claims.UserID), zap.String("request_user_id", userID))
			api.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		logMessage = "token refreshed successfully via JWT"
	}

	// Create new tokens and session
	response, sessionToken, err := h.createTokensAndSession(r.Context(), userID)
	if err != nil {
		log.Error("failed to create tokens and session", zap.String("user_id", userID), zap.Error(err))

		// Check if the error is because the user no longer exists
		if errors.Is(err, repository.ErrNotFound) || strings.Contains(err.Error(), "resource not found") {
			log.Warn("user account no longer exists during token refresh", zap.String("user_id", userID))
			api.WriteError(w, http.StatusUnauthorized, "User account no longer exists")
			return
		}

		api.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set new session cookie
	h.setCookie(w, r, sessionToken, response.ExpiresAt)

	log.Info(logMessage,
		zap.String("username", response.User.Username),
		zap.String("email", response.User.Email),
	)

	api.WriteSuccess(w, http.StatusOK, "Token refreshed successfully", response)
}

func (h *AuthHandler) authenticateUser(ctx context.Context, username, password string) (string, error) {
	user, err := h.users.GetByUsername(ctx, username)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", auth.ErrInvalidCredentials
	}

	return user.ID, nil
}
