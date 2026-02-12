package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hritikkanojiya/kvtxt/internal/cache"
	"github.com/hritikkanojiya/kvtxt/internal/constant"
	"github.com/hritikkanojiya/kvtxt/internal/crypto"
	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

type createRequest struct {
	Text        json.RawMessage `json:"text"`
	ContentType string          `json:"content_type"`
	TTLSeconds  *int64          `json:"ttl_seconds"`
}

type createResponse struct {
	Key       string `json:"key"`
	ExpiresAt *int64 `json:"expires_at,omitempty"`
}

func CreateKV(store *storage.Storage, crypt *crypto.Crypto, c *cache.Cache) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		if r.Method != http.MethodPost {
			return &APIError{
				Status:  http.StatusMethodNotAllowed,
				Code:    ErrBadRequest,
				Message: "Invalid Method",
			}
		}

		defer r.Body.Close()

		var req createRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if errors.Is(err, http.ErrBodyReadAfterClose) ||
				strings.Contains(err.Error(), "request body too large") {
				return &APIError{
					Status:  http.StatusRequestEntityTooLarge,
					Code:    ErrPayloadTooLarge,
					Message: "Request body exceeds allowed size",
				}
			}

			return &APIError{
				Status:  http.StatusBadRequest,
				Code:    ErrInvalidJSON,
				Message: "Invalid JSON body",
			}
		}

		if len(req.Text) == 0 {
			return &APIError{
				Status:  http.StatusBadRequest,
				Code:    ErrInvalidJSON,
				Message: "Text is required",
			}
		}

		if req.ContentType == "" {
			req.ContentType = "text/plain; charset=utf-8"
		}

		switch {
		case req.ContentType == "application/json":
			if !json.Valid(req.Text) {
				return &APIError{
					Status:  http.StatusBadRequest,
					Code:    ErrInvalidJSON,
					Message: "Invalid JSON body",
				}
			}

		case len(req.ContentType) >= 5 && req.ContentType[:5] == "text/":
			if !utf8.Valid(req.Text) {
				return &APIError{
					Status:  http.StatusBadRequest,
					Code:    ErrBadRequest,
					Message: "Invalid utf-8 text",
				}
			}
		}

		var ttl int64

		if req.TTLSeconds != nil {
			ttl = *req.TTLSeconds
		} else {
			ttl = constant.DefaultTTLSeconds
		}

		if ttl < constant.MinTTLSeconds {
			return &APIError{
				Status:  http.StatusBadRequest,
				Code:    ErrBadRequest,
				Message: "TTL must be greater than zero",
			}
		}

		if ttl > constant.MaxTTLSeconds {
			return &APIError{
				Status:  http.StatusBadRequest,
				Code:    ErrBadRequest,
				Message: "TTL exceeds maximum allowed",
			}
		}

		encrypted, err := crypt.Encrypt([]byte(req.Text))
		if err != nil {
			slog.Error("encryption failed", "error", err)
			return &APIError{
				Status:  http.StatusInternalServerError,
				Code:    ErrInternal,
				Message: "Encryption failed",
			}
		}

		now := time.Now().Unix()

		var expires sql.NullInt64
		expiryTime := time.Now().Add(time.Duration(ttl) * time.Second)

		expires = sql.NullInt64{
			Int64: expiryTime.Unix(),
			Valid: true,
		}

		var entry *storage.Entry

		const maxAttempts = 5
		for i := 0; i < maxAttempts; i++ {
			hash, err := storage.GenerateHash()
			if err != nil {
				slog.Error("hash generation failed", "error", err)
				return &APIError{
					Status:  http.StatusInternalServerError,
					Code:    ErrInternal,
					Message: "Hash generation failed",
				}
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
			return &APIError{
				Status:  http.StatusInternalServerError,
				Code:    ErrInternal,
				Message: "Storage Error",
			}
		}

		if entry == nil {
			slog.Error("hash collision retries exhausted")
			return &APIError{
				Status:  http.StatusInternalServerError,
				Code:    ErrInternal,
				Message: "Could not generate unique key",
			}
		}

		c.Set(entry.Hash, string(req.Text), entry.ContentType, entry.ExpiresAtPtr())

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(createResponse{
			Key:       entry.Hash,
			ExpiresAt: entry.ExpiresAtPtr(),
		})

		return nil
	}
}
