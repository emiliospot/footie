package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// Publisher handles publishing match events to Redis Streams and Pub/Sub.
type Publisher struct {
	redis  *redis.Client
	logger *logger.Logger
}

// MatchEvent represents a football match event.
type MatchEvent struct {
	ID                int32     `json:"id"`
	MatchID           int32     `json:"match_id"`
	TeamID            *int32    `json:"team_id,omitempty"`
	PlayerID          *int32    `json:"player_id,omitempty"`
	SecondaryPlayerID *int32    `json:"secondary_player_id,omitempty"`
	EventType         string    `json:"event_type"` // goal, shot, pass, card, substitution
	Minute            int       `json:"minute"`
	ExtraMinute       int       `json:"extra_minute,omitempty"`
	PositionX         *float64  `json:"position_x,omitempty"`
	PositionY         *float64  `json:"position_y,omitempty"`
	Description       string    `json:"description,omitempty"`
	Metadata          string    `json:"metadata,omitempty"` // JSON string with xG, pass completion, etc.
	Timestamp         time.Time `json:"timestamp"`
}

// ScoreUpdate represents a match score update.
type ScoreUpdate struct {
	MatchID       int32     `json:"match_id"`
	HomeTeamScore int       `json:"home_team_score"`
	AwayTeamScore int       `json:"away_team_score"`
	Timestamp     time.Time `json:"timestamp"`
}

// MatchStatusUpdate represents a match status change.
type MatchStatusUpdate struct {
	MatchID   int32     `json:"match_id"`
	Status    string    `json:"status"` // scheduled, live, finished, postponed, canceled
	Timestamp time.Time `json:"timestamp"`
}

// NewPublisher creates a new event publisher.
func NewPublisher(redis *redis.Client, logger *logger.Logger) *Publisher {
	return &Publisher{
		redis:  redis,
		logger: logger,
	}
}

// PublishMatchEvent publishes a match event to both Redis Stream and Pub/Sub.
func (p *Publisher) PublishMatchEvent(ctx context.Context, event *MatchEvent) error {
	event.Timestamp = time.Now()

	// Marshal event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 1. Add to Redis Stream for processing/analytics
	streamKey := fmt.Sprintf("match:%d:stream", event.MatchID)
	if err := p.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"event_type": event.EventType,
			"data":       string(eventJSON),
			"timestamp":  event.Timestamp.Unix(),
		},
	}).Err(); err != nil {
		p.logger.Error("Failed to add event to stream", "error", err, "match_id", event.MatchID)
		return fmt.Errorf("failed to add to stream: %w", err)
	}

	// 2. Publish to Pub/Sub for real-time WebSocket delivery
	channel := fmt.Sprintf("match:%d:events", event.MatchID)
	message := map[string]interface{}{
		"type":      "match_event",
		"match_id":  event.MatchID,
		"timestamp": event.Timestamp,
		"data":      event,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal pub/sub message: %w", err)
	}

	if err := p.redis.Publish(ctx, channel, messageJSON).Err(); err != nil {
		p.logger.Error("Failed to publish event", "error", err, "match_id", event.MatchID)
		return fmt.Errorf("failed to publish: %w", err)
	}

	p.logger.Info("Published match event",
		"match_id", event.MatchID,
		"event_type", event.EventType,
		"minute", event.Minute,
	)

	return nil
}

// PublishScoreUpdate publishes a score update.
func (p *Publisher) PublishScoreUpdate(ctx context.Context, update *ScoreUpdate) error {
	update.Timestamp = time.Now()

	channel := fmt.Sprintf("match:%d:events", update.MatchID)
	message := map[string]interface{}{
		"type":      "score_update",
		"match_id":  update.MatchID,
		"timestamp": update.Timestamp,
		"data":      update,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal score update: %w", err)
	}

	if err := p.redis.Publish(ctx, channel, messageJSON).Err(); err != nil {
		p.logger.Error("Failed to publish score update", "error", err, "match_id", update.MatchID)
		return fmt.Errorf("failed to publish score update: %w", err)
	}

	p.logger.Info("Published score update",
		"match_id", update.MatchID,
		"home_score", update.HomeTeamScore,
		"away_score", update.AwayTeamScore,
	)

	return nil
}

// PublishMatchStatusUpdate publishes a match status change.
func (p *Publisher) PublishMatchStatusUpdate(ctx context.Context, update *MatchStatusUpdate) error {
	update.Timestamp = time.Now()

	channel := fmt.Sprintf("match:%d:events", update.MatchID)
	message := map[string]interface{}{
		"type":      "match_status",
		"match_id":  update.MatchID,
		"timestamp": update.Timestamp,
		"data":      update,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal status update: %w", err)
	}

	if err := p.redis.Publish(ctx, channel, messageJSON).Err(); err != nil {
		p.logger.Error("Failed to publish status update", "error", err, "match_id", update.MatchID)
		return fmt.Errorf("failed to publish status update: %w", err)
	}

	p.logger.Info("Published match status update",
		"match_id", update.MatchID,
		"status", update.Status,
	)

	return nil
}

// InvalidateMatchCache invalidates cached match data.
func (p *Publisher) InvalidateMatchCache(ctx context.Context, matchID int32) error {
	keys := []string{
		fmt.Sprintf("match:%d", matchID),
		fmt.Sprintf("match:%d:events", matchID),
		fmt.Sprintf("match:%d:stats", matchID),
	}

	for _, key := range keys {
		if err := p.redis.Del(ctx, key).Err(); err != nil {
			p.logger.Error("Failed to invalidate cache", "error", err, "key", key)
		}
	}

	p.logger.Info("Invalidated match cache", "match_id", matchID)
	return nil
}

