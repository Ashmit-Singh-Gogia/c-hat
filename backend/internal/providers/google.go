package providers

import (
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
	c.Request = gothic.GetContextWithProvider(c.Request, "google")
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (g *GoogleProvider) Callback(c *gin.Context) (*services.UserDTO, error) {
	c.Request = gothic.GetContextWithProvider(c.Request, "google")
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
	c.Request = gothic.GetContextWithProvider(c.Request, "google")
	return gothic.Logout(c.Writer, c.Request)
}
