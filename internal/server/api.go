package server

import (
	"time"

	"mitsimi.dev/codeShare/internal/api"
	"mitsimi.dev/codeShare/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// setupAPIRoutes configures all API routes
func (s *Server) setupAPIRoutes(r chi.Router, authMiddleware *api.AuthMiddleware) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(15 * time.Second)) // Set a timeout for all API routes
		r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
		r.Use(middleware.AllowContentType("application/json"))

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			handler := handler.NewAuthHandler(s.users, s.sessions, s.secretKey)
			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
			r.Post("/logout", handler.Logout)
			r.Post("/refresh", handler.RefreshToken)
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			handler := handler.NewUserHandler(s.users, s.snippets, s.likes, s.bookmarks)
			r.Use(authMiddleware.RequireAuth) // Protect user routes

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetUser)                 // Get user by ID
				r.Get("/snippets", handler.GetUserSnippets) // Get user's snippets

				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.RequireSelfOrAdmin)      // Require self or admin access
					r.Get("/liked", handler.GetUserLikedSnippets) // Get user's liked snippets
					r.Get("/saved", handler.GetUserSavedSnippets) // Get user's saved snippets
					r.Patch("/", handler.UpdateProfile)
					r.Patch("/password", handler.UpdatePassword)
					r.Patch("/avatar", handler.UpdateAvatar)
				})
			})

			// /me routes - automatically use authenticated user's ID
			r.Route("/me", func(r chi.Router) {
				r.Get("/", handler.GetMe)                 // Get current user's profile
				r.Get("/snippets", handler.GetMySnippets) // Get current user's snippets

				r.Group(func(r chi.Router) {
					r.Get("/liked", handler.GetMyLikedSnippets)    // Get current user's liked snippets
					r.Get("/saved", handler.GetMySavedSnippets)    // Get current user's saved snippets
					r.Patch("/", handler.UpdateMyProfile)          // Update current user's profile
					r.Patch("/password", handler.UpdateMyPassword) // Update current user's password
					r.Patch("/avatar", handler.UpdateMyAvatar)     // Update current user's avatar
				})
			})
		})

		// Snippet routes
		r.Route("/snippets", func(r chi.Router) {
			handler := handler.NewSnippetHandler(s.snippets, s.likes, s.bookmarks, s.viewTracker, s.wsHub)

			// Public routes
			r.Group(func(r chi.Router) {
				r.Get("/", handler.GetSnippets)
				r.Get("/{id}", handler.GetSnippet)
			})

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Post("/", handler.CreateSnippet)

				r.Put("/{id}", handler.UpdateSnippet)
				r.Delete("/{id}", handler.DeleteSnippet)
				r.Patch("/{id}/like", handler.ToggleLikeSnippet)
				r.Patch("/{id}/save", handler.ToggleSaveSnippet)
			})
		})
	})
}
