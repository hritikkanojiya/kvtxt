package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

type createRequest struct {
	Text        json.RawMessage `json:"text"`
	ContentType string          `json:"content_type"`
	TTLSeconds  *int64          `json:"ttl_seconds"`
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
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
				return
			}

			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		if len(req.Text) == 0 {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}

		if req.ContentType == "" {
			req.ContentType = "text/plain; charset=utf-8"
		}

		switch {
		case req.ContentType == "application/json":
			if !json.Valid(req.Text) {
				http.Error(w, "invalid json payload", http.StatusBadRequest)
				return
			}

		case len(req.ContentType) >= 5 && req.ContentType[:5] == "text/":
			if !utf8.Valid(req.Text) {
				http.Error(w, "invalid utf-8 text", http.StatusBadRequest)
				return
			}
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
				Hash:        hash,
				Payload:     encrypted,
				ContentType: req.ContentType,
				CreatedAt:   now,
				ExpiresAt:   expires,
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

		c.Set(entry.Hash, string(req.Text), entry.ContentType, entry.ExpiresAtPtr())

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createResponse{
			Key: entry.Hash,
		})
	}
}
