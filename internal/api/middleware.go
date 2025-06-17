package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type contextKey string

const userIDKey contextKey = "user_id"

// AuthMiddleware is a middleware that checks for valid authentication
type AuthMiddleware struct {
	users     repository.UserRepository
	sessions  repository.SessionRepository
	logger    *zap.Logger
	secretKey string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(users repository.UserRepository, sessions repository.SessionRepository, secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		users:     users,
		sessions:  sessions,
		logger:    logger.Log,
		secretKey: secretKey,
	}
}

func (m *AuthMiddleware) TryAttachUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		log := m.logger.With(zap.String("request_id", requestID))

		var userID string

		// Try to get session from cookie first
		if cookie, err := r.Cookie("session"); err == nil {
			if session, err := m.sessions.GetByToken(r.Context(), cookie.Value); err == nil {
				if session.ExpiresAt > time.Now().Unix() {
					userID = session.UserID
				}
			}
		}

		// If no valid session, try JWT token
		if userID == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if user, err := auth.ValidateToken(token, m.secretKey); err == nil {
					userID = user.UserID
				}
			}
		}

		// Add user ID to context (even if empty)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		log.Debug("TryAuth completed", zap.String("user_id", userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware that requires valid authentication
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == "" {
			requestID := middleware.GetReqID(r.Context())
			log := m.logger.With(zap.String("request_id", requestID))
			log.Error("authentication required but not provided")
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireSelfOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == "" {
			requestID := middleware.GetReqID(r.Context())
			log := m.logger.With(zap.String("request_id", requestID))
			log.Error("authentication required but not provided")
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		// Check if the user is an admin or the same user
		if !auth.IsAdmin(userID) && chi.URLParam(r, "id") != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID gets the user ID from the context
func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(userIDKey).(string)
	return userID
}
