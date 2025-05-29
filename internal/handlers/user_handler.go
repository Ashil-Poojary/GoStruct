package handlers

import (
	"net/http"

	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: repo}
}

func (h *UserHandler) Me(c *gin.Context) {
	claims, ok := c.Get("user")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userClaims, ok := claims.(map[string]interface{})
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	// Usually JWT numeric claims come as float64
	idFloat, ok := userClaims["id"].(float64)
	if !ok {
		utils.Error(c, http.StatusBadRequest, "Invalid token data")
		return
	}
	userID := int64(idFloat)

	user, err := h.UserRepo.GetByID(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if user == nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, "user fetched successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// GetAllUsers
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserRepo.GetAll()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	if len(users) == 0 {
		utils.Success(c, "No users found", gin.H{"users": []string{}})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		})
	}

	utils.Success(c, "users fetched successfully", gin.H{"users": userList})
}
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}
	//convert id to int64
	idInt, err := utils.StringToInt64(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to Login")
		return
	}
	user, err := h.UserRepo.GetByID(idInt)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if user == nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, "user fetched successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}
