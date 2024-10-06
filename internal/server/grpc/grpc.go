package grpc

import (
	"fmt"
	"net"

	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/server/grpc/handlers"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Transport struct {
	addr         string
	grpcServer   *grpc.Server
	storage      storage.Storage
	userHandler  *handlers.GRPCUserHandler
	vaultHandler *handlers.GRPCVaultHandler
}

func NewGRPCTransport(storage storage.Storage, addr string) *Transport {
	s := grpc.NewServer()
	userHandler := handlers.NewGRPCUserHandler(storage)
	vaultHandler := handlers.NewGRPCVaultHandler(storage)
	grpcTransport := &Transport{
		addr:         addr,
		grpcServer:   s,
		storage:      storage,
		userHandler:  userHandler,
		vaultHandler: vaultHandler,
	}
	pb.RegisterUserServiceServer(grpcTransport.grpcServer, grpcTransport.userHandler)
	pb.RegisterVaultServiceServer(grpcTransport.grpcServer, grpcTransport.vaultHandler)
	reflection.Register(grpcTransport.grpcServer)
	return grpcTransport
}

func (s *Transport) Run() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen tcp has failed: %w", err)
	}
	return s.grpcServer.Serve(listen)
}

func (s *Transport) Stop() error {
	s.grpcServer.GracefulStop()
	return nil
}
