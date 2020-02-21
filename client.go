package cosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zhevron/cosmos/api"
)

const (
	apiVersion     = "2018-12-31"
	httpRetryAfter = 449
)

type Client struct {
	MaxRetries int
	client     *http.Client
	endpoint   *url.URL
	key        Key
}

func Dial(ctx context.Context, endpoint string, key string) (*Client, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	k, err := ParseKey(key)
	if err != nil {
		return nil, err
	}

	client := &Client{
		MaxRetries: 5,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoint: url,
		key:      k,
	}

	return client, nil
}

func (c Client) Database(ctx context.Context, id string) (*Database, error) {
	var db api.Database
	if _, err := c.get(ctx, createDatabaseLink(id), &db, nil); err != nil {
		return nil, err
	}

	return &Database{
		Database: db,
		client:   &c,
	}, nil
}

func (c Client) ListDatabases(ctx context.Context) ([]*Database, error) {
	var res api.ListDatabasesResponse
	if _, err := c.get(ctx, createDatabaseLink(""), &res, nil); err != nil {
		return nil, err
	}

	databases := make([]*Database, len(res.Databases))
	for i, d := range res.Databases {
		databases[i] = &Database{
			Database: d,
			client:   &c,
		}
	}

	return databases, nil
}

// TODO: Client.CreateDatabase
// TODO: Client.DeleteDatabase

func (c Client) get(ctx context.Context, link string, out interface{}, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodGet, link, nil, out, headers)
}

func (c Client) post(ctx context.Context, link string, body interface{}, out interface{}, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodPost, link, body, out, headers)
}

func (c Client) put(ctx context.Context, link string, body interface{}, out interface{}, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodPut, link, body, out, headers)
}

func (c Client) delete(ctx context.Context, link string, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodDelete, link, nil, nil, headers)
}

func (c Client) request(ctx context.Context, method string, link string, body interface{}, out interface{}, headers map[string]string) (*http.Response, error) {
	uri, _ := url.Parse(c.endpoint.String())
	uri.Path = link

	var reader io.Reader
	if body != nil {
		bodyJSON, err := serialize(body)
		if err != nil {
			return nil, err
		}

		reader = bytes.NewBuffer(bodyJSON)
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), reader)
	if err != nil {
		return nil, err
	}

	applyDefaultHeaders(req)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	signRequest(c.key, req)
	return doRequest(c.client, req, out, 0, c.MaxRetries)
}

func applyDefaultHeaders(req *http.Request) {
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		req.Header.Set(api.HEADER_CONTENT_TYPE, "application/json")
	}

	req.Header.Set(api.HEADER_DATE, time.Now().UTC().Format(api.TIME_FORMAT))
	req.Header.Set(api.HEADER_MAX_ITEM_COUNT, "-1")
	req.Header.Set(api.HEADER_VERSION, apiVersion)
}

func signRequest(key Key, req *http.Request) {
	date := req.Header.Get(api.HEADER_DATE)
	resourceType, resourceID := resourceTypeFromLink(req.URL.Path)

	payload := strings.ToLower(req.Method) + "\n" +
		strings.ToLower(resourceType) + "\n" +
		resourceID + "\n" +
		strings.ToLower(date) + "\n\n"
	signedPayload := key.Sign([]byte(payload))

	tokenType := "master"
	tokenVersion := "1.0"
	header := "type=" + tokenType + "&ver=" + tokenVersion + "&sig=" + signedPayload
	req.Header.Add("Authorization", url.QueryEscape(header))
}

func resourceTypeFromLink(uri string) (string, string) {
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}

	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}

	parts := strings.Split(uri, "/")
	partsLen := len(parts)

	if partsLen%2 == 0 {
		return parts[partsLen-3], strings.Join(parts[1:partsLen-1], "/")
	} else {
		return parts[partsLen-2], strings.Join(parts[1:partsLen-2], "/")
	}
}

func doRequest(client *http.Client, req *http.Request, out interface{}, currentAttempt int, maxRetries int) (*http.Response, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if res.ContentLength == 0 || out == nil {
			return res, nil
		}

		return res, json.NewDecoder(res.Body).Decode(out)
	case http.StatusNoContent:
		return res, nil
	}

	if shouldRetry(res.StatusCode) && currentAttempt < maxRetries {
		time.Sleep(100 * time.Millisecond)
		return doRequest(client, req, out, currentAttempt+1, maxRetries)
	}

	return res, errorFromStatusCode(res.StatusCode)
}

func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == httpRetryAfter
}

func errorFromStatusCode(statusCode int) error {
	errorMap := map[int]error{
		http.StatusBadRequest:            errors.New("bad request"), // TODO
		http.StatusUnauthorized:          ErrUnauthorized,
		http.StatusForbidden:             errors.New("forbidden"), // TODO
		http.StatusNotFound:              ErrNotFound,
		http.StatusRequestTimeout:        ErrTimeout,
		http.StatusConflict:              ErrConflict,
		http.StatusPreconditionFailed:    errors.New("precondition failed"), // TODO
		http.StatusRequestEntityTooLarge: ErrDocumentTooLarge,
	}

	if err, found := errorMap[statusCode]; found {
		return err
	}

	return ErrInternalServerError
}

func serialize(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil

	case string:
		return []byte(v), nil

	default:
		return json.Marshal(v)
	}
}
