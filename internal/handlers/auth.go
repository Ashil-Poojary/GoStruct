package handlers

import (
	"net/http"
	"time"

	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo         *repository.UserRepository
	RefreshTokenRepo *repository.AuthRepo
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	// RegisterRequest struct for user registration
}
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

func NewAuthHandler(userRepo *repository.UserRepository, rtRepo *repository.AuthRepo) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo, RefreshTokenRepo: rtRepo}
}

// On successful login (simplified example)
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// if !utils.CheckPasswordHash(req.Password, user.Password) {
	// 	utils.Error(c, http.StatusUnauthorized, "invalid credentials")
	// 	return
	// }
	if user.Password != req.Password { // Simplified password check, replace with proper hash check
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(int(user.ID), user.Email, user.Role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(int(user.ID), user.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	// Save refresh token with user agent and IP
	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()
	rtExpiry := time.Now().Add(7 * 24 * time.Hour) // 7 days expiry for refresh token

	if err := h.RefreshTokenRepo.SaveRefreshToken(refreshToken, user.ID, userAgent, ip, rtExpiry); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save refresh token")
		return
	}

	utils.Success(c, "login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600, // access token expiry in seconds
	})
}

// Refresh token endpoint
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "refresh_token required")
		return
	}

	rt, err := h.RefreshTokenRepo.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	//get user details from the refresh token
	user, err := h.UserRepo.GetByID(rt.UserID)
	if err != nil || user == nil {
		utils.Error(c, http.StatusUnauthorized, "user not found")
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(int(user.ID), user.Email, user.Role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	// (Optional) generate a new refresh token and revoke old one
	// For simplicity, just return the same refresh token here.

	utils.Success(c, "token refreshed", gin.H{
		"access_token":  accessToken,
		"refresh_token": req.RefreshToken,
		"expires_in":    3600,
	})
}

// Logout and revoke refresh token
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "refresh_token required")
		return
	}

	if err := h.RefreshTokenRepo.RevokeRefreshToken(req.RefreshToken); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to revoke token")
		return
	}

	utils.Success(c, "logged out", nil)
}
