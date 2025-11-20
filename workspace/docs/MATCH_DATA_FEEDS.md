# ⚽ Match Data Feed Integration

> **How to ingest live match events from external providers**

This guide covers integrating external football data feeds (Opta, StatsBomb, API-Football, etc.) into your real-time analytics platform.

---

## 🎯 Architecture Options

### Option A: Simple (Current - Perfect for MVP)

**For:** < 100 events/sec, single region, getting started

```
┌─────────────────────────────────────────────────────────────┐
│              EXTERNAL DATA FEED                             │
│  (Opta / StatsBomb / API-Football / Custom)                 │
└─────────────────┬───────────────────────────────────────────┘
                  │ HTTP Webhook / Polling
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              YOUR GOLANG BACKEND                            │
│  POST /api/v1/matches/:id/events                           │
│                                                             │
│  MatchHandler.CreateMatchEvent()                           │
│         ↓                                                   │
│  1. Save to PostgreSQL (sqlc)                              │
│  2. Publish to Redis Streams (analytics)                   │
│  3. Publish to Redis Pub/Sub (WebSocket)                   │
└─────────────────┬───────────────────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌──────────────┐    ┌──────────────┐
│  PostgreSQL  │    │    Redis     │
│              │    │              │
│ • Permanent  │    │ • Streams    │
│   storage    │    │ • Pub/Sub    │
└──────────────┘    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ WebSocket Hub│
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Clients    │
                    └──────────────┘
```

**✅ Simple, fast, and perfect for getting started!**

---

### Option B: Production Scale (AWS-Native)

**For:** > 1000 events/sec, multiple regions, high availability

```
┌─────────────────────────────────────────────────────────────┐
│              EXTERNAL DATA FEEDS                            │
│  (Opta / StatsBomb / API-Football / Custom)                 │
└─────────────────┬───────────────────────────────────────────┘
                  │ HTTP Webhooks
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS API GATEWAY                                │
│  • Rate limiting                                            │
│  • Authentication                                           │
│  • Request validation                                       │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS LAMBDA                                     │
│  • Validate webhook signature                               │
│  • Transform payload                                        │
│  • Publish to Kinesis                                       │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS KINESIS DATA STREAMS                       │
│  • Event buffer (ordered, replay)                           │
│  • 1000s events/sec throughput                              │
│  • Partitioned by match_id                                  │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│              GO KINESIS CONSUMER                            │
│  • Processes events in order                                │
│  • Auto-scaling based on shard count                        │
│  • Checkpoint management                                    │
└─────────────────┬───────────────────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌──────────────┐    ┌──────────────┐
│ RDS Postgres │    │ElastiCache   │
│              │    │   Redis      │
│ • Permanent  │    │              │
│   storage    │    │ • Streams    │
│              │    │ • Pub/Sub    │
└──────────────┘    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ WebSocket Hub│
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Clients    │
                    └──────────────┘
```

**Why Kinesis?**

- ✅ Handles 1000s events/sec (betting companies use this)
- ✅ Ordered processing (critical for match events)
- ✅ Replay capability (reprocess events if needed)
- ✅ Auto-scaling (handles traffic spikes)
- ✅ Decouples ingestion from processing

**When to use:**

- High event volume (> 1000 events/sec)
- Multiple data sources
- Need for event replay
- Production-grade reliability

---

---

## 🔌 Integration Methods

### Method 1: Webhooks (Recommended - Real-Time)

**Best for:** Opta, StatsBomb, custom feeds

```
External Feed → Webhook → Your API → Database + Redis → WebSocket
```

**Pros:**

- ✅ Real-time (< 1 second)
- ✅ No polling overhead
- ✅ Push-based (efficient)

**Implementation:**

```go
// POST /api/v1/webhooks/match-events
func (h *WebhookHandler) ReceiveMatchEvent(c *gin.Context) {
    var payload ExternalEventPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(400, gin.H{"error": "Invalid payload"})
        return
    }

    // 1. Validate webhook signature (security)
    if !h.validateSignature(c.GetHeader("X-Signature"), payload) {
        c.JSON(401, gin.H{"error": "Invalid signature"})
        return
    }

    // 2. Transform external format to internal format
    event := h.transformEvent(payload)

    // 3. Save to database
    savedEvent, err := h.queries.CreateMatchEvent(ctx, event)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to save event"})
        return
    }

    // 4. Publish to Redis (async)
    go h.publisher.PublishMatchEvent(ctx, savedEvent)

    c.JSON(200, gin.H{"status": "received"})
}
```

### Method 2: Polling (Fallback - Near Real-Time)

**Best for:** APIs without webhooks

```
Cron Job → Poll API → Your API → Database + Redis → WebSocket
```

**Pros:**

- ✅ Works with any API
- ✅ You control rate
- ✅ Simple to implement

**Implementation:**

```go
// Worker that polls external API every 10 seconds
func (w *MatchFeedWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return

        case <-ticker.C:
            w.pollLiveMatches(ctx)
        }
    }
}

func (w *MatchFeedWorker) pollLiveMatches(ctx context.Context) {
    // 1. Get list of live matches
    liveMatches, _ := w.queries.GetLiveMatches(ctx)

    for _, match := range liveMatches {
        // 2. Poll external API for new events
        events, err := w.externalAPI.GetMatchEvents(match.ExternalID)
        if err != nil {
            continue
        }

        // 3. Process new events
        for _, event := range events {
            // Check if event already exists
            exists, _ := w.queries.EventExists(ctx, event.ExternalID)
            if exists {
                continue
            }

            // Save and publish
            savedEvent, _ := w.queries.CreateMatchEvent(ctx, event)
            go w.publisher.PublishMatchEvent(ctx, savedEvent)
        }
    }
}
```

### Method 3: WebSocket Feed (Advanced - Ultra Real-Time)

**Best for:** Premium data providers with WebSocket streams

```
External WebSocket → Your Backend → Database + Redis → Your WebSocket
```

**Pros:**

- ✅ Ultra real-time (< 100ms)
- ✅ Bidirectional
- ✅ Most efficient

**Implementation:**

```go
func (w *ExternalFeedClient) Connect(ctx context.Context) {
    conn, _, err := websocket.DefaultDialer.Dial(w.feedURL, nil)
    if err != nil {
        w.logger.Error("Failed to connect to feed", "error", err)
        return
    }
    defer conn.Close()

    for {
        select {
        case <-ctx.Done():
            return

        default:
            var event ExternalEvent
            if err := conn.ReadJSON(&event); err != nil {
                w.logger.Error("Failed to read event", "error", err)
                continue
            }

            // Process event
            w.processEvent(ctx, event)
        }
    }
}
```

---

---

## 🚀 Production Implementation (AWS Lambda + Kinesis)

### AWS Lambda Webhook Receiver

```javascript
// lambda/webhook-receiver/index.js
const AWS = require("aws-sdk");
const kinesis = new AWS.Kinesis();

exports.handler = async (event) => {
  try {
    // 1. Parse webhook payload
    const payload = JSON.parse(event.body);

    // 2. Validate signature
    const signature = event.headers["X-Signature"];
    if (!validateSignature(payload, signature)) {
      return {
        statusCode: 401,
        body: JSON.stringify({ error: "Invalid signature" }),
      };
    }

    // 3. Transform to internal format
    const matchEvent = {
      match_id: payload.match_id,
      event_type: payload.event_type,
      minute: payload.minute,
      player_id: payload.player_id,
      timestamp: new Date().toISOString(),
      metadata: payload,
    };

    // 4. Publish to Kinesis
    await kinesis
      .putRecord({
        StreamName: "match-events-stream",
        PartitionKey: payload.match_id, // Events for same match go to same shard
        Data: JSON.stringify(matchEvent),
      })
      .promise();

    return {
      statusCode: 200,
      body: JSON.stringify({ status: "received" }),
    };
  } catch (error) {
    console.error("Error processing webhook:", error);
    return {
      statusCode: 500,
      body: JSON.stringify({ error: "Internal error" }),
    };
  }
};

function validateSignature(payload, signature) {
  const crypto = require("crypto");
  const secret = process.env.WEBHOOK_SECRET;
  const hmac = crypto.createHmac("sha256", secret);
  hmac.update(JSON.stringify(payload));
  const expectedSignature = hmac.digest("hex");
  return signature === expectedSignature;
}
```

### Go Kinesis Consumer

```go
// internal/workers/kinesis_consumer.go
package workers

import (
    "context"
    "encoding/json"

    "github.com/aws/aws-sdk-go-v2/service/kinesis"
    "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type KinesisConsumer struct {
    client    *kinesis.Client
    queries   *sqlc.Queries
    publisher *events.Publisher
    logger    *logger.Logger
}

func NewKinesisConsumer(
    client *kinesis.Client,
    queries *sqlc.Queries,
    publisher *events.Publisher,
    logger *logger.Logger,
) *KinesisConsumer {
    return &KinesisConsumer{
        client:    client,
        queries:   queries,
        publisher: publisher,
        logger:    logger,
    }
}

func (kc *KinesisConsumer) Start(ctx context.Context) error {
    streamName := "match-events-stream"

    // Get shard iterator
    shardIterator, err := kc.client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
        StreamName:        &streamName,
        ShardId:           aws.String("shardId-000000000000"),
        ShardIteratorType: types.ShardIteratorTypeLatest,
    })
    if err != nil {
        return err
    }

    iterator := shardIterator.ShardIterator

    for {
        select {
        case <-ctx.Done():
            return nil

        default:
            // Get records from Kinesis
            output, err := kc.client.GetRecords(ctx, &kinesis.GetRecordsInput{
                ShardIterator: iterator,
                Limit:         aws.Int32(100),
            })
            if err != nil {
                kc.logger.Error("Failed to get records", "error", err)
                continue
            }

            // Process each record
            for _, record := range output.Records {
                kc.processRecord(ctx, record)
            }

            // Update iterator for next batch
            iterator = output.NextShardIterator

            // Sleep if no records
            if len(output.Records) == 0 {
                time.Sleep(1 * time.Second)
            }
        }
    }
}

func (kc *KinesisConsumer) processRecord(ctx context.Context, record types.Record) {
    var event MatchEvent
    if err := json.Unmarshal(record.Data, &event); err != nil {
        kc.logger.Error("Failed to unmarshal event", "error", err)
        return
    }

    // 1. Save to PostgreSQL
    savedEvent, err := kc.queries.CreateMatchEvent(ctx, sqlc.CreateMatchEventParams{
        MatchID:   event.MatchID,
        EventType: event.EventType,
        Minute:    event.Minute,
        PlayerID:  event.PlayerID,
        // ... other fields
    })
    if err != nil {
        kc.logger.Error("Failed to save event", "error", err)
        return
    }

    // 2. Publish to Redis for WebSocket
    go kc.publisher.PublishMatchEvent(ctx, savedEvent)

    kc.logger.Info("Processed event", "match_id", event.MatchID, "type", event.EventType)
}
```

### Terraform Configuration

```hcl
# infra/terraform/kinesis.tf
resource "aws_kinesis_stream" "match_events" {
  name             = "match-events-stream"
  shard_count      = 2
  retention_period = 24

  shard_level_metrics = [
    "IncomingBytes",
    "IncomingRecords",
    "OutgoingBytes",
    "OutgoingRecords",
  ]

  tags = {
    Environment = var.environment
    Application = "footie"
  }
}

# Lambda function
resource "aws_lambda_function" "webhook_receiver" {
  filename      = "lambda/webhook-receiver.zip"
  function_name = "footie-webhook-receiver"
  role          = aws_iam_role.lambda_role.arn
  handler       = "index.handler"
  runtime       = "nodejs18.x"

  environment {
    variables = {
      KINESIS_STREAM_NAME = aws_kinesis_stream.match_events.name
      WEBHOOK_SECRET      = var.webhook_secret
    }
  }
}

# API Gateway
resource "aws_apigatewayv2_api" "webhook_api" {
  name          = "footie-webhook-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id           = aws_apigatewayv2_api.webhook_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.webhook_receiver.invoke_arn
}

resource "aws_apigatewayv2_route" "webhook" {
  api_id    = aws_apigatewayv2_api.webhook_api.id
  route_key = "POST /webhooks/match-events"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}
```

---

## 📋 Simple Implementation Guide (No AWS)

### Step 1: Create Webhook Handler

```go
// internal/api/handlers/webhook.go
package handlers

type WebhookHandler struct {
    *BaseHandler
    validator *WebhookValidator
}

type ExternalEventPayload struct {
    MatchID     string    `json:"match_id"`
    EventType   string    `json:"event_type"`
    Minute      int       `json:"minute"`
    PlayerID    string    `json:"player_id"`
    TeamID      string    `json:"team_id"`
    PositionX   float64   `json:"position_x"`
    PositionY   float64   `json:"position_y"`
    Metadata    string    `json:"metadata"`
    Timestamp   time.Time `json:"timestamp"`
    Signature   string    `json:"signature"`
}

func NewWebhookHandler(base *BaseHandler) *WebhookHandler {
    return &WebhookHandler{
        BaseHandler: base,
        validator:   NewWebhookValidator(base.cfg.Webhook.Secret),
    }
}

func (h *WebhookHandler) ReceiveMatchEvent(c *gin.Context) {
    var payload ExternalEventPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        h.logger.Error("Invalid webhook payload", "error", err)
        c.JSON(400, gin.H{"error": "Invalid payload"})
        return
    }

    // Validate signature
    if !h.validator.Validate(payload) {
        h.logger.Warn("Invalid webhook signature", "match_id", payload.MatchID)
        c.JSON(401, gin.H{"error": "Invalid signature"})
        return
    }

    // Transform to internal format
    event := h.transformEvent(payload)

    // Save to database
    savedEvent, err := h.queries.CreateMatchEvent(c.Request.Context(), event)
    if err != nil {
        h.logger.Error("Failed to save event", "error", err)
        c.JSON(500, gin.H{"error": "Failed to save event"})
        return
    }

    // Publish to Redis (async)
    go func() {
        if err := h.publisher.PublishMatchEvent(context.Background(), savedEvent); err != nil {
            h.logger.Error("Failed to publish event", "error", err)
        }
    }()

    h.logger.Info("Webhook event received",
        "match_id", payload.MatchID,
        "event_type", payload.EventType,
    )

    c.JSON(200, gin.H{"status": "received", "event_id": savedEvent.ID})
}

func (h *WebhookHandler) transformEvent(payload ExternalEventPayload) sqlc.CreateMatchEventParams {
    // Map external IDs to internal IDs
    matchID := h.getInternalMatchID(payload.MatchID)
    playerID := h.getInternalPlayerID(payload.PlayerID)
    teamID := h.getInternalTeamID(payload.TeamID)

    var posX, posY pgtype.Numeric
    _ = posX.Scan(payload.PositionX)
    _ = posY.Scan(payload.PositionY)

    return sqlc.CreateMatchEventParams{
        MatchID:   matchID,
        PlayerID:  &playerID,
        TeamID:    &teamID,
        EventType: payload.EventType,
        Minute:    int32(payload.Minute),
        PositionX: posX,
        PositionY: posY,
        Metadata:  []byte(payload.Metadata),
    }
}
```

### Step 2: Add Webhook Routes

```go
// internal/api/router.go

// Add webhook handler
webhookHandler := handlers.NewWebhookHandler(baseHandler)

// Webhook routes (no auth - validated by signature)
webhooks := router.Group("/api/v1/webhooks")
webhooks.POST("/match-events", webhookHandler.ReceiveMatchEvent)
webhooks.POST("/match-status", webhookHandler.ReceiveMatchStatus)
```

### Step 3: Configure Webhook Secret

```go
// internal/config/config.go

type Config struct {
    // ... existing fields
    Webhook WebhookConfig
}

type WebhookConfig struct {
    Secret string
}

func Load() (*Config, error) {
    // ... existing code

    cfg.Webhook = WebhookConfig{
        Secret: getEnv("WEBHOOK_SECRET", ""),
    }

    return cfg, nil
}
```

### Step 4: Add to Environment Variables

```bash
# .env
WEBHOOK_SECRET=your-secret-key-here
```

---

## 🔒 Security Best Practices

### 1. Webhook Signature Validation

```go
package handlers

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

type WebhookValidator struct {
    secret string
}

func NewWebhookValidator(secret string) *WebhookValidator {
    return &WebhookValidator{secret: secret}
}

func (v *WebhookValidator) Validate(payload ExternalEventPayload) bool {
    // Compute HMAC
    mac := hmac.New(sha256.New, []byte(v.secret))
    mac.Write([]byte(payload.String()))
    expectedSignature := hex.EncodeToString(mac.Sum(nil))

    // Compare signatures
    return hmac.Equal([]byte(expectedSignature), []byte(payload.Signature))
}
```

### 2. Rate Limiting

```go
// Limit webhooks to 1000 requests per minute
router.Use(middleware.RateLimit(1000, time.Minute))
```

### 3. IP Whitelisting

```go
func (h *WebhookHandler) validateIP(c *gin.Context) bool {
    clientIP := c.ClientIP()
    allowedIPs := []string{"203.0.113.0", "198.51.100.0"} // Provider IPs

    for _, ip := range allowedIPs {
        if clientIP == ip {
            return true
        }
    }
    return false
}
```

---

## 📊 Popular Data Providers

### 1. **Opta Sports** (Premium)

- **Type:** Webhooks + API
- **Coverage:** All major leagues
- **Latency:** < 1 second
- **Cost:** $$$$

### 2. **StatsBomb** (Premium)

- **Type:** API
- **Coverage:** Top leagues + free data
- **Latency:** Near real-time
- **Cost:** $$$ (free tier available)

### 3. **API-Football** (Affordable)

- **Type:** REST API
- **Coverage:** 1000+ leagues
- **Latency:** 10-30 seconds
- **Cost:** $ (free tier: 100 requests/day)

### 4. **Football-Data.org** (Free)

- **Type:** REST API
- **Coverage:** Major European leagues
- **Latency:** 1-2 minutes
- **Cost:** Free (limited)

### 5. **Custom Feed** (DIY)

- **Type:** Manual entry or scraping
- **Coverage:** Any
- **Latency:** Depends on you
- **Cost:** Free

---

## 🚀 Quick Start: API-Football Integration

### Step 1: Get API Key

```bash
# Sign up at https://www.api-football.com/
# Get your API key
```

### Step 2: Create Client

```go
// internal/infrastructure/external/api_football.go
package external

type APIFootballClient struct {
    apiKey     string
    httpClient *http.Client
    logger     *logger.Logger
}

func NewAPIFootballClient(apiKey string, logger *logger.Logger) *APIFootballClient {
    return &APIFootballClient{
        apiKey:     apiKey,
        httpClient: &http.Client{Timeout: 10 * time.Second},
        logger:     logger,
    }
}

func (c *APIFootballClient) GetMatchEvents(fixtureID int) ([]MatchEvent, error) {
    url := fmt.Sprintf("https://v3.football.api-sports.io/fixtures/events?fixture=%d", fixtureID)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("x-apisports-key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result APIFootballResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return c.transformEvents(result.Response), nil
}
```

### Step 3: Create Polling Worker

```go
// cmd/worker/main.go
package main

func main() {
    // Load config
    cfg, _ := config.Load()

    // Initialize dependencies
    pool, _ := database.NewPgxPool(context.Background(), cfg.Database)
    queries := sqlc.New(pool)

    // Create API client
    apiClient := external.NewAPIFootballClient(cfg.APIFootball.Key, logger)

    // Create worker
    worker := workers.NewMatchFeedWorker(queries, apiClient, logger)

    // Start polling
    worker.Start(context.Background())
}
```

---

## 📈 Monitoring & Observability

### Track Feed Health

```go
// Metrics to track
type FeedMetrics struct {
    EventsReceived    int64
    EventsProcessed   int64
    EventsFailed      int64
    AverageLatency    time.Duration
    LastEventTime     time.Time
}

func (h *WebhookHandler) recordMetrics(event ExternalEventPayload) {
    h.metrics.EventsReceived++
    h.metrics.LastEventTime = time.Now()

    // Calculate latency
    latency := time.Since(event.Timestamp)
    h.metrics.AverageLatency = (h.metrics.AverageLatency + latency) / 2
}
```

### Health Check Endpoint

```go
// GET /api/v1/feed/health
func (h *WebhookHandler) FeedHealth(c *gin.Context) {
    timeSinceLastEvent := time.Since(h.metrics.LastEventTime)

    status := "healthy"
    if timeSinceLastEvent > 5*time.Minute {
        status = "stale"
    }

    c.JSON(200, gin.H{
        "status":            status,
        "events_received":   h.metrics.EventsReceived,
        "events_processed":  h.metrics.EventsProcessed,
        "events_failed":     h.metrics.EventsFailed,
        "average_latency":   h.metrics.AverageLatency.String(),
        "last_event_time":   h.metrics.LastEventTime,
        "time_since_last":   timeSinceLastEvent.String(),
    })
}
```

---

## ✅ Your Architecture is Perfect!

You're absolutely right - you **DON'T need**:

- ❌ Complex event sourcing
- ❌ Change Data Capture (CDC)
- ❌ Microservice synchronization
- ❌ Event replay infrastructure

You **ALREADY HAVE**:

- ✅ Simple feed ingestion (webhook/polling)
- ✅ Single source of truth (PostgreSQL)
- ✅ Real-time delivery (Redis + WebSocket)
- ✅ Fast and efficient

**Your flow is perfect:**

```
Feed → Backend → PostgreSQL + Redis → WebSocket → Angular
```

---

## 🎯 Next Steps

1. **Choose a data provider** (API-Football for testing)
2. **Implement webhook handler** (30 minutes)
3. **Add signature validation** (15 minutes)
4. **Test with sample data** (1 hour)
5. **Deploy and monitor** (ongoing)

---

**Status:** ✅ Architecture Ready - Just add webhook endpoint!
**Complexity:** 🟢 Simple (exactly what you need)
**Time to Implement:** 2-4 hours
