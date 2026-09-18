package grpcserver

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type Container struct {
	application ports.ApiPort
	grpcServer  ports.GrpcService
}

func (c *Container) start(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("container recovered from panic")
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := c.grpcServer.Run(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error().Err(err).Msg("gRPC Server exited unexpectedly")
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigCh:
		log.Info().Str("signal", sig.String()).Msg("received shutdown signal")
		cancel()
	case <-ctx.Done():
		log.Info().Msg("context cancelled")
	}

	log.Info().Msg("shutting down grpc server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := c.grpcServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("grpc server shutdown failed")
	}

	log.Info().Msg("grpc server shutdown gracefully complete")
}
