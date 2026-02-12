package api

type ErrorCode string

const (
	ErrBadRequest       ErrorCode = "BAD_REQUEST"
	ErrInvalidJSON      ErrorCode = "INVALID_JSON"
	ErrPayloadTooLarge  ErrorCode = "PAYLOAD_TOO_LARGE"
	ErrMethodNotAllowed ErrorCode = "BAD_REQUEST"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrConflict         ErrorCode = "CONFLICT"
	ErrInternal         ErrorCode = "INTERNAL_ERROR"
)
