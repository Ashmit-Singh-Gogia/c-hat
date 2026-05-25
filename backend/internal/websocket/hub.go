package websocket

import "sync"

type Hub struct {
	// chatID → set of clients in that chat
	rooms      map[uint]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	ChatID  uint
	Payload []byte
	Sender  *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.ChatID] == nil {
				h.rooms[client.ChatID] = make(map[*Client]bool)
			}
			h.rooms[client.ChatID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.ChatID]; ok {
				if _, exists := room[client]; exists {
					delete(room, client)
					close(client.Send)
					if len(room) == 0 {
						delete(h.rooms, client.ChatID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.rooms[msg.ChatID] {
				if client == msg.Sender {
					continue
				}
				select {
				case client.Send <- msg.Payload:
				default:
					// Buffer full — client is too slow, drop them
					close(client.Send)
					delete(h.rooms[msg.ChatID], client)
				}
			}
			h.mu.RUnlock()
		}
	}
}
