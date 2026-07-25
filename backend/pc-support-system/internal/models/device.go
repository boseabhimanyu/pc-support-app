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

type DeviceState string

const (
	DeviceActive   DeviceState = "active"
	DeviceInactive DeviceState = "inactive"
)

type Device struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID bson.ObjectID `bson:"customer_id" json:"customerId"`

	Type DeviceType `bson:"type" json:"type"`

	Brand        string          `bson:"brand,omitempty" json:"brand,omitempty"`
	Model        string          `bson:"model,omitempty" json:"model,omitempty"`
	SerialNumber string          `bson:"serial_number,omitempty" json:"serial_number,omitempty"`
	Notes        string          `bson:"notes,omitempty" json:"notes,omitempty"`
	Condition    DeviceCondition `bson:"condition" json:"condition"`
	CreatedAt    time.Time       `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time       `bson:"updated_at" json:"updatedAt"`
}

type Component struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID bson.ObjectID `bson:"device_id" json:"deviceId"`

	Type         ComponentType `bson:"type" json:"type"`
	Manufacturer string        `bson:"manufacturer,omitempty" json:"manufacturer,omitempty"`
	Model        string        `bson:"model,omitempty" json:"model,omitempty"`
	SerialNumber string        `bson:"serial_number,omitempty" json:"serialNumber,omitempty"`

	InstalledAt time.Time `bson:"installed_at" json:"installedAt"`

	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
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

type DeviceCondition string

const (
	DeviceWorking          DeviceCondition = "working"
	DeviceNotWorking       DeviceCondition = "not_working"
	DevicePartiallyWorking DeviceCondition = "partially_working"
	DeviceUnknown          DeviceCondition = "unknown"
)

func (d DeviceType) IsValid() bool {
	switch d {
	case DeviceLaptop,
		DeviceDesktop,
		DevicePrinter,
		DeviceMonitor,
		DeviceUPS,
		DeviceRouter:
		return true
	default:
		return false
	}
}

func (d DeviceCondition) IsValid() bool {
	switch d {
	case DeviceWorking,
		DeviceNotWorking,
		DevicePartiallyWorking,
		DeviceUnknown:
		return true
	default:
		return false
	}
}
