package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceType string

const (
	DeviceLaptop  DeviceType = "laptop"
	DeviceDesktop DeviceType = "desktop"
	DevicePrinter DeviceType = "printer"
	DeviceMonitor DeviceType = "monitor"
	DeviceUPS     DeviceType = "ups"
	DeviceRouter  DeviceType = "router"
)

type DeviceStatus string

const (
	StatusActive    DeviceStatus = "active"
	StatusInRepair  DeviceStatus = "in_repair"
	StatusDelivered DeviceStatus = "delivered"
	StatusRetired   DeviceStatus = "retired"
)

type Device struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID bson.ObjectID `bson:"customer_id" json:"customerId"`

	Type   DeviceType   `bson:"type" json:"type"`
	Status DeviceStatus `bson:"status" json:"status"`

	Brand string `bson:"brand,omitempty" json:"brand,omitempty"`
	Model string `bson:"model,omitempty" json:"model,omitempty"`

	Notes string `bson:"notes,omitempty" json:"notes,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type Component struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	DeviceID bson.ObjectID `bson:"device_id"`

	Type         ComponentType `bson:"type"`
	Manufacturer string        `bson:"manufacturer,omitempty"`
	Model        string        `bson:"model,omitempty"`
	SerialNumber string        `bson:"serial_number,omitempty"`

	InstalledAt time.Time `bson:"installed_at"`
}

type ComponentType string

const (
	ComponentMotherboard ComponentType = "motherboard"
	ComponentRAM         ComponentType = "ram"
	ComponentSSD         ComponentType = "ssd"
	ComponentHDD         ComponentType = "hdd"
	ComponentPSU         ComponentType = "psu"
	ComponentGPU         ComponentType = "gpu"
	ComponentOptical     ComponentType = "optical_drive"
)
