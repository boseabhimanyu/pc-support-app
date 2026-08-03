package router

import (
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	handlers "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterJobRoutes(rg *gin.RouterGroup, jobHandler *handlers.JobHandler) {
	jobs := rg.Group("/jobs")
	{
		jobs.GET(
			"/customer/:customerId/",
			jobHandler.GetCustomerJobsByCustomerID,
		)
		jobs.POST(
			"",
			auth.RequireRoles(
				string(models.RoleReceptionist),
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
			),
			jobHandler.CreateJob,
		)
		jobs.GET(
			"/open",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			jobHandler.GetOpenJobs,
		)
		jobs.GET(
			"/in-progress",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			jobHandler.GetInProgressJobs,
		)
		jobs.GET(
			"/waiting-customer",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			jobHandler.GetWaitingCustomerJobs,
		)
		jobs.GET(
			"/resumed",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			jobHandler.GetResumedJobs,
		)
		jobs.PATCH(
			"/:jobId/assign",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
			),
			jobHandler.AssignJob,
		)
		jobs.GET(
			"/assigned",
			auth.RequireRoles(
				string(models.RoleHeadTechnician),
				string(models.RoleAdmin),
				string(models.RoleSuperAdmin),
			),
			jobHandler.GetAssignedJobs,
		)
		jobs.GET(
			"/my",
			auth.RequireRoles(
				string(models.RoleTechnician),
				string(models.RoleHeadTechnician),
			),
			jobHandler.GetMyJobs,
		)
		jobs.PATCH(
			"/:jobId/status",
			auth.RequireRoles(
				string(models.RoleTechnician),
				string(models.RoleHeadTechnician),
			),
			jobHandler.ChangeJobStatus,
		)
		jobs.POST(
			"/:jobId/notes",
			auth.RequireRoles(
				string(models.RoleReceptionist),
				string(models.RoleTechnician),
				string(models.RoleHeadTechnician),
			),
			jobHandler.AddJobNote,
		)
		jobs.GET("/:jobId/notes", jobHandler.GetJobNotes)
		jobs.GET("/:jobId", jobHandler.GetJobByID)
		jobs.POST("/:jobId/close", jobHandler.CloseJob)
		jobs.GET("/search", jobHandler.SearchJobs)
	}
}
