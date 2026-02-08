package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

type createRequest struct {
	Text       string `json:"text"`
	TTLSeconds *int64 `json:"ttl_seconds"`
}

type createResponse struct {
	Key string `json:"key"`
}

func CreateKV(store *storage.Storage, crypt *crypto.Crypto, c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		if req.Text == "" {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}

		if req.TTLSeconds != nil && *req.TTLSeconds <= 0 {
			http.Error(w, "ttl_seconds must be greater than zero", http.StatusBadRequest)
			return
		}

		encrypted, err := crypt.Encrypt([]byte(req.Text))
		if err != nil {
			slog.Error("encryption failed", "error", err)
			http.Error(w, "encryption failed", http.StatusInternalServerError)
			return
		}

		now := time.Now().Unix()

		var expires sql.NullInt64
		if req.TTLSeconds != nil {
			expires = sql.NullInt64{
				Int64: now + *req.TTLSeconds,
				Valid: true,
			}
		}

		var entry *storage.Entry

		const maxAttempts = 5
		for i := 0; i < maxAttempts; i++ {
			hash, err := storage.GenerateHash()
			if err != nil {
				slog.Error("hash generation failed", "error", err)
				http.Error(w, "hash generation failed", http.StatusInternalServerError)
				return
			}

			entry = &storage.Entry{
				Hash:      hash,
				Payload:   encrypted,
				CreatedAt: now,
				ExpiresAt: expires,
			}

			err = store.Insert(entry)
			if err == nil {
				break
			}

			if storage.IsUniqueConstraintError(err) {
				continue
			}

			slog.Error("insert failed", "error", err)
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		if entry == nil {
			slog.Error("hash collision retries exhausted")
			http.Error(w, "could not generate unique key", http.StatusInternalServerError)
			return
		}

		c.Set(entry.Hash, req.Text, entry.ExpiresAtPtr())

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createResponse{
			Key: entry.Hash,
		})
	}
}
