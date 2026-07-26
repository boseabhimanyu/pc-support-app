package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type UserResponse struct {
	ID        string           `json:"id"`
	FirstName string           `json:"firstName"`
	LastName  string           `json:"lastName"`
	Email     string           `json:"email"`
	Phone     string           `json:"phone"`
	Role      models.Role      `json:"role"`
	State     models.UserState `json:"state"`
}

type UpdateProfileRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone" binding:"omitempty,numeric"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

type CustomerSearchResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func ToCustomerSearchResponse(user *models.User) CustomerSearchResponse {
	return CustomerSearchResponse{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
	}
}

func ToCustomerSearchResponses(users []models.User) []CustomerSearchResponse {

	resp := make([]CustomerSearchResponse, 0, len(users))

	for _, user := range users {
		resp = append(resp, ToCustomerSearchResponse(&user))
	}

	return resp
}

type CreateCustomerRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
}

type CustomerResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func ToCustomerResponse(user *models.User) CustomerResponse {
	return CustomerResponse{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
	}
}

type SetCustomerPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateCustomerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}
