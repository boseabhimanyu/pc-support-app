package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterAuditRoutes(rg *gin.RouterGroup, auditHandler *handlers.AuditHandler) {
	audit := rg.Group("/audit-logs")
	{
		audit.GET(
			"",
			auth.RequireRoles(
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			auditHandler.GetAuditLogs,
		)
	}
}
