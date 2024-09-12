package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Response struct {
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *Error           `json:"error,omitempty"`
	*ID
	JSONRPC string `json:"jsonrpc"`
}

func (r Response) MarshalJSON() ([]byte, error) {
	request := map[string]any{
		"jsonrpc": "2.0",
	}
	if r.NullID {
		request["id"] = jsonNull
	} else {
		request["id"] = r.ID.ID
	}
	if r.Error != nil {
		request["error"] = r.Error
	}
	return json.Marshal(request)
}

func (r *Response) UnmarshalJSON(msg []byte) error {
	response := make(map[string]any)

	// used to tell apart params: null and params omitted
	emptyResult := &json.RawMessage{}
	response["result"] = emptyResult
	emptyError := &json.RawMessage{}
	response["error"] = emptyError

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()
	if err := decoder.Decode(&response); err != nil {
		return err
	}

	// parse jsonrpc
	rpc, ok := response["jsonrpc"].(string)
	if !ok {
		return fmt.Errorf("missing field jsonrpc")
	}
	if rpc != "2.0" {
		return fmt.Errorf("cannot decode jsonrpc: %s", rpc)
	}

	hasResult := response["result"] != nil && response["result"] != emptyResult
	hasError := response["error"] != nil && response["error"] != emptyError

	if hasResult == hasError {
		return fmt.Errorf("must provide exactly one of fields result or error")
	}

	var err error
	r.ID, err = parseID(response["id"])
	if err != nil {
		return err
	} else if r.NullID && hasResult {
		return fmt.Errorf("missing field id")
	}

	if hasResult {
		// we Marshal params because during the decode
		// it may have been decoded we want it to remain as raw json
		// ex: [42, 22] would be an array of ints instead of a string
		// thus, we Marshal back to a string and cast it
		params, err := json.Marshal(response["result"])
		if err != nil {
			return err
		}
		r.Result = (*json.RawMessage)(&params)
	}

	if hasError {
		respErr, err := json.Marshal(response["error"])
		if err != nil {
			return err
		}

		err = json.Unmarshal(respErr, &r.Error)
		if err != nil {
			return err
		}
	}

	return nil
}
