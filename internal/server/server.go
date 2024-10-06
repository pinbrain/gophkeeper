package server

import (
	"context"
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/pinbrain/gophkeeper/internal/server/grpc"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
)

type Server struct {
	storage   storage.Storage
	transport *grpc.Transport
}

func NewServer(cfg config.ServerConfig) (*Server, error) {
	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to run storage: %w", err)
	}
	transport := grpc.NewGRPCTransport(storage, cfg.ServerAddress)
	return &Server{
		storage:   storage,
		transport: transport,
	}, nil
}

func (s *Server) Run() error {
	return s.transport.Run()
}

func (s *Server) Stop() error {
	if err := s.transport.Stop(); err != nil {
		return err
	}
	s.storage.Close()
	return nil
}
