# Explore Service

Simple gRPC service to store like/pass decisions and list who liked a user. Data goes to MySQL, liked count is cached in Redis.

## Structure
- `cmd/explore-service/` entrypoint
- `internal/app/` wiring and gRPC server
- `internal/domain/decision/` entities and use cases
- `internal/repository/mysql/` MySQL repo
- `internal/repository/memory/` in-memory repo (dev/test)
- `internal/cache/redis/` Redis cache
- `internal/server/grpc/` gRPC handlers, pagination token, logging
- `proto/` protobuf; generated code in `pkg/proto/explore/proto`
- `migrations/` MySQL schema

## Run
1. Copy env (optional):
```bash
cp .env.example .env
```
2. Docker Compose (MySQL + Redis + app):
```bash
docker-compose up --build
```
gRPC on `localhost:50051`.

3. Local (no containers):
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

## Generate proto
Needs `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`:
```bash
protoc --go_out=. --go-grpc_out=. proto/explore-service.proto
```

## Notes
- Upsert by primary key `(actor_user_id, recipient_user_id)`.
- Cursor pagination order: `decision_unix_ts DESC, actor_id ASC`; token is base64 JSON.
- Redis cache for liked count with TTL; cache cleared on `PutDecision`.
- Page size limits, logging interceptor, graceful shutdown with signals.
