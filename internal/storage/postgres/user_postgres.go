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

var (
	ErrLoginTaken = errors.New("login or email is already taken")
	ErrNoUser     = errors.New("user not found in db")
)

func (pg *PGStorage) CreateUser(ctx context.Context, user *model.User) (int, error) {
	user.Login = strings.ToLower(user.Login)
	row := pg.pool.QueryRow(
		ctx,
		"INSERT INTO users(login, email, password_hash) VALUES($1, $2, $3) RETURNING id;",
		user.Login, user.Email, user.PasswordHash,
	)
	if err := row.Scan(&user.ID); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return 0, ErrLoginTaken
			}
		}
		return 0, fmt.Errorf("failed to create new user: %w", err)
	}
	return user.ID, nil
}

func (pg *PGStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	row := pg.pool.QueryRow(
		ctx,
		"SELECT id, email, password_hash FROM users WHERE login = $1;",
		login,
	)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoUser
		}
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
	user.Login = login
	return &user, nil
}
