package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/pinbrain/gophkeeper/internal/proto"
	"google.golang.org/grpc/status"
)

// addData реализует логику передачи объекта для сохранения на сервере.
func (s *Service) addData(ctx context.Context, data []byte, meta string, dataType model.DataType) error {
	_, err := s.grpcClient.VaultClient.AddData(ctx, &proto.AddDataReq{
		Item: &proto.Item{
			Data: data,
			Type: string(dataType),
			Meta: meta,
		},
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return fmt.Errorf("не удалось сохранить данные: %s", s.Message())
		}
		return err
	}
	return nil
}

// AddPassword сохраняет пароль в хранилище.
func (s *Service) AddPassword(ctx context.Context, data string, meta model.PasswordMeta) error {
	metaB, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать строку с мета данными: %w", err)
	}
	return s.addData(ctx, []byte(data), string(metaB), model.Password)
}

// AddText сохраняет произвольные текстовые данные в хранилище.
func (s *Service) AddText(ctx context.Context, data string, meta model.TextMeta) error {
	metaB, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать строку с мета данными: %w", err)
	}
	return s.addData(ctx, []byte(data), string(metaB), model.Text)
}

// AddBankCard сохраняет данные банковской карты в хранилище.
func (s *Service) AddBankCard(ctx context.Context, data model.BankCardData, meta model.BankCardMeta) error {
	metaB, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать строку с мета данными: %w", err)
	}
	dataB, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать строку с данными банковской карты: %w", err)
	}
	return s.addData(ctx, dataB, string(metaB), model.BankCard)
}

// AddFile сохраняет файл в хранилище.
func (s *Service) AddFile(ctx context.Context, file string, comment string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	meta := model.FileMeta{
		Name:      fileInfo.Name(),
		Extension: filepath.Ext(file),
		Comment:   comment,
	}
	if !filepath.IsAbs(file) && !filepath.IsLocal(file) {
		return errors.New("невалидное полное имя файла")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	metaB, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать строку с мета данными: %w", err)
	}
	return s.addData(ctx, data, string(metaB), model.File)
}

// GetData загружает данные из хранилища.
func (s *Service) GetData(ctx context.Context, id string) (model.DataType, any, error) {
	res, err := s.grpcClient.VaultClient.GetData(ctx, &proto.GetDataReq{
		Id: id,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return "", nil, fmt.Errorf("не удалось получить данные: %s", s.Message())
		}
		return "", nil, err
	}
	resItem := res.GetItem()
	itemType := resItem.GetType()
	switch itemType {
	case string(model.Password):
		meta := &model.PasswordMeta{}
		err = json.Unmarshal([]byte(resItem.GetMeta()), meta)
		if err != nil {
			return "", nil, fmt.Errorf("не удалось прочитать мета данные: %w", err)
		}
		return model.Password, &model.PasswordItem{
			Type: model.Password,
			Meta: *meta,
			Data: string(resItem.GetData()),
		}, nil

	case string(model.BankCard):
		meta := &model.BankCardMeta{}
		err = json.Unmarshal([]byte(resItem.GetMeta()), meta)
		if err != nil {
			return "", nil, fmt.Errorf("не удалось прочитать мета данные: %w", err)
		}
		data := &model.BankCardData{}
		err = json.Unmarshal(resItem.GetData(), data)
		if err != nil {
			return "", nil, fmt.Errorf("не удалось прочитать данные: %w", err)
		}
		return model.BankCard, &model.BankCardItem{
			Type: model.BankCard,
			Meta: *meta,
			Data: model.BankCardData{
				Number:     data.Number,
				Holder:     data.Holder,
				CSV:        data.CSV,
				ValidMonth: data.ValidMonth,
				ValidYear:  data.ValidYear,
			},
		}, nil

	case string(model.Text):
		meta := &model.TextMeta{}
		err = json.Unmarshal([]byte(resItem.GetMeta()), meta)
		if err != nil {
			return "", nil, fmt.Errorf("не удалось прочитать мета данные: %w", err)
		}
		return model.Text, &model.TextItem{
			Type: model.Text,
			Meta: *meta,
			Data: string(resItem.GetData()),
		}, nil

	case string(model.File):
		meta := &model.FileMeta{}
		err = json.Unmarshal([]byte(resItem.GetMeta()), meta)
		if err != nil {
			return "", nil, fmt.Errorf("не удалось прочитать мета данные: %w", err)
		}
		file, fileErr := os.Create(meta.Name)
		if fileErr != nil {
			return "", nil, fmt.Errorf("не удалось сохранить файл: %w", fileErr)
		}
		defer file.Close()

		_, err = file.Write(resItem.GetData())
		if err != nil {
			return "", nil, fmt.Errorf("не удалось записать данные в файл: %w", err)
		}
		return model.File, &model.FileItem{
			Type: model.File,
			Meta: *meta,
		}, nil
	}
	return "", nil, fmt.Errorf("неизвестный тип данных: %s", itemType)
}

// GetAllByType получает список данных из хранилища по типу.
func (s *Service) GetAllByType(ctx context.Context, dataType model.DataType) ([]model.ItemInfo, error) {
	res, err := s.grpcClient.VaultClient.GetAllByType(ctx, &proto.GetAllByTypeReq{
		Type: string(dataType),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, fmt.Errorf("не удалось получить данные: %s", s.Message())
		}
		return nil, err
	}
	items := res.GetItems()
	result := []model.ItemInfo{}

	for _, item := range items {
		var meta any
		switch dataType {
		case model.Password:
			meta = &model.PasswordMeta{}
		case model.BankCard:
			meta = &model.BankCardMeta{}
		case model.Text:
			meta = &model.TextMeta{}
		case model.File:
			meta = &model.FileMeta{}
		default:
			return nil, fmt.Errorf("неизвестный тип данных: %s", dataType)
		}
		err = json.Unmarshal([]byte(item.GetMeta()), meta)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать мета данные: %w", err)
		}
		result = append(result, model.ItemInfo{
			ID:   item.GetId(),
			Meta: meta,
		})
	}

	return result, nil
}

// DeleteData удаляет данные из хранилища.
func (s *Service) DeleteData(ctx context.Context, id string) error {
	_, err := s.grpcClient.VaultClient.DeleteData(ctx, &proto.DeleteDataReq{
		Id: id,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return fmt.Errorf("не удалось удалить данные: %s", s.Message())
		}
		return err
	}
	return nil
}
