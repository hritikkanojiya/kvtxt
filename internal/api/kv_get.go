// GetKV retrieves a stored value by key.
// Flow:
// 1. Validate key
// 2. Fetch from storage
// 3. Decrypt (if required)
// 4. Return response

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

func GetKV(store *storage.Storage, crypt *crypto.Crypto, c *cache.Cache) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		if r.Method != http.MethodGet {
			return &APIError{
				Status:  http.StatusMethodNotAllowed,
				Code:    ErrBadRequest,
				Message: "Invalid Method",
			}
		}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 {
			return &APIError{
				Status:  http.StatusNotFound,
				Code:    ErrNotFound,
				Message: "Not found",
			}
		}

		hash := parts[2]
		if hash == "" {
			return &APIError{
				Status:  http.StatusNotFound,
				Code:    ErrNotFound,
				Message: "Not found",
			}
		}

		if val, ct, ok := c.Get(hash); ok {
			w.Header().Set("Content-Type", ct)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return nil
		}

		entry, err := store.Get(hash)
		if err != nil {
			slog.Error("storage error", "error", err)
			return &APIError{
				Status:  http.StatusInternalServerError,
				Code:    ErrInternal,
				Message: "Storage error",
			}
		}
		if entry == nil {
			return &APIError{
				Status:  http.StatusNotFound,
				Code:    ErrNotFound,
				Message: "Not found",
			}
		}

		now := time.Now().Unix()
		if entry.ExpiresAt.Valid && entry.ExpiresAt.Int64 <= now {
			return &APIError{
				Status:  http.StatusGone,
				Code:    ErrConflict,
				Message: "Key expired",
			}
		}

		plaintext, err := crypt.Decrypt(entry.Payload)
		if err != nil {
			slog.Error("decryption failed", "error", err)
			return &APIError{
				Status:  http.StatusInternalServerError,
				Code:    ErrInternal,
				Message: "Decryption failed",
			}
		}

		c.Set(entry.Hash, string(plaintext), entry.ContentType, entry.ExpiresAtPtr())

		w.Header().Set("Content-Type", entry.ContentType)
		w.WriteHeader(http.StatusOK)
		w.Write(plaintext)

		return nil
	}

}
