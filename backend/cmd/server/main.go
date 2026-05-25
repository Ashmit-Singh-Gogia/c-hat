package main

import (
	"fmt"
	"net/http"

	"github.com/ashmit-singh-gogia/c-hat/internal/config"
	"github.com/ashmit-singh-gogia/c-hat/internal/database"
	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/ashmit-singh-gogia/c-hat/internal/providers"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
	"github.com/ashmit-singh-gogia/c-hat/internal/routes"
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/ashmit-singh-gogia/c-hat/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func main() {
	router := gin.Default()

	Cfg := config.LoadConfig() // the env vars are initialized inside this method
	config.InitOAuth(Cfg)      // the oauth config is initialized inside this method

	// --- NEW CHROME COOKIE FIX START ---
	// Create a custom session store for Goth and relax the SameSite rules
	store := sessions.NewCookieStore([]byte(Cfg.JWT_SECRET))
	store.MaxAge(86400)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false                  // Must be false for localhost HTTP
	store.Options.SameSite = http.SameSiteLaxMode // The magic flag for Chrome!
	gothic.Store = store
	// --- NEW CHROME COOKIE FIX END ---

	database.ConnectDB(Cfg)
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	chatRepo := repositories.NewChatRepository(database.DB)
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	messageRepo := repositories.NewMessageRepository(database.DB)
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService)

	hub := websocket.NewHub()
	go hub.Run() // runs forever in background

	wsHandler := handlers.NewWSHandler(hub, messageService)

	authService := services.NewAuthService(userRepo)
	authManager := services.NewAuthManager()
	googleProvider := providers.NewGoogleProvider()
	authManager.RegisterProvider("google", googleProvider)
	authHandler := handlers.NewAuthHandler(authService, authManager, Cfg)

	routes.LoadRoutes(router, userHandler, chatHandler, messageHandler, authHandler, wsHandler)
	fmt.Println("Server running on port", Cfg.PORT)
	if err := router.Run(":" + Cfg.PORT); err != nil {
		panic(err)
	}
}
