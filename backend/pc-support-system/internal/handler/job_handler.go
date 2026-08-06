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

func (h *JobHandler) GetCustomerJobs(c *gin.Context) {

	userID := c.GetString("userID")

	jobs, err := h.jobService.GetCustomerJobs(
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

func (h *JobHandler) ChangeJobStatus(c *gin.Context) {

	jobID := c.Param("jobId")
	userID := c.GetString("userID")

	var req dto.UpdateJobStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	job, err := h.jobService.ChangeJobStatus(
		c.Request.Context(),
		jobID,
		userID,
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

func (h *JobHandler) GetInProgressJobs(c *gin.Context) {

	jobs, err := h.jobService.GetInProgressJobs(
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

func (h *JobHandler) GetWaitingCustomerJobs(c *gin.Context) {

	jobs, err := h.jobService.GetWaitingCustomerJobs(
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

func (h *JobHandler) GetResumedJobs(c *gin.Context) {

	jobs, err := h.jobService.GetResumedJobs(
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

func (h *JobHandler) AddJobNote(c *gin.Context) {

	jobID := c.Param("jobId")
	userID := c.GetString("userID")

	var req dto.AddJobNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.jobService.AddJobNote(
		c.Request.Context(),
		jobID,
		userID,
		req,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "job note added successfully",
	})
}

func (h *JobHandler) GetJobNotes(c *gin.Context) {

	jobID := c.Param("jobId")
	userID := c.GetString("userID")

	notes, err := h.jobService.GetJobNotes(
		c.Request.Context(),
		jobID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func (h *JobHandler) GetJobByID(c *gin.Context) {

	jobID := c.Param("jobId")
	userID := c.GetString("userID")

	job, err := h.jobService.GetJobByID(
		c.Request.Context(),
		jobID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) GetCustomerJobsByCustomerID(c *gin.Context) {

	customerID := c.Param("customerId")
	userID := c.GetString("userID")

	jobs, err := h.jobService.GetCustomerJobsByCustomerID(
		c.Request.Context(),
		customerID,
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

func (h *JobHandler) CloseJob(c *gin.Context) {

	jobID := c.Param("jobId")

	var req dto.CloseJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	job, err := h.jobService.CloseJob(
		c.Request.Context(),
		jobID,
		c.GetString("userID"),
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

func (h *JobHandler) SearchJobs(c *gin.Context) {

	query := c.Query("q")

	jobs, err := h.jobService.SearchJobs(
		c.Request.Context(),
		query,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(jobs),
		"jobs":  jobs,
	})
}

func (h *JobHandler) GetJobByNumber(c *gin.Context) {

	jobNumber := c.Param("jobNumber")
	userID := c.GetString("userID")

	job, err := h.jobService.GetJobByNumber(
		c.Request.Context(),
		jobNumber,
		userID,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, job)
}
