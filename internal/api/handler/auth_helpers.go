package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/api/dto"
	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/constants"
	"mitsimi.dev/codeShare/internal/domain"
)

// createTokensAndSession generates new tokens, creates a session, and returns the auth response
func (h *AuthHandler) createTokensAndSession(ctx context.Context, userID string) (*dto.AuthResponse, string, error) {
	// Generate access token
	accessTokenResp, err := auth.GenerateToken(userID, h.secretKey, false)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenResp, err := auth.GenerateToken(userID, h.secretKey, true)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session token
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		Token:        sessionToken,
		RefreshToken: refreshTokenResp.Token,
		ExpiresAt:    refreshTokenResp.ExpiresAt,
	}

	// Create session in storage
	if err := h.sessions.Create(ctx, session); err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	// Get user details
	user, err := h.users.GetByID(ctx, userID)
	if err != nil {
		// Clean up the session we just created since user doesn't exist
		if deleteErr := h.sessions.Delete(ctx, sessionToken); deleteErr != nil {
			// Log the cleanup failure but don't override the original error
			h.logger.Warn("failed to clean up session after user lookup failure", zap.Error(deleteErr))
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	response := &dto.AuthResponse{
		Token:        accessTokenResp.Token,
		RefreshToken: refreshTokenResp.Token,
		User:         dto.ToUserResponse(user),
		ExpiresAt:    accessTokenResp.ExpiresAt,
	}

	return response, sessionToken, nil
}

// setCookie sets the session cookie with proper security settings
func (h *AuthHandler) setCookie(w http.ResponseWriter, r *http.Request, sessionToken string, expiresAt int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sessionToken,
		Path:     constants.CookiePath,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(expiresAt, 0),
	})
}

// isValidEmail checks if the email is valid
func isValidEmail(email string) bool {
	// Basic email validation regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

// validatePassword checks if the password meets the requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
