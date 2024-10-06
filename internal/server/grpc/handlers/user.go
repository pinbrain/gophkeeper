package handlers

import (
	"context"
	"errors"

	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/model"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/pinbrain/gophkeeper/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCUserHandler struct {
	pb.UnimplementedUserServiceServer
	storage storage.Storage
}

func NewGRPCUserHandler(storage storage.Storage) *GRPCUserHandler {
	return &GRPCUserHandler{
		storage: storage,
	}
}

func (h *GRPCUserHandler) Register(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterRes, error) {
	if in.GetEmail() == "" || in.GetLogin() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Некорректные входные данные")
	}
	passwordHash, err := utils.GeneratePasswordHash(in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	user := &model.User{
		Login:        in.GetLogin(),
		Email:        in.GetEmail(),
		PasswordHash: passwordHash,
	}
	_, err = h.storage.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrLoginTaken):
			return nil, status.Error(codes.AlreadyExists, "Пользователь с таким логином или email уже существует")
		default:
			logger.Log.WithError(err).Error("Error while creating new user")
			return nil, status.Error(codes.Internal, "Не удалось создать пользователя")
		}
	}

	jwt, err := utils.BuildJWTSting(user)
	if err != nil {
		logger.Log.WithError(err).Error("Error while creating new user")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	response := &pb.RegisterRes{
		Token: jwt,
	}
	return response, nil
}

func (h *GRPCUserHandler) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginRes, error) {
	if in.GetLogin() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Некорректные входные данные")
	}
	user, err := h.storage.GetUserByLogin(ctx, in.GetLogin())
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoUser):
			return nil, status.Error(codes.NotFound, "Пользователь с таким логином не найден")
		default:
			logger.Log.WithError(err).Error("Error while login user")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}
	if isPwdOk := utils.ComparePwdAndHash(in.GetPassword(), user.PasswordHash); !isPwdOk {
		return nil, status.Error(codes.Unauthenticated, "Неверные логин/пароль")
	}
	jwt, err := utils.BuildJWTSting(user)
	if err != nil {
		logger.Log.WithError(err).Error("Error while login user")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	response := &pb.LoginRes{
		Token: jwt,
	}
	return response, nil
}
