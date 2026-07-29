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

func (h *JobHandler) GetOpenJobs(c *gin.Context) {

	jobs, err := h.jobService.GetOpenJobs(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) AssignJob(c *gin.Context) {

	jobID := c.Param("jobId")

	var req dto.AssignJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	assignedBy := c.GetString("userID")

	job, err := h.jobService.AssignJob(
		c.Request.Context(),
		jobID,
		assignedBy,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) GetAssignedJobs(c *gin.Context) {

	jobs, err := h.jobService.GetAssignedJobs(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetMyJobs(c *gin.Context) {

	userID := c.GetString("userID")

	jobs, err := h.jobService.GetMyJobs(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jobs)
}
