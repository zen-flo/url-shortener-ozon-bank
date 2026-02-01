package postgres

import (
	"fmt"
	"log"
)

// PostgresStorage is a stub for Postgres that implements the Storage interface
type PostgresStorage struct {
	// ...
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	log.Printf("Postgres storage initialized with DSN: %s\n", dsn)
	return &PostgresStorage{}, nil
}

func (p *PostgresStorage) Save(originalURL string) (string, error) {
	// ...
	return "", fmt.Errorf("postgres storage Save not implemented yet")
}

func (p *PostgresStorage) Get(shortCode string) (string, error) {
	// ...
	return "", fmt.Errorf("postgres storage Get not implemented yet")
}
