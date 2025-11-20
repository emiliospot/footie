package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// UserHandler handles user-related endpoints.
type UserHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewUserHandler creates a new user handler.
func NewUserHandler(db *gorm.DB, log *logger.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: log,
	}
}

// @Router /users/me [get].
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Router /users/me [put].
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only allow updating certain fields
	allowedFields := []string{"first_name", "last_name", "avatar", "organization"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if value, exists := updateData[field]; exists {
			updates[field] = value
		}
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		h.logger.Error("Failed to update user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Router /users/{id} [get].
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Router /admin/users [get].
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	h.db.Model(&models.User{}).Count(&total)

	if err := h.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		h.logger.Error("Failed to fetch users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Router /admin/users/{id}/role [put].
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required,oneof=user analyst admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = req.Role
	if err := h.db.Save(&user).Error; err != nil {
		h.logger.Error("Failed to update user role", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Router /admin/users/{id} [delete].
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
