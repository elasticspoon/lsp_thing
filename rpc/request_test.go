package rpc

import (
	"encoding/json"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	t.Run("decode basic request", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1 }`)

		var request Request
		err := request.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if request.Method != "subtract" {
			t.Fatalf("expected method: subtract, got: %s", request.Method)
		}
		if request.NullID != false {
			t.Fatalf("expected request NullID: false, got: %v", request.NullID)
		}
		if request.ID.ID != 1 {
			t.Fatalf("expect ID: 1 got: %d", request.ID.ID)
		}

		wantParams := `[42,23]`
		if string(*request.Params) != wantParams {
			t.Fatalf("expected params: %s, got: %s", wantParams, *request.Params)
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

	t.Run("deals with null params", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "jsonrpc": "2.0", "method": "subtract", "params": null, "id": 1 }`)

		var request Request
		err := request.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if string(*request.Params) != string(jsonNull) {
			t.Fatalf("expected params: null, got: %s", *request.Params)
		}
	})

	t.Run("returns correct errors", func(t *testing.T) {
		tests := []struct {
			request string
			want    string
		}{
			{request: `{ "method": "subtract", "params": [42, 23], "id": 3  }`, want: "missing field jsonrpc"},
			{request: `{ "jsonrpc": "3.0", "method": "subtract", "params": [42, 23], "id": 3  }`, want: "cannot decode jsonrpc: 3.0"},
			{request: `{ "jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": "3"  }`, want: "cannot decode ID type: string"},
			{request: `{ "jsonrpc": "2.0", "params": [42, 23], "id": 3  }`, want: "missing field method"},
		}

		for _, test := range tests {
			var request Request
			err := request.UnmarshalJSON([]byte(test.request))

			if err == nil {
				t.Fatalf("expected error")
			}
			if err.Error() != test.want {
				t.Errorf("want error: %s, got: %s", test.want, err.Error())
			}
		}
	})
}

func TestEncodeRequest(t *testing.T) {
	params, _ := json.Marshal([]int{21, 23})
	request := Request{
		Params: (*json.RawMessage)(&params),
		Method: "test",
		ID: &ID{
			ID:     3,
			NullID: false,
		},
	}

	str, err := request.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var got Request
	if err = got.UnmarshalJSON(str); err != nil {
		t.Fatalf("error unmarshalling %s: %s ", str, err)
	}

	if got.Method != "test" {
		t.Fatalf("expected method: test, got: %s", got.Method)
	}
	if got.NullID != false {
		t.Fatalf("expected NullID: false, got: %v", got.NullID)
	}
	if got.ID.ID != 3 {
		t.Fatalf("expected ID: 3, got: %d", got.ID.ID)
	}

	wantParams := string(json.RawMessage(`[21,23]`))
	if string(*got.Params) != wantParams {
		t.Fatalf("expected params: %s, got: %s", wantParams, *got.Params)
	}
}
