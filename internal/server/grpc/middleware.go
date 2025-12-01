package grpc

import (
	"context"
	"log"
	"time"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(logger *log.Logger) gogrpc.UnaryServerInterceptor {
	if logger == nil {
		logger = log.Default()
	}

	return func(
		ctx context.Context,
		req any,
		info *gogrpc.UnaryServerInfo,
		handler gogrpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)
		logger.Printf("method=%s code=%s latency=%s", info.FullMethod, codeToString(code), time.Since(start))
		return resp, err
	}
}

func codeToString(code codes.Code) string {
	if code == codes.OK {
		return "OK"
	}
	return code.String()
}
