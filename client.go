package clova

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	config ClientConfig
}

type Response interface {
	SetHeader(http.Header)
}

type httpHeader http.Header

func (h *httpHeader) SetHeader(header http.Header) {
	*h = httpHeader(header)
}

func (h *httpHeader) Header() http.Header {
	return http.Header(*h)
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

func withSetHeader(key string, value string) requestOption {
	return func(args *requestOptions) {
		args.header.Set(key, value)
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

	request, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Replace request headers with provided headers (if any)
	if args.header != nil {
		request.Header = args.header
	}

	c.setCommonHeaders(request)

	return request, nil
}

func (c *Client) sendRequest(request *http.Request, v Response) error {
	response, err := c.config.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if v != nil {
		v.SetHeader(response.Header)
	}

	if isFailureStatusCode(response) {
		return c.handleErrorResp(response)
	}

	return decodeResponse(response.Body, v)
}

func sendRequestStream[T streamable](client *Client, request *http.Request) (*streamReader[T], error) {
	response, err := client.config.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	if isFailureStatusCode(response) {
		return nil, client.handleErrorResp(response)
	}

	return &streamReader[T]{
		isFinished:         false,
		emptyMessagesLimit: client.config.EmptyMessagesLimit,
		reader:             bufio.NewReader(response.Body),
		response:           response,
		unmarshaler:        &JSONUnmarshaler{},
		httpHeader:         httpHeader(response.Header),
	}, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	if c.config.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.apiKey))
	}
}

func (c *Client) handleErrorResp(response *http.Response) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error, reading response body: %w", err)
	}
	var errRes ErrorResponse
	if err = json.Unmarshal(body, &errRes); err != nil {
		return err
	}

	errRes.ErrStatus.HTTPStatus = response.Status
	errRes.ErrStatus.HTTPStatusCode = response.StatusCode
	return errRes.ErrStatus
}

func isFailureStatusCode(response *http.Response) bool {
	return response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}

type fullURLOptions struct {
	apiVersion string
}

type fullURLOption func(*fullURLOptions)

func withFullURLAPIVersion(apiVersion string) fullURLOption {
	return func(args *fullURLOptions) {
		args.apiVersion = apiVersion
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

	if args.apiVersion != "" {
		baseURL = fmt.Sprintf("%s/%s", baseURL, args.apiVersion)
	}

	return fmt.Sprintf("%s%s", baseURL, suffix)
}
