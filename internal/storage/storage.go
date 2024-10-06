package storage

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/model"
)

type Storage interface {
	Close() error

	UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *model.User) (id int, err error)
	GetUserByLogin(ctx context.Context, login string) (user *model.User, err error)
}
