package interceptor

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	return resp, err
}

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("паника в %s: %v\n%s", info.FullMethod, r, debug.Stack())
		}
	}()

	return handler(ctx, req)
}
