// APIError represents a structured HTTP error response.
// It standardizes error handling across the entire API.
//
// Design decision:
// - Errors are returned, not written directly
// - Centralized error-to-response mapping

package api

// APIError contains HTTP status, internal error code,
// and a human-readable message.
type APIError struct {
	Status  int
	Code    ErrorCode
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
