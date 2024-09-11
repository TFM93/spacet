package grpc

import (
	"context"
	"fmt"
	"spacet/pkg/logger"

	"google.golang.org/grpc"
)

// Setup creates a grpcServer, configures the necessary interceptors and registers the following services:
// - RegisterSpaceTServiceServer
func Setup(l logger.Interface) (*grpc.Server, error) {
	if l == nil {
		return nil, fmt.Errorf("invalid input parameters: logger must not be nil")
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(loggerInterceptor(l)))
	// todo: gen.RegisterSpaceTServiceServer(server, nil)
	return server, nil
}

func loggerInterceptor(l logger.Interface) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			l.Info("gRPC method: %s, error: %v", info.FullMethod, err)
		} else {
			l.Debug("gRPC method: %s, ok", info.FullMethod)
		}
		return resp, err
	}
}
