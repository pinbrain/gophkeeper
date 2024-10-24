// Package storage содержит реализацию хранилища сервера.
package storage

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/model"
)

// Storage описывает интерфейс хранилища приложения.
type Storage interface {
	Close() error

	UserStorage
	VaultStorage
}

// UserStorage описывает методы хранилища в части работы с пользователем.
type UserStorage interface {
	CreateUser(ctx context.Context, user *model.User) (id string, err error)
	GetUserByLogin(ctx context.Context, login string) (user *model.User, err error)
	GetUserByID(ctx context.Context, id string) (user *model.User, err error)
}

// VaultStorage описывает методы хранилища в части работы с данными.
type VaultStorage interface {
	CreateItem(ctx context.Context, userID string, item *model.VaultItem) (string, error)
	GetItem(ctx context.Context, id string, userID string) (*model.VaultItem, error)
	DeleteItem(ctx context.Context, id string, userID string) error
	GetItemsByType(ctx context.Context, dataType string, userID string) ([]model.VaultItem, error)
	UpdateItem(ctx context.Context, id string, userID string, item *model.VaultItem) error
}
