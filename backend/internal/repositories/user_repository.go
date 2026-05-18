package repositories

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindByProvider(provider string, providerID string) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id uint) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil { // first already searches for primary key no need of where clause
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
