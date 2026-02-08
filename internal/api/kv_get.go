package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

type getResponse struct {
	Text string `json:"text"`
}

func GetKV(store *storage.Storage, crypt *crypto.Crypto) http.HandlerFunc {
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

		entry, err := store.Get(hash)
		if err != nil {
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
			http.Error(w, "decryption failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getResponse{
			Text: string(plaintext),
		})
	}
}
