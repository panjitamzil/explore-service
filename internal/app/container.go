package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	rediscache "explore-service/internal/cache/redis"
	"explore-service/internal/config"
	"explore-service/internal/domain/decision"
	mysqlrepo "explore-service/internal/repository/mysql"
	grpcserver "explore-service/internal/server/grpc"
	explorepb "explore-service/pkg/proto/explore/proto"

	"github.com/redis/go-redis/v9"
	gogrpc "google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
)

type Container struct {
	Config       config.Config
	DB           *sql.DB
	RedisClient  *redis.Client
	GRPCServer   *gogrpc.Server
	DecisionSvc  *decision.Service
	DecisionRepo decision.Repository
}

func NewContainer(ctx context.Context, cfg config.Config) (*Container, error) {
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr(),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	repo := mysqlrepo.NewDecisionRepository(db)
	cache := rediscache.NewLikedCountCache(rdb, "liked_count:")
	svc := decision.NewService(repo, cache)

	grpcOpts := []gogrpc.ServerOption{
		gogrpc.ChainUnaryInterceptor(
			grpcserver.UnaryLoggingInterceptor(log.Default()),
		),
	}
	grpcServer := gogrpc.NewServer(grpcOpts...)
	handler := grpcserver.NewExploreHandler(svc)
	explorepb.RegisterExploreServiceServer(grpcServer, handler)

	return &Container{
		Config:       cfg,
		DB:           db,
		RedisClient:  rdb,
		GRPCServer:   grpcServer,
		DecisionSvc:  svc,
		DecisionRepo: repo,
	}, nil
}

func (c *Container) Close() error {
	var err error
	if c.RedisClient != nil {
		if e := c.RedisClient.Close(); e != nil {
			err = e
		}
	}
	if c.DB != nil {
		if e := c.DB.Close(); e != nil && err == nil {
			err = e
		}
	}
	return err
}
