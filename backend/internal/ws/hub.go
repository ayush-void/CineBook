package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool // map[showID]map[conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request, showID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	if h.clients[showID] == nil {
		h.clients[showID] = make(map[*websocket.Conn]bool)
	}
	h.clients[showID][conn] = true
	h.mu.Unlock()

	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients[showID], conn)
			h.mu.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// BroadcastSeatUpdate pushes a seat status change to every client
// watching a given show. Call this after Redis holds a seat or
// Postgres confirms/cancels a booking.
func (h *Hub) BroadcastSeatUpdate(showID, seatID, status string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := map[string]string{"seatId": seatID, "status": status}

	for conn := range h.clients[showID] {
		conn.WriteJSON(msg)
	}
}
