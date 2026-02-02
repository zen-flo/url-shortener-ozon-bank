package postgres

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"url-shortener-ozon-bank/internal/storage"
)

var testStorage *PostgresStorage

func TestMain(m *testing.M) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://shortener_user:secret@localhost:5433/shortener_test?sslmode=disable"
	}

	var err error
	testStorage, err = NewPostgresStorage(dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestPostgresStorage_Save(t *testing.T) {
	cleanDB(t)

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

	var firstCode string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortCode, err := testStorage.Save(tt.inputURL)

			assert.ErrorIs(t, err, tt.expectErr)

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

func TestPostgresStorage_Get(t *testing.T) {
	cleanDB(t)

	code, err := testStorage.Save("https://example.com")
	require.NoError(t, err)

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
			url, err := testStorage.Get(tt.shortCode)

			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expectURL, url)
		})
	}
}

func cleanDB(t *testing.T) {
	_, err := testStorage.db.Exec("TRUNCATE TABLE urls")
	require.NoError(t, err)
}
