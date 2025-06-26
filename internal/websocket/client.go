package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/logger"
)

// Client represents a WebSocket client
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string

	// Subscriptions
	userActionsSubscribed    bool
	snippetUpdatesSubscribed map[string]bool // snippetID -> subscribed
	listUpdatesSubscribed    bool            // Global list updates subscription

	mutex  sync.RWMutex
	logger *zap.Logger
}

// NewClient creates a new client
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:                      hub,
		conn:                     conn,
		send:                     make(chan []byte, 256),
		userID:                   userID,
		userActionsSubscribed:    false,
		snippetUpdatesSubscribed: make(map[string]bool),
		listUpdatesSubscribed:    false,
		mutex:                    sync.RWMutex{},
		logger:                   logger.With(zap.String("websocket", "client"), zap.String("user_id", userID)),
	}
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(message WebSocketMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("Error marshaling message", zap.Error(err))
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		c.hub.unregister <- c
	}
}

// HandleSubscription handles subscription requests
func (c *Client) HandleSubscription(subReq SubscriptionRequest) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch subReq.Type {
	case SubTypeUserActions:
		if c.userID == "anonymous" {
			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeError,
				Data:      "Anonymous users cannot subscribe to user_actions",
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if !c.userActionsSubscribed {
			c.userActionsSubscribed = true
			c.hub.mutex.Lock()
			c.hub.userClients[c.userID] = append(c.hub.userClients[c.userID], c)
			c.hub.mutex.Unlock()

			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeSuccess,
				Data:      "Subscribed to user_actions",
				Timestamp: time.Now().Unix(),
			})
			c.logger.Info("User subscribed to user_actions")
		}

	case SubTypeSnippetUpdates:
		if subReq.SnippetID != nil {
			snippetID := *subReq.SnippetID
			if !c.snippetUpdatesSubscribed[snippetID] {
				c.snippetUpdatesSubscribed[snippetID] = true
				c.hub.mutex.Lock()
				c.hub.snippetUpdateClients[snippetID] = append(
					c.hub.snippetUpdateClients[snippetID], c)
				c.hub.mutex.Unlock()

				c.SendMessage(WebSocketMessage{
					Type:      MessageTypeSuccess,
					Data:      "Subscribed to snippet_updates for " + snippetID,
					SnippetID: &snippetID,
					Timestamp: time.Now().Unix(),
				})
				c.logger.Info("User subscribed to snippet_updates", zap.String("snippet_id", snippetID))
			}
		}

	case SubTypeListUpdates:
		if !c.listUpdatesSubscribed {
			c.listUpdatesSubscribed = true
			c.hub.mutex.Lock()
			c.hub.listUpdateClients = append(c.hub.listUpdateClients, c)
			c.hub.mutex.Unlock()

			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeSuccess,
				Data:      "Subscribed to list_updates",
				Timestamp: time.Now().Unix(),
			})
			c.logger.Info("User subscribed to list_updates")
		}
	}
}

// HandleUnsubscription handles unsubscription requests
func (c *Client) HandleUnsubscription(subReq SubscriptionRequest) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch subReq.Type {
	case SubTypeUserActions:
		if c.userActionsSubscribed {
			c.userActionsSubscribed = false
			c.hub.removeFromUserClients(c)

			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeSuccess,
				Data:      "Unsubscribed from user_actions",
				Timestamp: time.Now().Unix(),
			})
			c.logger.Info("User unsubscribed from user_actions")
		}

	case SubTypeSnippetUpdates:
		if subReq.SnippetID != nil {
			snippetID := *subReq.SnippetID
			if c.snippetUpdatesSubscribed[snippetID] {
				delete(c.snippetUpdatesSubscribed, snippetID)
				// Remove from hub's snippet update clients
				c.hub.mutex.Lock()
				if clients, exists := c.hub.snippetUpdateClients[snippetID]; exists {
					for i, client := range clients {
						if client == c {
							c.hub.snippetUpdateClients[snippetID] = append(
								clients[:i], clients[i+1:]...)
							if len(c.hub.snippetUpdateClients[snippetID]) == 0 {
								delete(c.hub.snippetUpdateClients, snippetID)
							}
							break
						}
					}
				}
				c.hub.mutex.Unlock()

				c.SendMessage(WebSocketMessage{
					Type:      MessageTypeSuccess,
					Data:      "Unsubscribed from snippet_updates for " + snippetID,
					SnippetID: &snippetID,
					Timestamp: time.Now().Unix(),
				})
				c.logger.Info("User unsubscribed from snippet_updates", zap.String("snippet_id", snippetID))
			}
		}

	case SubTypeListUpdates:
		if c.listUpdatesSubscribed {
			c.listUpdatesSubscribed = false
			c.hub.removeFromListClients(c)

			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeSuccess,
				Data:      "Unsubscribed from list_updates",
				Timestamp: time.Now().Unix(),
			})
			c.logger.Info("User unsubscribed from list_updates")
		}
	}
}

// ReadPump handles reading messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		var message WebSocketMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			c.SendMessage(WebSocketMessage{
				Type:      MessageTypeError,
				Data:      "Invalid message format",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		switch message.Type {
		case MessageTypeSubscribe:
			if subData, ok := message.Data.(map[string]any); ok {
				subReq := SubscriptionRequest{}
				if subType, exists := subData["type"].(string); exists {
					subReq.Type = SubscriptionType(subType)
				}
				if snippetID, exists := subData["snippet_id"].(string); exists {
					subReq.SnippetID = &snippetID
				}
				c.HandleSubscription(subReq)
			}

		case MessageTypeUnsubscribe:
			if subData, ok := message.Data.(map[string]any); ok {
				subReq := SubscriptionRequest{}
				if subType, exists := subData["type"].(string); exists {
					subReq.Type = SubscriptionType(subType)
				}
				if snippetID, exists := subData["snippet_id"].(string); exists {
					subReq.SnippetID = &snippetID
				}
				c.HandleUnsubscription(subReq)
			}
		}
	}
}

// WritePump handles writing messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.logger.Debug("Sending message", zap.String("message", string(message)))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			c.logger.Debug("Message sent", zap.String("message", string(message)))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
