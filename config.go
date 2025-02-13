package clova

import "net/http"

const (
	clovaAPIURL                    = "https://clovastudio.stream.ntruss.com"
	defaultEmptyMessagesLimit uint = 300
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	apiKey string

	BaseURL    string
	HTTPClient HTTPDoer

	EmptyMessagesLimit uint
}

func DefaultConfig(apiKey string) ClientConfig {
	return ClientConfig{
		apiKey:  apiKey,
		BaseURL: clovaAPIURL,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
