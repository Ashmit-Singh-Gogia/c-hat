package services

import (
	"errors"

	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
)

type ChatService struct {
	repo *repositories.ChatRepository
}

func NewChatService(repo *repositories.ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}
func (c *ChatService) CreateDirectChat(user1ID uint, user2ID uint) (models.Chat, error) {
	if user1ID == 0 {
		return models.Chat{}, errors.New("id is 0 for user 1")
	}
	if user2ID == 0 {
		return models.Chat{}, errors.New("id is 0 for user 2")
	}
	if user1ID == user2ID {
		return models.Chat{}, errors.New("cannot create chat with yourself")
	}

	tx := c.repo.DB.Begin()
	chat, err := c.repo.CreateChat(tx, false)
	if err != nil {
		tx.Rollback()
		return models.Chat{}, err
	}

	userIDs := []uint{user1ID, user2ID}
	err = c.repo.AddParticipants(tx, chat.ID, userIDs)
	if err != nil {
		tx.Rollback()
		return models.Chat{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.Chat{}, err
	}

	// Reload chat with participants after commit
	if err := c.repo.DB.Preload("Participants").First(&chat, chat.ID).Error; err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

// Add this inside chat_servicce.go

func (c *ChatService) GetUserChats(userID uint) ([]models.Chat, error) {
	return c.repo.FindChatsByUserID(userID)
}
