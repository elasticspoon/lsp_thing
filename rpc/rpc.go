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

func DecodeMessage(msg []byte) (*BaseMessage, error) {
	header, _, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return nil, errors.New("header not found")
	}

	// NOTE: header => Content-Length: <number>
	// we already removed the /r/n/r/n bit
	contentLengthBytes := header[HEADER_PREAMBLE_LEN:]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return nil, err
	}

	_ = contentLength
	return nil, nil
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}
