package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// setupAPIRoutes configures all API routes
func (s *Server) setupAPIRoutes(r chi.Router) {
	// Create auth middleware
	authMiddleware := api.NewAuthMiddleware(s.storage, s.secretKey)

	r.Use(authMiddleware.TryAttachUserID)       // Attach user ID to context
	r.Use(middleware.Timeout(15 * time.Second)) // Set a timeout for all API routes
	r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
	r.Use(middleware.AllowContentType("application/json"))

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",         // Development
				"https://codeshare.mitsimi.dev", // Production
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		handler := api.NewAuthHandler(s.storage, s.secretKey)
		r.Post("/signup", handler.Signup)
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/refresh", handler.RefreshToken)

		// Protected profile routes
		r.Route("/me", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/", handler.GetCurrentUser)
			r.Patch("/", handler.UpdateProfile)
			r.Patch("/password", handler.UpdatePassword)
			r.Patch("/avatar", handler.UpdateAvatar)
		})
	})

	// Snippet routes
	r.Route("/snippets", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",         // Development
				"https://codeshare.mitsimi.dev", // Production
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		handler := api.NewSnippetHandler(s.storage)

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

			r.Get("/liked", func(w http.ResponseWriter, r *http.Request) {
				s.logger.Info("API called: Get liked snippets",
					zap.String("request_id", middleware.GetReqID(r.Context())),
				)
			})
			r.Get("/saved", func(w http.ResponseWriter, r *http.Request) {
				s.logger.Info("API called: Get liked snippets",
					zap.String("request_id", middleware.GetReqID(r.Context())),
				)
			})
		})
	})
}
