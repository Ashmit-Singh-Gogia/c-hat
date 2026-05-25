package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type IncomingMessage struct {
	Content string `json:"content"`
}

type ChatMessage struct {
	ID        uint   `json:"id"`
	SenderID  uint   `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type OutgoingMessage struct {
	Type    string      `json:"type"`
	ChatID  uint        `json:"chat_id"`
	Payload ChatMessage `json:"payload"`
}

type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         uint
	ChatID         uint
	MessageService *services.MessageService
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(rawMsg, &incoming); err != nil {
			log.Printf("invalid json from client: %v", err)
			continue
		}

		savedMsg, err := c.MessageService.SendMessage(c.ChatID, c.UserID, incoming.Content)
		if err != nil {
			log.Printf("failed to save message: %v", err)
			continue
		}
		message := ChatMessage{
			ID:        savedMsg.ID,
			SenderID:  savedMsg.SenderID,
			Content:   savedMsg.Content,
			CreatedAt: savedMsg.CreatedAt.Format(time.RFC3339),
		}

		outgoing, _ := json.Marshal(OutgoingMessage{
			Type:    "new_message",
			ChatID:  c.ChatID,
			Payload: message,
		})

		c.Hub.Broadcast <- &BroadcastMessage{
			ChatID:  c.ChatID,
			Payload: outgoing,
			Sender:  c,
		}

		c.Send <- outgoing
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
