// Package crypto provides encryption and decryption utilities
// used to protect stored values.
//
// All encryption details must remain encapsulated in this package.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Crypto struct {
	aead cipher.AEAD
}

func New(keyB64 string) (*Crypto, error) {
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, errors.New("invalid base64 encryption key")
	}

	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes (AES-256)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Crypto{aead: aead}, nil
}

// Encrypt secures plaintext value before persistence.
func (c *Crypto) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := c.aead.Seal(nil, nonce, plaintext, nil)

	// nonce || ciphertext
	return append(nonce, ciphertext...), nil
}

// Decrypt restores original value before returning to client.
func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.aead.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return c.aead.Open(nil, nonce, ciphertext, nil)
}
