# ⏱️ Duration Handling in Match Events

> **Where and how duration is stored for matches and events**

## Current State

**Duration is NOT currently stored** in the database schema. Here's what we have:

### Match Events

- ✅ `minute` - When the event occurred (0-120)
- ✅ `extra_minute` - Extra time period (0-30)
- ❌ `duration` - How long the event lasted (missing)

### Matches

- ✅ `match_date` - When the match started
- ✅ `status` - Current state (scheduled, live, finished)
- ❌ `duration` - Total match duration (missing)
- ❌ `stoppage_time` - Injury/stoppage time (missing)

---

## Where Duration Should Be Stored

### Option 1: Event Duration (in `metadata` JSONB)

**Best for:** Event-specific durations (VAR reviews, injury time, etc.)

```json
{
  "matchId": 12345,
  "eventType": "var_review",
  "minute": 78,
  "metadata": {
    "duration_seconds": 120, // VAR review took 2 minutes
    "review_type": "penalty"
  }
}
```

**Pros:**

- ✅ Flexible - different events can have different duration meanings
- ✅ No schema changes needed
- ✅ Already supported via `metadata` JSONB field

**Cons:**

- ❌ Not queryable directly (need JSONB queries)
- ❌ No validation

### Option 2: Add `duration` Column to `match_events`

**Best for:** Standardized event durations

```sql
ALTER TABLE match_events
ADD COLUMN duration_seconds INTEGER;  -- Duration in seconds
```

**Pros:**

- ✅ Queryable directly
- ✅ Type-safe
- ✅ Indexable

**Cons:**

- ❌ Requires migration
- ❌ Not all events have duration (nullable)

### Option 3: Match Duration (in `matches` table)

**Best for:** Total match duration

```sql
ALTER TABLE matches
ADD COLUMN duration_minutes INTEGER DEFAULT 90,  -- Standard match duration
ADD COLUMN stoppage_time_first_half INTEGER DEFAULT 0,
ADD COLUMN stoppage_time_second_half INTEGER DEFAULT 0;
```

---

## Recommended Approach

### For Events: Use `metadata` JSONB

Store event-specific durations in metadata:

```go
// In webhook payload
{
  "matchId": 12345,
  "eventType": "var_review",
  "minute": 78,
  "metadata": {
    "duration_seconds": 120,
    "review_type": "penalty",
    "decision": "overturned"
  }
}
```

### For Matches: Add Duration Fields

Add to `matches` table:

```sql
ALTER TABLE matches
ADD COLUMN duration_minutes INTEGER DEFAULT 90,
ADD COLUMN stoppage_time_first_half INTEGER DEFAULT 0,
ADD COLUMN stoppage_time_second_half INTEGER DEFAULT 0,
ADD COLUMN extra_time_minutes INTEGER DEFAULT 0;
```

---

## Implementation Examples

### 1. Webhook Payload with Duration

```go
// GenericProvider payload
type GenericPayload struct {
    MatchID           int32                  `json:"matchId"`
    EventType         string                 `json:"eventType"`
    Minute            int32                   `json:"minute"`
    ExtraMinute       *int32                  `json:"extraMinute,omitempty"`
    Duration          *int32                  `json:"duration,omitempty"`  // NEW: Duration in seconds
    // ... other fields
    Metadata          map[string]interface{} `json:"metadata,omitempty"`
}
```

### 2. Extract Duration in Provider

```go
func (p *GenericProvider) ExtractEvent(ctx context.Context, payload []byte) (*infraEvents.MatchEvent, error) {
    var genericPayload GenericPayload
    json.Unmarshal(payload, &genericPayload)

    // Store duration in metadata if provided
    metadata := make(map[string]interface{})
    if genericPayload.Metadata != nil {
        metadata = genericPayload.Metadata
    }

    if genericPayload.Duration != nil {
        metadata["duration_seconds"] = *genericPayload.Duration
    }

    // ... rest of extraction
}
```

### 3. Query Events with Duration

```sql
-- Get VAR reviews with duration > 60 seconds
SELECT
    id,
    minute,
    metadata->>'duration_seconds' as duration
FROM match_events
WHERE event_type = 'var_review'
AND (metadata->>'duration_seconds')::int > 60;
```

---

## Use Cases

### 1. VAR Review Duration

```json
{
  "eventType": "var_review",
  "minute": 78,
  "metadata": {
    "duration_seconds": 180,
    "review_type": "offside",
    "decision": "goal_disallowed"
  }
}
```

### 2. Injury Time

```json
{
  "eventType": "injury",
  "minute": 45,
  "metadata": {
    "duration_seconds": 240,
    "player_id": 42,
    "injury_type": "head"
  }
}
```

### 3. Substitution (Player Time on Field)

```json
{
  "eventType": "substitution_off",
  "minute": 67,
  "playerId": 10,
  "metadata": {
    "minutes_played": 67,
    "substitution_reason": "tactical"
  }
}
```

### 4. Match Duration

```sql
-- Update match duration after completion
UPDATE matches
SET duration_minutes = 90,
    stoppage_time_first_half = 3,
    stoppage_time_second_half = 5
WHERE id = 12345 AND status = 'finished';
```

---

## Current Payload Structure

### Generic Provider

```json
{
  "matchId": 12345,
  "eventType": "var_review",
  "minute": 78,
  "extraMinute": 0,
  "teamId": 10,
  "playerId": 42,
  "positionX": 0.75,
  "positionY": 0.5,
  "description": "VAR review for penalty",
  "metadata": {
    "duration_seconds": 120, // ← Duration goes here
    "review_type": "penalty",
    "decision": "awarded"
  }
}
```

### Opta Provider

```json
{
  "event": {
    "type": "VarReview",
    "minute": 78,
    "qualifiers": [
      {
        "type": "duration",
        "value": 120 // ← Duration in qualifiers
      },
      {
        "type": "reviewType",
        "value": "penalty"
      }
    ]
  }
}
```

### StatsBomb Provider

```json
{
  "type": "Var Review",
  "minute": 78,
  "duration": 120, // ← Direct field
  "review_type": "penalty"
}
```

---

## Summary

**Current Status:**

- ❌ No duration field in schema
- ✅ Duration can be stored in `metadata` JSONB
- ✅ Providers can extract duration from payloads

**Recommendation:**

1. **Events**: Store duration in `metadata` JSONB (flexible, no migration)
2. **Matches**: Add duration columns if needed (for analytics)

**Next Steps:**

1. Update provider adapters to extract duration from payloads
2. Store duration in metadata JSONB
3. Add match duration columns if analytics require it
