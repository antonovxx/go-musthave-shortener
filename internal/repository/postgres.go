package repository

import (
	"antonovxx/go-musthave-shortener/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const dbTimeout = 3 * time.Second

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveIfNotExists(id, originalURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO urls (short_id, original_url) VALUES ($1, $2)
		 ON CONFLICT (short_id) DO NOTHING`,
		id, originalURL)

	if err != nil {
		return fmt.Errorf("insert url: %w", err)
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return model.ErrIDExists
	}

	return nil
}

func (r *PostgresRepository) Get(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var originalURL string

	err := r.db.QueryRowContext(ctx,
		`SELECT original_url FROM urls WHERE short_id = $1`, id).Scan(&originalURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", model.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("select url: %w", err)
	}

	return originalURL, nil
}
