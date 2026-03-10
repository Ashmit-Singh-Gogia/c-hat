package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DATABASE_URL string
	JWT_SECRET   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error : %q", err)
	}
	config := &Config{
		PORT:         os.Getenv("PORT"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
	}
	return config
}
