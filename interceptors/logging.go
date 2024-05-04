package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func (i *InterceptorImpl) Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("Request : %s", info.FullMethod)

	return handler(ctx, req)
}
