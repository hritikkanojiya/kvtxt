package main

import (
	"log"

	"github.com/hritikkanojiya/kvtxt/internal/api"
	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/config"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"

	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	crypt, err := crypto.New(cfg.EncryptionKey)
	if err != nil {
		slog.Error("crypto init error", "error", err)
		os.Exit(1)
	}

	c := cache.New(1000)

	store, err := storage.Open(cfg.DBPath)
	if err != nil {
		slog.Error("storage init failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Handle("/v1/kv", api.CreateKV(store, crypt, c))

	mux.Handle("/v1/kv/", api.GetKV(store, crypt))

	log.Printf("kvtxt starting on %s\n", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, api.Logging(mux)); err != nil {
		slog.Error("server stopped", "error", err)
	}
}
