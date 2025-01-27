package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	logger "github.com/mizmorr/loggerm"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutDownTimeout time.Duration
	stop            chan interface{}
	logger          *logger.Logger
}

func New(handler http.Handler, host, port string, shutDownTimeout time.Duration) *Server {
	address := net.JoinHostPort(host, port)

	httpServer := &http.Server{
		Handler: handler,
		Addr:    address,
	}

	return &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutDownTimeout: shutDownTimeout,
		stop:            make(chan interface{}),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger = logger.GetLoggerFromContext(ctx)

	go func() {
		s.logger.Info().Msg("HTTP server is starting..")
		s.notify <- s.server.ListenAndServe()
	}()

	go func() {
		s.keepAlive(ctx)
	}()

	return nil
}

func (s *Server) keepAlive(ctx context.Context) {
	s.logger.Debug().Msg("HTTP server keep-alive is running..")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Msg("Keep alive http-worker is stopping (context canceled)...")
			return
		case <-s.stop:
			s.logger.Info().Msg("Keep alive http-worker is stopped correct")
			return
		case err := <-s.notify:
			s.logger.Warn().Msg("Got that err: " + err.Error() + "while running server, restarting..")
			s.restart()
		}
	}
}

func (s *Server) restart() {
	go func() {
		s.notify <- s.server.ListenAndServe()
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping http-server...")

	stopCtx, cancel := context.WithTimeout(ctx, s.shutDownTimeout)
	defer cancel()

	close(s.stop)

	if err := s.server.Shutdown(stopCtx); err != nil {
		s.logger.Error().Msg("Shutting down http-server failed")
		return err
	}

	return nil
}
