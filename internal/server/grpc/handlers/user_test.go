package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/model"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/pinbrain/gophkeeper/internal/server/jwt"
	"github.com/pinbrain/gophkeeper/internal/server/utils"
	"github.com/pinbrain/gophkeeper/internal/storage/mocks"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	jwtService := jwt.NewJWTService(config.JWTConfig{
		LifeTime:  360,
		SecretKey: "some_secret_key",
		MetaKey:   "jwt",
	})
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCUserHandler(masterKey, mockStorage, jwtService, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err    error
		userID string
	}

	tests := []struct {
		name    string
		request *pb.RegisterReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			request: &pb.RegisterReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err:    nil,
				userID: "1",
			},
			wantErr: false,
		},
		{
			name: "Логин уже занят",
			request: &pb.RegisterReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err: postgres.ErrLoginTaken,
			},
			wantErr: true,
			errCode: codes.AlreadyExists,
		},
		{
			name: "Ошибка БД",
			request: &pb.RegisterReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err: errors.New("db error"),
			},
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.store != nil {
				mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(tt.store.userID, tt.store.err)
			} else {
				mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			}

			response, err := handler.Register(context.Background(), tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
				userData, err := jwtService.GetJWTClaims(response.GetToken())
				require.NoError(t, err)
				assert.Equal(t, tt.request.GetLogin(), userData.Login)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	jwtService := jwt.NewJWTService(config.JWTConfig{
		LifeTime:  360,
		SecretKey: "some_secret_key",
		MetaKey:   "jwt",
	})
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCUserHandler(masterKey, mockStorage, jwtService, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err      error
		login    string
		user     *model.User
		calcHash bool
	}
	tests := []struct {
		name    string
		request *pb.LoginReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			request: &pb.LoginReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err:   nil,
				login: "user",
				user: &model.User{
					ID:              "1",
					Login:           "user",
					EncryptedSecret: "secret",
				},
				calcHash: true,
			},
			wantErr: false,
		},
		{
			name: "Неверный пароль",
			request: &pb.LoginReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err:   nil,
				login: "user",
				user: &model.User{
					ID:              "1",
					Login:           "user",
					EncryptedSecret: "secret",
				},
				calcHash: false,
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name: "Неверный логин",
			request: &pb.LoginReq{
				Login:    "user",
				Password: "password",
			},
			store: &Store{
				err:      postgres.ErrNoUser,
				login:    "user",
				calcHash: false,
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.store != nil {
				if tt.store.calcHash {
					hash, err := utils.GeneratePasswordHash(tt.request.GetPassword())
					require.NoError(t, err)
					tt.store.user.PasswordHash = hash
				}
				mockStorage.EXPECT().GetUserByLogin(gomock.Any(), tt.store.login).Times(1).Return(tt.store.user, tt.store.err)
			} else {
				mockStorage.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Times(0)
			}

			response, err := handler.Login(context.Background(), tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
				userData, err := jwtService.GetJWTClaims(response.GetToken())
				require.NoError(t, err)
				assert.Equal(t, tt.request.GetLogin(), userData.Login)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}
