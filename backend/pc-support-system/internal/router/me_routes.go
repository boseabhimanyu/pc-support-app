package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterMeRoutes(
	protected *gin.RouterGroup,
	jobHandler *handlers.JobHandler,
	deviceHandler *handlers.DeviceHandler,
) {
	me := protected.Group("/me")
	{
		me.GET(
			"/jobs",
			auth.RequireRoles(string(models.RoleCustomer)),
			jobHandler.GetCustomerJobs,
		)

		me.GET(
			"/devices",
			auth.RequireRoles(string(models.RoleCustomer)),
			deviceHandler.GetMyDevices,
		)
	}
}
