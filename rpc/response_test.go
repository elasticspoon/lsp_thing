package rpc

import (
	"encoding/json"
	"testing"
)

func TestUnmarshallResponse(t *testing.T) {
	t.Run("unmarshall basic success response", func(t *testing.T) {
		incomingJsonRpcResponse := []byte(`{ "jsonrpc": "2.0", "result": [42, 23], "id": 1 }`)

		var response Response
		err := response.UnmarshalJSON(incomingJsonRpcResponse)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if response.NullID != false {
			t.Fatalf("expect NullID: false, got: %v", response.NullID)
		}
		if response.ID.ID != 1 {
			t.Fatalf("expect ID: 1 got: %d", response.ID.ID)
		}

		wantResult := `[42,23]`
		if string(*response.Result) != wantResult {
			t.Fatalf("expected result: %s, got: %s", wantResult, *response.Result)
		}
	})

	t.Run("unmarshall basic error response", func(t *testing.T) {
		incomingJsonRpcResponse := []byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error", "data": [21, 34]}, "id": null}`)

		var response Response
		err := response.UnmarshalJSON(incomingJsonRpcResponse)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if response.NullID != true {
			t.Fatalf("expect NullID: true, got: %v", response.NullID)
		}

		if response.Result != nil {
			t.Fatalf("expected result: nil, got: %s", response.Result)
		}

		if response.Error == nil {
			t.Fatalf("expected response to contain error")
		}

		wantErrMessage := "Parse error"
		if response.Error.Message != wantErrMessage {
			t.Fatalf("expect error message: %s got: %s", wantErrMessage, response.Error.Message)
		}
		if response.Error.Code != -32700 {
			t.Fatalf("expect error code: -32700 got: %d", response.Error.Code)
		}
		wantErrData := `[21,34]`
		if response.Error.Data == nil || string(*response.Error.Data) != wantErrData {
			t.Fatalf("expected data: %s, got: %s", wantErrData, *response.Error.Data)
		}
	})

	t.Run("returns correct errors", func(t *testing.T) {
		tests := []struct {
			response string
			want     string
		}{
			{response: `{"error": {}, "id": null}`, want: "missing field jsonrpc"},
			{response: `{"jsonrpc": "2.0", "result": [], "id": null}`, want: "missing field id"},
			{response: `{"jsonrpc": "2.0", "error": {}, "result": [], "id": null}`, want: "must provide exactly one of fields result or error"},
			{response: `{"jsonrpc": "2.0", "id": null}`, want: "must provide exactly one of fields result or error"},
		}

		for _, test := range tests {
			var response Response
			err := response.UnmarshalJSON([]byte(test.response))

			if err == nil {
				t.Fatalf("expected error: %s", test.want)
			}
			if err.Error() != test.want {
				t.Errorf("want error: %s, got: %s", test.want, err.Error())
			}
		}
	})
}

func TestEncodeResponse(t *testing.T) {
	t.Run("encode basic error response", func(t *testing.T) {
		data := json.RawMessage(`{"test": 33}`)
		response := Response{
			Error: &Error{
				Data:    &data,
				Message: "test message",
				Code:    42,
			},
			ID:      &ID{ID: 3, NullID: false},
			JSONRPC: "2.0",
		}

		str, err := response.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		var got Response
		if err = got.UnmarshalJSON(str); err != nil {
			t.Fatalf("error unmarshalling \"%s\": %s", str, err)
		}

		if got.NullID != false {
			t.Fatalf("expected NullID: false, got: %v", got.NullID)
		}
		if got.ID.ID != 3 {
			t.Fatalf("expected ID: 3, got: %d", got.ID.ID)
		}
		if got.Error == nil {
			t.Fatalf("expected Error: got: nil")
		}
		if got.Error.Code != 42 {
			t.Fatalf("expected Error.Code: 42, got: %d", got.Error.Code)
		}
		if got.Error.Message != "test message" {
			t.Fatalf("expected Error.Message: test message, got: %s", got.Error.Message)
		}
		if string(*got.Error.Data) != string(data) {
			t.Fatalf(`expected Error.Data: %v, got: %v`, data, *got.Error.Data)
		}
		if string(*got.Error.Data) != string(data) {
			t.Fatalf(`expected Error.Data: "{"test":33}", got: "%s"`, *got.Error.Data)
		}
	})
}
