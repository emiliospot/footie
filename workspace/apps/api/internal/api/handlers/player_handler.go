package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// PlayerHandler handles player-related endpoints.
type PlayerHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewPlayerHandler creates a new player handler.
func NewPlayerHandler(db *gorm.DB, log *logger.Logger) *PlayerHandler {
	return &PlayerHandler{
		db:     db,
		logger: log,
	}
}

// @Router /players [get].
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	teamID := c.Query("team_id")
	position := c.Query("position")

	query := h.db.Model(&models.Player{}).Preload("Team")
	if teamID != "" {
		query = query.Where("team_id = ?", teamID)
	}
	if position != "" {
		query = query.Where("position = ?", position)
	}

	var total int64
	query.Count(&total)

	var players []models.Player
	if err := query.Offset(offset).Limit(limit).Find(&players).Error; err != nil {
		h.logger.Error("Failed to fetch players", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"players": players,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Router /players/{id} [get].
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")

	var player models.Player
	if err := h.db.Preload("Team").First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// @Router /players [post].
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&player).Error; err != nil {
		h.logger.Error("Failed to create player", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player"})
		return
	}

	h.db.Preload("Team").First(&player, player.ID)
	c.JSON(http.StatusCreated, player)
}

// @Router /players/{id} [put].
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")

	var player models.Player
	if err := h.db.First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&player).Error; err != nil {
		h.logger.Error("Failed to update player", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
		return
	}

	h.db.Preload("Team").First(&player, player.ID)
	c.JSON(http.StatusOK, player)
}

// @Router /players/{id} [delete].
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Player{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Router /players/{id}/statistics [get].
func (h *PlayerHandler) GetPlayerStatistics(c *gin.Context) {
	id := c.Param("id")
	season := c.Query("season")
	competition := c.Query("competition")

	query := h.db.Where("player_id = ?", id)
	if season != "" {
		query = query.Where("season = ?", season)
	}
	if competition != "" {
		query = query.Where("competition = ?", competition)
	}

	var stats []models.PlayerStatistics
	if err := query.Preload("Player").Find(&stats).Error; err != nil {
		h.logger.Error("Failed to fetch statistics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
