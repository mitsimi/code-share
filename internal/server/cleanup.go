package server

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// startSessionCleanup starts a background goroutine to periodically clean up expired sessions
func (s *Server) startSessionCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Run cleanup every hour
		defer ticker.Stop()

		for range ticker.C {
			if err := s.repos.Sessions.DeleteExpired(context.Background()); err != nil {
				s.logger.Error("Failed to delete expired sessions", zap.Error(err))
			} else {
				s.logger.Debug("Successfully cleaned up expired sessions")
			}
		}
	}()
}

// startViewCleanup starts a background goroutine to periodically clean up old view tracking records
func (s *Server) startViewCleanup() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Run cleanup daily
		defer ticker.Stop()

		for range ticker.C {
			if err := s.viewTracker.CleanupOldViews(context.Background()); err != nil {
				s.logger.Error("Failed to clean up old view tracking records", zap.Error(err))
			} else {
				s.logger.Debug("Successfully cleaned up old view tracking records")
			}
		}
	}()
}
