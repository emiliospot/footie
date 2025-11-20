# 🔍 OpenSearch Integration Guide

> **Advanced search and real-time analytics for football data**

Amazon OpenSearch Service (formerly Elasticsearch) is perfect for football analytics because it excels at complex queries, aggregations, and real-time dashboards.

---

## 🎯 Why OpenSearch for Football Analytics?

### The Problem with PostgreSQL Alone

```sql
-- This query is SLOW in PostgreSQL:
SELECT player_id,
       COUNT(*) as progressive_passes,
       AVG(distance) as avg_distance
FROM passes
WHERE pass_type = 'progressive'
  AND match_date >= NOW() - INTERVAL '10 games'
  AND success = true
GROUP BY player_id
HAVING COUNT(*) > 20
ORDER BY progressive_passes DESC;
```

**Issues:**

- ❌ Slow on millions of events
- ❌ Complex aggregations are expensive
- ❌ Not optimized for analytics
- ❌ Hard to do fuzzy search
- ❌ Time-series queries are complex

### The OpenSearch Solution

```json
POST /match_events/_search
{
  "query": {
    "bool": {
      "must": [
        { "term": { "event_type": "pass" }},
        { "term": { "pass_type": "progressive" }},
        { "term": { "success": true }},
        { "range": { "timestamp": { "gte": "now-10d" }}}
      ]
    }
  },
  "aggs": {
    "by_player": {
      "terms": { "field": "player_id" },
      "aggs": {
        "avg_distance": { "avg": { "field": "distance" }}
      }
    }
  }
}
```

**Benefits:**

- ✅ Sub-second response
- ✅ Scales to billions of events
- ✅ Built for aggregations
- ✅ Real-time indexing
- ✅ Fuzzy search included

---

## 🏗️ Architecture: The Perfect Trio

```
┌─────────────────────────────────────────────────────────────────┐
│                    MATCH EVENT FLOW                             │
└─────────────────────────────────────────────────────────────────┘

External Feed / API Call
    ↓
POST /api/v1/matches/:id/events
    ↓
┌─────────────────────────────────────────────────────────────────┐
│                  GOLANG BACKEND                                 │
│                                                                 │
│  1. Save to PostgreSQL (source of truth)                       │
│  2. Publish to Redis Streams (event log)                       │
│  3. Publish to Redis Pub/Sub (WebSocket)                       │
└─────────────────────────────────────────────────────────────────┘
    │           │               │
    ▼           ▼               ▼
┌──────────┐ ┌──────────┐ ┌──────────────┐
│PostgreSQL│ │  Redis   │ │  WebSocket   │
│          │ │ Streams  │ │     Hub      │
│Source of │ │          │ │              │
│  Truth   │ └────┬─────┘ │ Real-time    │
│          │      │       │ Broadcasts   │
└──────────┘      │       └──────────────┘
                  │
                  ▼
         ┌────────────────┐
         │ Analytics      │
         │ Worker         │
         │ (Go Service)   │
         └────────┬───────┘
                  │
                  ▼
         ┌────────────────┐
         │   OpenSearch   │
         │                │
         │ • Full-text    │
         │ • Analytics    │
         │ • Aggregations │
         │ • Dashboards   │
         └────────────────┘
```

### The Three-Database Pattern

| Database       | Purpose             | Use Cases                                 |
| -------------- | ------------------- | ----------------------------------------- |
| **PostgreSQL** | Source of truth     | CRUD, transactions, relationships         |
| **Redis**      | Real-time messaging | WebSocket, caching, pub/sub               |
| **OpenSearch** | Analytics & search  | Complex queries, aggregations, dashboards |

---

## 🚀 Use Cases for Football Analytics

### 1. Advanced Player Search

```json
// Find players similar to "Pedri"
POST /players/_search
{
  "query": {
    "more_like_this": {
      "fields": ["playing_style", "position", "stats"],
      "like": [{ "_id": "pedri_id" }],
      "min_term_freq": 1,
      "max_query_terms": 12
    }
  }
}
```

### 2. Heat Maps

```json
// Get all shots in specific area
POST /match_events/_search
{
  "query": {
    "bool": {
      "must": [
        { "term": { "event_type": "shot" }},
        { "term": { "match_id": 123 }},
        {
          "geo_bounding_box": {
            "position": {
              "top_left": { "lat": 80, "lon": 40 },
              "bottom_right": { "lat": 100, "lon": 60 }
            }
          }
        }
      ]
    }
  }
}
```

### 3. xG Trends

```json
// Average xG per minute bucket
POST /match_events/_search
{
  "query": { "term": { "event_type": "shot" }},
  "aggs": {
    "xg_per_minute": {
      "histogram": {
        "field": "minute",
        "interval": 5
      },
      "aggs": {
        "avg_xg": { "avg": { "field": "metadata.xG" }}
      }
    }
  }
}
```

### 4. Pass Networks

```json
// Most common pass combinations
POST /match_events/_search
{
  "query": { "term": { "event_type": "pass" }},
  "aggs": {
    "pass_combinations": {
      "composite": {
        "sources": [
          { "from": { "terms": { "field": "player_id" }}},
          { "to": { "terms": { "field": "secondary_player_id" }}}
        ]
      }
    }
  }
}
```

### 5. Fuzzy Player Search

```json
// "Messy" → finds "Messi"
POST /players/_search
{
  "query": {
    "fuzzy": {
      "name": {
        "value": "Messy",
        "fuzziness": "AUTO"
      }
    }
  }
}
```

---

## 💻 Implementation

### Step 1: Create Analytics Worker

```go
// internal/workers/analytics_worker.go
package workers

import (
    "context"
    "encoding/json"

    "github.com/opensearch-project/opensearch-go/v2"
    "github.com/redis/go-redis/v9"
)

type AnalyticsWorker struct {
    redis      *redis.Client
    opensearch *opensearch.Client
    logger     *logger.Logger
}

func NewAnalyticsWorker(
    redis *redis.Client,
    opensearch *opensearch.Client,
    logger *logger.Logger,
) *AnalyticsWorker {
    return &AnalyticsWorker{
        redis:      redis,
        opensearch: opensearch,
        logger:     logger,
    }
}

func (w *AnalyticsWorker) Start(ctx context.Context) {
    // Subscribe to Redis Streams
    for {
        select {
        case <-ctx.Done():
            return

        default:
            // Read from Redis Streams
            streams, err := w.redis.XRead(ctx, &redis.XReadArgs{
                Streams: []string{"match:*:stream", "0"},
                Block:   time.Second,
            }).Result()

            if err != nil {
                continue
            }

            // Process each event
            for _, stream := range streams {
                for _, message := range stream.Messages {
                    w.indexEvent(ctx, message)
                }
            }
        }
    }
}

func (w *AnalyticsWorker) indexEvent(ctx context.Context, msg redis.XMessage) {
    // Parse event
    var event MatchEvent
    json.Unmarshal([]byte(msg.Values["data"].(string)), &event)

    // Transform for OpenSearch
    doc := map[string]interface{}{
        "match_id":    event.MatchID,
        "player_id":   event.PlayerID,
        "team_id":     event.TeamID,
        "event_type":  event.EventType,
        "minute":      event.Minute,
        "position":    map[string]float64{
            "lat": *event.PositionX,
            "lon": *event.PositionY,
        },
        "metadata":    event.Metadata,
        "timestamp":   event.Timestamp,
    }

    // Index to OpenSearch
    body, _ := json.Marshal(doc)
    _, err := w.opensearch.Index(
        "match_events",
        bytes.NewReader(body),
        w.opensearch.Index.WithContext(ctx),
    )

    if err != nil {
        w.logger.Error("Failed to index event", "error", err)
    }
}
```

### Step 2: Create OpenSearch Client

```go
// internal/infrastructure/opensearch/client.go
package opensearch

import (
    "crypto/tls"
    "net/http"

    "github.com/opensearch-project/opensearch-go/v2"
)

func NewClient(endpoint, username, password string) (*opensearch.Client, error) {
    cfg := opensearch.Config{
        Addresses: []string{endpoint},
        Username:  username,
        Password:  password,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
            },
        },
    }

    return opensearch.NewClient(cfg)
}
```

### Step 3: Create Search Service

```go
// internal/application/services/search_service.go
package services

type SearchService struct {
    opensearch *opensearch.Client
    logger     *logger.Logger
}

func (s *SearchService) SearchPlayers(query string) ([]Player, error) {
    // Build fuzzy search query
    searchBody := map[string]interface{}{
        "query": map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":     query,
                "fields":    []string{"name^2", "team", "position"},
                "fuzziness": "AUTO",
            },
        },
    }

    body, _ := json.Marshal(searchBody)

    res, err := s.opensearch.Search(
        s.opensearch.Search.WithIndex("players"),
        s.opensearch.Search.WithBody(bytes.NewReader(body)),
    )

    // Parse results...
    return players, nil
}

func (s *SearchService) GetPlayerHeatMap(playerID int, matchID int) (HeatMap, error) {
    // Query all events for player in match
    searchBody := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {"term": map[string]interface{}{"player_id": playerID}},
                    {"term": map[string]interface{}{"match_id": matchID}},
                },
            },
        },
        "aggs": map[string]interface{}{
            "position_grid": map[string]interface{}{
                "geohash_grid": map[string]interface{}{
                    "field":     "position",
                    "precision": 5,
                },
            },
        },
    }

    // Execute and return heat map data
    return heatMap, nil
}
```

### Step 4: Add Configuration

```go
// internal/config/config.go

type Config struct {
    // ... existing fields
    OpenSearch OpenSearchConfig
}

type OpenSearchConfig struct {
    Endpoint string
    Username string
    Password string
}

func Load() (*Config, error) {
    // ... existing code

    cfg.OpenSearch = OpenSearchConfig{
        Endpoint: getEnv("OPENSEARCH_ENDPOINT", ""),
        Username: getEnv("OPENSEARCH_USERNAME", "admin"),
        Password: getEnv("OPENSEARCH_PASSWORD", ""),
    }

    return cfg, nil
}
```

---

## 🔧 AWS Setup

### 1. Create OpenSearch Domain

```bash
aws opensearch create-domain \
  --domain-name footie-analytics \
  --engine-version OpenSearch_2.11 \
  --cluster-config \
    InstanceType=t3.small.search,\
    InstanceCount=2 \
  --ebs-options \
    EBSEnabled=true,\
    VolumeType=gp3,\
    VolumeSize=20 \
  --access-policies '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Principal": {"AWS": "*"},
      "Action": "es:*",
      "Resource": "arn:aws:es:us-east-1:ACCOUNT:domain/footie-analytics/*"
    }]
  }'
```

### 2. Create Index Templates

```json
PUT _index_template/match_events_template
{
  "index_patterns": ["match_events*"],
  "template": {
    "settings": {
      "number_of_shards": 2,
      "number_of_replicas": 1
    },
    "mappings": {
      "properties": {
        "match_id": { "type": "integer" },
        "player_id": { "type": "integer" },
        "team_id": { "type": "integer" },
        "event_type": { "type": "keyword" },
        "minute": { "type": "integer" },
        "position": { "type": "geo_point" },
        "metadata": { "type": "object" },
        "timestamp": { "type": "date" }
      }
    }
  }
}
```

---

## 📊 Cost Estimation

### AWS OpenSearch Pricing (us-east-1)

| Instance Type    | vCPU | RAM   | Storage | Cost/Month |
| ---------------- | ---- | ----- | ------- | ---------- |
| t3.small.search  | 2    | 2 GB  | 20 GB   | ~$35       |
| t3.medium.search | 2    | 4 GB  | 50 GB   | ~$70       |
| r6g.large.search | 2    | 16 GB | 100 GB  | ~$180      |

**For Development:** t3.small (2 nodes) = ~$70/month
**For Production:** r6g.large (3 nodes) = ~$540/month

---

## 🎯 When to Add OpenSearch?

### ✅ Add OpenSearch When:

- You have > 1 million events
- You need complex analytics queries
- You want real-time dashboards
- You need full-text search
- You're doing ML/AI on event data

### ⏳ Don't Add Yet If:

- You have < 100k events (PostgreSQL is fine)
- You only need simple queries
- Budget is very tight
- Team is small (< 3 developers)

---

## 📚 Resources

- [Amazon OpenSearch Service](https://aws.amazon.com/opensearch-service/)
- [OpenSearch Go Client](https://github.com/opensearch-project/opensearch-go)
- [OpenSearch Query DSL](https://opensearch.org/docs/latest/query-dsl/)
- [Geo Queries](https://opensearch.org/docs/latest/query-dsl/geo-and-xy/)

---

**Status:** 🟡 Future Enhancement (Phase 3)  
**Complexity:** 🟠 Medium (2-3 weeks)  
**Cost:** 💰 $70-500/month (AWS)  
**Value:** 🚀 High (for analytics-heavy platforms)
