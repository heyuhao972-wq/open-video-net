# Recommendation Platform Dev Guide

This service provides recommendation APIs, behavior tracking, follow graph, favorites, notifications, and stats.
It is designed so that algorithm changes are isolated to a small set of files.

## Project Layout

```
recommendation-platform/
  api/                 HTTP router wiring
  config/              env config (ports, index base, jwt secret)
  handler/             HTTP handlers
  model/               data models
  pipeline/            retrieval/scoring pipeline components
  repository/          persistence (SQLite) + in-memory helpers
  service/             business logic (recommendation algorithms live here)
  main.go              service bootstrap
```

## Run Locally

From `recommendation-platform/recommendation-platform`:

```
$env:CONTENT_PLATFORMS="platformA=http://localhost:8080"
$env:DB_PATH="./data/recommendation.db"
go run .
```

Health endpoint:

```
GET http://localhost:8082/health
```

## Environment Variables

These are read in `config/config.go`:

- `RECOMMEND_PORT` (default `8082`)
- `INDEX_BASE` (default `http://localhost:8083`)  
  Used by pipeline to load index videos. If you do not run an index service, you can set this to an empty string.
- `JWT_SECRET` (default `dev-secret`)  
  Must match content platform for auth token validation.

Behavior and follow endpoints use the JWT in the `Authorization: Bearer <token>` header.

## Core Endpoints

Recommendation:

- `GET /recommend?user=<id>&type=default|hot|latest|following&page=1&limit=10`

Behavior & social:

- `POST /behavior`
- `POST /follow`
- `POST /unfollow`
- `POST /favorite`
- `POST /unfavorite`

User data:

- `GET /me/likes`
- `GET /me/favorites`
- `GET /me/follows`
- `GET /me/followers`
- `GET /me/history`

Notifications:

- `GET /notifications`
- `POST /notifications/read`

Stats:

- `GET /video/:id/stats`

## Where To Change Recommendation Algorithms

All algorithm changes are isolated to:

- `service/algorithms.go`  
  This is the only file you should edit for algorithm logic.

How it works:

1. The `Algorithm` interface defines:
   ```
   Name() string
   Recommend(userID string, limit int) []model.Video
   ```
2. `AlgorithmRegistry` maps `type` query to an algorithm.
3. `RecommendService` delegates to the registry.

### Edit the default algorithm

In `service/algorithms.go`:

- Modify `DefaultAlgorithm.Recommend(...)`
- Tune constants in `defaultAlgoConfig()`

### Add a new algorithm

1. Add a new struct that implements `Algorithm`.
2. Register it in `NewAlgorithmRegistry()`.
3. Call it by name using `type=<name>` in `/recommend`.

## Pipeline Components (Advanced)

The default algorithm uses a pipeline:

- Sources: `pipeline/sources.go`
- Filters: `pipeline/filters.go`
- Scorers: `pipeline/scorers.go`
- Selector: `pipeline/selector.go`

If you want to add new retrieval logic or scoring, add a new component in the `pipeline` package and wire it into `DefaultAlgorithm`.

## Storage

The recommendation platform stores:

- behaviors (likes/shares/watches/comments)
- follows
- favorites
- notifications
- watch history

The repository layer writes to SQLite using `DB_PATH`.

## Common Debug Checks

- `/recommend` returns empty:
  - index service not reachable or index has no videos
  - set `INDEX_BASE` correctly or ensure content platform `/videos` is reachable

- `/me/*` endpoints return 401:
  - JWT invalid or mismatch in `JWT_SECRET`

- Missing stats/notifications:
  - check `/behavior` is called from client
  - confirm DB file is writable under `DB_PATH`

## Minimal Test Flow

1. Start content platform (8080).
2. Start recommendation platform (8082).
3. Upload a video from the client.
4. Open `/recommend` and verify it returns `video://platformA/<id>`.

