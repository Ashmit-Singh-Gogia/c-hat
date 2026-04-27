package handlers

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

type chatStruct struct {
	Uid uint `json:"uid"`
}

func (handler *ChatHandler) CreateDirectChat(c *gin.Context) {
	userId, ok := c.Get("userID") // main users id, from the token
	if !ok {
		c.JSON(401, gin.H{
			"Error": "Unauthorized",
		})
		return
	}
	id, ok := userId.(uint)
	if !ok {
		c.JSON(400, gin.H{
			"Error": "Invalid user ID",
		})
		return
	}

	chatStr := chatStruct{} // to get the 2nd user's id
	err := c.ShouldBindJSON(&chatStr)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	chat, err := handler.service.CreateDirectChat(id, chatStr.Uid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"chat created": chat,
		"message":      "success",
	})

}

// Add this inside chat_handler.go

func (handler *ChatHandler) GetChats(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	chats, err := handler.service.GetUserChats(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the raw array so the React frontend maps it correctly
	c.JSON(200, chats)
}
