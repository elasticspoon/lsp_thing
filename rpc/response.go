package rpc

import (
	"bytes"
	"encoding/json"
)

type Response struct {
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
	JSONRPC string           `json:"jsonrpc"`
	ID      int              `json:"id"`
}

type Error struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

func (r Response) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (r *Response) UnmarshalJSON(msg []byte) error {
	response := make(map[string]any)
	//
	// used to tell apart params: null and params omitted
	emptyResult := &json.RawMessage{}
	response["params"] = emptyResult

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()
	if err := decoder.Decode(&response); err != nil {
		return err
	}

	return nil
}
