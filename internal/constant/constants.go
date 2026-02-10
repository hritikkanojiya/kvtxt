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
	MinEncryptionKeyLength  = 16
)
