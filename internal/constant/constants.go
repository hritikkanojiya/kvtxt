package constant

import "time"

// Server configuration
const (
	DefaultPort     = ":8080"
	ShutdownTimeout = 10 * time.Second
	ReadTimeout     = 10 * time.Second
	WriteTimeout    = 10 * time.Second
	IdleTimeout     = 30 * time.Second
)

// Payload configuration
const (
	DefaultMaxPayloadSizeMB = 50
	MinMaxPayloadSizeMB     = 1
	MaxMaxPayloadSizeMB     = 200
	MB                      = int64(1 << 20)
)

// Cache configuration
const (
	DefaultCacheSize = 1000
	CleanupInterval  = 1800 * time.Second
)

// Time-to-live configuration
const (
	MinTTL     = 1
	DefaultTTL = 5184000 * time.Second
	MaxTTL     = 31536000 * time.Second
)

// Security configuration
const (
	MinEncryptionKeyLength = 16
	Base62Characters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Application metadata
const (
	RequestIdKey = "request_id"
	AppVersion   = "1.0.0"
)
