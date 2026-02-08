package storage

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 16

func GenerateHash() (string, error) {
	result := make([]byte, keyLength)
	max := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < keyLength; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[n.Int64()]
	}

	return string(result), nil
}
