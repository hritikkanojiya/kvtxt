package api

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

func GetKV(store *storage.Storage, crypt *crypto.Crypto, c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 {
			http.NotFound(w, r)
			return
		}

		hash := parts[2]
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		if val, ct, ok := c.Get(hash); ok {
			w.Header().Set("Content-Type", ct)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
		}

		entry, err := store.Get(hash)
		if err != nil {
			slog.Error("storage error", "error", err)
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}
		if entry == nil {
			http.NotFound(w, r)
			return
		}

		now := time.Now().Unix()
		if entry.ExpiresAt.Valid && entry.ExpiresAt.Int64 <= now {
			w.WriteHeader(http.StatusGone)
			return
		}

		plaintext, err := crypt.Decrypt(entry.Payload)
		if err != nil {
			slog.Error("decryption failed", "error", err)
			http.Error(w, "decryption failed", http.StatusInternalServerError)
			return
		}

		c.Set(entry.Hash, string(plaintext), entry.ContentType, entry.ExpiresAtPtr())

		w.Header().Set("Content-Type", entry.ContentType)
		w.WriteHeader(http.StatusOK)
		w.Write(plaintext)
	}
}
