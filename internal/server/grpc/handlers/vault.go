package handlers

import (
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/storage"
)

type GRPCVaultHandler struct {
	pb.UnimplementedVaultServiceServer
	storage storage.Storage
}

func NewGRPCVaultHandler(storage storage.Storage) *GRPCVaultHandler {
	return &GRPCVaultHandler{storage: storage}
}
