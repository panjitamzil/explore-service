FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o explore-service ./cmd/explore-service

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/explore-service .

ENV GRPC_PORT=50051

EXPOSE 50051

ENTRYPOINT ["./explore-service"]
