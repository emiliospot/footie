package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
	"github.com/emiliospot/footie/api/pkg/auth"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	db     *gorm.DB
	cfg    *config.Config
	logger *logger.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(cfg *config.Config, db *gorm.DB, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		cfg:    cfg,
		logger: log,
	}
}

// RegisterRequest represents a registration request.
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Organization string `json:"organization"`
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents an authentication response.
type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

// @Router /auth/register [post].
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create user
	user := models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Organization: req.Organization,
		Role:         "user",
		IsActive:     true,
	}

	if createErr := h.db.Create(&user).Error; createErr != nil {
		h.logger.Error("Failed to create user", "error", createErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWT.Secret, h.cfg.JWT.ExpiryHours)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWT.Secret, h.cfg.JWT.RefreshExpiryHours)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// @Router /auth/login [post].
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is disabled"})
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWT.Secret, h.cfg.JWT.ExpiryHours)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWT.Secret, h.cfg.JWT.RefreshExpiryHours)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// @Router /auth/refresh [post].
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Extract token
	tokenString := authHeader[7:] // Remove "Bearer "

	// Generate new token
	newToken, err := auth.RefreshToken(tokenString, h.cfg.JWT.Secret, h.cfg.JWT.ExpiryHours)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}
