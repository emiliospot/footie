package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients per match.
	clients map[int32]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Redis client for pub/sub.
	redis *redis.Client

	// Logger.
	logger *logger.Logger

	// Mutex for thread-safe operations.
	mu sync.RWMutex
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Match ID this client is subscribed to.
	matchID int32

	// User ID (optional, for authentication).
	userID int32
}

// Message represents a real-time event message.
type Message struct {
	Type      string      `json:"type"`       // "match_event", "score_update", "match_status"
	MatchID   int32       `json:"match_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// NewHub creates a new Hub instance.
func NewHub(redis *redis.Client, logger *logger.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int32]map[*Client]bool),
		redis:      redis,
		logger:     logger,
	}
}

// Run starts the hub's main loop.
func (h *Hub) Run(ctx context.Context) {
	// Start Redis Pub/Sub listener
	go h.listenToRedis(ctx)

	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Hub shutting down")
			return

		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.matchID] == nil {
				h.clients[client.matchID] = make(map[*Client]bool)
			}
			h.clients[client.matchID][client] = true
			h.mu.Unlock()
			h.logger.Info("Client registered", "match_id", client.matchID, "total_clients", len(h.clients[client.matchID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.matchID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.matchID)
					}
				}
			}
			h.mu.Unlock()
			h.logger.Info("Client unregistered", "match_id", client.matchID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.MatchID]
			h.mu.RUnlock()

			messageBytes, err := json.Marshal(message)
			if err != nil {
				h.logger.Error("Failed to marshal message", "error", err)
				continue
			}

			for client := range clients {
				select {
				case client.send <- messageBytes:
				default:
					close(client.send)
					h.mu.Lock()
					delete(h.clients[message.MatchID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// listenToRedis subscribes to Redis pub/sub channels for match updates.
func (h *Hub) listenToRedis(ctx context.Context) {
	pubsub := h.redis.PSubscribe(ctx, "match:*:events")
	defer pubsub.Close()

	h.logger.Info("Started Redis pub/sub listener")

	for {
		select {
		case <-ctx.Done():
			return

		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				h.logger.Error("Redis pub/sub error", "error", err)
				time.Sleep(time.Second)
				continue
			}

			// Parse the message and broadcast to WebSocket clients
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				h.logger.Error("Failed to unmarshal Redis message", "error", err)
				continue
			}

			h.broadcast <- &message
		}
	}
}

// BroadcastToMatch sends a message to all clients watching a specific match.
func (h *Hub) BroadcastToMatch(matchID int32, msgType string, data interface{}) {
	message := &Message{
		Type:      msgType,
		MatchID:   matchID,
		Timestamp: time.Now(),
		Data:      data,
	}
	h.broadcast <- message
}

// GetClientCount returns the number of clients watching a match.
func (h *Hub) GetClientCount(matchID int32) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[matchID])
}

