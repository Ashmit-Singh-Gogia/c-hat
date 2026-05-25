package handlers

import (
	"net/http"
	"strconv"

	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	ws "github.com/ashmit-singh-gogia/c-hat/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub            *ws.Hub
	messageService *services.MessageService
}

func NewWSHandler(hub *ws.Hub, messageService *services.MessageService) *WSHandler {
	return &WSHandler{hub: hub, messageService: messageService}
}

func (h *WSHandler) HandleWS(c *gin.Context) {

	chatIDStr := c.Param("chatId")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// UserID was set by AuthMiddleware
	userID := c.GetUint("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Hub:            h.hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		ChatID:         uint(chatID),
		MessageService: h.messageService,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
