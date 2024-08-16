package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

const HEADER_PREAMBLE_LEN = len("Content-Length: ")

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return "", nil, errors.New("header not found")
	}

	// NOTE: header => Content-Length: <number>
	// we already removed the /r/n/r/n bit
	contentLengthBytes := header[HEADER_PREAMBLE_LEN:]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (int, []byte, error) {
	header, content, found := bytes.Cut(data, []byte("\r\n\r\n"))
	if !found {
		return 0, nil, nil // we are waitng
	}

	contentLengthBytes := header[HEADER_PREAMBLE_LEN:]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}

	// Wait for more
	if len(content) < contentLength {
		return 0, nil, nil
	}

	// header + /r/n/r/n + content
	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
