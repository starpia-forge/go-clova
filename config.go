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
	authToken string

	BaseURL    string
	HTTPClient HTTPDoer

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   clovaAPIURL,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
