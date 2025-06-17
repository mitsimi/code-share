package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/auth"
	"mitsimi.dev/codeShare/internal/models"
	"mitsimi.dev/codeShare/internal/storage"
	sqlitestorage "mitsimi.dev/codeShare/internal/storage/sqlite"
)

func setupTestSnippetHandler(
	t *testing.T,
	store storage.Storage,
) (*chi.Mux, func()) {
	// Create a test logger
	testLogger := zap.NewNop()

	// Create handler with test logger
	handler := &SnippetHandler{
		storage: store,
		logger:  testLogger,
	}

	authMiddleware := NewAuthMiddleware(store, "test-secret-key")
	authMiddleware.logger = testLogger // Set the test logger for the middleware

	// Create and configure router
	r := chi.NewRouter()
	r.Use(authMiddleware.RequireAuth)
	r.Patch("/api/snippets/{id}/like", handler.ToggleLikeSnippet)
	r.Get("/api/snippets/{id}", handler.GetSnippet)
	r.Get("/api/snippets", handler.GetSnippets)

	return r, func() {}
}

// createTestToken creates a valid JWT token for testing
func createTestToken(userID string, secretKey string) string {
	tokenResp, _ := auth.GenerateToken(userID, secretKey, false)
	return tokenResp.Token
}

func runToggleLikeSnippetTests(t *testing.T, store storage.Storage) {
	router, cleanup := setupTestSnippetHandler(t, store)
	defer cleanup()

	// Create test users
	user1, err := store.CreateUser(
		"testuser1",
		"testuser1@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	user2, err := store.CreateUser(
		"testuser2",
		"testuser2@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// Create test tokens
	user1Token := createTestToken(user1.ID, "test-secret-key")
	user2Token := createTestToken(user2.ID, "test-secret-key")

	// Create a test snippet
	snippet := models.Snippet{
		Title:    "Test Snippet",
		Content:  "Test Content",
		Language: "text",
		Author:   user1.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	tests := []struct {
		name         string
		userID       string
		snippetID    string
		action       string
		setupAuth    func(r *http.Request)
		wantStatus   int
		wantLikes    int
		wantIsLiked  bool
		wantError    bool
		errorMessage string
	}{
		{
			name:      "like snippet as authenticated user",
			userID:    user1.ID,
			snippetID: snippetID,
			action:    "like",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+user1Token)
			},
			wantStatus:  http.StatusOK,
			wantLikes:   1,
			wantIsLiked: true,
			wantError:   false,
		},
		{
			name:      "unlike snippet as authenticated user",
			userID:    user1.ID,
			snippetID: snippetID,
			action:    "unlike",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+user1Token)
			},
			wantStatus:  http.StatusOK,
			wantLikes:   0,
			wantIsLiked: false,
			wantError:   false,
		},
		{
			name:      "like non-existent snippet",
			userID:    user1.ID,
			snippetID: "non-existent-id",
			action:    "like",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+user1Token)
			},
			wantStatus:   http.StatusNotFound,
			wantError:    true,
			errorMessage: "snippet not found",
		},
		{
			name:      "like snippet as unauthenticated user",
			userID:    "",
			snippetID: snippetID,
			action:    "like",
			setupAuth: func(r *http.Request) {
				// No auth header
			},
			wantStatus:   http.StatusUnauthorized,
			wantError:    true,
			errorMessage: "Not authenticated",
		},
		{
			name:      "like snippet as another user",
			userID:    user2.ID,
			snippetID: snippetID,
			action:    "like",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+user2Token)
			},
			wantStatus:  http.StatusOK,
			wantLikes:   1,
			wantIsLiked: true,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(
				"PATCH",
				"/api/snippets/"+tt.snippetID+"/like?action="+tt.action,
				nil,
			)

			// Setup authentication
			tt.setupAuth(req)

			// Create response recorder
			w := httptest.NewRecorder()

			// Setup chi router context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.snippetID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Add user ID to context if authenticated
			if tt.userID != "" {
				req = req.WithContext(context.WithValue(req.Context(), userIDKey, tt.userID))
			}

			// Call handler through the router
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("ToggleLikeSnippet() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// For successful requests, verify the response
			if !tt.wantError {
				var response models.Snippet
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Verify response fields
				if response.Likes != tt.wantLikes {
					t.Errorf("Likes = %v, want %v", response.Likes, tt.wantLikes)
				}
				if response.IsLiked != tt.wantIsLiked {
					t.Errorf("IsLiked = %v, want %v", response.IsLiked, tt.wantIsLiked)
				}
			} else if tt.errorMessage != "" {
				// For error responses, verify the error message
				if !bytes.Contains(w.Body.Bytes(), []byte(tt.errorMessage)) {
					t.Errorf("Error message = %v, want %v", w.Body.String(), tt.errorMessage)
				}
			}
		})
	}
}

func TestToggleLikeSnippet(t *testing.T) {
	t.Run("with memory storage", func(t *testing.T) {
		store := storage.NewMemoryStorage()
		runToggleLikeSnippetTests(t, store)
	})

	t.Run("with SQLite storage", func(t *testing.T) {
		// Create a temporary SQLite database for testing
		store, err := sqlitestorage.NewSQLiteStorage(":memory:")
		if err != nil {
			t.Fatalf("Failed to create SQLite storage: %v", err)
		}
		defer store.Close()

		runToggleLikeSnippetTests(t, store)
	})
}

func runLikeStateConsistencyTests(t *testing.T, store storage.Storage) {
	router, cleanup := setupTestSnippetHandler(t, store)
	defer cleanup()

	// Create test user
	user, err := store.CreateUser(
		"testuser",
		"testuser@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create test token
	userToken := createTestToken(user.ID, "test-secret-key")

	// Create a test snippet
	snippet := models.Snippet{
		Title:    "Test Snippet",
		Content:  "Test Content",
		Language: "text",
		Author:   user.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Helper function to get snippet
	getSnippet := func() models.Snippet {
		req := httptest.NewRequest("GET", "/api/snippets/"+snippetID, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", snippetID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req = req.WithContext(context.WithValue(req.Context(), userIDKey, user.ID))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GetSnippet() status = %v, want %v", w.Code, http.StatusOK)
		}

		var response models.Snippet
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		return response
	}

	// Helper function to like/unlike snippet
	toggleLike := func(action string) models.Snippet {
		req := httptest.NewRequest(
			"PATCH",
			"/api/snippets/"+snippetID+"/like?action="+action,
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", snippetID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req = req.WithContext(context.WithValue(req.Context(), userIDKey, user.ID))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ToggleLikeSnippet() status = %v, want %v", w.Code, http.StatusOK)
		}

		var response models.Snippet
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		return response
	}

	t.Run("initial state has 0 likes and is not liked", func(t *testing.T) {
		initialSnippet := getSnippet()
		if initialSnippet.Likes != 0 {
			t.Errorf("Initial likes = %v, want 0", initialSnippet.Likes)
		}
		if initialSnippet.IsLiked {
			t.Error("Initial IsLiked = true, want false")
		}
	})

	t.Run("liking a snippet increments likes and sets IsLiked to true", func(t *testing.T) {
		likedSnippet := toggleLike("like")
		if likedSnippet.Likes != 1 {
			t.Errorf("Liked snippet likes = %v, want 1", likedSnippet.Likes)
		}
		if !likedSnippet.IsLiked {
			t.Error("Liked snippet IsLiked = false, want true")
		}
	})

	t.Run("getting the snippet after liking reflects the new state", func(t *testing.T) {
		gotLikedSnippet := getSnippet()
		if gotLikedSnippet.Likes != 1 {
			t.Errorf("Got liked snippet likes = %v, want 1", gotLikedSnippet.Likes)
		}
		if !gotLikedSnippet.IsLiked {
			t.Error("Got liked snippet IsLiked = false, want true")
		}
	})

	t.Run("unliking a snippet decrements likes and sets IsLiked to false", func(t *testing.T) {
		unlikedSnippet := toggleLike("unlike")
		if unlikedSnippet.Likes != 0 {
			t.Errorf("Unliked snippet likes = %v, want 0", unlikedSnippet.Likes)
		}
		if unlikedSnippet.IsLiked {
			t.Error("Unliked snippet IsLiked = true, want false")
		}
	})

	t.Run("getting the snippet after unliking reflects the final state", func(t *testing.T) {
		gotUnlikedSnippet := getSnippet()
		if gotUnlikedSnippet.Likes != 0 {
			t.Errorf("Got unliked snippet likes = %v, want 0", gotUnlikedSnippet.Likes)
		}
		if gotUnlikedSnippet.IsLiked {
			t.Error("Got unliked snippet IsLiked = true, want false")
		}
	})
}

func TestLikeStateConsistency(t *testing.T) {
	t.Run("with memory storage", func(t *testing.T) {
		store := storage.NewMemoryStorage()
		runLikeStateConsistencyTests(t, store)
	})

	t.Run("with SQLite storage", func(t *testing.T) {
		store, err := sqlitestorage.NewSQLiteStorage(":memory:")
		if err != nil {
			t.Fatalf("Failed to create SQLite storage: %v", err)
		}
		defer store.Close()

		runLikeStateConsistencyTests(t, store)
	})
}
