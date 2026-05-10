// Manages the providers map to be passed into the handlers
// This is done here to provide a separation of concerns as the methods here are to just serve the handler methods
// so we dont need to make them there

package services

import "errors"

type AuthManager struct {
	Providers map[string]AuthProvider
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		Providers: make(map[string]AuthProvider),
	}
}

func (auth_man *AuthManager) RegisterProvider(name string, p AuthProvider) {
	auth_man.Providers[name] = p
}

func (auth_man *AuthManager) GetProvider(name string) (AuthProvider, error) {
	provider, ok := auth_man.Providers[name]
	if !ok {
		return nil, errors.New("provider not supported")
	}
	return provider, nil
}
