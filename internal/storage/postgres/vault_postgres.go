package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pinbrain/gophkeeper/internal/model"
)

var (
	ErrNoData = errors.New("data not found in db")
)

func (pg *PGStorage) CreateItem(ctx context.Context, userID string, item *model.VaultItem) (string, error) {
	row := pg.pool.QueryRow(
		ctx,
		`INSERT INTO user_data(user_id, encrypt_data, meta, data_type) VALUES($1, $2, $3, $4) RETURNING id;`,
		userID, item.EncryptData, item.Meta, item.Type,
	)
	if err := row.Scan(&item.ID); err != nil {
		return "", fmt.Errorf("failed to create new item: %w", err)
	}
	return item.ID, nil
}

func (pg *PGStorage) GetItem(ctx context.Context, id string, userID string) (*model.VaultItem, error) {
	var item model.VaultItem
	row := pg.pool.QueryRow(
		ctx,
		`SELECT encrypt_data, meta, data_type, created_at, updated_at FROM user_data WHERE id = $1 AND user_id = $2;`,
		id, userID,
	)
	if err := row.Scan(
		&item.EncryptData, &item.Meta, &item.Type, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoData
		}
		return nil, fmt.Errorf("failed to get data from db: %w", err)
	}
	item.ID = id
	item.UserID = userID
	return &item, nil
}

func (pg *PGStorage) DeleteItem(ctx context.Context, id string, userID string) error {
	res, err := pg.pool.Exec(ctx, `DELETE FROM user_data WHERE id = $1 AND user_id = $2;`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNoData
	}
	return nil
}
