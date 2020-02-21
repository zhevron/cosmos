package cosmos

import (
	"context"

	"github.com/zhevron/cosmos/api"
	"github.com/zhevron/cosmos/query"
)

type Collection struct {
	api.Collection

	database *Database
}

func (c Collection) Get(ctx context.Context, partitionKey string, id string, out interface{}) error {
	headers := map[string]string{
		api.HEADER_PARTITION_KEY: partitionKey,
	}

	_, err := c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, id), out, headers)
	return err
}

func (c Collection) Query(ctx context.Context, partitionKey string, query string, params ...query.QueryParameter) (*DocumentIterator, error) {
	headers := map[string]string{
		api.HEADER_CONTENT_TYPE: "application/query+json",
		api.HEADER_IS_QUERY:     "True",
	}

	if len(partitionKey) == 0 {
		headers[api.HEADER_QUERY_CROSSPARTITION] = "True"
	} else {
		headers[api.HEADER_PARTITION_KEY] = partitionKey
	}

	queryParams := make([]api.QueryParameter, len(params))
	for i, v := range params {
		queryParams[i] = api.QueryParameter{
			Name:  v.Name,
			Value: v.ValueAsString(),
		}
	}

	apiQuery := api.Query{
		Query:      query,
		Parameters: queryParams,
	}

	var queryResult api.QueryDocumentsResponse
	res, err := c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, ""), apiQuery, &queryResult, headers)
	if err != nil {
		return nil, err
	}

	return newDocumentIterator(ctx, c.database.Client(), res, apiQuery, queryResult), nil
}

func (c Collection) Database() *Database {
	return c.database
}
