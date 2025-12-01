package app

import (
	"context"
	"log"

	"explore-service/internal/config"
	"explore-service/internal/domain/decision"
	memrepo "explore-service/internal/repository/memory"
	grpcserver "explore-service/internal/server/grpc"
	explorepb "explore-service/pkg/proto/explore/proto"

	gogrpc "google.golang.org/grpc"
)

type Container struct {
	Config       config.Config
	GRPCServer   *gogrpc.Server
	DecisionSvc  *decision.Service
	DecisionRepo decision.Repository
}

func NewContainer(ctx context.Context, cfg config.Config) (*Container, error) {
	_ = ctx

	repo := memrepo.NewDecisionRepository()
	svc := decision.NewService(repo, nil)

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
		GRPCServer:   grpcServer,
		DecisionSvc:  svc,
		DecisionRepo: repo,
	}, nil
}

func (c *Container) Close() error {
	return nil
}
