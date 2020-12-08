package cosmos

import (
	"net/url"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

type DialOption func(*Client) error

func WithAccountName(accountName string) DialOption {
	return func(c *Client) error {
		endpoint, err := url.Parse("https://" + accountName + ".documents.azure.com:443/")
		if err != nil {
			return err
		}

		c.endpoint = endpoint
		return nil
	}
}

func WithConnectionString(connectionString string) DialOption {
	return func(c *Client) error {
		endpoint, key, err := parseConnectionString(connectionString)
		if err != nil {
			return err
		}

		c.endpoint = endpoint
		c.key = key
		return nil
	}
}

func WithKey(key string) DialOption {
	return func(c *Client) error {
		k, err := ParseKey(key)
		if err != nil {
			return err
		}

		c.key = k
		return nil
	}
}

func WithEndpoint(endpoint *url.URL) DialOption {
	return func(c *Client) error {
		c.endpoint = endpoint
		return nil
	}
}

func WithRetries(retries int) DialOption {
	return func(c *Client) error {
		c.MaxRetries = retries
		return nil
	}
}

func WithRetryForStatusCode(statusCodes ...int) DialOption {
	return func(c *Client) error {
		if c.retryOnStatus == nil {
			c.retryOnStatus = []int{}
		}
		c.retryOnStatus = append(c.retryOnStatus, statusCodes...)
		return nil
	}
}

func WithTimeout(timeout time.Duration) DialOption {
	return func(c *Client) error {
		c.client.Timeout = timeout
		return nil
	}
}

func WithTracer(tracerInstance opentracing.Tracer) DialOption {
	return func(c *Client) error {
		c.tracer = tracerInstance
		return nil
	}
}

func parseConnectionString(connectionString string) (endpoint *url.URL, key Key, err error) {
	parts := strings.Split(connectionString, ";")
	for _, part := range parts {
		pair := strings.Split(part, "=")
		if len(pair) < 2 {
			continue
		}

		name := strings.ToUpper(pair[0])
		value := strings.Join(pair[1:], "=")

		switch name {
		case "ACCOUNTENDPOINT":
			endpoint, err = url.Parse(value)
			if err != nil {
				return nil, nil, err
			}

		case "ACCOUNTKEY":
			key, err = ParseKey(value)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return endpoint, key, nil
}
