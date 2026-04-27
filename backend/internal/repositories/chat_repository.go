package repositories

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		DB: db,
	}
}

func (c *ChatRepository) CreateChat(tx *gorm.DB, isgroup bool) (models.Chat, error) {
	chat := models.Chat{
		IsGroup: isgroup,
	}
	result := tx.Create(&chat)
	if result.Error != nil {
		return models.Chat{}, result.Error
	}
	return chat, nil
}

func (c *ChatRepository) FindChatsByUserID(userID uint) ([]models.Chat, error) {
	var chats []models.Chat

	// Find all chats where the user's ID is in the chat_participants table
	err := c.DB.Model(&models.Chat{}).
		Joins("JOIN chat_participants ON chat_participants.chat_id = chats.id").
		Where("chat_participants.user_id = ?", userID).
		Preload("Participants.User"). // Preloads the participants array
		Find(&chats).Error

	return chats, err
}

func (c *ChatRepository) AddParticipants(tx *gorm.DB, chatID uint, userIDs []uint) error {
	var participants []models.ChatParticipant
	for _, uid := range userIDs {
		participant := models.ChatParticipant{
			ChatID: chatID,
			UserID: uid,
		}
		participants = append(participants, participant)
	}
	result := tx.Create(&participants)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
