package rpc

import (
	"encoding/json"
	"fmt"
)

func parseOptionalJson(optJson any, emptyJson *json.RawMessage) (*json.RawMessage, error) {
	switch optJson {
	case nil:
		fmt.Println("nil val")
		return &jsonNull, nil
	case emptyJson:
		return nil, nil
	default:
		// we Marshal optJson because during the decode
		// it may have been decoded we want it to remain as raw json
		// ex: [42, 22] would be an array of ints instead of a string
		// thus, we Marshal back to a string and cast it
		data, err := json.Marshal(optJson)
		if err != nil {
			return nil, err
		}
		return (*json.RawMessage)(&data), nil
	}
}

func parseID(id any) (*int, error) {
	switch rawID := id.(type) {
	case json.Number:
		num, err := rawID.Int64()
		if err != nil {
			return nil, err
		}
		num2 := int(num)
		return &num2, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot decode ID type: %T", rawID)
	}
}
