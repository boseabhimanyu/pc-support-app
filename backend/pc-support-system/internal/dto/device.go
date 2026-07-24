package dto

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"

type AddDeviceRequest struct {
	CustomerID string            `json:"customerId" binding:"required"`
	Type       models.DeviceType `json:"type" binding:"required"`

	Brand        string `json:"brand"`
	Model        string `json:"model"`
	SerialNumber string `json:"serialNumber"`
	Notes        string `json:"notes"`
}

type DeviceResponse struct {
	ID           string            `json:"id"`
	CustomerID   string            `json:"customerId"`
	Type         models.DeviceType `json:"type"`
	Brand        string            `json:"brand,omitempty"`
	Model        string            `json:"model,omitempty"`
	SerialNumber string            `json:"serialNumber,omitempty"`
	Notes        string            `json:"notes,omitempty"`
}

func ToDeviceResponse(device *models.Device) DeviceResponse {
	return DeviceResponse{
		ID:           device.ID.Hex(),
		CustomerID:   device.CustomerID.Hex(),
		Type:         device.Type,
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
