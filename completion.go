package clova

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	ModelHCX003     = "HCX-003"
	ModelHCXDASH001 = "HCX-DASH-001"
)

type CompletionMessageRole string

const (
	CompletionMessageRoleSystem    CompletionMessageRole = "system"
	CompletionMessageRoleUser      CompletionMessageRole = "user"
	CompletionMessageRoleAssistant CompletionMessageRole = "assistant"
)

type CompletionMessage struct {
	Role    CompletionMessageRole `json:"role"`
	Content string                `json:"content"`
}

type CompletionRequest struct {
	Messages         []CompletionMessage `json:"messages"`
	Temperature      float64             `json:"temperature"` // TODO - is it Double?
	TopK             int                 `json:"topK"`
	TopP             float64             `json:"topP"`
	RepeatPenalty    float64             `json:"repeatPenalty"`
	StopBefore       []string            `json:"stopBefore"`
	MaxTokens        int                 `json:"maxTokens"`
	IncludeAIFilters bool                `json:"includeAiFilters"`
	Seed             int                 `json:"seed"`
}

type CompletionResponse struct {
	Status       string     `json:"status"`
	Result       Result     `json:"result"`
	StopReason   string     `json:"stopReason"`
	InputLength  int        `json:"inputLength"`
	OutputLength int        `json:"outputLength"`
	Seed         int        `json:"seed"`
	AIFilter     []AIFilter `json:"aiFilter"`

	httpHeader
}

type Result struct {
	Message CompletionMessage `json:"message"`
}

type AIFilter struct {
	GroupName string `json:"groupName"`
	Name      string `json:"name"`
	Score     string `json:"score"`
	Result    string `json:"result"`
}

func (c *Client) CreateCompletion(
	ctx context.Context,
	model string,
	request CompletionRequest,
) (CompletionResponse, error) {
	if model == "" {
		return CompletionResponse{}, errors.New("model cannot be empty")
	}

	suffix := fmt.Sprintf("/chat-completions/%s", model)

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(suffix, withFullURLAPIVersion("v1")),
		withBody(request),
		withContentType("application/json"),
		withAccept("application/json"),
	)
	if err != nil {
		return CompletionResponse{}, err
	}

	// Send Request
	res := CompletionResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return CompletionResponse{}, err
	}

	return res, nil
}
