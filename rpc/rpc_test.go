package rpc

import "testing"

func TestDecodeMessage(t *testing.T) {
	t.Run("Decodes Empty Line", func(t *testing.T) {
		want := ""
		got := DecodeMessage([]byte("Content-Length: 3\r\n\r\n\"\""))

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
	t.Run("Encodes Basic Struct", func(t *testing.T) {
		testStruct := struct {
			Testing bool
		}{Testing: true}
		want := testStruct
		got := DecodeMessage([]byte("Content-Length: 16\r\n\r\n{\"Testing\":true}"))

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
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
