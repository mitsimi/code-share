package api

import (
	"fmt"
	"net/http"
	"time"

	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/models"
)

// createTokensAndSession generates new tokens, creates a session, and returns the auth response
func (h *AuthHandler) createTokensAndSession(userID string) (*models.AuthResponse, string, error) {
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

	// Create session in storage
	if err := h.storage.CreateSession(userID, sessionToken, refreshTokenResp.Token, refreshTokenResp.ExpiresAt); err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	// Get user details
	user, err := h.storage.GetUser(userID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	response := &models.AuthResponse{
		Token:        accessTokenResp.Token,
		RefreshToken: refreshTokenResp.Token,
		User:         models.FromDBUser(user),
		ExpiresAt:    accessTokenResp.ExpiresAt,
	}

	return response, sessionToken, nil
}

// setCookie sets the session cookie with proper security settings
func (h *AuthHandler) setCookie(w http.ResponseWriter, r *http.Request, sessionToken string, expiresAt int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(expiresAt, 0),
	})
}
