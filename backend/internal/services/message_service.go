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
		return models.Message{}, errors.New("message must not be empty")
	}
	message, err := service.repo.CreateMessage(chatId, senderId, content)
	if err != nil {
		return models.Message{}, err
	}

	return message, nil
}

func (service *MessageService) GetMessagesByChat(chatID, userID uint) ([]models.Message, error) {
	// I have to bring that chat here first
	// Then check the participants of that chat and see if the user is a participant or not
	chat, err := service.repo.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	isParticipant := false
	for _, participant := range chat.Participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant of this chat")
	}
	messages, err := service.repo.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
