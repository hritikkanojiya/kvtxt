package storage

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateHash() (string, error) {
	buf := make([]byte, 18) // ~24 chars base64url
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
