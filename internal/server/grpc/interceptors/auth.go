package interceptors

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	pb "github.com/pinbrain/gophkeeper/internal/proto"
	appCtx "github.com/pinbrain/gophkeeper/internal/server/context"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	"github.com/pinbrain/gophkeeper/internal/server/utils"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor описывает структуру перехватчика для авторизации и аутентификации.
type AuthInterceptor struct {
	masterKey  string
	storage    storage.Storage
	jwtService jwt.ServiceI

	protectedServices map[string]bool

	log *logrus.Entry
}

// NewAuthInterceptor создает обработчик авторизации и аутентификации.
func NewAuthInterceptor(
	masterKey string, storage storage.Storage, jwtService jwt.ServiceI, log *logrus.Entry,
) *AuthInterceptor {
	return &AuthInterceptor{
		masterKey:  masterKey,
		storage:    storage,
		jwtService: jwtService,
		protectedServices: map[string]bool{
			pb.VaultService_ServiceDesc.ServiceName: true,
		},
		log: log,
	}
}

// AuthenticateUser аутентифицирует пользователя запроса.
func (i *AuthInterceptor) AuthenticateUser(
	ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(i.jwtService.GetMdJWTKey())
		if len(values) > 0 {
			jwt := values[0]
			userData, err := i.jwtService.GetJWTClaims(jwt)
			if err != nil {
				i.log.WithError(err).Error("failed to get claims from jwt")
				return nil, status.Error(codes.Unauthenticated, "Invalid jwt")
			}
			user, err := i.storage.GetUserByLogin(ctx, userData.Login)
			if err != nil {
				switch {
				case errors.Is(err, postgres.ErrNoUser):
					return nil, status.Error(codes.NotFound, "Пользователь не найден")
				default:
					i.log.WithError(err).Error("error while authenticating user by jwt")
					return nil, status.Error(codes.Internal, "Не удалось получить данные пользователя из БД")
				}
			}
			encUserSecretB, err := hex.DecodeString(user.EncryptedSecret)
			if err != nil {
				i.log.WithError(err).Error("error while decoding user secret key")
				return nil, status.Error(codes.Internal, "Internal Server Error")
			}
			userSecret, err := utils.Decrypt(encUserSecretB, i.masterKey)
			if err != nil {
				i.log.WithError(err).Error("error while decrypting user secret key")
				return nil, status.Error(codes.Internal, "Internal Server Error")
			}
			ctx = appCtx.CtxWithUser(ctx, &appCtx.CtxUser{
				ID:     user.ID,
				Login:  user.Login,
				Secret: hex.EncodeToString(userSecret),
			})
		}
	}

	return handler(ctx, req)
}

// RequireUser проверяет что пользователь авторизован.
// В противном случае прерывает обработку запроса и возвращает ошибку Unauthorized.
func (i *AuthInterceptor) RequireUser(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if !i.protectedServices[strings.Split(info.FullMethod, "/")[1]] {
		return handler(ctx, req)
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil || user.ID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}
	return handler(ctx, req)
}
