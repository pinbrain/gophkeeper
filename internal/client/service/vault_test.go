package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/pinbrain/gophkeeper/internal/proto"
	"github.com/pinbrain/gophkeeper/internal/proto/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	type request struct {
		data string
		meta model.PasswordMeta
	}
	tests := []struct {
		name    string
		resErr  error
		request *request
	}{
		{
			name:   "Успешный запрос",
			resErr: nil,
			request: &request{
				data: "some_password",
				meta: model.PasswordMeta{
					Resource: "some_resource",
					Login:    "user",
					Comment:  "some_comment",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metaB, err := json.Marshal(tt.request.meta)
			require.NoError(t, err)
			vaultSrvGRPCMock.EXPECT().AddData(gomock.Any(), &proto.AddDataReq{
				Item: &proto.Item{
					Data: []byte(tt.request.data),
					Type: string(model.Password),
					Meta: string(metaB),
				},
			}).Times(1).Return(&proto.AddDataRes{}, tt.resErr)

			err = service.AddPassword(context.Background(), tt.request.data, tt.request.meta)
			if tt.resErr == nil {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	type request struct {
		data string
		meta model.TextMeta
	}
	tests := []struct {
		name    string
		resErr  error
		request *request
	}{
		{
			name:   "Успешный запрос",
			resErr: nil,
			request: &request{
				data: "some_password",
				meta: model.TextMeta{
					Name:    "text name",
					Comment: "some_comment",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metaB, err := json.Marshal(tt.request.meta)
			require.NoError(t, err)
			vaultSrvGRPCMock.EXPECT().AddData(gomock.Any(), &proto.AddDataReq{
				Item: &proto.Item{
					Data: []byte(tt.request.data),
					Type: string(model.Text),
					Meta: string(metaB),
				},
			}).Times(1).Return(&proto.AddDataRes{}, tt.resErr)

			err = service.AddText(context.Background(), tt.request.data, tt.request.meta)
			if tt.resErr == nil {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAddBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	type request struct {
		data model.BankCardData
		meta model.BankCardMeta
	}
	tests := []struct {
		name    string
		resErr  error
		request *request
	}{
		{
			name:   "Успешный запрос",
			resErr: nil,
			request: &request{
				data: model.BankCardData{
					Number:     "0000",
					ValidMonth: 12,
					ValidYear:  24,
					Holder:     "user",
					CSV:        "123",
				},
				meta: model.BankCardMeta{
					Bank:    "some bank",
					Comment: "some comment",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metaB, err := json.Marshal(tt.request.meta)
			require.NoError(t, err)
			dataB, err := json.Marshal(tt.request.data)
			require.NoError(t, err)
			vaultSrvGRPCMock.EXPECT().AddData(gomock.Any(), &proto.AddDataReq{
				Item: &proto.Item{
					Data: dataB,
					Type: string(model.BankCard),
					Meta: string(metaB),
				},
			}).Times(1).Return(&proto.AddDataRes{}, tt.resErr)

			err = service.AddBankCard(context.Background(), tt.request.data, tt.request.meta)
			if tt.resErr == nil {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAddFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	tests := []struct {
		name        string
		createFile  bool
		file        string
		comment     string
		isFileError bool
	}{
		{
			name:        "Успешный запрос",
			createFile:  true,
			file:        "test_file",
			comment:     "some comment",
			isFileError: false,
		},
		{
			name:        "Невалидное имя файла",
			createFile:  false,
			file:        ".",
			comment:     "some comment",
			isFileError: true,
		},
		{
			name:        "Файл не найден",
			createFile:  false,
			file:        "test_file",
			comment:     "some comment",
			isFileError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile {
				tmpFile, err := os.Create(tt.file)
				require.NoError(t, err)
				defer os.Remove(tmpFile.Name())
			}
			if !tt.isFileError {
				metaB, err := json.Marshal(model.FileMeta{
					Name:    tt.file,
					Comment: tt.comment,
				})
				require.NoError(t, err)
				vaultSrvGRPCMock.EXPECT().AddData(gomock.Any(), &proto.AddDataReq{
					Item: &proto.Item{
						Data: []byte{},
						Type: string(model.File),
						Meta: string(metaB),
					},
				})
			}
			err := service.AddFile(context.Background(), tt.file, tt.comment)
			if !tt.isFileError {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	type want struct {
		dataType model.DataType
		item     any
		err      error
	}
	tests := []struct {
		name    string
		id      string
		grpcRes *proto.GetDataRes
		grpcErr error
		want    want
	}{
		{
			name: "Получение пароля",
			id:   "1",
			grpcRes: &proto.GetDataRes{
				Id: "1",
				Item: &proto.Item{
					Data: []byte("password"),
					Type: string(model.Password),
					Meta: `
						{
							"resource": "some_resource",
							"login": "user",
							"comment":"some comment"
						}
					`,
				},
			},
			grpcErr: nil,
			want: want{
				dataType: model.Password,
				item: &model.PasswordItem{
					Type: model.Password,
					Meta: model.PasswordMeta{
						Resource: "some_resource",
						Login:    "user",
						Comment:  "some comment",
					},
					Data: "password",
				},
				err: nil,
			},
		},
		{
			name: "Получение банковской карты",
			id:   "1",
			grpcRes: &proto.GetDataRes{
				Id: "1",
				Item: &proto.Item{
					Data: []byte(`
						{
							"number": "0000",
							"validMonth": 12,
							"validYear": 24,
							"holder": "user",
							"csv": "123"
						}
					`),
					Type: string(model.BankCard),
					Meta: `
						{
							"bank": "some bank",
							"comment":"some comment"
						}
					`,
				},
			},
			grpcErr: nil,
			want: want{
				dataType: model.BankCard,
				item: &model.BankCardItem{
					Type: model.BankCard,
					Meta: model.BankCardMeta{
						Bank:    "some bank",
						Comment: "some comment",
					},
					Data: model.BankCardData{
						Number:     "0000",
						ValidMonth: 12,
						ValidYear:  24,
						Holder:     "user",
						CSV:        "123",
					},
				},
				err: nil,
			},
		},
		{
			name: "Получение текстовые данные",
			id:   "1",
			grpcRes: &proto.GetDataRes{
				Id: "1",
				Item: &proto.Item{
					Data: []byte("some text data"),
					Type: string(model.Text),
					Meta: `
						{
							"name": "text name",
							"comment":"some comment"
						}
					`,
				},
			},
			grpcErr: nil,
			want: want{
				dataType: model.Text,
				item: &model.TextItem{
					Type: model.Text,
					Meta: model.TextMeta{
						Name:    "text name",
						Comment: "some comment",
					},
					Data: "some text data",
				},
				err: nil,
			},
		},
		{
			name:    "Ошибка получения данных",
			id:      "1",
			grpcRes: &proto.GetDataRes{},
			grpcErr: errors.New("grpc error"),
			want: want{
				err: fmt.Errorf("не удалось получить данные: %s", errors.New("grpc error")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vaultSrvGRPCMock.EXPECT().GetData(gomock.Any(), &proto.GetDataReq{Id: tt.id}).
				Times(1).Return(tt.grpcRes, tt.grpcErr)

			dataType, item, err := service.GetData(context.Background(), tt.id)
			if tt.want.err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want.dataType, dataType)
				assert.Equal(t, tt.want.item, item)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vaultSrvGRPCMock := mocks.NewMockVaultServiceClient(ctrl)
	service := NewService(&grpc.Client{VaultClient: vaultSrvGRPCMock})

	vaultSrvGRPCMock.EXPECT().DeleteData(gomock.Any(), &proto.DeleteDataReq{Id: "1"}).
		Times(1).Return(&proto.DeleteDataRes{}, nil)
	err := service.DeleteData(context.Background(), "1")
	require.NoError(t, err)
}
