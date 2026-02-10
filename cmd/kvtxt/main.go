package main

import (
	"errors"

	"github.com/hritikkanojiya/kvtxt/internal/api"
	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/config"
	"github.com/hritikkanojiya/kvtxt/internal/constant"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"

	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	c := cache.New(constant.DefaultCacheSize)

	store, err := storage.Open(cfg.DatabaseFilePath)
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

	mux.Handle("/v1/kv/", api.GetKV(store, crypt, c))

	maxSizeMB := cfg.MaxPayloadSize
	if maxSizeMB <= 0 {
		slog.Warn("invalid max payload size, using default", "value", cfg.MaxPayloadSize)
		maxSizeMB = constant.DefaultMaxPayloadSizeMB
	}

	maxPayloadSize := int64(maxSizeMB) * constant.MB

	var handler http.Handler = mux
	handler = api.MaxBodySize(maxPayloadSize)(handler)
	handler = api.Logging(handler)

	srv := &http.Server{
		Addr:    cfg.AppPort,
		Handler: handler,
	}

	go func() {
		slog.Info("server listening", "addr", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}
