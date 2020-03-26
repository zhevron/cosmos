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

func IsNoDocumentID(err error) bool {
	return isErrorCode(err, ErrNoDocumentID)
}

func IsInvalidKey(err error) bool {
	return isErrorCode(err, ErrInvalidKey)
}

func IsTimeout(err error) bool {
	return isErrorCode(err, ErrTimeout)
}

func InUnauthorized(err error) bool {
	return isErrorCode(err, ErrUnauthorized)
}

func IsForbidden(err error) bool {
	return isErrorCode(err, ErrForbidden)
}

func IsBadRequest(err error) bool {
	return isErrorCode(err, ErrBadRequest)
}

func IsNotFound(err error) bool {
	return isErrorCode(err, ErrNotFound)
}

func IsConflict(err error) bool {
	return isErrorCode(err, ErrConflict)
}

func IsConcurrency(err error) bool {
	return isErrorCode(err, ErrConcurrency)
}

func IsDocumentTooLarge(err error) bool {
	return isErrorCode(err, ErrDocumentTooLarge)
}

func IsInternalServerError(err error) bool {
	return isErrorCode(err, ErrInternalServerError)
}

func isErrorCode(err error, code ErrorCode) bool {
	if cerr, ok := err.(*CosmosError); ok {
		return cerr.Code == code
	}
	return false
}
