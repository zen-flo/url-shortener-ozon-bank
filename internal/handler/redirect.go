package handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"url-shortener-ozon-bank/internal/storage"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		log.Printf("missing code in request")
		http.NotFound(w, r)
		return
	}

	url, err := h.shortener.Resolve(shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Printf("short code %q not found: %s", shortCode, err)
			http.NotFound(w, r)
			return
		}
		log.Printf("failed to resolve short code %q: %v", shortCode, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
	log.Printf("redirecting %q to %q", shortCode, url)
}
