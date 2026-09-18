package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err)

		event := log.Debug()
		if err != nil {
			event = log.Error().Err(err)
		}

		event.
			Str("method", info.FullMethod).
			Str("code", code.String()).
			Dur("duration", duration).
			Msg("grpc request completed")

		return resp, err
	}
}
