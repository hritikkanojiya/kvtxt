package main

import (
	"github.com/hritikkanojiya/kvtxt/internal/api"
	"github.com/hritikkanojiya/kvtxt/internal/config"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"

	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	crypt, err := crypto.New(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("crypto init error: %v", err)
	}

	store, err := storage.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("storage init error: %v", err)
	}
	defer store.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Handle("/v1/kv", api.CreateKV(store, crypt))

	log.Printf("kvtxt starting on %s\n", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
