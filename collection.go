package cosmos

import (
	"context"

	"github.com/zhevron/cosmos/api"
)

type Collection struct {
	api.Collection

	database *Database
}

func (c Collection) Get(ctx context.Context, partitionKey string, id string, out interface{}) error {
	headers := map[string]string{
		api.HEADER_PARTITION_KEY: partitionKey,
	}

	return c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, id), out, headers)
}

func (c Collection) Query(ctx context.Context, partitionKey string, query string, params ...QueryParameter) (*DocumentIterator, error) {
	headers := map[string]string{
		api.HEADER_CONTENT_TYPE: "application/query+json",
		api.HEADER_IS_QUERY:     "True",
	}

	if len(partitionKey) == 0 {
		headers[api.HEADER_QUERY_CROSSPARTITION] = "True"
	}

	return &DocumentIterator{}, nil
}

func (c Collection) Database() *Database {
	return c.database
}
