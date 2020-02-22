package cosmos

import (
	"context"
	"reflect"

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

func (c Collection) CreateDocument(ctx context.Context, document interface{}, upsert bool) error {
	id, err := getDocumentID(document)
	if err != nil {
		return err
	}

	headers := map[string]string{}
	if upsert {
		headers[api.HEADER_IS_UPSERT] = "True"
	}

	_, err = c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, id), document, nil, headers)
	return err
}

func (c Collection) ReplaceDocument(ctx context.Context, partitionKey interface{}, document interface{}) error {
	id, err := getDocumentID(document)
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
	id, err := getDocumentID(document)
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

	if params == nil {
		params = []api.QueryParameter{}
	}

	apiQuery := &api.Query{
		Query:      query,
		Parameters: params,
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

func getDocumentID(document interface{}) (string, error) {
	rv := reflect.ValueOf(document)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", ErrNoDocumentID
	}

	rt := rv.Type()
	numField := rt.NumField()
	for i := 0; i < numField; i++ {
		if rt.Field(i).Tag.Get("json") == "id" {
			if id, ok := rv.Field(i).Interface().(string); ok && id != "" {
				return id, nil
			}

			return "", ErrNoDocumentID
		}
	}

	return "", ErrNoDocumentID
}
