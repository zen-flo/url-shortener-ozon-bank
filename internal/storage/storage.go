package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
)

type Storage interface {
	Save(originalURL string) (string, error)
	Get(shortCode string) (string, error)
}
