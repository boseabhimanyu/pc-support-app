package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RoleSuperAdmin     Role = "super_admin"
	RoleAdmin          Role = "admin"
	RoleTechnician     Role = "technician"
	RoleHeadTechnician Role = "head_technician"
	RoleReceptionist   Role = "receptionist"
	RoleCustomer       Role = "customer"
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	// Basic Details
	FirstName string `bson:"first_name" json:"firstName"`
	LastName  string `bson:"last_name" json:"lastName"`

	Email string `bson:"email" json:"email"`
	Phone string `bson:"phone" json:"phone"`

	PasswordHash string `bson:"password_hash" json:"-"`

	Role Role `bson:"role" json:"role"`

	// Employee / Customer State
	State UserState `bson:"state" json:"state"`

	// Audit
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`

	CreatedByID           *bson.ObjectID `bson:"created_by_id,omitempty" json:"-"`
	UpdatedByID           *bson.ObjectID `bson:"updated_by_id,omitempty" json:"-"`
	CurrentRefreshToken   string         `bson:"current_refresh_token,omitempty" json:"-"`
	RefreshTokenExpiresAt *time.Time     `bson:"refresh_token_expires_at,omitempty" json:"-"`
}

func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin,
		RoleAdmin,
		RoleTechnician,
		RoleHeadTechnician,
		RoleReceptionist,
		RoleCustomer:
		return true
	default:
		return false
	}
}

type CreateCustomerRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
}

type UserState string

const (
	UserActive   UserState = "active"
	UserInactive UserState = "inactive"
)

func (s UserState) IsValid() bool {
	switch s {
	case UserActive, UserInactive:
		return true
	default:
		return false
	}
}
