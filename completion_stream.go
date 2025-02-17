package clova

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	CompletionMessageStreamEventToken  = "token"
	CompletionMessageStreamEventResult = "result"
	CompletionMessageStreamEventSignal = "signal"
)

type ChatCompletionStreamResponse struct {
	Message      CompletionMessage `json:"message"`
	Index        int               `json:"index"`
	InputLength  int               `json:"inputLength"`
	OutputLength int               `json:"outputLength"`
	StopReason   string            `json:"stopReason"`
}

type ChatCompletionStream struct {
	*streamReader[ChatCompletionStreamResponse]
}

func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	model string,
	request CompletionRequest,
) (*ChatCompletionStream, error) {
	if model == "" {
		return nil, errors.New("model cannot be empty")
	}

	suffix := fmt.Sprintf("/chat-completions/%s", model)

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(suffix, withFullURLAPIVersion("v1")),
		withBody(request),
		withSetHeader("Content-Type", "application/json"),
		withSetHeader("Accept", "text/event-stream"),
		withSetHeader("Cache-Control", "no-cache"),
		withSetHeader("Connection", "keep-alive"),
	)
	if err != nil {
		return nil, err
	}

	// Send Request Stream
	response, err := sendRequestStream[ChatCompletionStreamResponse](c, req)
	if err != nil {
		return nil, err
	}
	stream := &ChatCompletionStream{
		streamReader: response,
	}

	return stream, nil
}
