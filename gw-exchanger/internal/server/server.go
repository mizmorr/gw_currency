package server

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	server          *grpc.Server
	shutDownTimeout time.Duration
	notify          chan error
	listener        net.Listener
}

func New(host, port string) (*Server, error) {
	server := grpc.NewServer()

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
	}, nil
}

func (s *Server) Start() error {
	pb.RegisterCurrencyExchangeServiceServer(server, &CurrencyExchangeServer{})

	go func() {
		s.notify <- s.server.Serve(s.listener)
	}()

	return nil
}
