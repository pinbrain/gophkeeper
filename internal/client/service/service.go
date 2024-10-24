package service

import (
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
)

// Service описывает структуру сервиса.
type Service struct {
	grpcClient *grpc.Client
}

// NewService создает и возвращает новый сервис.
func NewService(grpcClient *grpc.Client) *Service {
	return &Service{
		grpcClient: grpcClient,
	}
}
