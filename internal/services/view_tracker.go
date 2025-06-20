package services

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
)

// ViewTracker handles view counting with debouncing logic
type ViewTracker struct {
	viewRepo repository.ViewRepository
	logger   *zap.Logger

	// Configuration
	ViewCooldownMinutes int // Time before counting another view from same viewer
}

// NewViewTracker creates a new view tracker with default settings
func NewViewTracker(viewRepo repository.ViewRepository) *ViewTracker {
	return &ViewTracker{
		viewRepo:            viewRepo,
		logger:              logger.Log,
		ViewCooldownMinutes: 10,
	}
}

// ViewerIdentifier extracts a unique identifier for the viewer
func (vt *ViewTracker) ViewerIdentifier(r *http.Request, userID string) string {
	// Use userID if authenticated, otherwise use session cookie or create temporary identifier
	if userID != "" {
		return "user:" + userID
	}

	// Try to get session cookie
	if cookie, err := r.Cookie("session"); err == nil {
		return "session:" + cookie.Value
	}

	// Fallback to IP address (less reliable but better than nothing)
	return "ip:" + GetClientIP(r)
}

// GetClientIP extracts the real client IP from request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return cleanIPAddress(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return cleanIPAddress(xri)
	}

	// Fallback to RemoteAddr (remove port if present)
	return cleanIPAddress(r.RemoteAddr)
}

// cleanIPAddress removes port from IP address if present
func cleanIPAddress(addr string) string {
	// Handle IPv6 addresses like [::1]:64011
	if strings.HasPrefix(addr, "[") {
		if idx := strings.LastIndex(addr, "]:"); idx != -1 {
			return addr[1:idx] // Extract IP without brackets and port
		}
		return addr // Return as-is if no port found
	}

	// Handle IPv4 addresses like 192.168.1.1:8080
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx] // Extract IP without port
	}

	return addr // Return as-is if no port found
}

// ShouldCountView determines if a view should be counted based on debouncing rules
func (vt *ViewTracker) ShouldCountView(ctx context.Context, snippetID, viewerIdentifier string) (bool, error) {
	vt.logger.Debug("checking if view should be counted",
		zap.String("snippet_id", snippetID),
		zap.String("viewer_identifier", viewerIdentifier),
		zap.Int("cooldown_minutes", vt.ViewCooldownMinutes),
	)

	recentView, err := vt.viewRepo.CheckRecentView(ctx, snippetID, viewerIdentifier)
	if err != nil {
		if err == repository.ErrNotFound {
			// No previous view found, should count
			vt.logger.Debug("no previous view found, counting view",
				zap.String("snippet_id", snippetID),
				zap.String("viewer_identifier", viewerIdentifier),
			)
			return true, nil
		}
		vt.logger.Error("error checking recent view",
			zap.Error(err),
			zap.String("snippet_id", snippetID),
			zap.String("viewer_identifier", viewerIdentifier),
		)
		return false, err
	}

	// Check if enough time has passed since last view
	cooldownSeconds := int64(vt.ViewCooldownMinutes * 60)
	shouldCount := recentView.SecondsSinceLastView >= cooldownSeconds

	vt.logger.Debug("recent view found, checking cooldown",
		zap.String("snippet_id", snippetID),
		zap.String("viewer_identifier", viewerIdentifier),
		zap.Int64("seconds_since_last_view", recentView.SecondsSinceLastView),
		zap.Int64("cooldown_seconds", cooldownSeconds),
		zap.Bool("should_count", shouldCount),
		zap.Time("last_viewed_at", recentView.LastViewedAt),
	)

	return shouldCount, nil
}

// TrackView records a view and increments the counter if appropriate
func (vt *ViewTracker) TrackView(ctx context.Context, r *http.Request, snippetID, userID string) error {
	viewerIdentifier := vt.ViewerIdentifier(r, userID)
	clientIP := GetClientIP(r)

	vt.logger.Debug("tracking view",
		zap.String("snippet_id", snippetID),
		zap.String("viewer_identifier", viewerIdentifier),
		zap.String("client_ip", clientIP),
		zap.String("user_id", userID),
	)

	shouldCount, err := vt.ShouldCountView(ctx, snippetID, viewerIdentifier)
	if err != nil {
		vt.logger.Error("failed to check if view should be counted",
			zap.Error(err),
			zap.String("snippet_id", snippetID),
			zap.String("viewer_identifier", viewerIdentifier),
		)
		return err
	}

	// Always record the view attempt (for analytics)
	if err := vt.viewRepo.RecordView(ctx, snippetID, viewerIdentifier, clientIP); err != nil {
		vt.logger.Error("failed to record view",
			zap.Error(err),
			zap.String("snippet_id", snippetID),
			zap.String("viewer_identifier", viewerIdentifier),
		)
		return err
	}

	// Only increment the public view count if cooldown has passed
	if shouldCount {
		if err := vt.viewRepo.IncrementViewCount(ctx, snippetID); err != nil {
			vt.logger.Error("failed to increment view count",
				zap.Error(err),
				zap.String("snippet_id", snippetID),
			)
			return err
		}

		vt.logger.Debug("view counted",
			zap.String("snippet_id", snippetID),
			zap.String("viewer_identifier", viewerIdentifier),
		)
	} else {
		vt.logger.Debug("view not counted (cooldown active)",
			zap.String("snippet_id", snippetID),
			zap.String("viewer_identifier", viewerIdentifier),
		)
	}

	return nil
}

// CleanupOldViews removes old view tracking records (should be called periodically)
func (vt *ViewTracker) CleanupOldViews(ctx context.Context) error {
	return vt.viewRepo.CleanupOldViews(ctx)
}

// SetCooldownMinutes allows customizing the view cooldown period
func (vt *ViewTracker) SetCooldownMinutes(minutes int) {
	vt.ViewCooldownMinutes = minutes
}
