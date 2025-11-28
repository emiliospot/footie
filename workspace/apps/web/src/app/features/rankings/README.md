# Competition Rankings Feature

This feature implements the "Competition Rankings" card component as specified in the technical assignment. It displays team and player rankings across different performance categories.

## Overview

The Competition Rankings component shows:

- **Team Rankings**: Top 5 teams for various metrics (xG, Shots, Crosses, etc.)
- **Player Rankings**: Top 5 players for various metrics with avatars
- **Categories**: Attacking, Defending, Distribution, Goalkeeper, Insights
- **Filters**: Championship and Season selection

## Architecture

### Frontend (Angular)

- **Component**: `competition-rankings.component.ts`
- **Service**: `rankings.service.ts` (calls backend API)
- **Models**: `rankings.model.ts` (TypeScript interfaces)
- **Route**: `/dashboard` (embedded in dashboard) or `/rankings` (standalone)

### Backend (Go)

- **Handler**: `rankings.go` (serves mock data)
- **Endpoint**: `GET /api/v1/rankings`
- **Query Parameters**:
  - `type`: `team` or `player` (default: `team`)
  - `category`: `attacking`, `defending`, `distribution`, `goalkeeper`, `insights` (default: `attacking`)
  - `championship`: Championship name (optional)
  - `season`: Season (optional)

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL (optional - rankings uses mock data)
- Redis (optional - for real-time features)

### Starting the Backend API

**Important:** All npm commands must be run from the `workspace` directory, not the project root.

**Option 1: Using npm script (recommended - from workspace directory):**

```bash
# Navigate to workspace directory first
cd workspace

# Then run the API
npm run api
```

This uses Nx (`nx run api:serve`) to run the API with Air hot-reload.

**Option 2: Using VSCode Tasks (if using VSCode):**

1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
2. Type "Tasks: Run Task"
3. Select "Dev: Backend Only"

This runs `npm run api` which uses Nx internally.

**Option 3: Using Make (optional - from API directory):**

Make is usually pre-installed on Linux/Mac. If you don't have it, skip to Option 4.

**Note:** `make dev` requires Air to be installed. Install it first with:

```bash
cd workspace/apps/api
make install  # Installs Air and other Go tools
```

Then you can use:

```bash
cd workspace/apps/api

# Run directly (no Air needed)
make run

# Run with hot-reload (requires Air)
make dev
```

**Option 4: Using Go directly (from API directory):**

```bash
cd workspace/apps/api
go run cmd/api/main.go
```

**Note:** Option 1 (npm script) is recommended as it works consistently across all platforms and uses Nx for task orchestration.

**Note:** The rankings endpoint uses mock data, so you can run the API without a database by setting environment variables:

```bash
SKIP_DB=true SKIP_REDIS=true npm run api
```

**The API will start on:** `http://localhost:8088`

**Test the health endpoint:**

```bash
curl http://localhost:8088/health
```

**Test the rankings endpoint:**

```bash
curl "http://localhost:8088/api/v1/rankings?type=team&category=attacking"
```

### Starting the Frontend (Angular)

**Important:** All npm commands must be run from the `workspace` directory, not the project root.

**Option 1: Using npm script (from workspace directory):**

```bash
# Navigate to workspace directory first
cd workspace

# Install dependencies (if not already done)
npm install

# Start the development server
npm run web
```

This uses Nx (`nx run web:serve`) to run the Angular dev server.

**Option 2: Using VSCode Tasks (if using VSCode):**

1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
2. Type "Tasks: Run Task"
3. Select "Dev: Frontend Only"

**The app will start on:** `http://localhost:4200`

**Navigate to:** `http://localhost:4200/dashboard` to see the Competition Rankings component

**Note:** The dashboard route is currently public (auth guard disabled for development).

### Using Docker Compose (Optional)

If you prefer to run PostgreSQL and Redis in Docker:

```bash
# Start database and Redis
cd workspace/infra/docker
docker compose up -d postgres redis

# Or from project root
docker compose -f workspace/infra/docker/docker-compose.yml up -d postgres redis
```

**Note:** The rankings endpoint uses mock data, so the database is optional for testing this feature.

## API Endpoints

### Health Check

```
GET /health
```

**Response:**

```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

### Get Competition Rankings

```
GET /api/v1/rankings?type=team&category=attacking&championship=Cyprus%20U19%20League%20Division%201&season=2025/2026
```

**Query Parameters:**

- `type` (optional): `team` or `player` (default: `team`)
- `category` (optional): `attacking`, `defending`, `distribution`, `goalkeeper`, `insights` (default: `attacking`)
- `championship` (optional): Championship name
- `season` (optional): Season

**Response:**

```json
{
  "type": "team",
  "category": "attacking",
  "categories": [
    {
      "title": "xG - Expected Goals",
      "unit": "/90'",
      "rankings": [
        {
          "rank": 1,
          "name": "Anorthosis U19",
          "value": 2.42
        },
        ...
      ]
    },
    ...
  ]
}
```

## Component Structure

```
rankings/
├── competition-rankings.component.ts    # Main component logic
├── competition-rankings.component.html # Template
├── competition-rankings.component.scss # Styles
├── rankings.routes.ts                  # Routing configuration
└── README.md                           # This file
```

## Features

- ✅ Team and Player rankings
- ✅ Multiple performance categories
- ✅ Championship and Season filters
- ✅ Responsive design
- ✅ Loading and error states
- ✅ Player avatars with initials
- ✅ TrackBy functions for performance
- ✅ Type-safe TypeScript interfaces
- ✅ Mock data from backend

## Development Notes

- The component is currently embedded in the dashboard (`/dashboard`)
- It's also available as a standalone route (`/rankings`)
- Mock data is served from the Go backend
- Database is optional for this feature (uses mock data)
- Auth guard is disabled on dashboard for development

## Future Enhancements

- Connect to real database instead of mock data
- Add more performance categories
- Implement filtering by multiple championships
- Add export functionality (CSV, PDF)
- Add charts and visualizations
- Implement real-time updates via WebSocket
