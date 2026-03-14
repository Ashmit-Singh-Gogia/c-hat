package database

import (
	"log"

	"github.com/ashmit-singh-gogia/c-hat/internal/config"
	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DATABASE_URL,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error : %q", err)
	}
	DB = db
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.ChatParticipant{})
	DB.AutoMigrate(&models.Chat{})
	DB.AutoMigrate(&models.Message{})
}
