// Package server содержит реализацию сервера.
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

// Server определяет структуру сервера.
type Server struct {
	storage   storage.Storage
	transport *grpc.Transport

	log *logrus.Entry
}

// NewServer создает и возвращает новый сервер.
func NewServer(ctx context.Context, cfg *config.ServerConfig, logger *logrus.Logger) (*Server, error) {
	log := logger.WithField("instance", "server")

	storage, err := postgres.NewStorage(ctx, cfg.DSN, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to run storage: %w", err)
	}

	jwtService := jwt.NewJWTService(cfg.JWT)

	transport, err := grpc.NewGRPCTransport(grpc.TransportConfig{
		MasterKey:     cfg.MasterKey,
		ServerAddress: cfg.ServerAddress,
	}, storage, jwtService, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc transport: %w", err)
	}

	return &Server{
		storage:   storage,
		transport: transport,
		log:       log,
	}, nil
}

// Run запускает сервер.
func (s *Server) Run() error {
	return s.transport.Run()
}

// Stop останавливает сервер.
func (s *Server) Stop() error {
	if err := s.transport.Stop(); err != nil {
		s.log.Errorf("an error occurred during grpc server shutdown: %v", err)
	}
	s.log.Info("gRPC server stopped")
	s.storage.Close()
	s.log.Info("Storage closed")
	return nil
}
