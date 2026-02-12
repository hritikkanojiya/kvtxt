package storage

import (
	"crypto/rand"
	"math/big"

	"github.com/hritikkanojiya/kvtxt/internal/constant"
)

const keyLength = 16

func GenerateHash() (string, error) {
	result := make([]byte, keyLength)
	maxValue := big.NewInt(int64(len(constant.Base62Characters)))

	for i := 0; i < keyLength; i++ {
		n, err := rand.Int(rand.Reader, maxValue)
		if err != nil {
			return "", err
		}
		result[i] = constant.Base62Characters[n.Int64()]
	}

	return string(result), nil
}
