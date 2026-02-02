package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"url-shortener-ozon-bank/internal/shortcode"
	"url-shortener-ozon-bank/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		short_code VARCHAR(10) PRIMARY KEY,
	    original_url TEXT UNIQUE NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) Save(originalURL string) (string, error) {
	var existingCode string
	err := p.db.QueryRow("SELECT short_code FROM urls WHERE original_url = $1", originalURL).Scan(&existingCode)
	if err == nil {
		return existingCode, storage.ErrURLExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to query existing url: %w", err)
	}

	for {
		shortCode, err := shortcode.GenerateCode()
		if err != nil {
			return "", fmt.Errorf("failed to generate code: %w", err)
		}

		_, err = p.db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, originalURL)
		if err != nil {
			if isUniqueViolation(err) {
				continue // the shortCode already exists, creating a new one
			}
			return "", fmt.Errorf("failed to insert url: %w", err)
		}

		return shortCode, nil
	}
}

func (p *PostgresStorage) Get(shortCode string) (string, error) {
	var originalURL string
	err := p.db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	} else if err != nil {
		return "", fmt.Errorf("failed to query existing url: %w", err)
	}

	return originalURL, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if ok := errors.As(err, &pgErr); ok {
		return pgErr.Code == "23505" // not unique
	}
	return false
}
