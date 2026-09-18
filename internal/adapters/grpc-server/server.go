package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	usersv1 "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/gen/users/v1"
	usersHandler "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/handlers/users"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/interceptors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

const Addr = "0.0.0.0:9090"

type GrpcServer struct {
	ctx    context.Context
	server *grpc.Server
	api    ports.ApiPort
}

func NewGrpcServer(ctx context.Context, api ports.ApiPort) ports.GrpcService {
	return &GrpcServer{
		ctx: ctx,
		api: api,
	}
}

func (s *GrpcServer) registerServices() *grpc.Server {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.UnaryLoggingInterceptor()))

	usersv1.RegisterUsersServiceServer(server, usersHandler.NewServer(s.api))
	reflection.Register(server)

	return server
}

func (s *GrpcServer) Run() error {
	s.server = s.registerServices()

	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s: %w", Addr, err)
	}

	serverCtx, serverStopCtx := context.WithCancel(s.ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		s.server.GracefulStop()
		serverStopCtx()
	}()

	log.Info().Str("addr", Addr).Msg("starting grpc server")
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("grpc server exited unexpectedly: %w", err)
	}

	<-serverCtx.Done()

	return nil
}

func (s *GrpcServer) Shutdown(_ context.Context) error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	return nil
}
