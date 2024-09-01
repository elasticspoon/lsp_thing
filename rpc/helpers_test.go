package rpc

import (
	"encoding/json"
	"testing"
)

func TestParseOptionJson(t *testing.T) {
	t.Run("returns nil for omitted json", func(t *testing.T) {
		emptyJson := &json.RawMessage{}

		result, err := parseOptionalJson(emptyJson, emptyJson)
		if err != nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Errorf("wanted: nil, got: %s", result)
		}
	})
	t.Run("returns jsonNull for null value", func(t *testing.T) {
		message := json.RawMessage("null")
		emptyJson := &json.RawMessage{}

		result, err := parseOptionalJson(message, emptyJson)
		if err != nil {
			t.Fatal(err)
		}

		if string(*result) != string(jsonNull) {
			t.Errorf("wanted: %s, got: %s", *result, jsonNull)
		}
	})
	t.Run("returns json if given anything else", func(t *testing.T) {
		message := []int{1, 2}
		want := `[1,2]`
		emptyJson := &json.RawMessage{}

		result, err := parseOptionalJson(message, emptyJson)
		if err != nil {
			t.Fatal(err)
		}

		if string(*result) != want {
			t.Errorf("wanted: %s, got: %s", want, *result)
		}
	})
}

func TestParseId(t *testing.T) {
	t.Run("pasrses json.Number", func(t *testing.T) {
		result, err := parseID(json.Number("11"))
		if err != nil {
			t.Fatal(err)
		}

		if *result != 11 {
			t.Errorf("wanted: 11, got: %d", *result)
		}
	})

	t.Run("returns error otherwise", func(t *testing.T) {
		_, err := parseID(11)

		wantError := "cannot decode ID type: int"

		if err == nil {
			t.Fatalf("wanted error: %s, got: nil", wantError)
		}

		if err.Error() != wantError {
			t.Errorf("expected error: %s, got: %s", wantError, err.Error())
		}
	})

	t.Run("parses nil", func(t *testing.T) {
		result, err := parseID(nil)
		if err != nil {
			t.Fatal(err)
		}

		if result != nil {
			t.Errorf("wanted: nil, got: %d", result)
		}
	})
}
