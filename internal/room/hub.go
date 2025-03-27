package room

import (
	"maps"
	"sync"
	"time"

	"github.com/escalopa/peer-cast/internal/domain"
	"github.com/gorilla/websocket"
)

const hubCleanupInterval = 5 * time.Minute

// Hub handles WebRTC signaling for multiple rooms
type Hub struct {
	rooms map[string]*room
	mutex sync.RWMutex
}

// NewHub creates a new WebRTCHandler
func NewHub() *Hub {
	h := &Hub{rooms: make(map[string]*room)}
	go h.cleanup()
	return h
}

func (h *Hub) getOrCreateRoom(roomID string) *room {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	r, ok := h.rooms[roomID]
	if !ok {
		r = newRoom()
		h.rooms[roomID] = r
		go r.run()
	}

	return r
}

func (h *Hub) Handle(user *domain.User, roomID string, conn *websocket.Conn) {
	r := h.getOrCreateRoom(roomID)
	r.join(newUser(user, conn))
}

func (h *Hub) cleanup() {
	ticker := time.NewTicker(hubCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		h.mutex.RLock()
		rs := maps.Clone(h.rooms)
		h.mutex.RUnlock()

		for roomID, r := range rs {
			if len(r.users) == 0 {
				r.stop()
				h.mutex.Lock()
				delete(h.rooms, roomID)
				h.mutex.Unlock()
			}
		}
	}
}
