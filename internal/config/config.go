package config

import (
	"errors"
	"os"
)

type Config struct {
	Addr          string
	DBPath        string
	EncryptionKey string
}

func Load() (*Config, error) {
	cfg := &Config{
		Addr:          os.Getenv("KVTXT_ADDR"),
		DBPath:        os.Getenv("KVTXT_DB_PATH"),
		EncryptionKey: os.Getenv("KVTXT_ENCRYPTION_KEY"),
	}

	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}

	if cfg.DBPath == "" {
		return nil, errors.New("KVTXT_DB_PATH is required")
	}

	if cfg.EncryptionKey == "" {
		return nil, errors.New("KVTXT_ENCRYPTION_KEY is required")
	}

	return cfg, nil
}
