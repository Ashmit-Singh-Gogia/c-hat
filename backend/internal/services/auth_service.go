package services

import (
	"errors"

	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository // ← holds the repo instance
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo} // ← initialize the service with the repo
}

func (s *AuthService) FindOrCreateUser(googleID, email, name, avatar string) (*models.User, error) {
	user, err := s.userRepo.FindByGoogleId(googleID)
	if err == nil {
		return &user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		createdUser, err := s.userRepo.CreateUser(name)
		if err != nil {
			return nil, err
		}
		createdUser.GoogleID = googleID
		createdUser.Email = email
		createdUser.Avatar = avatar
		if err := s.userRepo.DB.Save(&createdUser).Error; err != nil {
			return nil, err
		}
		return &createdUser, nil
	}
	return nil, errors.New("failed to find or create user")
}
