package handlers

import (
	"context"
	"errors"

	"github.com/pinbrain/gophkeeper/internal/model"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	appCtx "github.com/pinbrain/gophkeeper/internal/server/context"
	"github.com/pinbrain/gophkeeper/internal/server/utils"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCVaultHandler struct {
	pb.UnimplementedVaultServiceServer
	masterKey string
	storage   storage.Storage
	log       *logrus.Entry
}

func NewGRPCVaultHandler(masterKey string, storage storage.Storage, log *logrus.Entry) *GRPCVaultHandler {
	return &GRPCVaultHandler{
		masterKey: masterKey,
		storage:   storage,
		log:       log,
	}
}

func (h *GRPCVaultHandler) AddData(ctx context.Context, in *pb.AddDataReq) (*pb.AddDataRes, error) {
	reqItem := in.GetItem()
	if reqItem == nil {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует объект для сохранения")
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil {
		h.log.Error("failed to get user from context")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	dataType := reqItem.GetType()
	if !isValidDataType(dataType) {
		return nil, status.Error(codes.InvalidArgument, "Неизвестный тип данных")
	}

	encData, err := utils.Encrypt(reqItem.GetData(), user.Secret)
	if err != nil {
		h.log.WithError(err).Error("Error while encrypting user data")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	item := &model.VaultItem{
		UserID:      user.ID,
		Meta:        reqItem.GetMeta(),
		Type:        model.DataType(dataType),
		EncryptData: encData,
	}
	_, err = h.storage.CreateItem(ctx, user.ID, item)
	if err != nil {
		h.log.WithError(err).Error("Error while saving data")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &pb.AddDataRes{}, nil
}

func (h *GRPCVaultHandler) GetData(ctx context.Context, in *pb.GetDataReq) (*pb.GetDataRes, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует id данных")
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil {
		h.log.Error("failed to get user from context")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	data, err := h.storage.GetItem(ctx, in.GetId(), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoData):
			return nil, status.Error(codes.NotFound, "Данные не найдены")
		default:
			h.log.WithError(err).Error("Error while getting item")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}
	decData, err := utils.Decrypt(data.EncryptData, user.Secret)
	if err != nil {
		h.log.WithError(err).Error("Error while decrypting user data")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	response := &pb.GetDataRes{
		Id: data.ID,
		Item: &pb.Item{
			Data: decData,
			Type: string(data.Type),
			Meta: data.Meta,
		},
	}
	return response, nil
}

func (h *GRPCVaultHandler) DeleteData(ctx context.Context, in *pb.DeleteDataReq) (*pb.DeleteDataRes, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует id данных")
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil {
		h.log.Error("failed to get user from context")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	err := h.storage.DeleteItem(ctx, in.GetId(), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoData):
			return nil, status.Error(codes.NotFound, "Данные для удаления не найдены")
		default:
			h.log.WithError(err).Error("Error while deleting item")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}
	return &pb.DeleteDataRes{}, nil
}

func (h *GRPCVaultHandler) GetAllByType(ctx context.Context, in *pb.GetAllByTypeReq) (*pb.GetAllByTypeRes, error) {
	dataType := in.GetType()
	if !isValidDataType(dataType) {
		return nil, status.Error(codes.InvalidArgument, "Неизвестный тип данных")
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil {
		h.log.Error("failed to get user from context")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	items, err := h.storage.GetItemsByType(ctx, dataType, user.ID)
	if err != nil {
		h.log.WithError(err).Error("Error while deleting item")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	var responseItems []*pb.GetAllByTypeRes_TypeItem
	for _, item := range items {
		responseItems = append(responseItems, &pb.GetAllByTypeRes_TypeItem{
			Id:   item.ID,
			Meta: item.Meta,
		})
	}
	return &pb.GetAllByTypeRes{
		Items: responseItems,
	}, nil
}

func (h *GRPCVaultHandler) UpdateData(ctx context.Context, in *pb.UpdateDataReq) (*pb.UpdateDataRes, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует id данных")
	}
	user := appCtx.GetCtxUser(ctx)
	if user == nil {
		h.log.Error("failed to get user from context")
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	encData, err := utils.Encrypt(in.GetData(), user.Secret)
	if err != nil {
		h.log.WithError(err).Error("Error while encrypting user data")
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	item := &model.VaultItem{
		ID:          in.GetId(),
		UserID:      user.ID,
		Meta:        in.GetMeta(),
		EncryptData: encData,
	}
	err = h.storage.UpdateItem(ctx, in.GetId(), user.ID, item)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNoData):
			return nil, status.Error(codes.NotFound, "Данные для обновления не найдены")
		default:
			h.log.WithError(err).Error("Error while updating item")
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}
	return &pb.UpdateDataRes{}, nil
}

func isValidDataType(dataType string) bool {
	switch model.DataType(dataType) {
	case model.Password, model.Text, model.BankCard, model.File:
		return true
	}
	return false
}
