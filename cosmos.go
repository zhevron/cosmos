package cosmos

import (
	"errors"

	"github.com/zhevron/cosmos/api"
	"github.com/zhevron/cosmos/query"
)

var (
	ErrInvalidKey   = errors.New("invalid key")
	ErrNoDocumentID = errors.New("unable to find document ID from struct")
	// TODO: ErrBadRequest?
	ErrUnauthorized = errors.New("unauthorized or invalid key")
	// TODO: ErrForbidden?
	ErrNotFound = errors.New("resource not found")
	ErrTimeout  = errors.New("request timed out")
	ErrConflict = errors.New("resource already exists")
	// TODO: ErrConcurrency?
	ErrDocumentTooLarge    = errors.New("document size exceeds maximum")
	ErrInternalServerError = errors.New("internal server error")
)

type DateTime = api.DateTime
type Document = api.Document
type QueryParameter = api.QueryParameter

func Select(fields ...string) query.Query {
	return query.Select(fields...)
}
