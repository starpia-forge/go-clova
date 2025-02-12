package clova

import "context"

const (
	ModelHCX003     = "HCX-003"
	ModelHCXDASH001 = "HCX-DASH-001"
)

type CompletionRequest struct {
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature"` // TODO - is it Double?
	TopK             int       `json:"topK"`
	TopP             float64   `json:"topP"`
	RepeatPenalty    float64   `json:"repeatPenalty"`
	StopBefore       []string  `json:"stopBefore"`
	MaxTokens        int       `json:"maxTokens"`
	IncludeAIFilters bool      `json:"includeAiFilters"`
	Seed             int       `json:"seed"`
}

type CompletionResponse struct {
	Status       string     `json:"status"`
	Result       Result     `json:"result"`
	StopReason   string     `json:"stopReason"`
	InputLength  int        `json:"inputLength"`
	OutputLength int        `json:"outputLength"`
	Seed         int        `json:"seed"`
	AIFilter     []AIFilter `json:"aiFilter"`
}

type Result struct {
	Message Message `json:"message"`
}

type AIFilter struct {
	GroupName string `json:"groupName"`
	Name      string `json:"name"`
	Score     string `json:"score"`
	Result    string `json:"result"`
}

func (c *Client) CreateCompletion(
	ctx context.Context,
	request CompletionRequest,
) (CompletionResponse, error) {
	return CompletionResponse{}, nil
}
