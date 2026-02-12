package main

import (
	"errors"

	"github.com/hritikkanojiya/kvtxt/internal/api"
	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/config"
	"github.com/hritikkanojiya/kvtxt/internal/constant"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
	"github.com/hritikkanojiya/kvtxt/internal/worker"

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

	ctx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	worker.StartCleanupWorker(
		ctx,
		store,
		constant.CleanupInterval,
	)

	mux := http.NewServeMux()

	api.RegisterRoute(
		mux,
		"/health",
		http.MethodGet,
		api.Health(),
	)

	api.RegisterRoute(
		mux,
		"/ready",
		http.MethodGet,
		api.Readiness(store),
	)

	mux.Handle(
		"/v1/kv",
		api.Adapt(
			api.AllowMethods(http.MethodPost)(
				api.CreateKV(store, crypt, c),
			),
		),
	)

	mux.Handle(
		"/v1/kv/",
		api.Adapt(
			api.AllowMethods(http.MethodGet)(
				api.GetKV(store, crypt, c),
			),
		),
	)

	maxSizeMB := cfg.MaxPayloadSize
	if maxSizeMB <= 0 {
		slog.Warn("invalid max payload size, using default", "value", cfg.MaxPayloadSize)
		maxSizeMB = constant.DefaultMaxPayloadSizeMB
	}

	maxPayloadSize := int64(maxSizeMB) * constant.MB

	var handler http.Handler = mux
	handler = api.MaxBodySize(maxPayloadSize)(handler)
	handler = api.Logging(handler)
	handler = api.RequestID(handler)

	srv := &http.Server{
		Addr:         cfg.AppPort,
		Handler:      handler,
		ReadTimeout:  constant.ReadTimeout,
		WriteTimeout: constant.WriteTimeout,
		IdleTimeout:  constant.IdleTimeout,
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

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		constant.ShutdownTimeout,
	)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}

	signal.Stop(stop)
	close(stop)
}
