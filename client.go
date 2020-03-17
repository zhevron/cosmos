package cosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

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
	cache      *cache.Cache
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
		cache:    cache.New(5*time.Minute, 10*time.Minute),
	}

	return client, nil
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
			cache:    cache.New(5*time.Minute, 10*time.Minute),
		}
		c.cache.Set(d.ID, databases[i], cache.DefaultExpiration)
	}

	return databases, nil
}

func (c Client) GetDatabase(ctx context.Context, id string) (*Database, error) {
	if database, found := c.cache.Get(id); found {
		return database.(*Database), nil
	}

	var db api.Database
	if _, err := c.get(ctx, createDatabaseLink(id), &db, nil); err != nil {
		return nil, err
	}

	database := &Database{
		Database: db,
		client:   &c,
		cache:    cache.New(5*time.Minute, 10*time.Minute),
	}
	c.cache.Set(db.ID, database, cache.DefaultExpiration)

	return database, nil
}

func (c Client) CreateDatabase(ctx context.Context, id string) (*Database, error) {
	req := api.CreateDatabaseRequest{
		ID: id,
	}

	var db api.Database
	_, err := c.post(ctx, createDatabaseLink(""), req, &db, nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		Database: db,
		client:   &c,
		cache:    cache.New(5*time.Minute, 10*time.Minute),
	}, nil
}

func (c Client) DeleteDatabase(ctx context.Context, id string) error {
	_, err := c.delete(ctx, createDatabaseLink(id), nil)
	return err
}

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

	return res, errorFromResponse(res)
}

func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == httpRetryAfter
}

func errorFromResponse(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusBadRequest:
		message, err := errorMessageFromBody(res.Body)
		if err != nil {
			return err
		}
		return &CosmosError{Code: ErrBadRequest, Message: message}

	case http.StatusUnauthorized:
		return &CosmosError{Code: ErrUnauthorized, Message: res.Status} // TODO: Message from response?

	case http.StatusForbidden:
		return &CosmosError{Code: ErrForbidden, Message: res.Status} // TODO: Message from response?

	case http.StatusNotFound:
		return &CosmosError{Code: ErrNotFound, Message: res.Status}

	case http.StatusRequestTimeout:
		return &CosmosError{Code: ErrTimeout, Message: res.Status} // TODO: Message from response?

	case http.StatusConflict:
		return &CosmosError{Code: ErrConflict, Message: res.Status} // TODO: Message from response?

	case http.StatusPreconditionFailed:
		return &CosmosError{Code: ErrConcurrency, Message: res.Status} // TODO: Message from response?

	case http.StatusRequestEntityTooLarge:
		return &CosmosError{Code: ErrDocumentTooLarge, Message: res.Status} // TODO: Message from response?
	}

	return &CosmosError{Code: ErrInternalServerError, Message: "internal server error"}
}

func errorMessageFromBody(bodyReader io.ReadCloser) (string, error) {
	var body struct {
		Code    string
		Message string
	}

	if err := json.NewDecoder(bodyReader).Decode(&body); err != nil {
		return "", err
	}

	var errors struct {
		Errors []struct {
			Severity string
			Code     string
			Message  string
		}
	}

	errorsJSON := strings.TrimSpace(strings.TrimPrefix(strings.Split(strings.Replace(body.Message, "\r\n", "\n", -1), "\n")[0], "Message:"))
	if err := json.Unmarshal([]byte(errorsJSON), &errors); err != nil {
		return "", err
	}

	if len(errors.Errors) == 0 {
		return "", nil
	}

	return errors.Errors[0].Message, nil
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

func makePartitionKeyHeaderValue(partitionKey interface{}) string {
	v, err := json.Marshal([]interface{}{partitionKey})
	if err != nil {
		return ""
	}
	return string(v)
}
