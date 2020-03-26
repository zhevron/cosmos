package cosmos

import (
	"context"
	"strings"

	"github.com/zhevron/cosmos/api"
)

type Collection struct {
	api.Collection

	database *Database
}

func (c Collection) ListDocuments(ctx context.Context) (*DocumentIterator, error) {
	var listResult api.ListDocumentsResponse
	res, err := c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, ""), &listResult, nil)
	if err != nil {
		return nil, err
	}

	return newDocumentIterator(ctx, c.database.Client(), res, nil, listResult), nil
}

func (c Collection) GetDocument(ctx context.Context, partitionKey interface{}, id string, out interface{}) error {
	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err := c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, id), out, headers)
	return err
}

func (c Collection) CreateDocument(ctx context.Context, partitionKey interface{}, document interface{}, upsert bool) error {
	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}
	if upsert {
		headers[api.HEADER_IS_UPSERT] = "True"
	}

	_, err := c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, ""), document, nil, headers)
	return err
}

func (c Collection) ReplaceDocument(ctx context.Context, partitionKey interface{}, document interface{}) error {
	id, err := DocumentID(document)
	if err != nil {
		return err
	}

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err = c.database.Client().put(ctx, createDocumentLink(c.database.ID, c.ID, id), document, nil, headers)
	return err
}

func (c Collection) DeleteDocument(ctx context.Context, partitionKey interface{}, document interface{}) error {
	id, err := DocumentID(document)
	if err != nil {
		return err
	}

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err = c.database.Client().delete(ctx, createDocumentLink(c.database.ID, c.ID, id), headers)
	return err
}

func (c Collection) QueryDocuments(ctx context.Context, partitionKey interface{}, query string, params ...api.QueryParameter) (*DocumentIterator, error) {
	headers := map[string]string{
		api.HEADER_CONTENT_TYPE: "application/query+json",
		api.HEADER_IS_QUERY:     "True",
	}

	if partitionKey == nil {
		headers[api.HEADER_QUERY_CROSSPARTITION] = "True"
	} else {
		headers[api.HEADER_PARTITION_KEY] = makePartitionKeyHeaderValue(partitionKey)
	}

	queryParams := []api.QueryParameter{}
	for _, p := range params {
		if strings.Contains(query, p.Name) {
			queryParams = append(queryParams, p)
		}
	}

	apiQuery := &api.Query{
		Query:      query,
		Parameters: queryParams,
	}

	var queryResult api.ListDocumentsResponse
	res, err := c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, ""), apiQuery, &queryResult, headers)
	if err != nil {
		return nil, err
	}

	return newDocumentIterator(ctx, c.database.Client(), res, apiQuery, queryResult), nil
}

func (c Collection) Database() *Database {
	return c.database
}
