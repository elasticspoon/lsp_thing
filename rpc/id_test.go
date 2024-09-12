package rpc

import (
	"testing"
)

func TestDecodeID(t *testing.T) {
	t.Run("decode basic id string", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "id": "222" }`)

		var id ID
		err := id.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if id.NullID != false {
			t.Fatalf("expect NullID: false got: %v", id.NullID)
		}
		if id.ID != 222 {
			t.Fatalf("expect ID: 22 got: %d", id.ID)
		}
	})

	t.Run("decode basic id num", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "id": 222 }`)

		var id ID
		err := id.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if id.NullID != false {
			t.Fatalf("expect NullID: false got: %v", id.NullID)
		}
		if id.ID != 222 {
			t.Fatalf("expect ID: 22 got: %d", id.ID)
		}
	})

	t.Run("decode null id", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "id": null }`)

		var id ID
		err := id.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if id.NullID != true {
			t.Fatalf("expect NullID: true got: %v", id.NullID)
		}
	})

	t.Run("deals with omitted params", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "jsonrpc": "2.0", "method": "subtract", "id": 1 }`)

		var request Request
		err := request.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if request.Params != nil {
			t.Fatalf("expected params: nil, got: %T", request.Params)
		}
	})
}

func TestEncodeID(t *testing.T) {
	t.Run("encodes id with num", func(t *testing.T) {
		request := ID{ID: 22, NullID: false}

		str, err := request.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		want := `{"id":22}`
		if string(str) != want {
			t.Errorf("wanted %s, got: %s", want, str)
		}
	})
	t.Run("encodes omitted id", func(t *testing.T) {
		request := ID{NullID: true}

		str, err := request.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		want := `{"id":null}`
		if string(str) != want {
			t.Errorf("wanted %s, got: %s", want, str)
		}
	})
}
