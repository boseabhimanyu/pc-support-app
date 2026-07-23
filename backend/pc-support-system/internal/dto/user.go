package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type UserResponse struct {
	ID        string      `json:"id"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	Role      models.Role `json:"role"`
	Active    bool        `json:"active"`
}
