package services

import (
	"errors"

	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
)

type MessageService struct {
	repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (service *MessageService) SendMessage(chatId uint, senderId uint, content string) (models.Message, error) {
	if content == "" {
		return models.Message{}, errors.New("Message must not be empty")
	}
	message, err := service.repo.CreateMessage(chatId, senderId, content)
	if err != nil {
		return models.Message{}, err
	}

	return message, nil
}

func (service *MessageService) GetMessagesByChat(chatID uint) ([]models.Message, error) {
	messages, err := service.repo.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
