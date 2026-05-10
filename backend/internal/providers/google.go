package providers

import (
	"context"

	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type GoogleProvider struct {
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}

func (g *GoogleProvider) Login(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (g *GoogleProvider) Callback(c *gin.Context) (*services.UserDTO, error) {
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		return nil, err
	}
	return &services.UserDTO{
		Provider:   "google",
		ProviderID: user.UserID,
		Email:      user.Email,
		Username:   user.Name,
		Avatar:     user.AvatarURL,
	}, nil
}

func (g *GoogleProvider) Logout(c *gin.Context) error {
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)
	return gothic.Logout(c.Writer, c.Request)
}
