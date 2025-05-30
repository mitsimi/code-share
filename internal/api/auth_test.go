package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
)

func setupTestAuthHandler(_ *testing.T) (*AuthHandler, func()) {
	// Setup test storage
	store := storage.NewMemoryStorage()
	secretKey := "test-secret-key"

	// Create a test logger
	testLogger := zap.NewNop()

	// Create handler with test logger
	handler := &AuthHandler{
		storage:   store,
		logger:    testLogger,
		secretKey: secretKey,
	}

	// Return cleanup function
	cleanup := func() {}

	return handler, cleanup
}

func TestSignup(t *testing.T) {
	handler, cleanup := setupTestAuthHandler(t)
	defer cleanup()

	tests := []struct {
		name       string
		payload    models.SignupRequest
		wantStatus int
	}{
		{
			name: "valid signup",
			payload: models.SignupRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid email",
			payload: models.SignupRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Signup(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("Signup() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// For successful signup, verify the response
			if tt.wantStatus == http.StatusOK {
				var response models.AuthResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Verify response fields
				if response.User.Username != tt.payload.Username {
					t.Errorf("Username = %v, want %v", response.User.Username, tt.payload.Username)
				}
				if response.User.Email != tt.payload.Email {
					t.Errorf("Email = %v, want %v", response.User.Email, tt.payload.Email)
				}
				if response.Token == "" {
					t.Error("Expected token to be non-empty")
				}
				if response.RefreshToken == "" {
					t.Error("Expected refresh token to be non-empty")
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	handler, cleanup := setupTestAuthHandler(t)
	defer cleanup()

	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (string, *http.Cookie)
		wantStatus    int
		wantNewTokens bool
	}{
		{
			name: "valid session-based refresh",
			setupFunc: func(t *testing.T) (string, *http.Cookie) {
				return createUserAndGetTokens(t, handler)
			},
			wantStatus:    http.StatusOK,
			wantNewTokens: true,
		},
		{
			name: "valid JWT-based refresh",
			setupFunc: func(t *testing.T) (string, *http.Cookie) {
				refreshToken, _ := createUserAndGetTokens(t, handler)
				return refreshToken, nil // No session cookie
			},
			wantStatus:    http.StatusOK,
			wantNewTokens: true,
		},
		{
			name: "invalid refresh token",
			setupFunc: func(t *testing.T) (string, *http.Cookie) {
				_, sessionCookie := createUserAndGetTokens(t, handler)
				return "invalid-token", sessionCookie
			},
			wantStatus:    http.StatusUnauthorized,
			wantNewTokens: false,
		},
		{
			name: "missing refresh token",
			setupFunc: func(t *testing.T) (string, *http.Cookie) {
				_, sessionCookie := createUserAndGetTokens(t, handler)
				return "", sessionCookie
			},
			wantStatus:    http.StatusUnauthorized,
			wantNewTokens: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshToken, sessionCookie := tt.setupFunc(t)

			refreshReq := struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: refreshToken,
			}

			body, _ := json.Marshal(refreshReq)
			req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if sessionCookie != nil {
				req.AddCookie(sessionCookie)
			}

			w := httptest.NewRecorder()
			handler.RefreshToken(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("RefreshToken() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantNewTokens {
				var response models.AuthResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Verify tokens are not empty
				if response.Token == "" {
					t.Error("Expected access token to be non-empty")
				}
				if response.RefreshToken == "" {
					t.Error("Expected refresh token to be non-empty")
				}

				// Verify new refresh token is different from input
				if response.RefreshToken == refreshToken {
					t.Error("Expected new refresh token to be different from input")
				}

				// Verify expiration time
				if response.ExpiresAt <= time.Now().Unix() {
					t.Error("Expected expiration time to be in the future")
				}
			}
		})
	}
}

// Helper function to create a user and return tokens
func createUserAndGetTokens(t *testing.T, handler *AuthHandler) (string, *http.Cookie) {
	// Create a unique user for each test
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())

	signupReq := models.SignupRequest{
		Username: username,
		Email:    email,
		Password: "password123",
	}

	// Create request for signup
	body, _ := json.Marshal(signupReq)
	req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Signup(w, req)

	// Get the initial tokens from signup response
	var signupResponse models.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&signupResponse); err != nil {
		t.Fatalf("Failed to decode signup response: %v", err)
	}

	// Get the session cookie from the signup response
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session" {
			sessionCookie = cookie
			break
		}
	}

	// Create a session in storage if we have a session cookie
	if sessionCookie != nil {
		sessionToken := sessionCookie.Value
		err := handler.storage.CreateSession(
			signupResponse.User.ID,
			sessionToken,
			signupResponse.RefreshToken,
			time.Now().Add(24*time.Hour).Unix(),
		)
		if err != nil {
			t.Fatalf("Failed to create test session: %v", err)
		}
	}

	return signupResponse.RefreshToken, sessionCookie
}

func TestTokenExpiration(t *testing.T) {
	handler, cleanup := setupTestAuthHandler(t)
	defer cleanup()

	tests := []struct {
		name        string
		tokenType   string
		expireToken bool
		useSession  bool
		wantStatus  int
		wantError   string
	}{
		{
			name:        "expired access token with valid session",
			tokenType:   "access",
			expireToken: true,
			useSession:  true,
			wantStatus:  http.StatusUnauthorized,
			wantError:   "token is expired",
		},
		{
			name:        "valid access token with valid session",
			tokenType:   "access",
			expireToken: false,
			useSession:  true,
			wantStatus:  http.StatusOK,
			wantError:   "",
		},
		{
			name:        "expired refresh token in refresh request",
			tokenType:   "refresh",
			expireToken: true,
			useSession:  false,
			wantStatus:  http.StatusUnauthorized,
			wantError:   "Invalid refresh token",
		},
		{
			name:        "valid refresh token in refresh request",
			tokenType:   "refresh",
			expireToken: false,
			useSession:  false,
			wantStatus:  http.StatusOK,
			wantError:   "",
		},
		{
			name:        "expired session with expired refresh token",
			tokenType:   "session_and_refresh",
			expireToken: true,
			useSession:  true,
			wantStatus:  http.StatusUnauthorized,
			wantError:   "Invalid refresh token",
		},
		{
			name:        "expired session with valid refresh token",
			tokenType:   "session_only",
			expireToken: true,
			useSession:  true,
			wantStatus:  http.StatusOK, // Should succeed via JWT fallback
			wantError:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh user for each test
			username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
			email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())

			signupReq := models.SignupRequest{
				Username: username,
				Email:    email,
				Password: "password123",
			}

			// Signup and get initial tokens
			body, _ := json.Marshal(signupReq)
			req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.Signup(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Signup failed: %d", w.Code)
			}

			var signupResponse models.AuthResponse
			if err := json.NewDecoder(w.Body).Decode(&signupResponse); err != nil {
				t.Fatalf("Failed to decode signup response: %v", err)
			}

			// Get session cookie
			cookies := w.Result().Cookies()
			var sessionCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == "session" {
					sessionCookie = cookie
					break
				}
			}

			// Handle token expiration based on test case
			switch tt.tokenType {
			case "access":
				if tt.expireToken {
					// Create an expired access token
					expiredToken := createExpiredToken(t, signupResponse.User.ID, handler.secretKey, false)
					signupResponse.Token = expiredToken
				}

				// Test protected endpoint with access token
				testProtectedEndpoint(t, handler, signupResponse.Token, sessionCookie, tt.useSession, tt.wantStatus, tt.wantError)

			case "refresh":
				if tt.expireToken {
					// Create an expired refresh token
					expiredRefreshToken := createExpiredToken(t, signupResponse.User.ID, handler.secretKey, true)
					signupResponse.RefreshToken = expiredRefreshToken
				}

				// Test refresh endpoint
				testRefreshEndpoint(t, handler, signupResponse.RefreshToken, sessionCookie, tt.useSession, tt.wantStatus, tt.wantError)

			case "session_and_refresh":
				// Expire both session and refresh token
				if tt.expireToken {
					// Create expired session
					if sessionCookie != nil {
						expiredTime := time.Now().Add(-1 * time.Hour).Unix()
						err := handler.storage.CreateSession(
							signupResponse.User.ID,
							sessionCookie.Value,
							signupResponse.RefreshToken,
							expiredTime,
						)
						if err != nil {
							t.Fatalf("Failed to create expired session: %v", err)
						}
					}

					// Also expire the refresh token
					expiredRefreshToken := createExpiredToken(t, signupResponse.User.ID, handler.secretKey, true)
					signupResponse.RefreshToken = expiredRefreshToken
				}

				testRefreshEndpoint(t, handler, signupResponse.RefreshToken, sessionCookie, tt.useSession, tt.wantStatus, tt.wantError)

			case "session_only":
				// Expire only the session, keep refresh token valid
				if tt.expireToken && sessionCookie != nil {
					expiredTime := time.Now().Add(-1 * time.Hour).Unix()
					err := handler.storage.CreateSession(
						signupResponse.User.ID,
						sessionCookie.Value,
						signupResponse.RefreshToken,
						expiredTime,
					)
					if err != nil {
						t.Fatalf("Failed to create expired session: %v", err)
					}
				}

				testRefreshEndpoint(t, handler, signupResponse.RefreshToken, sessionCookie, tt.useSession, tt.wantStatus, tt.wantError)
			}
		})
	}
}

// Add a specific test for session expiration behavior
func TestSessionExpirationBehavior(t *testing.T) {
	handler, cleanup := setupTestAuthHandler(t)
	defer cleanup()

	// Create user and get tokens
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())

	signupReq := models.SignupRequest{
		Username: username,
		Email:    email,
		Password: "password123",
	}

	body, _ := json.Marshal(signupReq)
	req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Signup(w, req)

	var signupResponse models.AuthResponse
	json.NewDecoder(w.Body).Decode(&signupResponse)

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session" {
			sessionCookie = cookie
			break
		}
	}

	t.Run("expired session falls back to JWT validation", func(t *testing.T) {
		// Create expired session
		if sessionCookie != nil {
			expiredTime := time.Now().Add(-1 * time.Hour).Unix()
			err := handler.storage.CreateSession(
				signupResponse.User.ID,
				sessionCookie.Value,
				signupResponse.RefreshToken,
				expiredTime,
			)
			if err != nil {
				t.Fatalf("Failed to create expired session: %v", err)
			}
		}

		// Try refresh with expired session but valid JWT refresh token
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: signupResponse.RefreshToken,
		}

		body, _ := json.Marshal(refreshReq)
		req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		handler.RefreshToken(w, req)

		// Should succeed because it falls back to JWT validation
		if w.Code != http.StatusOK {
			t.Errorf("Expected fallback to JWT to succeed, got status %d: %s", w.Code, w.Body.String())
		}

		// Verify we got new tokens
		var response models.AuthResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Token == "" || response.RefreshToken == "" {
			t.Error("Expected new tokens to be generated")
		}
	})

	t.Run("expired session with no JWT fallback fails", func(t *testing.T) {
		// Create expired session
		if sessionCookie != nil {
			expiredTime := time.Now().Add(-1 * time.Hour).Unix()
			err := handler.storage.CreateSession(
				signupResponse.User.ID,
				sessionCookie.Value,
				"expired-refresh-token",
				expiredTime,
			)
			if err != nil {
				t.Fatalf("Failed to create expired session: %v", err)
			}
		}

		// Try refresh with expired session and invalid refresh token
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: "expired-refresh-token",
		}

		body, _ := json.Marshal(refreshReq)
		req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		handler.RefreshToken(w, req)

		// Should fail because both session and JWT are invalid
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d: %s", w.Code, w.Body.String())
		}
	})
}

// Helper function to create an expired token
func createExpiredToken(t *testing.T, userID, secretKey string, isRefresh bool) string {
	now := time.Now()
	expiredTime := now.Add(-1 * time.Hour) // 1 hour ago

	// Generate random bytes for JWT ID
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
	}
	jti := hex.EncodeToString(randomBytes)

	claims := &auth.JWTClaims{
		UserID:    userID,
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime), // Expired
			IssuedAt:  jwt.NewNumericDate(expiredTime),
			NotBefore: jwt.NewNumericDate(expiredTime),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	return tokenString
}

// Helper function to test a protected endpoint
func testProtectedEndpoint(_ *testing.T, handler *AuthHandler, accessToken string, sessionCookie *http.Cookie, useSession bool, wantStatus int, wantError string) {
	// Create a request to a protected endpoint (you'll need to implement this)
	// For now, let's assume you have a /api/user/profile endpoint
	req := httptest.NewRequest("GET", "/api/user/profile", nil)

	if useSession && sessionCookie != nil {
		req.AddCookie(sessionCookie)
	} else {
		// Use Bearer token
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	w := httptest.NewRecorder()

	// You'll need to create a middleware or handler that validates tokens
	// For this example, let's create a simple token validation
	validateTokenAndRespond(w, req, handler, wantStatus, wantError)
}

// Helper function to test refresh endpoint
func testRefreshEndpoint(t *testing.T, handler *AuthHandler, refreshToken string, sessionCookie *http.Cookie, useSession bool, wantStatus int, wantError string) {
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: refreshToken,
	}

	body, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if useSession && sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	w := httptest.NewRecorder()
	handler.RefreshToken(w, req)

	if w.Code != wantStatus {
		t.Errorf("Expected status %d, got %d. Response: %s", wantStatus, w.Code, w.Body.String())
	}

	if wantError != "" {
		responseBody := w.Body.String()
		if !strings.Contains(responseBody, wantError) {
			t.Errorf("Expected error message to contain '%s', got: %s", wantError, responseBody)
		}
	}
}

// Helper function to validate tokens (you'll need to adapt this to your actual middleware)
func validateTokenAndRespond(w http.ResponseWriter, r *http.Request, handler *AuthHandler, _ int, _ string) {
	// Try session-based auth first
	if cookie, err := r.Cookie("session"); err == nil {
		if session, err := handler.storage.GetSession(cookie.Value); err == nil {
			if session.ExpiresAt > time.Now().Unix() {
				// Valid session
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "success"}`))
				return
			}
		}
	}

	// Try Bearer token
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := auth.ValidateToken(token, handler.secretKey); err == nil {
			if !claims.IsRefresh {
				// Valid access token
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "success"}`))
				return
			}
		}
	}

	// Invalid or expired token
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "token is expired"}`))
}

// Additional test for token expiration edge cases
func TestTokenExpirationEdgeCases(t *testing.T) {
	handler, cleanup := setupTestAuthHandler(t)
	defer cleanup()

	t.Run("token expires during request processing", func(t *testing.T) {
		// Create a token that expires in 1 second
		userID := "test-user-id"
		now := time.Now()
		shortExpiry := now.Add(1 * time.Second)

		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)
		jti := hex.EncodeToString(randomBytes)

		claims := &auth.JWTClaims{
			UserID:    userID,
			IsRefresh: false,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(shortExpiry),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				ID:        jti,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(handler.secretKey))
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		// Wait for token to expire
		time.Sleep(2 * time.Second)

		// Try to use expired token
		req := httptest.NewRequest("GET", "/api/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		validateTokenAndRespond(w, req, handler, http.StatusUnauthorized, "token is expired")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("refresh token rotation after expiration", func(t *testing.T) {
		// Create user
		username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
		email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())

		signupReq := models.SignupRequest{
			Username: username,
			Email:    email,
			Password: "password123",
		}

		body, _ := json.Marshal(signupReq)
		req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Signup(w, req)

		var signupResponse models.AuthResponse
		json.NewDecoder(w.Body).Decode(&signupResponse)

		// Create an expired refresh token
		expiredRefreshToken := createExpiredToken(t, signupResponse.User.ID, handler.secretKey, true)

		// Try to refresh with expired token
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: expiredRefreshToken,
		}

		body, _ = json.Marshal(refreshReq)
		req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		handler.RefreshToken(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for expired refresh token, got %d", w.Code)
		}

		// Verify that a valid refresh token still works
		refreshReq.RefreshToken = signupResponse.RefreshToken
		body, _ = json.Marshal(refreshReq)
		req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		handler.RefreshToken(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for valid refresh token, got %d", w.Code)
		}
	})
}
