# ⚽ Event Types System

> **How event types are defined, validated, and extended**

## Overview

Football matches generate **thousands of different event types** (goals, shots, passes, tackles, cards, substitutions, etc.). This system provides a flexible, extensible way to handle all event types while maintaining type safety and validation.

---

## 🏗️ Architecture

### 1. **Domain Layer** (`internal/domain/events/types.go`)

Defines the event type system with:

- **EventType**: String-based type (allows 1000s of types)
- **EventCategory**: Groups related types (goal, shot, pass, etc.)
- **Validation**: Format checking (lowercase, alphanumeric + underscores)
- **Normalization**: Converts external formats to internal format
- **Helper methods**: `IsGoal()`, `IsCard()`, `IsShot()`, etc.

### 2. **Database Schema**

```sql
event_type VARCHAR(50) NOT NULL  -- Flexible, allows any valid type
```

**Why VARCHAR instead of ENUM?**

- ✅ Supports 1000s of event types
- ✅ Easy to add new types without migrations
- ✅ Provider-specific types (Opta, StatsBomb) can be normalized
- ✅ No database lock when adding new types

### 3. **Provider Adapters**

Each provider (Opta, StatsBomb, Generic) normalizes their event types:

```go
// Opta sends: "Goal"
// StatsBomb sends: "Shot"
// Generic sends: "GOAL"

// All normalized to: "goal"
eventType := events.Normalize(externalType)
```

---

## 📋 Event Type Categories

Events are grouped into categories for analytics:

| Category         | Examples                                                    |
| ---------------- | ----------------------------------------------------------- |
| **Goal**         | `goal`, `own_goal`, `penalty`, `penalty_goal`               |
| **Card**         | `yellow_card`, `red_card`, `second_yellow_card`             |
| **Substitution** | `substitution`, `substitution_on`, `substitution_off`       |
| **Shot**         | `shot`, `shot_on_target`, `shot_off_target`, `shot_blocked` |
| **Pass**         | `pass`, `pass_completed`, `key_pass`, `assist`, `cross`     |
| **Defensive**    | `tackle`, `interception`, `clearance`, `block`              |
| **Duel**         | `duel`, `aerial_duel`, `ground_duel`                        |
| **Foul**         | `foul`, `foul_committed`, `foul_won`, `offside`             |
| **Goalkeeper**   | `save`, `save_penalty`, `punch`, `claim`                    |
| **VAR**          | `var_review`, `var_goal`, `var_penalty`                     |
| **Match State**  | `kick_off`, `half_time`, `full_time`                        |
| **Other**        | Any type not in above categories                            |

---

## 🔧 Usage Examples

### Normalizing Event Types

```go
import "github.com/emiliospot/footie/api/internal/domain/events"

// External provider sends: "GOAL"
normalized := events.Normalize("GOAL")
// Result: events.EventType("goal")

// Validate
if !events.IsValid(normalized) {
    return fmt.Errorf("invalid event type")
}

// Check category
category := events.GetCategory(normalized)
// Result: events.CategoryGoal

// Use helper methods
if normalized.IsGoal() {
    // Handle goal event
}
```

### In Webhook Providers

```go
func (p *OptaProvider) ExtractEvent(ctx context.Context, payload []byte) (*infraEvents.MatchEvent, error) {
    // Parse Opta payload
    var optaPayload OptaPayload
    json.Unmarshal(payload, &optaPayload)

    // Normalize Opta's event type
    eventType := events.Normalize(optaPayload.Event.Type)

    // Validate
    if !events.IsValid(eventType) {
        return nil, fmt.Errorf("invalid event type: %s", optaPayload.Event.Type)
    }

    // Use normalized type
    return &infraEvents.MatchEvent{
        EventType: eventType.String(),
        // ...
    }, nil
}
```

### Querying by Category

```sql
-- Get all goal-related events
SELECT * FROM match_events
WHERE event_type IN ('goal', 'own_goal', 'penalty', 'penalty_goal')
AND match_id = $1;

-- Get all shot-related events
SELECT * FROM match_events
WHERE event_type IN ('shot', 'shot_on_target', 'shot_off_target', 'shot_blocked')
AND match_id = $1;
```

---

## ➕ Adding New Event Types

### Step 1: Add Constant (Optional)

If it's a common type, add to `types.go`:

```go
const (
    EventTypeNewType EventType = "new_type"
)
```

### Step 2: Add Category (If Needed)

```go
func GetCategory(eventType EventType) EventCategory {
    switch eventType {
    case EventTypeNewType:
        return CategoryNewCategory
    // ...
    }
}
```

### Step 3: Use in Code

That's it! The system automatically accepts any valid event type:

```go
// This works immediately, no migration needed
event := &MatchEvent{
    EventType: "new_custom_type",  // Automatically accepted
    // ...
}
```

**Validation rules:**

- ✅ Lowercase
- ✅ Alphanumeric + underscores only
- ✅ Max 50 characters
- ✅ Non-empty

---

## 🔍 Provider-Specific Types

Different providers use different event type names. The system normalizes them:

| Provider     | External Type  | Normalized     |
| ------------ | -------------- | -------------- |
| Opta         | `Goal`         | `goal`         |
| StatsBomb    | `Shot`         | `shot`         |
| API-Football | `GOAL`         | `goal`         |
| Custom       | `penalty-kick` | `penalty_kick` |

**Normalization rules:**

1. Convert to lowercase
2. Replace hyphens with underscores
3. Trim whitespace

---

## 📊 Analytics Queries

### Count Events by Category

```sql
SELECT
    CASE
        WHEN event_type IN ('goal', 'own_goal', 'penalty') THEN 'goal'
        WHEN event_type IN ('shot', 'shot_on_target') THEN 'shot'
        -- ...
    END as category,
    COUNT(*) as count
FROM match_events
WHERE match_id = $1
GROUP BY category;
```

### Get All Goal Events

```go
// Using domain helpers
goalTypes := []string{
    events.EventTypeGoal.String(),
    events.EventTypeOwnGoal.String(),
    events.EventTypePenalty.String(),
}

events, _ := queries.GetMatchEventsByType(ctx, matchID, goalTypes...)
```

---

## 🎯 Best Practices

1. **Always normalize** external event types before storing
2. **Use categories** for analytics, not individual types
3. **Validate** event types before processing
4. **Use constants** for common types (better IDE support)
5. **Document** provider-specific types in provider adapters

---

## 🔐 Type Safety

While event types are strings (for flexibility), we maintain safety through:

- ✅ **Validation**: Format checking
- ✅ **Normalization**: Consistent format
- ✅ **Categories**: Grouped analysis
- ✅ **Constants**: Common types are typed
- ✅ **Database constraints**: VARCHAR(50) length limit

---

## 📚 Related Files

- `internal/domain/events/types.go` - Event type definitions
- `internal/infrastructure/webhooks/providers/` - Provider adapters
- `migrations/000001_init_schema.up.sql` - Database schema
- `internal/domain/models/match_event.go` - Domain model
