package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func DecodeMessage(msg []byte) error {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return errors.New("header not found")
	}

	return nil
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}
