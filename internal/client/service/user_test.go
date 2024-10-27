package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/proto/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSrvGRPCMock := mocks.NewMockUserServiceClient(ctrl)
	service := NewService(&grpc.Client{UserClient: userSrvGRPCMock})

	type want struct {
		err error
		jwt string
	}
	tests := []struct {
		name     string
		login    string
		password string
		response *pb.RegisterRes
		resErr   error
		want     want
	}{
		{
			name:     "Успешный запрос",
			login:    "user",
			password: "password",
			response: &pb.RegisterRes{
				Token: "some_jwt",
			},
			resErr: nil,
			want: want{
				err: nil,
				jwt: "some_jwt",
			},
		},
		{
			name:     "Ошибка запроса",
			login:    "user",
			password: "password",
			resErr:   errors.New("grpc res error"),
			want: want{
				err: errors.New("grpc res error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(".", "jsonDB_*.json")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			viper.SetConfigFile(tmpFile.Name())

			userSrvGRPCMock.EXPECT().Register(gomock.Any(), &pb.RegisterReq{Login: tt.login, Password: tt.password}).
				Times(1).Return(tt.response, tt.resErr)

			res, err := service.Register(context.Background(), tt.login, tt.password)
			if tt.want.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want.jwt, res)
				assert.Equal(t, tt.want.jwt, config.GetJWT())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSrvGRPCMock := mocks.NewMockUserServiceClient(ctrl)
	service := NewService(&grpc.Client{UserClient: userSrvGRPCMock})

	type want struct {
		err error
		jwt string
	}
	tests := []struct {
		name     string
		login    string
		password string
		response *pb.LoginRes
		resErr   error
		want     want
	}{
		{
			name:     "Успешный запрос",
			login:    "user",
			password: "password",
			response: &pb.LoginRes{
				Token: "some_jwt",
			},
			resErr: nil,
			want: want{
				err: nil,
				jwt: "some_jwt",
			},
		},
		{
			name:     "Ошибка запроса",
			login:    "user",
			password: "password",
			resErr:   errors.New("grpc res error"),
			want: want{
				err: errors.New("grpc res error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(".", "jsonDB_*.json")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			viper.SetConfigFile(tmpFile.Name())

			userSrvGRPCMock.EXPECT().Login(gomock.Any(), &pb.LoginReq{Login: tt.login, Password: tt.password}).
				Times(1).Return(tt.response, tt.resErr)

			res, err := service.Login(context.Background(), tt.login, tt.password)
			if tt.want.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want.jwt, res)
				assert.Equal(t, tt.want.jwt, config.GetJWT())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
