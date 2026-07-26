package router

import (
	"net/http"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database, cfg config.Config) *gin.Engine {
	r := gin.Default()

	// 1. Global/Public Endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":     true,
			"status": "healthy pc service",
		})
	})

	userRepo := repository.NewUserRepository(database)

	authService := services.NewAuthService(userRepo)

	authHandler := handlers.NewAuthHandler(authService, cfg)

	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	deviceRepo := repository.NewDeviceRepository(database)

	deviceService := services.NewDeviceService(
		deviceRepo,
		userRepo,
	)

	deviceHandler := handlers.NewDeviceHandler(deviceService)

	jobRepo := repository.NewJobRepository(database)

	jobService := services.NewJobService(
		jobRepo,
		userRepo,
		deviceRepo,
	)

	jobHandler := handlers.NewJobHandler(jobService)
	// API Version
	api := r.Group("/api/v1")
	{
		// Authentication Routes
		open := api.Group("/auth")
		{
			open.POST("/register", authHandler.Register)
			open.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/me", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"userID": c.GetString("userID"),
					"role":   c.GetString("role"),
				})
			})
			protected.GET("/users/me", userHandler.GetProfile)
			protected.PATCH("/users/me", userHandler.UpdateProfile)
			protected.POST("/logout", authHandler.Logout)
			protected.POST(
				"/devices",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				deviceHandler.AddDevice,
			)
			protected.GET(
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
			protected.GET(
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
			protected.GET(
				"/customers-search",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleHeadTechnician),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				userHandler.SearchCustomers,
			)
			protected.POST(
				"/customers",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleHeadTechnician),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				userHandler.CreateCustomer,
			)
			protected.PATCH(
				"/customers/:customerId/password",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				userHandler.SetCustomerPassword,
			)
			protected.POST(
				"/jobs",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleHeadTechnician),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				jobHandler.CreateJob,
			)
			protected.PATCH(
				"/customers/:customerId",
				auth.RequireRoles(
					string(models.RoleReceptionist),
					string(models.RoleAdmin),
					string(models.RoleSuperAdmin),
				),
				userHandler.UpdateCustomer,
			)
			// protected.PATCH(
			// 	"PATCH /customers/:customerId/state", auth.RequireRoles(
			// 		string(models.RoleReceptionist),
			// 		string(models.RoleAdmin),
			// 		string(models.RoleSuperAdmin),
			// 	),
			// 	userHandler.CustomerState,
			// )
		}
	}

	return r
}
