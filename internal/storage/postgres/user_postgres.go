package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pinbrain/gophkeeper/internal/model"
)

// Ошибки, возвращаемые хранилищем.
var (
	ErrLoginTaken = errors.New("login or email is already taken")
	ErrNoUser     = errors.New("user not found in db")
)

// CreateUser создает нового пользователя.
func (pg *PGStorage) CreateUser(ctx context.Context, user *model.User) (string, error) {
	user.Login = strings.ToLower(user.Login)
	row := pg.pool.QueryRow(
		ctx,
		"INSERT INTO users(login, password_hash, encrypt_secret) VALUES($1, $2, $3) RETURNING id;",
		user.Login, user.PasswordHash, user.EncryptedSecret,
	)
	if err := row.Scan(&user.ID); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return "", ErrLoginTaken
			}
		}
		return "", fmt.Errorf("failed to create new user: %w", err)
	}
	return user.ID, nil
}

// GetUserByLogin возвращает данные пользователя по логину.
func (pg *PGStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	row := pg.pool.QueryRow(
		ctx,
		"SELECT id, password_hash, encrypt_secret FROM users WHERE login = $1;",
		login,
	)
	if err := row.Scan(&user.ID, &user.PasswordHash, &user.EncryptedSecret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoUser
		}
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
	user.Login = login
	return &user, nil
}

// GetUserByID возвращает данные пользователя по ID.
func (pg *PGStorage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	row := pg.pool.QueryRow(
		ctx,
		"SELECT login, password_hash, encrypt_secret FROM users WHERE id = $1;",
		id,
	)
	if err := row.Scan(&user.Login, &user.PasswordHash, &user.EncryptedSecret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoUser
		}
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
	user.ID = id
	return &user, nil
}
