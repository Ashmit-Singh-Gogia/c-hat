package database

import (
	"fmt"
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
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	if err := DB.AutoMigrate(&models.ChatParticipant{}); err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	if err := DB.AutoMigrate(&models.Chat{}); err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	if err := DB.AutoMigrate(&models.Message{}); err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
}
