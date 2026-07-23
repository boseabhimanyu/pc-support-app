package router

import (
	"net/http"
	//"your-project/handlers" // Import your handlers package

	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database) *gin.Engine {
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

	authHandler := handlers.NewAuthHandler(authService)

	// API Routes
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
	}

	return r
}
