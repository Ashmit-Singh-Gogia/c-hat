package config

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Config struct {
	PORT                 string
	DATABASE_URL         string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
	JWT_SECRET           string
	SESSION_SECRET       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error : %q", err)
	}
	config := &Config{
		PORT:                 os.Getenv("PORT"),
		DATABASE_URL:         os.Getenv("DATABASE_URL"),
		JWT_SECRET:           os.Getenv("JWT_SECRET"),
		SESSION_SECRET:       os.Getenv("SESSION_SECRET"),
		GOOGLE_CLIENT_ID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_REDIRECT_URL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}
	return config
}

func InitOAuth(cfg *Config) {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	gothic.Store = store

	// Force gothic to bypass its default checks and just use Google
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "google", nil
	}
	goth.UseProviders(
		google.New(cfg.GOOGLE_CLIENT_ID, cfg.GOOGLE_CLIENT_SECRET, cfg.GOOGLE_REDIRECT_URL, "openid", "email", "profile"),
	)
}
