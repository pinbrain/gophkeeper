package service

import (
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
)

type Service struct {
	grpcClient *grpc.Client
}

func NewService(grpcClient *grpc.Client) *Service {
	return &Service{
		grpcClient: grpcClient,
	}
}
