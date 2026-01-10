package chiserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

const Addr = "0.0.0.0:8080"

type HttpServer struct {
	ctx    context.Context
	server *http.Server
	api    ports.ApiPort
}

func NewHttpServer(ctx context.Context, api ports.ApiPort) ports.HttpService {
	return &HttpServer{
		ctx:    ctx,
		server: &http.Server{},
		api:    api,
	}
}

func (s *HttpServer) Run() error {
	s.server = &http.Server{
		Addr:              Addr,
		Handler:           s.setRouter(),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	serverCtx, serverStopCtx := context.WithCancel(s.ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxCancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal().Err(shutdownCtx.Err()).Msg("graceful shutdown timed out, forcing shutdown")
			}
		}()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown server")
		}
		serverStopCtx()
	}()

	if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("could not listen on %s: %w", s.server.Addr, err)
	}

	<-serverCtx.Done()

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error { return s.server.Shutdown(ctx) }
