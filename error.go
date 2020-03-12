package cosmos

import (
	"fmt"
)

type ErrorCode int

const (
	ErrNoDocumentID        ErrorCode = 0
	ErrInvalidKey          ErrorCode = 1
	ErrTimeout             ErrorCode = 2
	ErrUnauthorized        ErrorCode = 3
	ErrForbidden           ErrorCode = 4
	ErrBadRequest          ErrorCode = 5
	ErrNotFound            ErrorCode = 6
	ErrConflict            ErrorCode = 7
	ErrConcurrency         ErrorCode = 8
	ErrDocumentTooLarge    ErrorCode = 9
	ErrInternalServerError ErrorCode = 10
)

type CosmosError struct {
	Code    ErrorCode
	Message string
}

func (e *CosmosError) Error() string {
	return fmt.Sprintf("cosmosdb error: code=%d message=%s", e.Code, e.Message)
}
