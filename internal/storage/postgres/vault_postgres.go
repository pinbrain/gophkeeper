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

func (pg *PGStorage) GetItemsByType(ctx context.Context, dataType string, userID string) ([]model.VaultItem, error) {
	var items []model.VaultItem
	rows, err := pg.pool.Query(ctx,
		`SELECT id, meta, created_at, updated_at FROM user_data WHERE user_id = $1 AND data_type = $2;`,
		userID, dataType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.VaultItem
		if err = rows.Scan(&item.ID, &item.Meta, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to read data from db - item row: %w", err)
		}
		item.UserID = userID
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get items by type: %w", err)
	}
	return items, nil
}

func (pg *PGStorage) UpdateItem(ctx context.Context, id string, userID string, item *model.VaultItem) error {
	res, err := pg.pool.Exec(ctx,
		`UPDATE user_data SET encrypt_data = $1, meta = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4;`,
		item.EncryptData, item.Meta, id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNoData
	}
	return nil
}
