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
	Uid1 uint `json:"uid1"`
	Uid2 uint `json:"uid2"`
}

func (handler *ChatHandler) CreateDirectChat(c *gin.Context) {
	chatStr := chatStruct{}
	err := c.ShouldBindJSON(&chatStr)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	chat, err := handler.service.CreateDirectChat(chatStr.Uid1, chatStr.Uid2)
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
