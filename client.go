package clova

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	body       any
	header     http.Header
	marshaller Marshaller
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

func withMarshaller(marshaller Marshaller) requestOption {
	return func(args *requestOptions) {
		args.marshaller = marshaller
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method, url string,
	opts ...requestOption,
) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:       nil,
		header:     make(http.Header),
		marshaller: &JSONMarshaller{},
	}

	for _, opt := range opts {
		opt(args)
	}

	// Create Body Reader
	var bodyReader io.Reader
	if args.body != nil {
		b, ok := args.body.(io.Reader)
		if ok {
			bodyReader = b
		} else {
			bodyBytes, err := args.marshaller.Marshal(args.body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewBuffer(bodyBytes)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Replace request headers with provided headers (if any)
	if args.header != nil {
		req.Header = args.header
	}

	return req, nil
}

type fullURLOptions struct {
	apiVersion string
	model      string
	taskID     string
}

type fullURLOption func(*fullURLOptions)

func withAPIVersion(apiVersion string) fullURLOption {
	return func(args *fullURLOptions) {
		args.apiVersion = apiVersion
	}
}

func withModel(model string) fullURLOption {
	return func(args *fullURLOptions) {
		args.model = model
	}
}

func withTaskID(taskID string) fullURLOption {
	return func(args *fullURLOptions) {
		args.taskID = taskID
	}
}

func (c *Client) fullURL(suffix string, opts ...fullURLOption) string {
	baseURL := strings.TrimRight(c.config.BaseURL, "/")
	args := fullURLOptions{
		apiVersion: "v1",
	}

	for _, opt := range opts {
		opt(&args)
	}

	if len(args.apiVersion) > 0 {
		baseURL = fmt.Sprintf("%s/%s", baseURL, args.apiVersion)
	}

	if len(args.taskID) > 0 {
		baseURL = fmt.Sprintf("%s/tasks/%s", baseURL, args.taskID)
	}

	if len(args.model) > 0 {
		suffix = fmt.Sprintf("%s/%s", suffix, args.model)
	}

	return fmt.Sprintf("%s%s", baseURL, suffix)
}
