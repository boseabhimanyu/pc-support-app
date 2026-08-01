package handlers

import (
	"net/http"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditHandler struct {
	auditService *services.AuditService
}

func NewAuditHandler(
	auditService *services.AuditService,
) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {

	filter := repository.AuditFilter{}

	// entity=user/job/customer/device
	if entity := c.Query("entity"); entity != "" {
		filter.Entity = entity
	}

	// entityId filter
	if entityID := c.Query("entityId"); entityID != "" {

		id, err := bson.ObjectIDFromHex(entityID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid entity id",
			})
			return
		}

		filter.EntityID = &id
	}

	// performedBy filter
	if userID := c.Query("userId"); userID != "" {

		id, err := bson.ObjectIDFromHex(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id",
			})
			return
		}

		filter.PerformedBy = &id
	}

	// action filter
	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	logs, err := h.auditService.FindAuditLogs(
		c.Request.Context(),
		filter,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(logs),
		"logs":  logs,
	})
}
