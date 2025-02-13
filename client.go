package clova

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
	config ClientConfig
}

func NewClient(apiKey string) *Client {
	config := DefaultConfig(apiKey)
	return NewClientWithConfig(config)
}

func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func withContentType(contentType string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}

func withAccept(accept string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Accept", accept)
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method, url string,
	opts ...requestOption,
) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}

	for _, opt := range opts {
		opt(args)
	}

	var bodyReader io.Reader
	if args.body != nil {
		bodyReader = args.body.(io.Reader)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if args.header != nil {
		req.Header = args.header
	}

	return req, nil
}
