package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"url-shortener-ozon-bank/internal/storage"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		log.Printf("missing url in request body")
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	shortCode, err := h.shortener.Shorten(req.URL)
	if err != nil && !errors.Is(err, storage.ErrURLExists) {
		log.Printf("failed to shorten url %q: %v", req.URL, err)
		http.Error(w, "failed to shorten url", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortCode: shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
