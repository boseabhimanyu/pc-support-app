package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, userHandler *handlers.UserHandler) {
	// Current user profile
	rg.GET("/users/me", userHandler.GetProfile)
	rg.PATCH("/users/me", userHandler.UpdateProfile)

	// Customer management
	rg.GET(
		"/customers-search",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleHeadTechnician),
			string(models.RoleAdmin),
			string(models.RoleSuperAdmin),
		),
		userHandler.SearchCustomers,
	)
	rg.POST(
		"/customers",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleHeadTechnician),
			string(models.RoleAdmin),
		),
		userHandler.CreateCustomer,
	)
	rg.PATCH(
		"/customers/:customerId/password",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleAdmin),
		),
		userHandler.SetCustomerPassword,
	)
	rg.PATCH(
		"/customers/:customerId",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleAdmin),
		),
		userHandler.UpdateCustomer,
	)
	rg.GET("customers/:customerId",
		auth.RequireRoles(
			string(models.RoleReceptionist),
			string(models.RoleAdmin),
			string(models.RoleHeadTechnician),
			string(models.RoleSuperAdmin),
		), userHandler.GetCustomerByID,
	)

	// Staff management group
	staff := rg.Group("/staff")
	{
		staff.POST(
			"",
			auth.RequireRoles(
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			userHandler.CreateStaff,
		)
		staff.PATCH(
			"/:staffId/password",
			auth.RequireRoles(
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			userHandler.SetStaffPassword,
		)
		staff.GET(
			"/search",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			userHandler.SearchStaff,
		)
		staff.PATCH(
			"/:staffId",
			auth.RequireRoles(
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			userHandler.UpdateStaff,
		)
		staff.GET(
			"/:staffId",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			userHandler.GetStaffByID,
		)
	}
}
