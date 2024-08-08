package rpc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	t.Run("Decodes Length of Header", func(t *testing.T) {
		want := BaseMessage{Method: "hello"}
		testMsg, _ := json.Marshal(want)

		got, err := DecodeMessage([]byte(fmt.Sprintf("Content-Length: 16/r/n/r/n%s", testMsg)))
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if got != &want {
			t.Errorf("expected %d bytes, got %d", want, got)
		}
	})

	t.Run("Decodes Length of Header", func(t *testing.T) {
		want := 2
		bytes, err := DecodeMessage([]byte("Content-Length: 2\r\n\r\n\"\""))
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if bytes != nil {
			t.Errorf("expected %d bytes, got %d", want, bytes)
		}
	})

	t.Run("Decodes Empty Line", func(t *testing.T) {
		_, err := DecodeMessage([]byte(""))
		if err == nil {
			t.Errorf("wanted an error, got: %s", err)
		}
	})
}

func TestEncodeMessage(t *testing.T) {
	t.Run("Encodes Empty Line", func(t *testing.T) {
		message := ""
		want := "Content-Length: 2\r\n\r\n\"\""
		got := EncodeMessage(message)

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
	t.Run("Encodes Basic Struct", func(t *testing.T) {
		testStruct := struct {
			Testing bool
		}{Testing: true}
		want := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
		got := EncodeMessage(testStruct)

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
}
