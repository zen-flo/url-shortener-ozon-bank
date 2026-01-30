package service

import (
	"url-shortener-ozon-bank/internal/storage"
)

type ShortenerService struct {
	storage storage.Storage
}

func NewShortenerService(storage storage.Storage) *ShortenerService {
	return &ShortenerService{
		storage: storage,
	}
}

func (s *ShortenerService) Shorten(originalURL string) (string, error) {
	return s.storage.Save(originalURL)
}
