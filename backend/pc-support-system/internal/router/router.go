package router

import (
	"net/http"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database, cfg config.Config) *gin.Engine {
	r := gin.Default()

	// -------------------------------------------------------------
	// CORS Configuration
	// -------------------------------------------------------------
	r.Use(cors.New(cors.Config{
		// Allow your Vite development server origin
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},

		// Explicitly allow HTTP methods used by your frontend
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// Explicitly allow headers sent by your frontend
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization", // Essential for JWT Bearer tokens
			"X-Requested-With",
		},

		// Headers exposed to the frontend in responses
		ExposeHeaders: []string{"Content-Length", "Authorization"},

		// Allow cookies / HTTP authentication headers if needed
		AllowCredentials: true,

		// Cache preflight OPTIONS request responses for 12 hours
		MaxAge: 12 * time.Hour,
	}))

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
	jobRepo := repository.NewJobRepository(database)
	deviceRepo := repository.NewDeviceRepository(database)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService, cfg)

	userService := services.NewUserService(userRepo, deviceRepo, jobRepo, auditService)
	userHandler := handlers.NewUserHandler(userService)

	deviceService := services.NewDeviceService(deviceRepo, userRepo, auditService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	jobService := services.NewJobService(jobRepo, userRepo, deviceRepo, auditService)
	jobHandler := handlers.NewJobHandler(jobService)

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
