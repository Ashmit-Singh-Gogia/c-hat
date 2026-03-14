package repositories

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		DB: db,
	}
}

func (repo *MessageRepository) CreateMessage(chatId, senderID uint, content string) (models.Message, error) {
	message := models.Message{ChatID: chatId, SenderID: senderID, Content: content}
	result := repo.DB.Create(&message)
	if result.Error != nil {
		return models.Message{}, result.Error
	}
	return message, nil
}

func (repo *MessageRepository) GetMessagesByChatID(chatID uint) ([]models.Message, error) {

	var messages []models.Message

	result := repo.DB.
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}
