// Package ws_hub 推送實時數據給前端
package ws_hub

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Hub 管理所有 WebSocket 客戶端，廣播訊息
type Hub struct {
	clients map[*client]struct{}
	register chan *client
	unregister chan *client
	broadcast  chan []byte
	log        *zap.Logger
	mu         sync.RWMutex
}

type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// New 建立 Hub
func New(log *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*client]struct{}),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan []byte, 256),
		log:        log,
	}
}

// Run 在單一 goroutine 中處理註冊、註銷與廣播
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()
			h.log.Info("ws client connected", zap.Int("total", len(h.clients)))

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastJSON 將 v 序列化後廣播給所有連線
func (h *Hub) BroadcastJSON(v interface{}) {
	payload, err := json.Marshal(v)
	if err != nil {
		h.log.Warn("broadcast json marshal", zap.Error(err))
		return
	}
	select {
	case h.broadcast <- payload:
	default:
		h.log.Warn("broadcast channel full, drop message")
	}
}

// Upgrader 供 HTTP 升級為 WebSocket 使用
func Upgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}

// Register 註冊新連線（由 HTTP handler 呼叫）
func (h *Hub) Register(conn *websocket.Conn) {
	c := &client{hub: h, conn: conn, send: make(chan []byte, 256)}
	h.register <- c
	go c.writePump()
	go c.readPump()
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (c *client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
