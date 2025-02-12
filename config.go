package clova

import "net/http"

const (
	clovaAPIURLv1                  = "https://clovastudio.stream.ntruss.com/v1"
	clovaAPIURLv2                  = "https://clovastudio.stream.ntruss.com/v2"
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
		BaseURL:   clovaAPIURLv1,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
