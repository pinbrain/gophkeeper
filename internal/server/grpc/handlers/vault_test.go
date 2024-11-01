package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/model"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	appCtx "github.com/pinbrain/gophkeeper/internal/server/context"
	"github.com/pinbrain/gophkeeper/internal/server/utils"
	"github.com/pinbrain/gophkeeper/internal/storage/mocks"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCVaultHandler(masterKey, mockStorage, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err error
	}
	tests := []struct {
		name    string
		user    *appCtx.CtxUser
		request *pb.AddDataReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "user",
				Secret: masterKey,
			},
			request: &pb.AddDataReq{
				Item: &pb.Item{
					Data: []byte("123"),
					Type: string(model.Password),
					Meta: "some meta info",
				},
			},
			store: &Store{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Отсутствуют данные",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "user",
				Secret: masterKey,
			},
			request: &pb.AddDataReq{
				Item: nil,
			},
			store:   nil,
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка получения данных пользователя",
			user: nil,
			request: &pb.AddDataReq{
				Item: &pb.Item{
					Data: []byte("123"),
					Type: string(model.Password),
					Meta: "some meta info",
				},
			},
			store:   nil,
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Некорректный тип данных",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "user",
				Secret: masterKey,
			},
			request: &pb.AddDataReq{
				Item: &pb.Item{
					Data: []byte("123"),
					Type: "INVALID",
					Meta: "some meta info",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка сохранения данных в БД",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "user",
				Secret: masterKey,
			},
			request: &pb.AddDataReq{
				Item: &pb.Item{
					Data: []byte("123"),
					Type: string(model.Password),
					Meta: "some meta info",
				},
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
				mockStorage.EXPECT().CreateItem(gomock.Any(), tt.user.ID, gomock.Any()).DoAndReturn(
					func(_ context.Context, userID string, item *model.VaultItem) (string, error) {
						if len(item.EncryptData) == 0 {
							t.Errorf("EncryptData is nil or empty")
						}
						if userID != tt.user.ID ||
							item.Meta != tt.request.GetItem().GetMeta() ||
							item.Type != model.DataType(tt.request.GetItem().GetType()) {
							t.Errorf("Unexpected VaultItem data: got %+v", item)
						}
						return "1", tt.store.err
					},
				)
			}

			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}

			_, err = handler.AddData(ctx, tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCVaultHandler(masterKey, mockStorage, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err     error
		resItem *model.VaultItem
	}
	tests := []struct {
		name    string
		user    *appCtx.CtxUser
		request *pb.GetDataReq
		data    []byte
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "password",
				Secret: masterKey,
			},
			request: &pb.GetDataReq{
				Id: "1",
			},
			data: []byte("some stored data"),
			store: &Store{
				err: nil,
				resItem: &model.VaultItem{
					ID:     "1",
					UserID: "1",
					Meta:   "some data meta",
					Type:   "PASSWORD",
				},
			},
			wantErr: false,
		},
		{
			name:    "Нет id в запросе",
			request: &pb.GetDataReq{},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка получения пользователя запроса",
			request: &pb.GetDataReq{
				Id: "1",
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Данные не найдены",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "password",
				Secret: masterKey,
			},
			request: &pb.GetDataReq{
				Id: "1",
			},
			store: &Store{
				err: postgres.ErrNoData,
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Ошибка БД",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "password",
				Secret: masterKey,
			},
			request: &pb.GetDataReq{
				Id: "1",
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
				mockStorage.EXPECT().GetItem(gomock.Any(), tt.request.GetId(), tt.user.ID).DoAndReturn(
					func(ctx context.Context, id string, userID string) (*model.VaultItem, error) {
						if tt.store.err != nil {
							return nil, tt.store.err
						}
						encData, err := utils.Encrypt(tt.data, tt.user.Secret)
						require.NoError(t, err)
						tt.store.resItem.EncryptData = encData
						return tt.store.resItem, nil
					},
				)
			}

			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}

			response, err := handler.GetData(ctx, tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.data, response.GetItem().GetData())
				assert.Equal(t, tt.store.resItem.Meta, response.GetItem().GetMeta())
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCVaultHandler(masterKey, mockStorage, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err error
	}
	tests := []struct {
		name    string
		user    *appCtx.CtxUser
		request *pb.DeleteDataReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "user",
			},
			request: &pb.DeleteDataReq{
				Id: "1",
			},
			store: &Store{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Нет id в запросе",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "user",
			},
			request: &pb.DeleteDataReq{},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка получения пользователя запроса",
			request: &pb.DeleteDataReq{
				Id: "1",
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Данные не найдены",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "password",
			},
			request: &pb.DeleteDataReq{
				Id: "1",
			},
			store: &Store{
				err: postgres.ErrNoData,
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Ошибка БД",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "password",
			},
			request: &pb.DeleteDataReq{
				Id: "1",
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
				mockStorage.EXPECT().DeleteItem(gomock.Any(), tt.request.GetId(), tt.user.ID).Times(1).Return(tt.store.err)
			}

			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}

			_, err = handler.DeleteData(ctx, tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestGetAllByType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCVaultHandler(masterKey, mockStorage, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err error
		res []model.VaultItem
	}
	tests := []struct {
		name    string
		user    *appCtx.CtxUser
		request *pb.GetAllByTypeReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "user",
			},
			request: &pb.GetAllByTypeReq{
				Type: string(model.Password),
			},
			store: &Store{
				err: nil,
				res: []model.VaultItem{
					{
						ID:   "1",
						Meta: "some meta",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Некорректный тип данных",
			request: &pb.GetAllByTypeReq{
				Type: "INVALID",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка получения пользователя запроса",
			request: &pb.GetAllByTypeReq{
				Type: string(model.Password),
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Ошибка БД",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "password",
			},
			request: &pb.GetAllByTypeReq{
				Type: string(model.Password),
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
				mockStorage.EXPECT().GetItemsByType(gomock.Any(), tt.request.GetType(), tt.user.ID).
					Times(1).Return(tt.store.res, tt.store.err)
			}

			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}

			response, err := handler.GetAllByType(ctx, tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
				for i, item := range response.GetItems() {
					assert.Equal(t, tt.store.res[i].ID, item.GetId())
					assert.Equal(t, tt.store.res[i].Meta, item.GetMeta())
				}
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}

func TestUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	log, err := logger.NewLogger("info")
	require.NoError(t, err)
	masterKey := "1d0e95ed9e11b59ba42200720c252f98d4cd440412926a0c15b6a95e03ab4480"
	handler := NewGRPCVaultHandler(masterKey, mockStorage, log.WithField("instance", "grpcTransport"))

	type Store struct {
		err error
	}
	tests := []struct {
		name    string
		user    *appCtx.CtxUser
		request *pb.UpdateDataReq
		store   *Store
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешный запрос",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "user",
				Secret: masterKey,
			},
			request: &pb.UpdateDataReq{
				Id:   "1",
				Data: []byte("some data"),
				Meta: "some meta",
			},
			store: &Store{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Нет id в запросе",
			user: &appCtx.CtxUser{
				ID:    "1",
				Login: "user",
			},
			request: &pb.UpdateDataReq{},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Ошибка получения данных пользователя",
			user: nil,
			request: &pb.UpdateDataReq{
				Id:   "1",
				Data: []byte("some data"),
				Meta: "some meta",
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "Данные не найдены",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "password",
				Secret: masterKey,
			},
			request: &pb.UpdateDataReq{
				Id:   "1",
				Data: []byte("some data"),
				Meta: "some meta",
			},
			store: &Store{
				err: postgres.ErrNoData,
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "Ошибка БД",
			user: &appCtx.CtxUser{
				ID:     "1",
				Login:  "password",
				Secret: masterKey,
			},
			request: &pb.UpdateDataReq{
				Id:   "1",
				Data: []byte("some data"),
				Meta: "some meta",
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
				mockStorage.EXPECT().UpdateItem(gomock.Any(), tt.request.GetId(), tt.user.ID, gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, userID string, item *model.VaultItem) error {
						if len(item.EncryptData) == 0 {
							t.Errorf("EncryptData is nil or empty")
						}
						if userID != tt.user.ID ||
							item.Meta != tt.request.GetMeta() {
							t.Errorf("Unexpected VaultItem data: got %+v", item)
						}
						return tt.store.err
					},
				)
			}

			ctx := context.Background()
			if tt.user != nil {
				ctx = appCtx.CtxWithUser(ctx, tt.user)
			}

			_, err = handler.UpdateData(ctx, tt.request)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				code, _ := status.FromError(err)
				assert.Equal(t, tt.errCode, code.Code())
			}
		})
	}
}
