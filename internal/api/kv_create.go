package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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

func CreateKV(store *storage.Storage, crypt *crypto.Crypto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		if req.Text == "" {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}

		encrypted, err := crypt.Encrypt([]byte(req.Text))
		if err != nil {
			http.Error(w, "encryption failed", http.StatusInternalServerError)
			return
		}

		hash, err := storage.GenerateHash()
		if err != nil {
			http.Error(w, "hash generation failed", http.StatusInternalServerError)
			return
		}

		now := time.Now().Unix()

		var expires sql.NullInt64
		if req.TTLSeconds != nil && *req.TTLSeconds > 0 {
			expires = sql.NullInt64{
				Int64: now + *req.TTLSeconds,
				Valid: true,
			}
		}

		entry := &storage.Entry{
			Hash:      hash,
			Payload:   encrypted,
			CreatedAt: now,
			ExpiresAt: expires,
		}

		if err := store.Insert(entry); err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(createResponse{
			Key: hash,
		})
	}
}
