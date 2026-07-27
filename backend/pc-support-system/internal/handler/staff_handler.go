package handlers

import (
	"net/http"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) CreateStaff(c *gin.Context) {

	var req dto.CreateStaffRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdBy := c.GetString("userID")

	staff, err := h.userService.CreateStaff(
		c.Request.Context(),
		createdBy,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, staff)
}

func (h *UserHandler) SetStaffPassword(c *gin.Context) {

	staffID := c.Param("staffId")
	updatedBy := c.GetString("userID")

	var req dto.SetStaffPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.userService.SetStaffPassword(
		c.Request.Context(),
		staffID,
		updatedBy,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "staff password set successfully",
	})
}

func (h *UserHandler) SearchStaff(c *gin.Context) {

	query := c.Query("q")

	staff, err := h.userService.SearchStaff(
		c.Request.Context(),
		query,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, staff)
}
