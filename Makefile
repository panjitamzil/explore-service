APP_NAME=explore-service

.PHONY: run test proto build docker-build compose-up compose-down

run:
	go run ./cmd/explore-service

test:
	go test ./...

proto:
	protoc --go_out=. --go-grpc_out=. proto/explore-service.proto

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -trimpath -o bin/$(APP_NAME) ./cmd/explore-service

docker-build:
	docker build -t $(APP_NAME):local .

compose-up:
	docker-compose up --build

compose-down:
	docker-compose down -v
