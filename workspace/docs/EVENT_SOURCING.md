# 🔄 Event Sourcing Preparation Guide

> **Preparing the Footie platform for event sourcing architecture**

Event sourcing is perfect for football analytics because matches are naturally event-driven (goals, shots, passes, cards). This guide shows how to prepare your current architecture for event sourcing.

---

## 🎯 What is Event Sourcing?

Instead of storing current state, you store **all events** that led to that state.

### Traditional Approach (Current)

```
Match Table:
┌────┬──────────┬──────────┬────────┐
│ ID │ Home     │ Away     │ Status │
├────┼──────────┼──────────┼────────┤
│ 1  │ Team A   │ Team B   │ Live   │
│    │ Score: 2 │ Score: 1 │        │
└────┴──────────┴──────────┴────────┘

Problem: Lost history of HOW we got to 2-1
```

### Event Sourcing Approach

```
Event Store:
┌────┬──────────┬────────────┬────────────────┐
│ ID │ Match ID │ Event Type │ Data           │
├────┼──────────┼────────────┼────────────────┤
│ 1  │ 1        │ MatchStart │ {...}          │
│ 2  │ 1        │ Goal       │ {team: A, ...} │
│ 3  │ 1        │ Goal       │ {team: B, ...} │
│ 4  │ 1        │ Goal       │ {team: A, ...} │
└────┴──────────┴────────────┴────────────────┘

Benefit: Can replay match, calculate stats at any point
```

---

## 🏗️ Current Architecture (Already Event-Ready!)

You're **80% there** already! Here's what you have:

### ✅ What's Already in Place

1. **Event Table** (`match_events`)

   ```sql
   CREATE TABLE match_events (
       id SERIAL PRIMARY KEY,
       match_id INT NOT NULL,
       event_type VARCHAR(50) NOT NULL,
       minute INT NOT NULL,
       metadata JSONB,
       created_at TIMESTAMPTZ DEFAULT NOW()
   );
   ```

2. **Event Publisher** (`internal/infrastructure/events/publisher.go`)

   ```go
   type Publisher struct {
       redis  *redis.Client
       logger *logger.Logger
   }

   func (p *Publisher) PublishMatchEvent(ctx context.Context, event *MatchEvent) error {
       // Publishes to Redis Streams + Pub/Sub
   }
   ```

3. **Redis Streams** (Event Log)

   ```
   Stream: match:123:stream
   - Ordered event log
   - Can replay events
   - Multiple consumers
   ```

4. **Immutable Events**
   - Events are never updated (only soft-deleted)
   - `created_at` timestamp preserved
   - Metadata stored as JSONB

---

## 🚀 Migration Path to Full Event Sourcing

### Phase 1: Event Store Pattern (Easy - 1 week)

**Goal:** Make events the source of truth for match state

#### 1.1 Add Event Versioning

```sql
-- Migration: Add version and aggregate_id
ALTER TABLE match_events
ADD COLUMN aggregate_id VARCHAR(100),  -- e.g., "match:123"
ADD COLUMN aggregate_version INT,      -- Event sequence number
ADD COLUMN causation_id UUID,          -- What caused this event
ADD COLUMN correlation_id UUID;        -- Request tracking

CREATE INDEX idx_match_events_aggregate
ON match_events(aggregate_id, aggregate_version);
```

#### 1.2 Update Event Publisher

```go
// internal/infrastructure/events/event_store.go
package events

type EventStore struct {
    queries   *sqlc.Queries
    publisher *Publisher
}

type StoredEvent struct {
    ID               int32
    AggregateID      string    // "match:123"
    AggregateVersion int32     // Sequential version
    EventType        string    // "goal", "shot", "pass"
    EventData        []byte    // JSON payload
    Metadata         map[string]string
    CausationID      string    // What caused this
    CorrelationID    string    // Request ID
    Timestamp        time.Time
}

func (es *EventStore) AppendEvent(ctx context.Context, event *StoredEvent) error {
    // 1. Validate version (optimistic locking)
    currentVersion, err := es.queries.GetAggregateVersion(ctx, event.AggregateID)
    if err != nil {
        return err
    }

    if event.AggregateVersion != currentVersion + 1 {
        return ErrConcurrencyConflict
    }

    // 2. Store event in database
    storedEvent, err := es.queries.AppendEvent(ctx, sqlc.AppendEventParams{
        AggregateID:      event.AggregateID,
        AggregateVersion: event.AggregateVersion,
        EventType:        event.EventType,
        EventData:        event.EventData,
        CausationID:      event.CausationID,
        CorrelationID:    event.CorrelationID,
    })
    if err != nil {
        return err
    }

    // 3. Publish to Redis Streams
    return es.publisher.PublishMatchEvent(ctx, storedEvent)
}

func (es *EventStore) GetEvents(ctx context.Context, aggregateID string) ([]StoredEvent, error) {
    return es.queries.GetEventsByAggregate(ctx, aggregateID)
}
```

#### 1.3 Create Projections

```go
// internal/application/projections/match_projection.go
package projections

type MatchProjection struct {
    queries *sqlc.Queries
}

// Rebuild match state from events
func (p *MatchProjection) Project(ctx context.Context, matchID int32) (*Match, error) {
    // 1. Load all events for match
    events, err := p.queries.GetEventsByAggregate(ctx, fmt.Sprintf("match:%d", matchID))
    if err != nil {
        return nil, err
    }

    // 2. Replay events to build current state
    match := &Match{ID: matchID}

    for _, event := range events {
        switch event.EventType {
        case "match_started":
            match.Status = "live"
            match.StartTime = event.Timestamp

        case "goal":
            var goalData GoalEvent
            json.Unmarshal(event.EventData, &goalData)

            if goalData.TeamID == match.HomeTeamID {
                match.HomeScore++
            } else {
                match.AwayScore++
            }

        case "match_finished":
            match.Status = "finished"
            match.EndTime = event.Timestamp
        }
    }

    return match, nil
}
```

---

### Phase 2: CQRS Pattern (Medium - 2 weeks)

**Goal:** Separate read and write models

```
┌─────────────────────────────────────────────┐
│              COMMAND SIDE                   │
│  (Write Model - Event Sourced)              │
│                                             │
│  POST /matches/123/events                   │
│         ↓                                   │
│  CommandHandler                             │
│         ↓                                   │
│  EventStore.AppendEvent()                   │
│         ↓                                   │
│  Redis Streams                              │
└─────────────────┬───────────────────────────┘
                  │
                  │ Event Published
                  │
┌─────────────────▼───────────────────────────┐
│              QUERY SIDE                     │
│  (Read Model - Optimized Views)             │
│                                             │
│  Projection Workers                         │
│         ↓                                   │
│  Update Read Models                         │
│  - matches (current state)                  │
│  - player_statistics (aggregated)           │
│  - team_statistics (aggregated)             │
│         ↓                                   │
│  GET /matches/123                           │
│  (Fast reads from materialized views)       │
└─────────────────────────────────────────────┘
```

#### 2.1 Command Handlers

```go
// internal/application/commands/create_match_event.go
package commands

type CreateMatchEventCommand struct {
    MatchID       int32
    EventType     string
    PlayerID      *int32
    TeamID        *int32
    Minute        int32
    PositionX     *float64
    PositionY     *float64
    Metadata      map[string]interface{}
    CorrelationID string
}

type CreateMatchEventHandler struct {
    eventStore *events.EventStore
    validator  *EventValidator
}

func (h *CreateMatchEventHandler) Handle(ctx context.Context, cmd *CreateMatchEventCommand) error {
    // 1. Validate command
    if err := h.validator.Validate(cmd); err != nil {
        return err
    }

    // 2. Load current aggregate version
    aggregateID := fmt.Sprintf("match:%d", cmd.MatchID)
    currentVersion, _ := h.eventStore.GetCurrentVersion(ctx, aggregateID)

    // 3. Create event
    event := &events.StoredEvent{
        AggregateID:      aggregateID,
        AggregateVersion: currentVersion + 1,
        EventType:        cmd.EventType,
        EventData:        marshalEventData(cmd),
        CorrelationID:    cmd.CorrelationID,
        Timestamp:        time.Now(),
    }

    // 4. Append to event store
    return h.eventStore.AppendEvent(ctx, event)
}
```

#### 2.2 Projection Workers

```go
// internal/infrastructure/workers/projection_worker.go
package workers

type ProjectionWorker struct {
    redis       *redis.Client
    projections []Projection
    logger      *logger.Logger
}

func (w *ProjectionWorker) Start(ctx context.Context) {
    // Subscribe to Redis Streams
    streams := []string{"match:*:stream"}

    for {
        select {
        case <-ctx.Done():
            return

        default:
            // Read events from streams
            events, err := w.redis.XRead(ctx, &redis.XReadArgs{
                Streams: streams,
                Block:   time.Second,
            }).Result()

            if err != nil {
                continue
            }

            // Process each event
            for _, stream := range events {
                for _, message := range stream.Messages {
                    w.processEvent(ctx, message)
                }
            }
        }
    }
}

func (w *ProjectionWorker) processEvent(ctx context.Context, msg redis.XMessage) {
    // Apply event to all projections
    for _, projection := range w.projections {
        if err := projection.Apply(ctx, msg); err != nil {
            w.logger.Error("Projection failed", "error", err, "projection", projection.Name())
        }
    }
}
```

---

### Phase 3: Full Event Sourcing (Advanced - 4 weeks)

**Goal:** Complete event-sourced system with snapshots

#### 3.1 Snapshots (Performance Optimization)

```sql
CREATE TABLE aggregate_snapshots (
    aggregate_id VARCHAR(100) PRIMARY KEY,
    aggregate_version INT NOT NULL,
    snapshot_data JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

```go
// Load aggregate with snapshot optimization
func (es *EventStore) LoadAggregate(ctx context.Context, aggregateID string) (*Match, error) {
    // 1. Try to load latest snapshot
    snapshot, err := es.queries.GetLatestSnapshot(ctx, aggregateID)

    var match *Match
    var fromVersion int32

    if err == nil {
        // Start from snapshot
        json.Unmarshal(snapshot.SnapshotData, &match)
        fromVersion = snapshot.AggregateVersion
    } else {
        // Start from scratch
        match = &Match{}
        fromVersion = 0
    }

    // 2. Load events after snapshot
    events, err := es.queries.GetEventsAfterVersion(ctx, aggregateID, fromVersion)
    if err != nil {
        return nil, err
    }

    // 3. Replay events
    for _, event := range events {
        match.Apply(event)
    }

    return match, nil
}
```

#### 3.2 Event Upcasting (Schema Evolution)

```go
// Handle event schema changes
type EventUpcaster struct {
    upcasters map[string]func([]byte) ([]byte, error)
}

func (u *EventUpcaster) Upcast(eventType string, version int, data []byte) ([]byte, error) {
    key := fmt.Sprintf("%s:v%d", eventType, version)

    if upcaster, exists := u.upcasters[key]; exists {
        return upcaster(data)
    }

    return data, nil
}

// Example: goal event v1 → v2
func upcastGoalEventV1ToV2(data []byte) ([]byte, error) {
    var v1 GoalEventV1
    json.Unmarshal(data, &v1)

    v2 := GoalEventV2{
        PlayerID:  v1.PlayerID,
        TeamID:    v1.TeamID,
        Minute:    v1.Minute,
        XG:        0.0,  // New field, default value
        ShotType:  "unknown",  // New field
    }

    return json.Marshal(v2)
}
```

---

## 🎯 Quick Wins (Start Here!)

### 1. Add Event Versioning (1 day)

```sql
-- Run this migration
ALTER TABLE match_events
ADD COLUMN aggregate_version INT DEFAULT 1,
ADD COLUMN correlation_id UUID;

CREATE INDEX idx_match_events_version
ON match_events(match_id, aggregate_version);
```

### 2. Create Event Replay Endpoint (1 day)

```go
// GET /api/v1/matches/:id/replay
func (h *MatchHandler) ReplayMatch(c *gin.Context) {
    matchID := c.Param("id")

    // Get all events
    events, err := h.queries.GetMatchEvents(ctx, matchID)

    // Replay to build state
    match := &Match{ID: matchID}
    timeline := []MatchState{}

    for _, event := range events {
        match.Apply(event)
        timeline = append(timeline, match.Clone())
    }

    c.JSON(200, gin.H{
        "match": match,
        "timeline": timeline,
    })
}
```

### 3. Add Event Metadata (1 day)

```go
// Enrich events with metadata
type EventMetadata struct {
    UserID        int32     `json:"user_id"`
    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    CorrelationID string    `json:"correlation_id"`
    Timestamp     time.Time `json:"timestamp"`
}

// Store in match_events.metadata JSONB column
```

---

## 📊 Benefits of Event Sourcing for Football Analytics

### 1. **Time Travel**

```go
// Get match state at minute 45
state := match.StateAtMinute(45)

// Compare team performance in first vs second half
firstHalf := match.StateAtMinute(45)
secondHalf := match.FinalState()
```

### 2. **Audit Trail**

```
Who scored? When? From where? What was the xG?
All questions answered from event log
```

### 3. **New Analytics Without Migration**

```
Want to add "pass completion rate"?
Just replay all pass events - no database migration needed!
```

### 4. **Debugging**

```
Bug in xG calculation? Replay events with fixed algorithm
```

### 5. **Machine Learning**

```
Feed event stream directly to ML models
Train on historical event sequences
```

---

## 🛠️ Tools & Libraries

### Event Store

- **PostgreSQL** (current) - Good for small-medium scale
- **EventStoreDB** - Purpose-built event store
- **Apache Kafka** - Distributed event log

### CQRS

- **Watermill** - Go library for CQRS/Event Sourcing
- **go-cqrs** - Lightweight CQRS framework

### Projections

- **Redis Streams** (current) - Already set up!
- **Kafka Streams** - For complex projections

---

## 📋 Migration Checklist

### Week 1: Foundation

- [ ] Add `aggregate_version` column to `match_events`
- [ ] Add `correlation_id` for request tracking
- [ ] Create `GetEventsByAggregate` sqlc query
- [ ] Implement event replay endpoint
- [ ] Add event metadata enrichment

### Week 2: Event Store

- [ ] Create `EventStore` struct
- [ ] Implement `AppendEvent` with version checking
- [ ] Add optimistic locking for concurrency
- [ ] Create event replay function
- [ ] Add unit tests for event store

### Week 3: Projections

- [ ] Create `MatchProjection` to rebuild state
- [ ] Create `StatisticsProjection` for aggregates
- [ ] Implement projection worker
- [ ] Subscribe to Redis Streams
- [ ] Handle projection errors gracefully

### Week 4: CQRS

- [ ] Separate command handlers
- [ ] Separate query handlers
- [ ] Create read models (materialized views)
- [ ] Update API to use CQRS pattern
- [ ] Performance testing

---

## 🎓 Learning Resources

- [Event Sourcing by Martin Fowler](https://martinfowler.com/eaaDev/EventSourcing.html)
- [CQRS Journey by Microsoft](<https://docs.microsoft.com/en-us/previous-versions/msp-n-p/jj554200(v=pandp.10)>)
- [Watermill Documentation](https://watermill.io/)
- [EventStoreDB](https://www.eventstore.com/)

---

## ⚠️ Considerations

### When NOT to Use Event Sourcing

- ❌ Simple CRUD applications
- ❌ Small datasets (< 10k events)
- ❌ Team unfamiliar with pattern
- ❌ Tight deadlines

### When TO Use Event Sourcing

- ✅ Audit trail required (football matches!)
- ✅ Time travel needed (replay matches)
- ✅ Complex domain logic (analytics)
- ✅ Multiple read models (stats, reports, ML)
- ✅ Event-driven by nature (sports!)

---

**Status:** 🟡 80% Ready - Event table + Redis Streams already in place!  
**Next Step:** Add event versioning and create replay endpoint  
**Estimated Time:** 1-4 weeks depending on scope
