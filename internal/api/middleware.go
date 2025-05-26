package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// AuthMiddleware is a middleware that checks for valid authentication
type AuthMiddleware struct {
	storage   storage.Storage
	logger    *zap.Logger
	secretKey string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(storage storage.Storage, secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		storage:   storage,
		logger:    logger.Log,
		secretKey: secretKey,
	}
}

// RequireAuth is a middleware that requires valid authentication
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		log := m.logger.With(zap.String("request_id", requestID))

		// Try to get session from cookie
		cookie, err := r.Cookie("session")
		if err == nil {
			// Validate session
			session, err := m.storage.GetSession(cookie.Value)
			if err == nil {
				if session.ExpiresAt > time.Now().Unix() {
					// Add user ID to context
					ctx := context.WithValue(r.Context(), "user_id", session.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				// Session expired
				log.Error("session expired", zap.Int64("expires_at", session.ExpiresAt))
				http.Error(w, "Not authenticated", http.StatusUnauthorized)
				return
			}
		}

		// Try to get JWT token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error("no authorization header")
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		// Check if header has Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Error("invalid authorization header format")
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		userID, err := auth.ValidateToken(token, m.secretKey)
		if err != nil {
			log.Error("invalid token", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID gets the user ID from the context
func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value("user_id").(string)
	return userID
}
