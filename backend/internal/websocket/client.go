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
		err := c.Conn.Close()
		if err != nil {
			log.Printf("failed to close websocket connection: %v", err)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("failed to set read deadline: %v", err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("failed to set read deadline: %v", err)
			return err
		}
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
		err := c.Conn.Close()
		if err != nil {
			log.Printf("failed to close websocket connection: %v", err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("failed to set write deadline: %v", err)
				return
			}
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("failed to write close message: %v", err)
				}
				return
			}
			err = c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("failed to write text message: %v", err)
				return
			}

		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("failed to set write deadline: %v", err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
