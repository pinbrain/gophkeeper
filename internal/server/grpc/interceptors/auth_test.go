package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/pinbrain/gophkeeper/internal/proto"
	appCtx "github.com/pinbrain/gophkeeper/internal/server/context"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	jwt_mocks "github.com/pinbrain/gophkeeper/internal/server/jwt/mocks"
	"github.com/pinbrain/gophkeeper/internal/storage/mocks"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockJWT := jwt_mocks.NewMockServiceI(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"

	authInterceptor := NewAuthInterceptor(masterKey, mockStorage, mockJWT, log.WithField("instance", "grpcTransport"))
	handler := func(_ context.Context, req any) (any, error) {
		return req, nil
	}

	type storage struct {
		user *model.User
		err  error
	}
	type jwtService struct {
		err      error
		userData *jwt.Claims
	}
	tests := []struct {
		name       string
		storage    *storage
		jwt        string
		jwtService jwtService
		wantErr    bool
		errCode    codes.Code
	}{
		{
			name: "Успешный запрос",
			storage: &storage{
				user: &model.User{
					ID:              "1",
					Login:           "user",
					EncryptedSecret: "ce4ef7c0df5d1738675b5f16d7c7bccf5e2267a09d6e8d3115c26fbab619aed088abd055ba50d550e8d9f578f14ed095804c5fe6014f44e4a4e40665",
				},
			},
			jwt: "some_jwt",
			jwtService: jwtService{
				err: nil,
				userData: &jwt.Claims{
					UserID: "1",
					Login:  "user",
				},
			},
			wantErr: false,
		},
		{
			name:    "Нет jwt в мете запроса",
			wantErr: false,
		},
		{
			name: "Невалидный jwt",
			jwt:  "some_jwt",
			jwtService: jwtService{
				err: errors.New("jwt service error"),
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name: "Пользователь не найден",
			storage: &storage{
				err: postgres.ErrNoUser,
			},
			jwt: "some_jwt",
			jwtService: jwtService{
				err: nil,
				userData: &jwt.Claims{
					UserID: "1",
					Login:  "user",
				},
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Ошибка чтения данных из БД",
			storage: &storage{
				err: errors.New("db error"),
			},
			jwt: "some_jwt",
			jwtService: jwtService{
				err: nil,
				userData: &jwt.Claims{
					UserID: "1",
					Login:  "user",
				},
			},
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(make(map[string]string))
			mockJWT.EXPECT().GetMdJWTKey().Times(1).Return("jwt")
			if tt.jwt != "" {
				mockJWT.EXPECT().GetJWTClaims(tt.jwt).Times(1).Return(tt.jwtService.userData, tt.jwtService.err)
				md.Set("jwt", tt.jwt)
			}
			if tt.storage != nil {
				mockStorage.EXPECT().GetUserByLogin(gomock.Any(), tt.jwtService.userData.Login).
					Times(1).Return(tt.storage.user, tt.storage.err)
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)
			info := &grpc.UnaryServerInfo{}
			_, err = authInterceptor.AuthenticateUser(ctx, nil, info, handler)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestRequireUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockJWT := jwt_mocks.NewMockServiceI(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"

	authInterceptor := NewAuthInterceptor(masterKey, mockStorage, mockJWT, log.WithField("instance", "grpcTransport"))
	handler := func(_ context.Context, req any) (any, error) {
		return req, nil
	}

	tests := []struct {
		name    string
		method  string
		user    *appCtx.CtxUser
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "Успешный запрос в защищенный метод",
			method:  proto.VaultService_AddData_FullMethodName,
			user:    &appCtx.CtxUser{ID: "1"},
			wantErr: false,
		},
		{
			name:    "Успешный запрос в незащищенный метод",
			method:  proto.UserService_Login_FullMethodName,
			wantErr: false,
		},
		{
			name:    "Ошибка",
			method:  proto.VaultService_AddData_FullMethodName,
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}
			info := &grpc.UnaryServerInfo{FullMethod: tt.method}
			_, err = authInterceptor.RequireUser(ctx, nil, info, handler)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}
