package grpc

import (
	"fmt"
	"net"

	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/server/grpc/handlers"
	"github.com/pinbrain/gophkeeper/internal/server/grpc/interceptors"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Transport struct {
	addr       string
	grpcServer *grpc.Server

	storage storage.Storage

	userHandler  *handlers.GRPCUserHandler
	vaultHandler *handlers.GRPCVaultHandler
	log          *logrus.Entry
}

type TransportConfig struct {
	MasterKey     string
	ServerAddress string
}

func NewGRPCTransport(
	cfg TransportConfig, storage storage.Storage, jwtService *jwt.Service, logger *logrus.Logger,
) (*Transport, error) {
	tlsCredentials, err := credentials.NewServerTLSFromFile("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	log := logger.WithField("instance", "grpcTransport")
	authInterceptor := interceptors.NewAuthInterceptor(cfg.MasterKey, storage, jwtService, log)
	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(
			interceptors.LoggerInterceptor(log),
			authInterceptor.AuthenticateUser,
			authInterceptor.RequireUser,
		),
	)
	userHandler := handlers.NewGRPCUserHandler(cfg.MasterKey, storage, jwtService, log)
	vaultHandler := handlers.NewGRPCVaultHandler(cfg.MasterKey, storage, log)
	grpcTransport := &Transport{
		addr:         cfg.ServerAddress,
		grpcServer:   s,
		storage:      storage,
		userHandler:  userHandler,
		vaultHandler: vaultHandler,
		log:          log,
	}
	pb.RegisterUserServiceServer(grpcTransport.grpcServer, grpcTransport.userHandler)
	pb.RegisterVaultServiceServer(grpcTransport.grpcServer, grpcTransport.vaultHandler)
	reflection.Register(grpcTransport.grpcServer)
	return grpcTransport, nil
}

func (s *Transport) Run() error {
	listen, err := net.Listen("tcp", s.addr)
	fmt.Println("going to listen:", s.addr)
	if err != nil {
		return fmt.Errorf("listen tcp has failed: %w", err)
	}
	return s.grpcServer.Serve(listen)
}

func (s *Transport) Stop() error {
	s.grpcServer.GracefulStop()
	return nil
}
