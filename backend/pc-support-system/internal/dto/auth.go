package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type RegisterRequest struct {
	FirstName string `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string `json:"lastName" binding:"omitempty,max=50"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	Phone     string `json:"phone" binding:"required,numeric,len=10"`
	Password  string `json:"password" binding:"required,min=8,max=64"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		State:     user.State,
	}
}
