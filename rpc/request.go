package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Request struct {
	Params *json.RawMessage `json:"params,omitempty"`
	Method string           `json:"method"`
	*ID
}

func (r Request) MarshalJSON() ([]byte, error) {
	request := map[string]any{
		"jsonrpc": "2.0",
		"method":  r.Method,
	}
	if r.NullID {
		request["id"] = jsonNull
	} else {
		request["id"] = r.ID.ID
	}
	if r.Params != nil {
		request["params"] = r.Params
	}
	return json.Marshal(request)
}

func (r *Request) UnmarshalJSON(msg []byte) error {
	request := make(map[string]any)

	// used to tell apart params: null and params omitted
	emptyParams := &json.RawMessage{}
	request["params"] = emptyParams

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()

	if err := decoder.Decode(&request); err != nil {
		return err
	}

	// parse method
	var ok bool
	if r.Method, ok = request["method"].(string); !ok {
		return fmt.Errorf("missing field method")
	}

	// parse jsonrpc
	rpc, ok := request["jsonrpc"].(string)
	if !ok {
		return fmt.Errorf("missing field jsonrpc")
	}
	if rpc != "2.0" {
		return fmt.Errorf("cannot decode jsonrpc: %s", rpc)
	}

	var err error
	if r.ID, err = parseID(request["id"]); err != nil {
		return err
	}

	params, err := parseOptionalJson(request["params"], emptyParams)
	if err != nil {
		return err
	}
	r.Params = params

	return nil
}
