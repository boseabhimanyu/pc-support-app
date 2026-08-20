package handlers

import (
	"net/http"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/auth"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/dto"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/services"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/utils"
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

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

	// Hash the refresh token before storing it in the database.
	refreshTokenHash := utils.HashRefreshToken(refreshToken)

	err = h.authService.UpdateRefreshToken(
		c.Request.Context(),
		user.ID,
		refreshTokenHash,
		expiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save refresh token",
		})
		return
	}

	// Apply SameSite policy to cookies created on this context.
	c.SetSameSite(h.cfg.CookieSameSite)

	c.SetCookie(
		"access_token",
		accessToken,
		h.cfg.JWTExpiryHours*60*60,
		"/",
		"",
		h.cfg.CookieSecure,
		true, // HttpOnly
	)

	// Keep the RAW refresh token in the browser cookie.
	// Only its hash is stored in the database.
	c.SetCookie(
		"refresh_token",
		refreshToken,
		h.cfg.RefreshTokenExpiryDays*24*60*60,
		"/",
		"",
		h.cfg.CookieSecure,
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetSameSite(h.cfg.CookieSameSite)

	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		h.cfg.CookieSecure,
		true,
	)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		h.cfg.CookieSecure,
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

	// Validate the JWT and extract the user ID.
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

	// Find the user associated with the refresh token.
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

	// Compare the raw token from the browser
	// against the hashed token stored in the database.
	if user.CurrentRefreshToken == "" ||
		!utils.VerifyRefreshToken(
			refreshToken,
			user.CurrentRefreshToken,
		) {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	// Check expiry stored in the database.
	if user.RefreshTokenExpiresAt == nil ||
		time.Now().After(*user.RefreshTokenExpiresAt) {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token expired",
		})
		return
	}

	// Generate a new access token.
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

	// Rotate refresh token.
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

	// Store only the hash in the database.
	newRefreshTokenHash := utils.HashRefreshToken(newRefreshToken)

	err = h.authService.UpdateRefreshToken(
		c.Request.Context(),
		user.ID,
		newRefreshTokenHash,
		expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update refresh token",
		})
		return
	}

	c.SetSameSite(h.cfg.CookieSameSite)

	// Send the raw access token to the browser.
	c.SetCookie(
		"access_token",
		accessToken,
		h.cfg.JWTExpiryHours*60*60,
		"/",
		"",
		h.cfg.CookieSecure,
		true,
	)

	// Send the raw refresh token to the browser.
	// The database contains only its hash.
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		h.cfg.RefreshTokenExpiryDays*24*60*60,
		"/",
		"",
		h.cfg.CookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed",
	})
}
