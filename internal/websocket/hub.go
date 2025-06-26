package ws

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
	"mitsimi.dev/codeShare/internal/logger"
)

// Hub maintains active clients and handles broadcasting
type Hub struct {
	// Client management
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client

	// Subscription maps
	userClients          map[string][]*Client // userID -> clients for user_actions
	snippetUpdateClients map[string][]*Client // snippetID -> clients for snippet_updates
	listUpdateClients    []*Client            // clients subscribed to list_updates

	// Broadcasting
	broadcast chan BroadcastMessage

	mutex  sync.RWMutex
	logger *zap.Logger
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:              make(map[*Client]struct{}),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
		userClients:          make(map[string][]*Client),
		snippetUpdateClients: make(map[string][]*Client),
		listUpdateClients:    make([]*Client, 0),
		broadcast:            make(chan BroadcastMessage),
		logger:               logger.With(zap.String("websocket", "hub")),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case broadcastMsg := <-h.broadcast:
			h.handleBroadcast(broadcastMsg)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = struct{}{}
	h.logger.Info("Client registered",
		zap.String("user_id", client.userID),
		zap.Int("total_clients", len(h.clients)))

	// Send success message
	client.SendMessage(WebSocketMessage{
		Type:      MessageTypeSuccess,
		Data:      "Connected successfully",
		Timestamp: time.Now().Unix(),
	})
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		// Remove from all subscription maps
		h.removeFromUserClients(client)
		h.removeFromSnippetClients(client)
		h.removeFromListClients(client)

		delete(h.clients, client)
		close(client.send)

		h.logger.Info("Client unregistered",
			zap.String("user_id", client.userID),
			zap.Int("total_clients", len(h.clients)))
	}
}

func (h *Hub) removeFromUserClients(client *Client) {
	if clients, exists := h.userClients[client.userID]; exists {
		for i, c := range clients {
			if c == client {
				h.userClients[client.userID] = append(clients[:i], clients[i+1:]...)
				if len(h.userClients[client.userID]) == 0 {
					delete(h.userClients, client.userID)
				}
				break
			}
		}
	}
}

func (h *Hub) removeFromSnippetClients(client *Client) {
	// Remove from snippet updates
	for snippetID := range client.snippetUpdatesSubscribed {
		if clients, exists := h.snippetUpdateClients[snippetID]; exists {
			for i, c := range clients {
				if c == client {
					h.snippetUpdateClients[snippetID] = append(clients[:i], clients[i+1:]...)
					if len(h.snippetUpdateClients[snippetID]) == 0 {
						delete(h.snippetUpdateClients, snippetID)
					}
					break
				}
			}
		}
	}
}

func (h *Hub) removeFromListClients(client *Client) {
	for i, c := range h.listUpdateClients {
		if c == client {
			h.listUpdateClients = append(h.listUpdateClients[:i], h.listUpdateClients[i+1:]...)
			break
		}
	}
}

func (h *Hub) handleBroadcast(broadcastMsg BroadcastMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	switch broadcastMsg.Target.Type {
	case BroadcastTargetTypeUser:
		if broadcastMsg.Target.UserID != nil {
			h.broadcastToUser(*broadcastMsg.Target.UserID, broadcastMsg.Message)
		}
	case BroadcastTargetTypeSnippetUpdates:
		if broadcastMsg.Target.SnippetID != nil {
			h.broadcastToSnippetUpdates(*broadcastMsg.Target.SnippetID, broadcastMsg.Message)
		}
	case BroadcastTargetTypeListUpdates:
		h.broadcastToListUpdates(broadcastMsg.Message)
	}
}

func (h *Hub) broadcastToUser(userID string, message WebSocketMessage) {
	clients := h.userClients[userID]
	messageBytes, _ := json.Marshal(message)

	for _, client := range clients {
		select {
		case client.send <- messageBytes:
		default:
			// Client buffer full, remove client
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) broadcastToSnippetUpdates(snippetID string, message WebSocketMessage) {
	clients := h.snippetUpdateClients[snippetID]
	messageBytes, _ := json.Marshal(message)

	for _, client := range clients {
		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) broadcastToListUpdates(message WebSocketMessage) {
	messageBytes, _ := json.Marshal(message)

	for _, client := range h.listUpdateClients {
		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]any {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]any{
		"total_clients":         len(h.clients),
		"user_subscriptions":    len(h.userClients),
		"snippet_subscriptions": len(h.snippetUpdateClients),
		"list_subscriptions":    len(h.listUpdateClients),
	}
}
