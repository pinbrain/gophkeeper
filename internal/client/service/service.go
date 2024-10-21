package service

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/client/grpc"
	"github.com/pinbrain/gophkeeper/internal/model"
)

type Service struct {
	grpcClient *grpc.Client

	UserService
	VaultService
}

type UserService interface {
	Register(ctx context.Context, login, password string) (token string, err error)
	Login(ctx context.Context, login, password string) (token string, err error)
}

type VaultService interface {
	AddPassword(ctx context.Context, data string, meta model.PasswordMeta) error
	AddText(ctx context.Context, data string, meta model.TextMeta) error
	AddBankCard(ctx context.Context, data model.BankCardData, meta model.BankCardMeta) error
	AddFile(ctx context.Context, file string, comment string) error
	GetData(ctx context.Context, id string) (model.DataType, any, error)
}

func NewService(grpcClient *grpc.Client) *Service {
	return &Service{
		grpcClient: grpcClient,
	}
}
