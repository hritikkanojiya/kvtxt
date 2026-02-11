package api

type APIError struct {
	Status  int
	Code    ErrorCode
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
