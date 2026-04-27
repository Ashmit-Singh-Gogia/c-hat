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
	ChatID  uint   `json:"chat_id"`
	Content string `json:"content"`
}

func NewMessageHandler(service *services.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (handler *MessageHandler) SendMessage(c *gin.Context) {
	senderID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id, ok := senderID.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user id type"})
		return
	}

	message := SendMessageRequest{}
	err := c.ShouldBindJSON(&message)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if message.Content == "" {
		c.JSON(400, gin.H{
			"error": "Content Should not be empty",
		})
		return
	}
	newMessage, err := handler.service.SendMessage(message.ChatID, id, message.Content)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{"message": newMessage})
}

func (handler *MessageHandler) GetMessages(c *gin.Context) {

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userid, ok := userID.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user id type"})
		return
	}
	chatID := c.Param("id")
	id, err := strconv.Atoi(chatID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	messagesResponse, err := handler.service.GetMessagesByChat(uint(id), userid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"messages": messagesResponse})
}
