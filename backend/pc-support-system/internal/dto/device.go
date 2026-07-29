package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type AddDeviceRequest struct {
	CustomerID   string                 `json:"customerId" binding:"required"`
	Type         models.DeviceType      `json:"type" binding:"required"`
	Condition    models.DeviceCondition `json:"condition" binding:"required"`
	Brand        string                 `json:"brand"`
	Model        string                 `json:"model"`
	SerialNumber string                 `json:"serialNumber"`
	Notes        string                 `json:"notes"`
}

type DeviceResponse struct {
	ID string `json:"id"`

	Customer CustomerSummary `json:"customer"`

	Type      models.DeviceType      `json:"type"`
	Condition models.DeviceCondition `json:"condition"`

	Brand        string `json:"brand,omitempty"`
	Model        string `json:"model,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type CustomerSummary struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

func ToDeviceResponse(device *models.Device) DeviceResponse {
	return DeviceResponse{
		ID:           device.ID.Hex(),
		Type:         device.Type,
		Condition:    device.Condition,
		Brand:        device.Brand,
		Model:        device.Model,
		SerialNumber: device.SerialNumber,
		Notes:        device.Notes,
	}
}

func ToDeviceResponses(devices []models.Device) []DeviceResponse {

	resp := make([]DeviceResponse, 0, len(devices))

	for _, device := range devices {
		resp = append(resp, ToDeviceResponse(&device))
	}

	return resp
}

type UpdateDeviceRequest struct {
	Type      *models.DeviceType      `json:"type"`
	Condition *models.DeviceCondition `json:"condition"`

	Brand        *string `json:"brand"`
	Model        *string `json:"model"`
	SerialNumber *string `json:"serialNumber"`
	Notes        *string `json:"notes"`
}

func ToCustomerSummary(user *models.User) CustomerSummary {
	return CustomerSummary{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}
}

func ToDeviceSummary(device *models.Device) DeviceSummary {

	return DeviceSummary{
		ID:           device.ID.Hex(),
		Type:         device.Type,
		Brand:        device.Brand,
		Model:        device.Model,
		SerialNumber: device.SerialNumber,
	}
}

type CustomerDevicesResponse struct {
	DevicesCount int                      `json:"devicesCount"`
	Devices      []CustomerDeviceResponse `json:"devices"`
}

type CustomerDeviceResponse struct {
	ID           string                 `json:"id"`
	Type         models.DeviceType      `json:"type"`
	Condition    models.DeviceCondition `json:"condition"`
	Brand        string                 `json:"brand,omitempty"`
	Model        string                 `json:"model,omitempty"`
	SerialNumber string                 `json:"serialNumber,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
}

func ToCustomerDeviceResponse(device *models.Device) CustomerDeviceResponse {
	return CustomerDeviceResponse{
		ID:           device.ID.Hex(),
		Type:         device.Type,
		Condition:    device.Condition,
		Brand:        device.Brand,
		Model:        device.Model,
		SerialNumber: device.SerialNumber,
		Notes:        device.Notes,
	}
}
