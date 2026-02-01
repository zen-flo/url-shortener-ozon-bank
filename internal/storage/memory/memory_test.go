package memory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"url-shortener-ozon-bank/internal/storage"
)

func TestInMemoryStorage_Save(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		expectErr error
	}{
		{
			name:      "Save new URL",
			inputURL:  "https://example.com",
			expectErr: nil,
		},
		{
			name:      "Save duplicate URL",
			inputURL:  "https://example.com",
			expectErr: storage.ErrURLExists,
		},
	}
	store := NewInMemoryStorage()

	var firstCode string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortCode, err := store.Save(tt.inputURL)

			assert.Equal(t, tt.expectErr, err)

			if tt.expectErr == nil {
				assert.Len(t, shortCode, 10)
				firstCode = shortCode
			}

			if errors.Is(tt.expectErr, storage.ErrURLExists) {
				assert.Equal(t, firstCode, shortCode)
			}
		})
	}
}

func TestInMemoryStorage_Get(t *testing.T) {
	store := NewInMemoryStorage()

	code, err := store.Save("https://example.com")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		shortCode string
		expectURL string
		expectErr error
	}{
		{
			name:      "Get existing URL",
			shortCode: code,
			expectURL: "https://example.com",
			expectErr: nil,
		},
		{
			name:      "Get non-existing URL",
			shortCode: "unknown",
			expectURL: "",
			expectErr: storage.ErrURLNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := store.Get(tt.shortCode)

			assert.Equal(t, tt.expectErr, err)
			assert.Equal(t, tt.expectURL, url)
		})
	}
}
