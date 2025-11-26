package events

import "strings"

// EventType represents a match event type.
// This is a string type to allow extensibility for 1000s of event types.
type EventType string

// Common event types - these are the most frequently used.
// For a complete list, see the event type registry below.
const (
	// Goals
	EventTypeGoal      EventType = "goal"
	EventTypeOwnGoal   EventType = "own_goal"
	EventTypePenalty   EventType = "penalty"
	EventTypePenaltyGoal EventType = "penalty_goal"
	EventTypePenaltyMiss EventType = "penalty_miss"

	// Cards
	EventTypeYellowCard EventType = "yellow_card"
	EventTypeRedCard    EventType = "red_card"
	EventTypeSecondYellow EventType = "second_yellow_card"

	// Substitutions
	EventTypeSubstitution EventType = "substitution"
	EventTypeSubstitutionOn  EventType = "substitution_on"
	EventTypeSubstitutionOff EventType = "substitution_off"

	// Shots
	EventTypeShot        EventType = "shot"
	EventTypeShotOnTarget EventType = "shot_on_target"
	EventTypeShotOffTarget EventType = "shot_off_target"
	EventTypeShotBlocked  EventType = "shot_blocked"
	EventTypeShotSaved    EventType = "shot_saved"
	EventTypeShotPost     EventType = "shot_post"
	EventTypeShotWoodwork EventType = "shot_woodwork"

	// Passes
	EventTypePass         EventType = "pass"
	EventTypePassCompleted EventType = "pass_completed"
	EventTypePassIncomplete EventType = "pass_incomplete"
	EventTypeKeyPass       EventType = "key_pass"
	EventTypeAssist        EventType = "assist"
	EventTypeThroughBall   EventType = "through_ball"
	EventTypeCross         EventType = "cross"
	EventTypeLongBall      EventType = "long_ball"
	EventTypeShortPass     EventType = "short_pass"

	// Defensive actions
	EventTypeTackle        EventType = "tackle"
	EventTypeTackleWon     EventType = "tackle_won"
	EventTypeTackleLost    EventType = "tackle_lost"
	EventTypeInterception  EventType = "interception"
	EventTypeClearance     EventType = "clearance"
	EventTypeBlock         EventType = "block"
	EventTypeBlockedShot   EventType = "blocked_shot"

	// Duels
	EventTypeDuel         EventType = "duel"
	EventTypeDuelWon       EventType = "duel_won"
	EventTypeDuelLost     EventType = "duel_lost"
	EventTypeAerialDuel   EventType = "aerial_duel"
	EventTypeAerialDuelWon EventType = "aerial_duel_won"
	EventTypeAerialDuelLost EventType = "aerial_duel_lost"
	EventTypeGroundDuel    EventType = "ground_duel"

	// Fouls
	EventTypeFoul         EventType = "foul"
	EventTypeFoulCommitted EventType = "foul_committed"
	EventTypeFoulWon      EventType = "foul_won"
	EventTypeOffside      EventType = "offside"

	// Goalkeeper events
	EventTypeSave         EventType = "save"
	EventTypeSavePenalty   EventType = "save_penalty"
	EventTypeSaveSixYardBox EventType = "save_six_yard_box"
	EventTypeSavePenaltyArea EventType = "save_penalty_area"
	EventTypeSaveOutOfBox EventType = "save_out_of_box"
	EventTypePunch        EventType = "punch"
	EventTypeClaim        EventType = "claim"
	EventTypeSweeperKeeper EventType = "sweeper_keeper"

	// VAR and reviews
	EventTypeVarReview    EventType = "var_review"
	EventTypeVarGoal      EventType = "var_goal"
	EventTypeVarPenalty   EventType = "var_penalty"
	EventTypeVarRedCard   EventType = "var_red_card"

	// Match state
	EventTypeKickOff      EventType = "kick_off"
	EventTypeHalfTime     EventType = "half_time"
	EventTypeFullTime     EventType = "full_time"
	EventTypeExtraTime    EventType = "extra_time"
	EventTypePenaltyShootout EventType = "penalty_shootout"
)

// EventCategory groups related event types for analytics and filtering.
type EventCategory string

const (
	CategoryGoal         EventCategory = "goal"
	CategoryCard         EventCategory = "card"
	CategorySubstitution EventCategory = "substitution"
	CategoryShot         EventCategory = "shot"
	CategoryPass         EventCategory = "pass"
	CategoryDefensive    EventCategory = "defensive"
	CategoryDuel         EventCategory = "duel"
	CategoryFoul         EventCategory = "foul"
	CategoryGoalkeeper   EventCategory = "goalkeeper"
	CategoryVar          EventCategory = "var"
	CategoryMatchState   EventCategory = "match_state"
	CategoryOther        EventCategory = "other"
)

// GetCategory returns the category for an event type.
func GetCategory(eventType EventType) EventCategory {
	switch eventType {
	// Goals
	case EventTypeGoal, EventTypeOwnGoal, EventTypePenalty, EventTypePenaltyGoal, EventTypePenaltyMiss:
		return CategoryGoal

	// Cards
	case EventTypeYellowCard, EventTypeRedCard, EventTypeSecondYellow:
		return CategoryCard

	// Substitutions
	case EventTypeSubstitution, EventTypeSubstitutionOn, EventTypeSubstitutionOff:
		return CategorySubstitution

	// Shots
	case EventTypeShot, EventTypeShotOnTarget, EventTypeShotOffTarget, EventTypeShotBlocked,
		EventTypeShotSaved, EventTypeShotPost, EventTypeShotWoodwork:
		return CategoryShot

	// Passes
	case EventTypePass, EventTypePassCompleted, EventTypePassIncomplete, EventTypeKeyPass,
		EventTypeAssist, EventTypeThroughBall, EventTypeCross, EventTypeLongBall, EventTypeShortPass:
		return CategoryPass

	// Defensive
	case EventTypeTackle, EventTypeTackleWon, EventTypeTackleLost, EventTypeInterception,
		EventTypeClearance, EventTypeBlock, EventTypeBlockedShot:
		return CategoryDefensive

	// Duels
	case EventTypeDuel, EventTypeDuelWon, EventTypeDuelLost, EventTypeAerialDuel,
		EventTypeAerialDuelWon, EventTypeAerialDuelLost, EventTypeGroundDuel:
		return CategoryDuel

	// Fouls
	case EventTypeFoul, EventTypeFoulCommitted, EventTypeFoulWon, EventTypeOffside:
		return CategoryFoul

	// Goalkeeper
	case EventTypeSave, EventTypeSavePenalty, EventTypeSaveSixYardBox, EventTypeSavePenaltyArea,
		EventTypeSaveOutOfBox, EventTypePunch, EventTypeClaim, EventTypeSweeperKeeper:
		return CategoryGoalkeeper

	// VAR
	case EventTypeVarReview, EventTypeVarGoal, EventTypeVarPenalty, EventTypeVarRedCard:
		return CategoryVar

	// Match state
	case EventTypeKickOff, EventTypeHalfTime, EventTypeFullTime, EventTypeExtraTime, EventTypePenaltyShootout:
		return CategoryMatchState

	default:
		return CategoryOther
	}
}

// IsValid checks if an event type is valid (non-empty and reasonable length).
// Note: We don't restrict to a fixed list to allow extensibility.
func IsValid(eventType EventType) bool {
	if eventType == "" {
		return false
	}
	// Reasonable length check (VARCHAR(50) in DB)
	if len(eventType) > 50 {
		return false
	}
	// Basic format check: lowercase, alphanumeric + underscores
	for _, r := range eventType {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// Normalize normalizes an event type string (lowercase, trim spaces).
func Normalize(eventType string) EventType {
	normalized := strings.ToLower(strings.TrimSpace(eventType))
	return EventType(normalized)
}

// IsGoal returns true if the event type is goal-related.
func (et EventType) IsGoal() bool {
	return GetCategory(et) == CategoryGoal
}

// IsCard returns true if the event type is card-related.
func (et EventType) IsCard() bool {
	return GetCategory(et) == CategoryCard
}

// IsShot returns true if the event type is shot-related.
func (et EventType) IsShot() bool {
	return GetCategory(et) == CategoryShot
}

// IsPass returns true if the event type is pass-related.
func (et EventType) IsPass() bool {
	return GetCategory(et) == CategoryPass
}

// String returns the string representation of the event type.
func (et EventType) String() string {
	return string(et)
}
