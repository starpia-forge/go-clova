package clova

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

var (
	idPrefix    = []byte(`id:`)
	eventPrefix = []byte(`event:`)
	headerData  = []byte(`data:`)
	errorPrefix = []byte(`data:{"status":`)
)

type StreamResponse[T streamable] struct {
	ID    string
	Event string
	Data  T
}

type streamable interface {
	ChatCompletionStreamResponse
}

type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool

	reader      *bufio.Reader
	response    *http.Response
	unmarshaler Unmarshaler

	httpHeader
}

func (stream *streamReader[T]) Recv() (response StreamResponse[T], err error) {
	for {
		var lineBytes []byte
		lineBytes, err = stream.processLines()
		if err != nil {
			return
		}
		if bytes.HasPrefix(lineBytes, idPrefix) {
			response.ID = string(bytes.TrimPrefix(lineBytes, idPrefix))
			continue
		}
		if bytes.HasPrefix(lineBytes, eventPrefix) {
			response.Event = string(bytes.TrimPrefix(lineBytes, eventPrefix))
			continue
		}
		err = stream.unmarshaler.Unmarshal(lineBytes, &response.Data)
		if err != nil {
			return
		}
		break
	}
	return
}

func (stream *streamReader[T]) processLines() ([]byte, error) {
	emptyMessageCount := uint(0)

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
		if bytes.HasPrefix(noSpaceLine, idPrefix) {
			// Return ID
			return noSpaceLine, nil
		}
		if bytes.HasPrefix(noSpaceLine, eventPrefix) {
			// Return Event
			return noSpaceLine, nil
		}
		if bytes.HasPrefix(noSpaceLine, headerData) == false {
			emptyMessageCount++
			if emptyMessageCount > stream.emptyMessagesLimit {
				return nil, ErrTooManyEmptyStreamMessages
			}
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
