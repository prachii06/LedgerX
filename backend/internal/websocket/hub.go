package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for simplicity, or we can check origins
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

type WSMessage struct {
	Type          string `json:"type"`
	TransactionID string `json:"transaction_id,omitempty"`
	Status        string `json:"status,omitempty"`
	Timestamp     string `json:"timestamp"`
	// Additional fields can be added if needed
	Details interface{} `json:"details,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 1024),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
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
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Hub) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	h.register <- client

	go client.writePump()

	// Read pump to handle disconnects (ping/pong)
	go func() {
		defer func() {
			h.unregister <- client
			conn.Close()
		}()
		conn.SetReadLimit(512)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// BroadcastMessage sends a JSON message to all clients without blocking the caller
func (h *Hub) BroadcastMessage(msgType string, transactionID string, status string, details interface{}) {
	msg := WSMessage{
		Type:          msgType,
		TransactionID: transactionID,
		Status:        status,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Details:       details,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling WS message:", err)
		return
	}

	select {
	case h.broadcast <- bytes:
	default:
		log.Println("WS broadcast channel full, dropping message")
	}
}
