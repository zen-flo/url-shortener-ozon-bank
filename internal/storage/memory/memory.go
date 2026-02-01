package memory

import (
	"sync"
	"url-shortener-ozon-bank/internal/shortcode"
	"url-shortener-ozon-bank/internal/storage"
)

type InMemoryStorage struct {
	mu        sync.RWMutex
	codeToURL map[string]string
	urlToCode map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (m *InMemoryStorage) Save(originalURL string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if shortCode, ok := m.urlToCode[originalURL]; ok {
		return shortCode, storage.ErrURLExists
	}

	for {
		shortCode, err := shortcode.GenerateCode()
		if err != nil {
			return "", err
		}

		if _, ok := m.codeToURL[shortCode]; ok {
			continue // the shortCode already exists, creating a new one
		}

		m.codeToURL[shortCode] = originalURL
		m.urlToCode[originalURL] = shortCode
		return shortCode, nil
	}
}

func (m *InMemoryStorage) Get(shortCode string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, ok := m.codeToURL[shortCode]
	if !ok {
		return "", storage.ErrURLNotFound
	}
	return url, nil
}
