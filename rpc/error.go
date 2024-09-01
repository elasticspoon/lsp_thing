package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Error struct {
	Data    *json.RawMessage `json:"data,omitempty"`
	Message string           `json:"message"`
	Code    int              `json:"code"`
}

func (r *Error) UnmarshalJSON(msg []byte) error {
	respErr := make(map[string]any)

	// used to tell apart params: null and params omitted
	emptyData := &json.RawMessage{}
	respErr["data"] = emptyData

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()
	if err := decoder.Decode(&respErr); err != nil {
		return err
	}

	switch rawCode := respErr["code"].(type) {
	case json.Number:
		code, err := rawCode.Int64()
		if err != nil {
			return err
		}
		r.Code = int(code)
	case nil:
		return fmt.Errorf("missing field code")
	default:
		return fmt.Errorf("cannot decode ID type: %T", rawCode)
	}

	var ok bool
	if r.Message, ok = respErr["message"].(string); !ok {
		return fmt.Errorf("missing field message")
	}

	data, err := parseOptionalJson(respErr["data"], emptyData)
	if err != nil {
		return err
	}
	r.Data = data

	return nil
}
