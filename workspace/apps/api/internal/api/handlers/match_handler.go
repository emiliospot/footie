package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

const (
	// ErrMatchNotFound is returned when a match is not found.
	ErrMatchNotFound = "Match not found"
	// ErrFailedToFetchMatches is returned when fetching matches fails.
	ErrFailedToFetchMatches = "Failed to fetch matches"
	// ErrFailedToCreateMatch is returned when match creation fails.
	ErrFailedToCreateMatch = "Failed to create match"
	// ErrFailedToUpdateMatch is returned when match update fails.
	ErrFailedToUpdateMatch = "Failed to update match"
	// ErrFailedToFetchEvents is returned when fetching match events fails.
	ErrFailedToFetchEvents = "Failed to fetch events"
	// ErrFailedToCreateEvent is returned when event creation fails.
	ErrFailedToCreateEvent = "Failed to create event"
	// ErrInvalidMatchID is returned when match ID is invalid.
	ErrInvalidMatchID = "Invalid match ID"
)

// MatchHandler handles match-related endpoints.
type MatchHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewMatchHandler creates a new match handler.
func NewMatchHandler(db *gorm.DB, log *logger.Logger) *MatchHandler {
	return &MatchHandler{
		db:     db,
		logger: log,
	}
}

// @Router /matches [get].
func (h *MatchHandler) ListMatches(c *gin.Context) {
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
	season := c.Query("season")
	competition := c.Query("competition")
	status := c.Query("status")

	query := h.db.Model(&models.Match{}).Preload("HomeTeam").Preload("AwayTeam")
	if teamID != "" {
		query = query.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}
	if season != "" {
		query = query.Where("season = ?", season)
	}
	if competition != "" {
		query = query.Where("competition = ?", competition)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var matches []models.Match
	if err := query.Order("match_date DESC").Offset(offset).Limit(limit).Find(&matches).Error; err != nil {
		h.logger.Error("Failed to fetch matches", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToFetchMatches})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Router /matches/{id} [get].
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id := c.Param("id")

	var match models.Match
	if err := h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("Events").First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrMatchNotFound})
		return
	}

	c.JSON(http.StatusOK, match)
}

// @Router /matches [post].
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var match models.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&match).Error; err != nil {
		h.logger.Error("Failed to create match", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToCreateMatch})
		return
	}

	h.db.Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)
	c.JSON(http.StatusCreated, match)
}

// @Router /matches/{id} [put].
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	id := c.Param("id")

	var match models.Match
	if err := h.db.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrMatchNotFound})
		return
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&match).Error; err != nil {
		h.logger.Error("Failed to update match", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToUpdateMatch})
		return
	}

	h.db.Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)
	c.JSON(http.StatusOK, match)
}

// @Router /matches/{id} [delete].
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Match{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrMatchNotFound})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Router /matches/{id}/events [get].
func (h *MatchHandler) GetMatchEvents(c *gin.Context) {
	id := c.Param("id")

	var events []models.MatchEvent
	if err := h.db.Where("match_id = ?", id).
		Preload("Player").
		Preload("Team").
		Preload("SecondaryPlayer").
		Order("minute ASC, extra_minute ASC").
		Find(&events).Error; err != nil {
		h.logger.Error("Failed to fetch events", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToFetchEvents})
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Router /matches/{id}/events [post].
func (h *MatchHandler) CreateMatchEvent(c *gin.Context) {
	matchID := c.Param("id")

	var event models.MatchEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set match ID from URL
	id, err := strconv.Atoi(matchID)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidMatchID})
		return
	}
	event.MatchID = uint(id)

	if err := h.db.Create(&event).Error; err != nil {
		h.logger.Error("Failed to create event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToCreateEvent})
		return
	}

	h.db.Preload("Player").Preload("Team").Preload("SecondaryPlayer").First(&event, event.ID)
	c.JSON(http.StatusCreated, event)
}
