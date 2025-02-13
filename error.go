package clova

import "fmt"

type ErrorResponse struct {
	ErrStatus *ErrorResponseStatus `json:"status"`
}

type ErrorResponseStatus struct {
	Status

	HTTPStatus     string `json:"-"`
	HTTPStatusCode int    `json:"-"`
}

func (e *ErrorResponseStatus) Error() string {
	return fmt.Sprintf("status code: %d, error code: %s, message: %s", e.HTTPStatusCode, e.Code, e.Message)
}
