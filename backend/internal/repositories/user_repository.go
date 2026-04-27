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

func (r *UserRepository) CreateUser(username string) (models.User, error) {
	user := models.User{Username: username}
	result := r.DB.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *UserRepository) FindByGoogleId(googleID string) (models.User, error) {
	var user = models.User{}
	if err := r.DB.Where("google_id = ?", googleID).First(&user).Error; err != nil {
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
