package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

const HEADER_PREAMBLE_LEN = len("Content-Length: ")

var jsonNull = json.RawMessage("null")

func DecodeMessage(msg []byte) (*Request, error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
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

	var request Request
	if err := json.Unmarshal(content[:contentLength], &request); err != nil {
		return nil, err
	}

	return &request, nil
}
