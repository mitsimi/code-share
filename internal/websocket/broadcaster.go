package ws

import (
	"time"

	"go.uber.org/zap"
)

// BroadcastUserAction syncs like/save status across user's devices
func (h *Hub) BroadcastUserAction(userID string, data UserActionData) {
	message := WebSocketMessage{
		Type:      MessageTypeUserActions,
		Data:      data,
		UserID:    &userID,
		SnippetID: &data.SnippetID,
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- BroadcastMessage{
		Message: message,
		Target: BroadcastTarget{
			Type:   BroadcastTargetTypeUser,
			UserID: &userID,
		},
	}

	h.logger.Debug("Broadcasting user action",
		zap.String("action", data.Action),
		zap.String("snippet_id", data.SnippetID),
		zap.String("user_id", userID))
}

// BroadcastSnippetUpdate sends content + stats updates to specific snippet subscribers
func (h *Hub) BroadcastSnippetUpdate(snippetID string, data SnippetUpdateData) {
	message := WebSocketMessage{
		Type:      MessageTypeSnippetUpdates,
		Data:      data,
		SnippetID: &snippetID,
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- BroadcastMessage{
		Message: message,
		Target: BroadcastTarget{
			Type:      BroadcastTargetTypeSnippetUpdates,
			SnippetID: &snippetID,
		},
	}

	h.logger.Debug("Broadcasting snippet update",
		zap.String("update_type", data.UpdateType),
		zap.String("snippet_id", snippetID))
}

// BroadcastListUpdate sends content changes to all list view subscribers
func (h *Hub) BroadcastListUpdate(data ListUpdateData) {
	message := WebSocketMessage{
		Type:      MessageTypeListUpdates,
		Data:      data,
		SnippetID: &data.SnippetID,
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- BroadcastMessage{
		Message: message,
		Target: BroadcastTarget{
			Type: BroadcastTargetTypeListUpdates,
		},
	}

	h.logger.Debug("Broadcasting list update", zap.String("snippet_id", data.SnippetID))
}

// BroadcastSnippetContentUpdate is a convenience method that broadcasts both
// snippet_updates and list_updates for content changes
func (h *Hub) BroadcastSnippetContentUpdate(snippetID string, title, content, language *string) {
	// Broadcast to specific snippet subscribers (detail view)
	h.BroadcastSnippetUpdate(snippetID, SnippetUpdateData{
		SnippetID:  snippetID,
		UpdateType: "content",
		Title:      title,
		Content:    content,
		Language:   language,
	})

	// Broadcast to list view subscribers
	h.BroadcastListUpdate(ListUpdateData{
		SnippetID: snippetID,
		Title:     title,
		Content:   content,
		Language:  language,
	})
}

// BroadcastSnippetStatsUpdate is a convenience method for stats-only updates
func (h *Hub) BroadcastSnippetStatsUpdate(snippetID string, viewCount, likeCount *int) {
	h.BroadcastSnippetUpdate(snippetID, SnippetUpdateData{
		SnippetID:  snippetID,
		UpdateType: "stats",
		ViewCount:  viewCount,
		LikeCount:  likeCount,
	})
}
