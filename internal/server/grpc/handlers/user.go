package handlers

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/pinbrain/gophkeeper/internal/model"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	"github.com/pinbrain/gophkeeper/internal/server/utils"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCUserHandler определяет структуру обработчика grpc запросов в части работы с пользователями.
type GRPCUserHandler struct {
	pb.UnimplementedUserServiceServer
	masterKey  string
	storage    storage.Storage
	jwtService jwt.ServiceI
	log        *logrus.Entry
}

// NewGRPCUserHandler создает и возвращает новый обработчик grpc запросов в части работы с пользователями.
func NewGRPCUserHandler(
	masterKey string, storage storage.Storage, jwtService jwt.ServiceI, log *logrus.Entry,
) *GRPCUserHandler {
	return &GRPCUserHandler{
		masterKey:  masterKey,
		storage:    storage,
		jwtService: jwtService,
		log:        log,
	}
}

// Register регистрирует нового пользователя.
func (h *GRPCUserHandler) Register(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterRes, error) {
	if in.GetLogin() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Некорректные входные данные")
	}
	passwordHash, err := utils.GeneratePasswordHash(in.GetPassword())
	if err != nil {
		h.log.WithError(err).Error("Error while creating new user - failed to generate password hash")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	secretKey, err := utils.GenerateUserKey()
	if err != nil {
		h.log.WithError(err).Error("Error while creating new user - failed to generate user secret key")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	encSecretKey, err := utils.Encrypt(secretKey, h.masterKey)
	if err != nil {
		h.log.WithError(err).Error("Error while creating new user - failed to encrypt user secret key")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user := &model.User{
		Login:           in.GetLogin(),
		PasswordHash:    passwordHash,
		EncryptedSecret: hex.EncodeToString(encSecretKey),
	}
	id, err := h.storage.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrLoginTaken):
			return nil, status.Error(codes.AlreadyExists, "Пользователь с таким логином уже существует")
		default:
			h.log.WithError(err).Error("Error while creating new user - failed to save user in DB")
			return nil, status.Error(codes.Internal, "Не удалось создать пользователя")
		}
	}
	user.ID = id

	jwt, err := h.jwtService.BuildJWTSting(user)
	if err != nil {
		h.log.WithError(err).Error("Error while creating new user - failed to generate jwt")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	response := &pb.RegisterRes{
		Token: jwt,
	}
	return response, nil
}

// Login аутентифицирует пользователя по логину и паролю.
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
			h.log.WithError(err).Error("Error while login user")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}
	if isPwdOk := utils.ComparePwdAndHash(in.GetPassword(), user.PasswordHash); !isPwdOk {
		return nil, status.Error(codes.Unauthenticated, "Неверные логин/пароль")
	}
	jwt, err := h.jwtService.BuildJWTSting(user)
	if err != nil {
		h.log.WithError(err).Error("Error while login user")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	response := &pb.LoginRes{
		Token: jwt,
	}
	return response, nil
}
