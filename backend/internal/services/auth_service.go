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

func (s *AuthService) FindOrCreateUser(dto *UserDTO) (*models.User, error) {
	user, err := s.userRepo.FindByProvider(dto.Provider, dto.ProviderID)
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// user doesn't exist → create a new one
	user = models.User{
		Provider:   dto.Provider,
		ProviderID: dto.ProviderID,
		Username:   dto.Username,
		Email:      dto.Email,
		Avatar:     dto.Avatar,
		IsVerified: true, // since we trust the provider, we can mark the user as verified
	}
	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, errors.New("failed to create new user")
	}
	return &user, nil
}

func (s *AuthService) ProcessUserLogin(dto *UserDTO) (*models.User, error) {
	return s.FindOrCreateUser(dto)
}
