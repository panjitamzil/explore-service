package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"explore-service/internal/app"
	"explore-service/internal/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := config.LoadDotEnv(); err != nil {
		log.Printf("warning: failed to load .env: %v", err)
	}

	cfg := config.FromEnv()

	container, err := app.NewContainer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init app container: %v", err)
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Printf("error closing resources: %v", err)
		}
	}()

	addr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	log.Printf("ExploreService gRPC server listening on %s", addr)

	go func() {
		if err := container.GRPCServer.Serve(lis); err != nil {
			log.Printf("gRPC server stopped with error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received, stopping gRPC server...")

	stopCh := make(chan struct{})
	go func() {
		container.GRPCServer.GracefulStop()
		close(stopCh)
	}()

	select {
	case <-stopCh:
		log.Println("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		log.Println("gRPC server shutdown timed out, forcing stop")
		container.GRPCServer.Stop()
	}
}
