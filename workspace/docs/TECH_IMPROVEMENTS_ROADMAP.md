# 🚀 Technical Improvements Roadmap

> **Beyond the MVP: Production-Ready Enhancements**
> **Status:** Phase 1.5 - 2.0 Planning
> **Timeline:** Post-MVP, Iterative Implementation

---

## 📋 Overview

This document outlines technical improvements beyond the current MVP stack to make the platform production-ready, scalable, and enterprise-grade. These enhancements focus on observability, security, performance, and developer experience.

---

## 🔧 Backend & Infrastructure Improvements

### 1. Observability: Production-Ready Monitoring

**Current State:** Basic logging with Go's `slog`

**Improvements:**

#### Structured Logging

- **Upgrade to:** `zap` or `zerolog` for high-performance structured logging
- **Add correlation IDs:**
  - Request ID (trace entire request lifecycle)
  - Match ID (track match-specific operations)
  - User ID (audit user actions)
- **JSON logs** for CloudWatch/OpenSearch ingestion
- **Log levels** with environment-based configuration

**Example:**

```go
logger.Info("match event created",
    zap.String("request_id", reqID),
    zap.Int("match_id", matchID),
    zap.String("event_type", "goal"),
    zap.Duration("db_latency", dbLatency),
)
```

#### Metrics

- **Prometheus + Grafana** or **AWS CloudWatch Metrics**
- **Track:**
  - Per-endpoint latency (p50, p95, p99)
  - WebSocket connection count
  - Redis cache hit/miss rate
  - Database connection pool usage
  - Event processing throughput

#### Distributed Tracing

- **OpenTelemetry (OTel)** + **AWS X-Ray** or **Tempo/Jaeger**
- **Trace flow:** HTTP Request → Handler → sqlc Query → PostgreSQL → Redis Publish → WebSocket Broadcast
- **Benefits:** Identify bottlenecks, debug slow requests, visualize dependencies

**Why PM Cares:**

- Faster incident resolution (find issues in minutes, not hours)
- SLOs/SLAs (prove 99.9% uptime)
- Proactive monitoring (catch issues before users complain)

**Timeline:** 2-3 weeks
**Priority:** High (essential for production)

---

### 2. API Hygiene & Versioning

**Current State:** `/api/v1/...` prefix exists

**Improvements:**

#### OpenAPI Specification

- **Generate OpenAPI 3.0 spec** from Go code (using Swag or manual)
- **Benefits:**
  - Auto-generate Angular client types
  - Interactive API docs (Swagger UI)
  - Client SDK generation for partners
  - Contract testing

**Example:**

```go
// @Summary Create match event
// @Description Creates a new match event and broadcasts to WebSocket clients
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param event body CreateMatchEventRequest true "Event details"
// @Success 201 {object} MatchEvent
// @Router /api/v1/matches/{id}/events [post]
func (h *MatchHandler) CreateMatchEvent(c *gin.Context) { ... }
```

#### API Versioning Strategy

- **Current:** `/api/v1/...`
- **Future:** `/api/v2/...` when breaking changes needed
- **Deprecation policy:** Support N-1 versions for 6 months

**Why PM Cares:**

- Easier partner integrations
- No breaking changes for existing clients
- Professional API documentation

**Timeline:** 1-2 weeks
**Priority:** Medium (nice to have before external partners)

---

### 3. Hardening Auth & Security

**Current State:** Basic JWT structure in place

**Improvements:**

#### JWT Best Practices

- **Short-lived access tokens** (15 minutes)
- **Refresh token rotation** (7 days, single-use)
- **Token blacklist** in Redis for logout
- **Secure cookie storage** (HttpOnly, Secure, SameSite)

#### Role-Based Access Control (RBAC)

```go
type Role string

const (
    RoleAnalyst  Role = "analyst"  // Read-only access
    RoleCoach    Role = "coach"    // Edit team/player data
    RoleAdmin    Role = "admin"    // Full access
    RoleOrgOwner Role = "org_owner" // Multi-tenant owner
)

// Middleware
func RequireRole(roles ...Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := GetUserFromContext(c)
        if !user.HasRole(roles...) {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Usage
router.POST("/teams", RequireRole(RoleCoach, RoleAdmin), teamHandler.Create)
```

#### Rate Limiting

- **Redis-based rate limiter**
- **Per IP:** 100 req/min (global)
- **Per user:** 1000 req/min (authenticated)
- **Per endpoint:** Auth endpoints (5 req/min to prevent brute force)

#### Input Validation

- **go-playground/validator** for struct validation
- **Sanitize inputs** (XSS prevention)
- **SQL injection protection** (already handled by sqlc)

**Why PM Cares:**

- Enterprise-ready security
- Compliance (GDPR, SOC 2)
- Prevents abuse and attacks

**Timeline:** 2-3 weeks
**Priority:** High (essential for production)

---

### 4. Data Model & Retention Strategy

**Current State:** Single `match_events` table

**Problem:** Will explode over time (millions of events)

**Improvements:**

#### Table Partitioning

```sql
-- Partition by season
CREATE TABLE match_events_2024_25 PARTITION OF match_events
    FOR VALUES FROM ('2024-08-01') TO ('2025-07-31');

CREATE TABLE match_events_2025_26 PARTITION OF match_events
    FOR VALUES FROM ('2025-08-01') TO ('2026-07-31');
```

**Benefits:**

- Faster queries (scan only relevant partition)
- Easier archival (drop old partitions)
- Better index performance

#### Retention Policy

- **Hot data:** Last 2 seasons in PostgreSQL (fast queries)
- **Warm data:** 2-5 years in compressed format
- **Cold data:** 5+ years archived to S3 (Parquet format for ML)

#### Pre-Computed Aggregates

```sql
-- Materialized view for player stats
CREATE MATERIALIZED VIEW player_season_stats AS
SELECT
    player_id,
    season,
    COUNT(*) FILTER (WHERE event_type = 'goal') as goals,
    COUNT(*) FILTER (WHERE event_type = 'assist') as assists,
    AVG((metadata->>'xG')::numeric) as avg_xg
FROM match_events
GROUP BY player_id, season;

-- Refresh nightly
REFRESH MATERIALIZED VIEW CONCURRENTLY player_season_stats;
```

**Why PM Cares:**

- Keeps system fast as data grows
- Marketing: "Designed for long-term, large-scale data"
- Cost-effective (archive old data to cheap storage)

**Timeline:** 3-4 weeks
**Priority:** Medium (before hitting 1M+ events)

---

### 5. Better CI/CD Story

**Current State:** GitHub Actions with linting and type-checking

**Improvements:**

#### Preview Environments

- **Per-PR preview deployments**
- **Subdomain:** `pr-123.preview.footie.com`
- **Stack:** Docker Compose or AWS ECS
- **Benefits:** Test features before merging

#### Automated DB Migrations

```yaml
# .github/workflows/deploy.yml
- name: Run database migrations
  run: |
    cd workspace/apps/api
    migrate -path migrations -database "$DATABASE_URL" up
```

#### Security Scanning

- **Go:** `gosec` (security linter)
- **Node:** `npm audit` / `pnpm audit`
- **Docker:** Trivy (container scanning)
- **Secrets:** GitGuardian or Gitleaks

#### Deployment Strategy

- **Blue-Green deployments** (zero downtime)
- **Canary releases** (gradual rollout)
- **Automatic rollback** on health check failure

**Why PM Cares:**

- Reduced deployment risk
- Faster feedback loop
- Professional development process

**Timeline:** 2-3 weeks
**Priority:** Medium (before production launch)

---

### 6. Abstraction Over Redis (EventBus Pattern)

**Current State:** Direct Redis Streams + Pub/Sub usage

**Problem:** Tightly coupled to Redis implementation

**Improvement:**

```go
// Define interface
type EventBus interface {
    PublishMatchEvent(ctx context.Context, evt MatchEvent) error
    SubscribeMatchEvents(matchID int64, handler func(MatchEvent)) error
    PublishScoreUpdate(ctx context.Context, matchID int64, score Score) error
}

// Redis implementation
type RedisEventBus struct {
    client *redis.Client
    logger *logger.Logger
}

func (r *RedisEventBus) PublishMatchEvent(ctx context.Context, evt MatchEvent) error {
    // Redis Streams + Pub/Sub implementation
}

// Future: NATS implementation
type NATSEventBus struct {
    conn *nats.Conn
}

// Future: AWS SNS/SQS implementation
type AWSEventBus struct {
    sns *sns.Client
    sqs *sqs.Client
}
```

**Benefits:**

- Swappable implementations (Redis → NATS → Kinesis)
- Easier testing (mock EventBus)
- Cloud-agnostic (not locked to Redis)

**Why PM Cares:**

- Future-proof architecture
- Flexibility to change infrastructure
- No vendor lock-in

**Timeline:** 1-2 weeks
**Priority:** Low (nice to have, not urgent)

---

## 🅰️ Angular Frontend Improvements

### 1. Nx Library Architecture

**Current State:** Monolithic `apps/web` structure

**Improvement:** Split into domain libraries

```
workspace/
├── apps/
│   └── web/                    # Shell app only
└── libs/
    ├── feature/
    │   ├── match-live/         # Live match view
    │   ├── players/            # Player management
    │   ├── teams/              # Team management
    │   └── analytics/          # Analytics dashboards
    ├── data-access/
    │   ├── api/                # HTTP services
    │   ├── websocket/          # WebSocket service
    │   └── state/              # State management
    ├── ui/
    │   ├── components/         # Shared components
    │   ├── charts/             # Chart components
    │   └── tables/             # Table components
    └── util/
        ├── shared/             # Pipes, helpers
        └── models/             # TypeScript interfaces
```

**Benefits:**

- **Faster builds:** Nx only rebuilds affected libs
- **Clear boundaries:** Feature teams can own specific libs
- **Reusability:** Shared components in `ui/` libs
- **Testing:** Test libs independently

**Why PM Cares:**

- Faster development (parallel work on features)
- Easier to scale team
- Better code organization

**Timeline:** 2-3 weeks
**Priority:** High (before codebase grows too large)

---

### 2. Real-Time State Management

**Current State:** Basic RxJS observables

**Improvement:** Robust WebSocket service with reconnection

```typescript
@Injectable({ providedIn: "root" })
export class WebSocketService {
  private socket$ = new Subject<WebSocketMessage>();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;

  constructor(private config: ConfigService) {
    this.connect();
  }

  private connect() {
    const ws = new WebSocket(this.config.wsUrl);

    ws.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;
    };

    ws.onmessage = (event) => {
      this.socket$.next(JSON.parse(event.data));
    };

    ws.onclose = () => {
      console.log("WebSocket closed, reconnecting...");
      this.reconnect();
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  private reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      setTimeout(() => {
        this.reconnectAttempts++;
        this.connect();
      }, delay);
    }
  }

  // Expose typed streams
  liveMatchEvents$(matchId: number): Observable<MatchEvent> {
    return this.socket$.pipe(
      filter((msg) => msg.type === "match_event" && msg.matchId === matchId),
      map((msg) => msg.data as MatchEvent),
    );
  }

  liveScore$(matchId: number): Observable<Score> {
    return this.socket$.pipe(
      filter((msg) => msg.type === "score_update" && msg.matchId === matchId),
      map((msg) => msg.data as Score),
    );
  }
}
```

#### Signals Integration (Angular 20)

```typescript
export class MatchLiveComponent {
  private matchId = input.required<number>();
  private wsService = inject(WebSocketService);

  // Convert Observable to Signal
  liveEvents = toSignal(this.wsService.liveMatchEvents$(this.matchId()), {
    initialValue: [],
  });

  liveScore = toSignal(this.wsService.liveScore$(this.matchId()), {
    initialValue: { home: 0, away: 0 },
  });
}
```

#### OnPush Change Detection

```typescript
@Component({
  selector: "app-match-events-list",
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div *ngFor="let event of events; trackBy: trackByEventId">
      {{ event.type }} - {{ event.minute }}'
    </div>
  `,
})
export class MatchEventsListComponent {
  @Input() events: MatchEvent[] = [];

  trackByEventId(index: number, event: MatchEvent): number {
    return event.id;
  }
}
```

**Why PM Cares:**

- Smooth UI even with many events streaming
- No lost connections (automatic reconnect)
- Better user experience

**Timeline:** 2 weeks
**Priority:** High (essential for real-time features)

---

### 3. API Typing & Client Generation

**Current State:** Hand-written TypeScript interfaces

**Improvement:** Generate from OpenAPI spec

```bash
# Install generator
npm install --save-dev @openapitools/openapi-generator-cli

# Generate Angular client
npx openapi-generator-cli generate \
  -i http://localhost:8088/swagger/doc.json \
  -g typescript-angular \
  -o libs/data-access/api/generated
```

**Benefits:**

- No drift between backend and frontend
- Auto-complete in IDE
- Compile-time type safety
- Faster development

**Alternative:** Shared types in Nx lib

```
libs/models/
├── src/
│   ├── match.model.ts
│   ├── player.model.ts
│   ├── team.model.ts
│   └── index.ts
```

**Why PM Cares:**

- Fewer bugs (types catch errors)
- Faster feature development
- Better developer experience

**Timeline:** 1 week
**Priority:** Medium (nice to have)

---

### 4. Premium UX for Analytics

**Current State:** Basic Material components

**Improvements:**

#### Skeleton Loaders

```typescript
<ng-container *ngIf="loading; else content">
  <app-skeleton-table [rows]="10" [columns]="5"></app-skeleton-table>
</ng-container>

<ng-template #content>
  <app-data-table [data]="matchEvents"></app-data-table>
</ng-template>
```

#### Virtual Scroll (for large lists)

```typescript
<cdk-virtual-scroll-viewport itemSize="50" class="events-viewport">
  <div *cdkVirtualFor="let event of matchEvents; trackBy: trackByEventId">
    <app-match-event-card [event]="event"></app-match-event-card>
  </div>
</cdk-virtual-scroll-viewport>
```

#### Live Indicators

```typescript
<div class="live-badge" *ngIf="isLive">
  <span class="pulse"></span>
  LIVE
</div>

<div class="last-updated">
  Updated {{ lastUpdated | timeAgo }}
</div>
```

```css
.pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
```

#### Replay Mode

```typescript
export class MatchReplayComponent {
  private events = input.required<MatchEvent[]>();

  currentMinute = signal(0);
  maxMinute = computed(() => Math.max(...this.events().map((e) => e.minute)));

  visibleEvents = computed(() =>
    this.events().filter((e) => e.minute <= this.currentMinute()),
  );

  onSliderChange(minute: number) {
    this.currentMinute.set(minute);
  }
}
```

**Why PM Cares:**

- Looks like a premium product, not a dev demo
- Better user engagement
- Competitive advantage (unique features)

**Timeline:** 3-4 weeks
**Priority:** Medium (after core features work)

---

### 5. Testing & Quality

**Current State:** Basic Playwright E2E tests

**Improvements:**

#### Component Test Harnesses

```typescript
import { MatTableHarness } from "@angular/material/table/testing";

it("should display player statistics", async () => {
  const table = await loader.getHarness(MatTableHarness);
  const rows = await table.getRows();

  expect(rows.length).toBe(11); // 11 players

  const firstRow = await rows[0].getCells();
  const nameCell = await firstRow[0].getText();
  expect(nameCell).toBe("Lionel Messi");
});
```

#### Playwright Real-Time Scenarios

```typescript
test("live match updates", async ({ page }) => {
  await page.goto("/matches/123");

  // Wait for WebSocket connection
  await page.waitForSelector(".live-badge");

  // Simulate event from backend
  await page.evaluate(() => {
    window.postMessage(
      {
        type: "match_event",
        data: { type: "goal", minute: 45, player: "Messi" },
      },
      "*",
    );
  });

  // Verify UI updated
  await expect(page.locator(".score")).toContainText("1-0");
  await expect(page.locator(".event-list")).toContainText("Goal - Messi");
});

test("reconnect after network loss", async ({ page, context }) => {
  await page.goto("/matches/123");

  // Simulate offline
  await context.setOffline(true);
  await expect(page.locator(".offline-indicator")).toBeVisible();

  // Simulate online
  await context.setOffline(false);
  await expect(page.locator(".live-badge")).toBeVisible();
  await expect(page.locator(".offline-indicator")).not.toBeVisible();
});
```

#### Performance Budgets

```json
// angular.json
{
  "budgets": [
    {
      "type": "initial",
      "maximumWarning": "500kb",
      "maximumError": "1mb"
    },
    {
      "type": "anyComponentStyle",
      "maximumWarning": "2kb",
      "maximumError": "4kb"
    }
  ]
}
```

**Why PM Cares:**

- Confidence in real-time behavior
- Catch regressions before production
- Maintain performance as features grow

**Timeline:** 2-3 weeks
**Priority:** High (essential for quality)

---

### 6. Theming & White-Labelling

**Current State:** Single Material theme

**Improvement:** Multi-theme architecture

```typescript
// libs/config/theme/src/themes.ts
export interface Theme {
  primary: string;
  accent: string;
  warn: string;
  logo: string;
}

export const themes: Record<string, Theme> = {
  default: {
    primary: "#1976d2",
    accent: "#ff4081",
    warn: "#f44336",
    logo: "/assets/logos/footie.svg",
  },
  "man-united": {
    primary: "#DA291C",
    accent: "#FBE122",
    warn: "#000000",
    logo: "/assets/logos/man-united.svg",
  },
  barcelona: {
    primary: "#004D98",
    accent: "#A50044",
    warn: "#EDBB00",
    logo: "/assets/logos/barcelona.svg",
  },
};
```

```scss
// styles/themes/_default.scss
@use "@angular/material" as mat;

$primary: mat.define-palette(mat.$indigo-palette);
$accent: mat.define-palette(mat.$pink-palette);
$warn: mat.define-palette(mat.$red-palette);

$theme: mat.define-light-theme(
  (
    color: (
      primary: $primary,
      accent: $accent,
      warn: $warn,
    ),
  )
);

@include mat.all-component-themes($theme);
```

**Why PM Cares:**

- Easy upsell: "We can white-label for clubs, leagues, broadcasters"
- Multi-tenant support (per-organization branding)
- Professional customization

**Timeline:** 2 weeks
**Priority:** Low (nice to have for enterprise sales)

---

## 🎯 Implementation Priority Matrix

| Improvement                                   | Impact | Effort | Priority        | Timeline  |
| --------------------------------------------- | ------ | ------ | --------------- | --------- |
| **Observability (logs, metrics, traces)**     | High   | Medium | 🔴 Critical     | 2-3 weeks |
| **Auth & Security (RBAC, rate limiting)**     | High   | Medium | 🔴 Critical     | 2-3 weeks |
| **Nx Library Architecture**                   | High   | Medium | 🔴 Critical     | 2-3 weeks |
| **Real-Time WebSocket Service**               | High   | Low    | 🔴 Critical     | 2 weeks   |
| **Testing & Quality (Harnesses, Playwright)** | High   | Medium | 🔴 Critical     | 2-3 weeks |
| **Data Retention & Partitioning**             | Medium | High   | 🟡 Important    | 3-4 weeks |
| **CI/CD Improvements**                        | Medium | Medium | 🟡 Important    | 2-3 weeks |
| **OpenAPI Spec & Client Generation**          | Medium | Low    | 🟡 Important    | 1-2 weeks |
| **Premium UX (skeletons, virtual scroll)**    | Medium | Medium | 🟡 Important    | 3-4 weeks |
| **EventBus Abstraction**                      | Low    | Low    | 🟢 Nice to Have | 1-2 weeks |
| **Theming & White-Labelling**                 | Low    | Low    | 🟢 Nice to Have | 2 weeks   |

---

## 📊 Phased Rollout

### Phase 1.5 (Post-MVP, Pre-Production) - 6-8 weeks

**Goal:** Make it production-ready

- ✅ Observability (logs, metrics, traces)
- ✅ Auth & Security hardening
- ✅ Nx library architecture
- ✅ Real-time WebSocket service
- ✅ Testing improvements

**Outcome:** Ready for beta users

---

### Phase 2.0 (Production Launch) - 8-10 weeks

**Goal:** Scale and polish

- ✅ Data retention strategy
- ✅ CI/CD improvements
- ✅ OpenAPI spec
- ✅ Premium UX features
- ✅ Performance optimization

**Outcome:** Ready for production launch

---

### Phase 2.5 (Enterprise Features) - 10-12 weeks

**Goal:** Enterprise-ready

- ✅ EventBus abstraction
- ✅ Theming & white-labelling
- ✅ Multi-tenancy
- ✅ Advanced analytics

**Outcome:** Ready for enterprise sales

---

## 🎤 TL;DR - Elevator Pitch

Beyond Redis 8 and the MVP stack, we're planning:

**Backend:**

- Full observability (structured logs, metrics, distributed tracing)
- Hardened auth (RBAC, rate limiting, JWT best practices)
- Long-term data strategy (partitioning, retention, archival)
- EventBus abstraction (cloud-agnostic, swappable)

**Frontend:**

- Nx domain libraries (faster builds, clear boundaries)
- Robust WebSocket layer (reconnect, offline handling)
- Premium UX (skeletons, virtual scroll, live badges, replay mode)
- Quality enforcement (Component Harnesses, Playwright flows, perf budgets)

**DevOps:**

- Preview environments per PR
- Automated security scanning
- Blue-green deployments

**Result:** Enterprise-grade platform ready for scale, not just an MVP.

---

**Status:** 🟡 Planning Phase
**Next Step:** Prioritize based on business goals
**Estimated Total Time:** 6-12 weeks (depending on priorities)
