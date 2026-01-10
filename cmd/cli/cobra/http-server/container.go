package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
	"github.com/rs/zerolog/log"
)

type Container struct {
	application ports.ApiPort
	httpServer  ports.HttpService
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
		if err := c.httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("HTTP Server exited unexpectedly")
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

	log.Info().Msg("shutting down http server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := c.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http server shutdown failed")
	}

	log.Info().Msg("http server shutdown gracefully complete")
}
