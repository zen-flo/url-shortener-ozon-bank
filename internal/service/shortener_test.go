package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"url-shortener-ozon-bank/internal/storage"
)

// fakeStorage implements the storage.Storage interface for tests
type fakeStorage struct {
	store map[string]string
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		store: make(map[string]string),
	}
}

func (f *fakeStorage) Save(originalURL string) (string, error) {
	if shortCode, ok := f.store[originalURL]; ok {
		return shortCode, storage.ErrURLExists
	}
	shortCode := "0123456789" // fixed shortCode for tests
	f.store[originalURL] = shortCode
	return shortCode, nil
}

func (f *fakeStorage) Get(shortCode string) (string, error) {
	for url, code := range f.store {
		if code == shortCode {
			return url, nil
		}
	}
	return "", storage.ErrURLNotFound
}

func TestShortenerService_ShortenAndResolve(t *testing.T) {
	store := newFakeStorage()
	serve := NewShortenerService(store)

	tests := []struct {
		name          string
		originalURL   string
		expectCode    string
		expectErr     error
		expectResolve string
	}{
		{
			name:          "Shorten new URL",
			originalURL:   "https://example.com",
			expectCode:    "0123456789",
			expectErr:     nil,
			expectResolve: "https://example.com",
		},
		{
			name:          "Shorten duplicate URL",
			originalURL:   "https://example.com",
			expectCode:    "0123456789",
			expectErr:     storage.ErrURLExists,
			expectResolve: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := serve.Shorten(tt.originalURL)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expectCode, code)

			if tt.expectResolve != "" {
				url, err := serve.Resolve(code)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResolve, url)
			} else {
				url, err := serve.Resolve(code)
				assert.ErrorIs(t, err, storage.ErrURLNotFound)
				assert.Equal(t, "", url)
			}
		})
	}

	t.Run("Resolve unknown code", func(t *testing.T) {
		url, err := serve.Resolve("nonexistent_code")
		assert.ErrorIs(t, err, storage.ErrURLNotFound)
		assert.Equal(t, "", url)
	})
}
