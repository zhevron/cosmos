package cosmos

import (
	"github.com/zhevron/cosmos/api"
	"github.com/zhevron/cosmos/query"
)

type DateTime = api.DateTime
type Document = api.Document
type QueryParameter = api.QueryParameter

func Select(fields ...string) query.Query {
	return query.Select(fields...)
}
