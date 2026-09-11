# CineBook

Movie ticketing platform: Go backend (Postgres + Redis) with a Gemini-powered
AI booking assistant, real-time seat availability over WebSockets, and a
React admin dashboard.

## Structure

```
cinebook/
├── backend/
│   ├── cmd/api/            # entry point (main.go)
│   └── internal/
│       ├── models/         # GORM structs
│       ├── database/       # Postgres + Redis connections
│       ├── services/       # circuit breaker, booking logic
│       ├── ai/             # Gemini orchestrator & tool declarations
│       └── ws/             # WebSocket hub for live seat updates
├── frontend/                # React + Vite admin dashboard
└── docker-compose.yml        # local Postgres + Redis
```

## Local setup

1. Start Postgres and Redis:
   ```bash
   docker-compose up -d
   ```
2. Fetch Go dependencies and run the API:
   ```bash
   cd backend
   go mod tidy
   go run ./cmd/api
   ```
   The server listens on `:8080` by default. Override with `PORT`,
   `DATABASE_URL`, and `REDIS_ADDR` env vars as needed.
3. Frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Status

Scaffolding + core primitives (models, Redis seat holds, circuit breaker,
WebSocket hub, AI tool-calling loop) are in place. Tool implementations in
`internal/ai/executor.go` and the booking sub-agent are stubbed pending
repository/service wiring.
