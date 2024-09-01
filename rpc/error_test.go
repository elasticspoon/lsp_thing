package rpc

import (
	"testing"
)

func TestUnmarshalError(t *testing.T) {
	t.Run("decode basic error", func(t *testing.T) {
		incomingJsonRpcRequest := []byte(`{ "code": 32000, "message": "subtract", "data": [42, 23] }`)

		var respError Error
		err := respError.UnmarshalJSON(incomingJsonRpcRequest)
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		if respError.Message != "subtract" {
			t.Fatalf("expected method: subtract, got: %s", respError.Message)
		}

		if respError.Code != 32000 {
			t.Fatalf("expect ID: 1 got: %d", respError.Code)
		}

		wantData := `[42,23]`
		if string(*respError.Data) != wantData {
			t.Fatalf("expected params: %s, got: %s", wantData, *respError.Data)
		}
	})

	t.Run("returns correct errors", func(t *testing.T) {
		tests := []struct {
			respError string
			want      string
		}{
			{respError: `{ "data": [42, 23], "message": "c"  }`, want: "missing field code"},
			{respError: `{ "code": 32000, "data": [42, 23] }`, want: "missing field message"},
		}

		for _, test := range tests {
			var respError Error
			err := respError.UnmarshalJSON([]byte(test.respError))

			if err == nil {
				t.Fatalf("expected error: %s", test.want)
			}
			if err.Error() != test.want {
				t.Errorf("want error: %s, got: %s", test.want, err.Error())
			}
		}
	})
}
