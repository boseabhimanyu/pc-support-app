package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterDeviceRoutes(rg *gin.RouterGroup, deviceHandler *handlers.DeviceHandler) {
	rg.POST(
		"/devices",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleAdmin),
			string(models.RoleHeadTechnician),
		),
		deviceHandler.AddDevice,
	)
	rg.GET(
		"/customers/:customerId/devices",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleTechnician),
			string(models.RoleHeadTechnician),
			string(models.RoleAdmin),
			string(models.RoleSuperAdmin),
		),
		deviceHandler.GetCustomerDevices,
	)
	rg.GET(
		"/devices/:deviceId",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleTechnician),
			string(models.RoleHeadTechnician),
			string(models.RoleAdmin),
			string(models.RoleSuperAdmin),
		),
		deviceHandler.GetDevice,
	)
	rg.PATCH(
		"/devices/:deviceId",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleHeadTechnician),
			string(models.RoleAdmin),
		),
		deviceHandler.UpdateDevice,
	)
}
