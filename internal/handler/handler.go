package handler

import "url-shortener-ozon-bank/internal/service"

type Handler struct {
	shortener *service.ShortenerService
}

func NewHandler(shortener *service.ShortenerService) *Handler {
	return &Handler{
		shortener: shortener,
	}
}
