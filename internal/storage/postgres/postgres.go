package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pinbrain/gophkeeper/internal/storage"
	"github.com/pinbrain/gophkeeper/internal/storage/postgres/migrations"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

// PGStorage описывает структуру хранилища БД.
type PGStorage struct {
	pool *pgxpool.Pool
	log  *logrus.Entry
}

// NewStorage создает и возвращает новое хранилище.
func NewStorage(ctx context.Context, dsn string, logger *logrus.Logger) (storage.Storage, error) {
	log := logger.WithField("instance", "pgStorage")
	if err := runMigrations(dsn, log); err != nil {
		return nil, fmt.Errorf("failed to run db migration: %w", err)
	}
	pool, err := initPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialized a db connection: %w", err)
	}
	return &PGStorage{pool: pool, log: log}, nil
}

// Close закрывает пулл и все соединения с БД.
func (pg *PGStorage) Close() error {
	pg.pool.Close()
	return nil
}

// initPool инициализация пула для соединения с БД.
func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return pool, nil
}

// runMigrations запускает миграции БД.
func runMigrations(dsn string, logger *logrus.Entry) error {
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(logger)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	sqlDB, err := goose.OpenDBWithDriver("postgres", dsn)
	defer func() {
		if err = sqlDB.Close(); err != nil {
			logger.WithField("err", err).Error("failed to close db connection while migration")
		}
	}()
	if err != nil {
		return err
	}
	if err = goose.Up(sqlDB, "."); err != nil {
		return err
	}
	return nil
}
