package handlers

import (
	"net/http"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) CreateJob(c *gin.Context) {

	var req dto.CreateJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	job, err := h.jobService.CreateJob(
		c.Request.Context(),
		c.GetString("userID"),
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, job)
}
