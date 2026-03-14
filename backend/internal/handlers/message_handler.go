package handlers

import (
	"strconv"

	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	service *services.MessageService
}
type SendMessageRequest struct {
	ChatID   uint   `json:"chat_id"`
	SenderID uint   `json:"sender_id"`
	Content  string `json:"content"`
}

func NewMessageHandler(service *services.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (handler *MessageHandler) SendMessage(c *gin.Context) {
	message := SendMessageRequest{}
	err := c.ShouldBindJSON(&message)
	if err != nil {
		c.JSON(400, gin.H{
			"error : ": err.Error(),
		})
		return
	}
	if message.Content == "" {
		c.JSON(400, gin.H{
			"error : ": "Content Should not be empty",
		})
		return
	}
	newMessage, err := handler.service.SendMessage(message.ChatID, message.SenderID, message.Content)
	if err != nil {
		c.JSON(400, gin.H{
			"error : ": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{"message : ": newMessage})
}

func (handler *MessageHandler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	id, err := strconv.Atoi(chatID)
	if err != nil {
		c.JSON(400, gin.H{
			"error : ": err.Error(),
		})
		return
	}
	messagesResponse, err := handler.service.GetMessagesByChat(uint(id))
	if err != nil {
		c.JSON(500, gin.H{
			"error : ": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"messages : ": messagesResponse})
}
