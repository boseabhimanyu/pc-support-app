package handlers

import (
	"net/http"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         config.Config
}

func NewAuthHandler(authService *services.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// CreateUser handles adding a new user
func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.authService.Register(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, err := auth.GenerateToken(
		user,
		h.cfg.JWTSecret,
		h.cfg.JWTExpiryHours,
	)

	refreshToken, expiresAt, err := auth.GenerateRefreshToken(
		user,
		h.cfg.JWTSecret,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate refresh token",
		})
		return
	}

	err = h.authService.UpdateRefreshToken(
		c.Request.Context(),
		user.ID,
		refreshToken,
		expiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save refresh token",
		})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		h.cfg.JWTExpiryHours*60*60,
		"/",
		"",
		false, // Secure = true in production (HTTPS)
		true,  // HttpOnly
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		30*24*60*60, // 30 days
		"/",
		"",
		false, // Secure = true in production
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {

	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDHex := c.GetString("userID")

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := h.userService.UpdateProfile(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token missing",
		})
		return
	}

	// Validate JWT
	_, err = auth.ValidateRefreshToken(
		refreshToken,
		h.cfg.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	// Find user by stored refresh token
	userID, err := auth.ValidateRefreshToken(
		refreshToken,
		h.cfg.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	user, err := h.authService.FindByID(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	// Check expiry stored in DB
	if user.RefreshTokenExpiresAt == nil ||
		time.Now().After(*user.RefreshTokenExpiresAt) {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token expired",
		})
		return
	}

	// Generate new access token
	accessToken, err := auth.GenerateToken(
		user,
		h.cfg.JWTSecret,
		h.cfg.JWTExpiryHours,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	// Rotate refresh token
	newRefreshToken, expiresAt, err := auth.GenerateRefreshToken(
		user,
		h.cfg.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate refresh token",
		})
		return
	}

	err = h.authService.UpdateRefreshToken(
		c.Request.Context(),
		user.ID,
		newRefreshToken,
		expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update refresh token",
		})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		h.cfg.JWTExpiryHours*60*60,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		30*24*60*60,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed",
	})
}
