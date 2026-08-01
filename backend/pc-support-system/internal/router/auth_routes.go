package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, cfg config.Config) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		// authGroup.POST("/register", middleware.RateLimit(), authHandler.Register)
		// authGroup.POST("/login", middleware.RateLimit(), authHandler.Login)
	}
}
