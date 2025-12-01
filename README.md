# Explore Service

The Explore Service implements the core logic behind “Liked You”, “New Liked You”, and “Mutual Likes” features commonly used in matchmaking applications.  
It is a clean, scalable, gRPC-based backend with cursor pagination, Redis caching, and MySQL storage.


## Features
- **PutDecision** - record like/pass and detect mutual likes  
- **ListLikedYou** — list users who liked the recipient (paginated)  
- **ListNewLikedYou** — list users you haven’t liked back  
- **CountLikedYou** — efficient counting with Redis caching  
- **Cursor-based pagination** (timestamp + actor_id)  
- **MySQL schema** optimized with composite indexes  
- **Comprehensive tests** for service, handler, repository, pagination, and cache

## Requirements
- Go 1.22+
- Docker & Docker Compose (for MySQL/Redis stack)
- `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`

## Run
1) Copy env:
```bash
cp .env.example .env
```
2) Docker Compose (MySQL + Redis + app):
```bash
docker-compose up --build
```
gRPC listens on `localhost:50051`.

3) Local without containers:
```bash
export GRPC_PORT=50051
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=password
export MYSQL_DB=explore
export REDIS_HOST=localhost
export REDIS_PORT=6379
go run ./cmd/explore-service
```

## API (quick samples)
- PutDecision (mutual only if both directions are liked):
```json
{ "actor_user_id":"u1", "recipient_user_id":"u2", "liked_recipient":true }
```
Response: `{"mutual_likes": false}`

- ListLikedYou:
```json
{ "recipient_user_id":"u2", "page_size":20 }
```
Response returns likers and `next_pagination_token` if more data.

- ListNewLikedYou: same request shape, but only users not yet liked back.

- CountLikedYou: `{"recipient_user_id":"u2"}` (cached in Redis).

Pagination token format: JSON `{"ts": <unix_ts>, "actor_id": "<string>"}` encoded with Base64 URL-safe.

## Database
Schema is in `migrations/001_create_decisions_table.sql` with indexes for pagination and mutual-like checks.

## Tests
```bash
go test ./...
```
Coverage includes domain service, pagination token, handlers, MySQL repo (mocked), and Redis cache (mocked).

## Project Structure
- `cmd/explore-service/` — main entrypoint
- `internal/app/` — wiring (MySQL, Redis, gRPC server)
- `internal/config/` — env loading and DSN builder
- `internal/domain/decision/` — entities, errors, service, repository interface
- `internal/repository/mysql/` — MySQL repository
- `internal/repository/memory/` — in-memory repository
- `internal/cache/redis/` — Redis liked-count cache
- `internal/server/grpc/` — handlers, pagination helpers, middleware
- `proto/` — proto definitions
- `migrations/` — SQL schema
- `Dockerfile`, `docker-compose.yml`, `Makefile`, `README.md`
