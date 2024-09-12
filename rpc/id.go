package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type ID struct {
	ID     int `json:"id"`
	NullID bool
}

func (r ID) MarshalJSON() ([]byte, error) {
	request := map[string]any{}
	if r.NullID {
		request["id"] = jsonNull
	} else {
		request["id"] = r.ID
	}
	return json.Marshal(request)
}

func (r *ID) UnmarshalJSON(msg []byte) error {
	id := make(map[string]any)

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()

	if err := decoder.Decode(&id); err != nil {
		return err
	}

	switch rawID := id["id"].(type) {
	case string:
		num, err := strconv.Atoi(rawID)
		if err != nil {
			return nil
		}
		r.ID, r.NullID = int(num), false
	case json.Number:
		num, err := rawID.Int64()
		if err != nil {
			return nil
		}
		r.ID, r.NullID = int(num), false
		return nil
	case nil:
		r.NullID = true
	default:
		return fmt.Errorf("cannot decode ID type: %T", rawID)
	}
	return nil
}

type (
	key int
)

var idKey key

func NewContext(ctx context.Context, id *ID) context.Context {
	return context.WithValue(ctx, idKey, id)
}

func FromContext(ctx context.Context) (*ID, bool) {
	id, ok := ctx.Value(idKey).(*ID)
	return id, ok
}
