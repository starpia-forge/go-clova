package clova

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

var (
	headerData  = []byte("data:")
	errorPrefix = []byte(`data:{"status":`)
	idPrefix    = []byte(`id: `)
)

type streamable interface {
	ChatCompletionStreamResponse
}

type streamReader[T streamable] struct {
	isFinished bool

	reader      *bufio.Reader
	response    *http.Response
	unmarshaler Unmarshaler

	httpHeader
}

func (stream *streamReader[T]) Recv() (response T, err error) {
	rawLine, err := stream.RecvRaw()
	if err != nil {
		return
	}

	err = stream.unmarshaler.Unmarshal(rawLine, &response)
	if err != nil {
		return
	}
	return response, nil
}

func (stream *streamReader[T]) RecvRaw() ([]byte, error) {
	if stream.isFinished {
		return nil, io.EOF
	}

	return stream.processLines()
}

func (stream *streamReader[T]) processLines() ([]byte, error) {
	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil {
			return nil, readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			errResponse := ErrorResponse{}
			if err := stream.unmarshaler.Unmarshal(noSpaceLine, &errResponse); err != nil {
				return nil, err
			}
			return nil, errResponse.ErrStatus
		}
		if bytes.HasPrefix(noSpaceLine, headerData) == false {
			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == `{"data":"[DONE]"}` {
			stream.isFinished = true
			return nil, io.EOF
		}

		return noPrefixLine, nil
	}
}

func (stream *streamReader[T]) Close() error {
	return stream.response.Body.Close()
}
