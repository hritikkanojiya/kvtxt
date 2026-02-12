package constant

import "time"

const (
	DefaultPort             = ":8080"
	DefaultMaxPayloadSizeMB = 50
	MinMaxPayloadSizeMB     = 1
	MaxMaxPayloadSizeMB     = 200
	DefaultCacheSize        = 1000
	MB                      = int64(1 << 20)
	ShutdownTimeout         = 10 * time.Second
	ReadTimeout             = 10 * time.Second
	WriteTimeout            = 10 * time.Second
	IdleTimeout             = 30 * time.Second
	MinEncryptionKeyLength  = 16
	RequestIdKey            = "request_id"
	MaxTTL                  = 31536000 * time.Second
	MinTTL                  = 1 * time.Second
	DefaultTTL              = 5184000 * time.Second
	CleanupInterval         = 1800 * time.Second
)
