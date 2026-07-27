package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type CreateStaffRequest struct {
	FirstName string      `json:"firstName" binding:"required"`
	LastName  string      `json:"lastName" binding:"required"`
	Phone     string      `json:"phone" binding:"required"`
	Email     string      `json:"email" binding:"required"`
	Role      models.Role `json:"role" binding:"required"`
}

type UpdateStaffRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
}

type SetStaffPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type ChangeStaffRoleRequest struct {
	Role models.Role `json:"role" binding:"required"`
}

type ChangeStaffStateRequest struct {
	State models.UserState `json:"state" binding:"required"`
}

type StaffSummary struct {
	ID        string           `json:"id"`
	FirstName string           `json:"firstName"`
	LastName  string           `json:"lastName"`
	Email     string           `json:"email"`
	Phone     string           `json:"phone"`
	Role      models.Role      `json:"role"`
	State     models.UserState `json:"state"`
}
