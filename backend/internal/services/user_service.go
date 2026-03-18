package services

import (
	"errors"

	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(username string) (models.User, error) {
	_, err := s.repo.GetUserByUsername(username)
	if err == nil { // error nil means already user exists
		return models.User{}, errors.New("user already exists")
	}
	// but if any other error except this there is an issue
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}
	user, err := s.repo.CreateUser(username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
