package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"url-shortener-ozon-bank/internal/config"
	hndlr "url-shortener-ozon-bank/internal/handler"
	srvc "url-shortener-ozon-bank/internal/service"
	strg "url-shortener-ozon-bank/internal/storage"
	"url-shortener-ozon-bank/internal/storage/memory"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var storage strg.Storage
	switch cfg.StoreType {
	case "memory":
		storage = memory.NewInMemoryStorage()
	case "postgres":
		// a stub for Postgres
		log.Printf("Postgres storage is not implemented yet, using in-memory as fallback")
		storage = memory.NewInMemoryStorage()
	}

	service := srvc.NewShortenerService(storage)
	handler := hndlr.NewHandler(service)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	r.Post("/shorten", handler.Shorten)
	r.Get("/{code}", handler.Redirect)

	log.Printf("Starting server on port %s with store type %s", cfg.Port, cfg.StoreType)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
