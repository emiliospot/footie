package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// TeamHandler handles team-related endpoints.
type TeamHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewTeamHandler creates a new team handler.
func NewTeamHandler(db *gorm.DB, log *logger.Logger) *TeamHandler {
	return &TeamHandler{
		db:     db,
		logger: log,
	}
}

// @Router /teams [get].
func (h *TeamHandler) ListTeams(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	country := c.Query("country")

	query := h.db.Model(&models.Team{})
	if country != "" {
		query = query.Where("country = ?", country)
	}

	var total int64
	query.Count(&total)

	var teams []models.Team
	if err := query.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		h.logger.Error("Failed to fetch teams", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Router /teams/{id} [get].
func (h *TeamHandler) GetTeam(c *gin.Context) {
	id := c.Param("id")

	var team models.Team
	if err := h.db.Preload("Players").First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// @Router /teams [post].
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&team).Error; err != nil {
		h.logger.Error("Failed to create team", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// @Router /teams/{id} [put].
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")

	var team models.Team
	if err := h.db.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&team).Error; err != nil {
		h.logger.Error("Failed to update team", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// @Router /teams/{id} [delete].
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Team{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Router /teams/{id}/players [get].
func (h *TeamHandler) GetTeamPlayers(c *gin.Context) {
	id := c.Param("id")

	var players []models.Player
	if err := h.db.Where("team_id = ?", id).Find(&players).Error; err != nil {
		h.logger.Error("Failed to fetch players", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// @Router /teams/{id}/statistics [get].
func (h *TeamHandler) GetTeamStatistics(c *gin.Context) {
	id := c.Param("id")
	season := c.Query("season")
	competition := c.Query("competition")

	query := h.db.Where("team_id = ?", id)
	if season != "" {
		query = query.Where("season = ?", season)
	}
	if competition != "" {
		query = query.Where("competition = ?", competition)
	}

	var stats []models.TeamStatistics
	if err := query.Find(&stats).Error; err != nil {
		h.logger.Error("Failed to fetch statistics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
