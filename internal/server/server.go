package server

import (
	"context"
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/pinbrain/gophkeeper/internal/server/grpc"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/sirupsen/logrus"
)

type Server struct {
	storage   storage.Storage
	transport *grpc.Transport

	log *logrus.Entry
}

func NewServer(cfg *config.ServerConfig, logger *logrus.Logger) (*Server, error) {
	log := logger.WithField("instance", "server")

	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, cfg.DSN, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to run storage: %w", err)
	}

	jwtService := jwt.NewJWTService(cfg.JWT)

	transport := grpc.NewGRPCTransport(grpc.TransportConfig{
		MasterKey:     cfg.MasterKey,
		ServerAddress: cfg.ServerAddress,
	}, storage, jwtService, logger)
	return &Server{
		storage:   storage,
		transport: transport,
		log:       log,
	}, nil
}

func (s *Server) Run() error {
	return s.transport.Run()
}

func (s *Server) Stop() error {
	if err := s.transport.Stop(); err != nil {
		s.log.Errorf("an error occurred during grpc server shutdown: %v", err)
	}
	s.storage.Close()
	return nil
}
