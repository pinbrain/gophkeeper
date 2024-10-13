package storage

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/model"
)

type Storage interface {
	Close() error

	UserStorage
	VaultStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *model.User) (id string, err error)
	GetUserByLogin(ctx context.Context, login string) (user *model.User, err error)
	GetUserByID(ctx context.Context, id string) (user *model.User, err error)
}

type VaultStorage interface {
	CreateItem(ctx context.Context, userID string, item *model.VaultItem) (string, error)
	GetItem(ctx context.Context, id string, userID string) (*model.VaultItem, error)
	DeleteItem(ctx context.Context, id string, userID string) error
}
