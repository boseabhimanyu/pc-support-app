package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RoleSuperAdmin   Role = "super_admin"
	RoleAdmin        Role = "admin"
	RoleTechnician   Role = "technician"
	RoleReceptionist Role = "receptionist"
	RoleCustomer     Role = "customer"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string        `bson:"email" json:"email"`
	Phone        string        `bson:"phone" json:"phone"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	FirstName    string        `bson:"first_name" json:"firstName"`
	LastName     string        `bson:"last_name" json:"lastName"`

	Role   Role `bson:"role" json:"role"`
	Active bool `bson:"active" json:"active"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin,
		RoleAdmin,
		RoleTechnician,
		RoleReceptionist,
		RoleCustomer:
		return true
	default:
		return false
	}
}
