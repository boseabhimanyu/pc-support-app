package handlers

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
