package grpc

import (
	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/client/grpc/interceptors"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	UserClient  pb.UserServiceClient
	VaultClient pb.VaultServiceClient
}

func NewGRPCConnection(cfg *config.ClientConfig) (*Client, error) {
	conn, err := grpc.NewClient(cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
