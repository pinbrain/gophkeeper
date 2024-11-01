// Package grpc содержит реализацию grpc клиента.
package grpc

import (
	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/client/grpc/interceptors"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client определяет структуру клиента.
type Client struct {
	UserClient  pb.UserServiceClient
	VaultClient pb.VaultServiceClient
}

// NewGRPCConnection создает и возвращает новый grpc клиент.
func NewGRPCConnection(cfg *config.ClientConfig) (*Client, error) {
	creds, err := credentials.NewClientTLSFromFile("cert/server-cert.pem", "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(cfg.ServerAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(interceptors.TokenInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		UserClient:  pb.NewUserServiceClient(conn),
		VaultClient: pb.NewVaultServiceClient(conn),
	}, nil
}
