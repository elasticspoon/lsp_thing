package rpc

import (
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

		if *response.ID != 1 {
			t.Fatalf("expect ID: 1 got: %d", response.ID)
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

		if response.ID != nil {
			t.Fatalf("expect ID: nil got: %d", response.ID)
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
