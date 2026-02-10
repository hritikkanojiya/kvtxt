package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hritikkanojiya/kvtxt/internal/constant"
)

type Config struct {
	AppPort          string
	DatabaseFilePath string
	EncryptionKey    string
	MaxPayloadSize   int
}

func Load() (*Config, error) {
	cfg := &Config{
		AppPort:          os.Getenv("KVTXT_PORT"),
		DatabaseFilePath: os.Getenv("KVTXT_DB_PATH"),
		EncryptionKey:    os.Getenv("KVTXT_ENCRYPTION_KEY"),
		MaxPayloadSize:   getEnvInt("KVTXT_MAX_PAYLOAD_SIZE", constant.DefaultMaxPayloadSizeMB),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = constant.DefaultPort
	}

	if !strings.HasPrefix(cfg.AppPort, ":") {
		return nil, fmt.Errorf("invalid KVTXT_PORT format: %s", cfg.AppPort)
	}

	port, err := strconv.Atoi(strings.TrimPrefix(cfg.AppPort, ":"))
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid KVTXT_PORT value: %s", cfg.AppPort)
	}

	if cfg.DatabaseFilePath == "" {
		return nil, errors.New("KVTXT_DB_PATH is required")
	}

	if _, err := os.Stat(cfg.DatabaseFilePath); err != nil {
		return nil, fmt.Errorf("invalid KVTXT_DB_PATH: %w", err)
	}

	if cfg.EncryptionKey == "" {
		return nil, errors.New("KVTXT_ENCRYPTION_KEY is required")
	}

	if len(cfg.EncryptionKey) < constant.MinEncryptionKeyLength {
		return nil, fmt.Errorf(
			"KVTXT_ENCRYPTION_KEY must be at least %d characters",
			constant.MinEncryptionKeyLength,
		)
	}

	if cfg.MaxPayloadSize < constant.MinMaxPayloadSizeMB ||
		cfg.MaxPayloadSize > constant.MaxMaxPayloadSizeMB {
		return nil, fmt.Errorf(
			"KVTXT_MAX_PAYLOAD_SIZE must be between %d and %d MB",
			constant.MinMaxPayloadSizeMB,
			constant.MaxMaxPayloadSizeMB,
		)
	}

	return cfg, nil
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}
