package router

import (
	"net/http"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database, cfg config.Config) *gin.Engine {
	r := gin.Default()

	// Global/Public Endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":     true,
			"status": "healthy pc service",
		})
	})

	// Repositories & Services Wiring
	userRepo := repository.NewUserRepository(database)
	auditRepo := repository.NewAuditRepository(database)
	auditService := services.NewAuditService(auditRepo)
	auditHandler := handlers.NewAuditHandler(auditService)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService, cfg)

	userService := services.NewUserService(userRepo, auditService)
	userHandler := handlers.NewUserHandler(userService)

	deviceRepo := repository.NewDeviceRepository(database)
	deviceService := services.NewDeviceService(deviceRepo, userRepo, auditService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	jobRepo := repository.NewJobRepository(database)
	jobService := services.NewJobService(jobRepo, userRepo, deviceRepo, auditService)
	jobHandler := handlers.NewJobHandler(jobService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Version Group
	api := r.Group("/api/v1")
	{
		// Public Auth Routes
		RegisterAuthRoutes(api, authHandler, cfg)

		// Protected Routes Group
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(cfg.JWTSecret))
		{
			// Base protected endpoints (/me, /logout, /refresh)
			protected.GET("/me", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"userID": c.GetString("userID"),
					"role":   c.GetString("role"),
				})
			})
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/refresh", authHandler.RefreshToken)

			// Modular Feature Routes
			RegisterUserRoutes(protected, userHandler)
			RegisterDeviceRoutes(protected, deviceHandler)
			RegisterJobRoutes(protected, jobHandler)
			RegisterAuditRoutes(protected, auditHandler)
			RegisterMeRoutes(protected, jobHandler, deviceHandler)
		}
	}

	return r
}
