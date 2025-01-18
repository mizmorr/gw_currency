package server

import (
	"context"
	"net"
	"time"

	logger "github.com/mizmorr/loggerm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type CurrencyController interface {
	Register(ctx context.Context, server *grpc.Server)
}

type Server struct {
	server          *grpc.Server
	shutDownTimeout time.Duration
	notify          chan error
	listener        net.Listener
	controller      CurrencyController
}

func New(ctx context.Context, controller CurrencyController, host, port string) (*Server, error) {
	logger := logger.GetLoggerFromContext(ctx)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	address := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to listen")
	}

	return &Server{
		server:          server,
		shutDownTimeout: 5 * time.Second,
		notify:          make(chan error),
		listener:        listener,
		controller:      controller,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.controller.Register(ctx, s.server)
	go func() {
		s.notify <- s.server.Serve(s.listener)
		close(s.notify)
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutDownTimeout)
	defer cancel()
	okCh := make(chan interface{})

	go func() {
		s.server.GracefulStop()
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-s.notify:
		return err
	case <-okCh:
		return nil
	}
}

func loggingInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		// Определяем IP клиента
		var clientIP string
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		logger.Info().
			Str("method", info.FullMethod).
			Str("client_ip", clientIP).
			Interface("request", req). // Логируем тело запроса
			Msg("Incoming gRPC request")

		resp, err = handler(ctx, req)

		logger.Info().
			Str("method", info.FullMethod).
			Dur("duration", time.Since(start)).
			Err(err). // Логируем ошибку, если есть
			Msg("Completed gRPC request")

		return resp, err
	}
}
