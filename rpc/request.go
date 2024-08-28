package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Request struct {
	Params *json.RawMessage `json:"params,omitempty"`
	Method string           `json:"method"`
	ID     int              `json:"id"`
}

func (r Request) MarshalJSON() ([]byte, error) {
	request := map[string]any{
		"jsonrpc": "2.0",
		"method":  r.Method,
		"id":      r.ID,
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

	// parse ID
	switch rawID := request["id"].(type) {
	case json.Number:
		num, err := rawID.Int64()
		if err != nil {
			return err
		}
		r.ID = int(num)
	default:
		return fmt.Errorf("cannot decode ID type: %T", rawID)
	}

	switch request["params"] {
	case nil:
		r.Params = &jsonNull
	case emptyParams:
		r.Params = nil
	default:
		// parse params
		// we Marshal params because during the decode
		// it may have been decoded we want it to remain as raw json
		// ex: [42, 22] would be an array of ints instead of a string
		// thus, we Marshal back to a string and cast it
		params, err := json.Marshal(request["params"])
		if err != nil {
			return err
		}
		r.Params = (*json.RawMessage)(&params)
	}

	return nil
}
