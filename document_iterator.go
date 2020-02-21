package cosmos

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zhevron/cosmos/api"
)

type DocumentIterator struct {
	ctx               context.Context
	client            *Client
	headers           map[string]string
	path              string
	query             api.Query
	continuationTokan string
	documents         []json.RawMessage
	total             int64
	current           int64
	err               error
}

func newDocumentIterator(ctx context.Context, client *Client, res *http.Response, query api.Query, queryResult api.QueryDocumentsResponse) *DocumentIterator {
	headers := map[string]string{}
	for k := range res.Request.Header {
		headers[k] = res.Request.Header.Get(k)
	}
	delete(headers, "Authorization")

	return &DocumentIterator{
		ctx:               ctx,
		client:            client,
		headers:           headers,
		path:              res.Request.URL.Path,
		query:             query,
		continuationTokan: res.Header.Get(api.HEADER_CONTINUATION),
		documents:         queryResult.Documents,
		total:             queryResult.Count,
		current:           0,
	}
}

func (it *DocumentIterator) Next(out interface{}) bool {
	if it.err != nil {
		return false
	}

	if it.current < int64(len(it.documents)) {
		it.err = json.Unmarshal(it.documents[it.current], out)
		it.current++
		return true
	} else if it.current < it.total {
		it.err = it.fetchNext()
		return it.Next(out)
	}

	return false
}

func (it *DocumentIterator) Reset() {
	it.current = 0
	it.err = nil
}

func (it *DocumentIterator) Count() int64 {
	return it.total
}

func (it *DocumentIterator) Err() error {
	return it.err
}

func (it *DocumentIterator) fetchNext() error {
	if it.continuationTokan == "" {
		return nil
	}
	it.headers[api.HEADER_CONTINUATION] = it.continuationTokan

	var queryResult api.QueryDocumentsResponse
	res, err := it.client.post(it.ctx, it.path, it.query, &queryResult, it.headers)
	if err != nil {
		return err
	}

	it.continuationTokan = res.Header.Get(api.HEADER_CONTINUATION)
	it.documents = append(it.documents, queryResult.Documents...)

	return nil
}
