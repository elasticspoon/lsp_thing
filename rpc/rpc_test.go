package rpc

import (
	"slices"
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	t.Run("Decodes the Method", func(t *testing.T) {
		incomingMsg := "Content-Length: 17\r\n\r\n{\"Method\":\"test\"}"
		msgContent := []byte("{\"Method\":\"test\"}")

		method, content, err := DecodeMessage([]byte(incomingMsg))
		contentLength := len(content)
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if contentLength != 17 {
			t.Errorf("expected Content-Length: 17, got: %d", contentLength)
		}

		if slices.Compare(content, msgContent) != 0 {
			t.Errorf("expected %s, got: %s", msgContent, content)
		}

		if method != "test" {
			t.Errorf("expected %s, got %s", "test", method)
		}
	})

	t.Run("Decodes Empty Line", func(t *testing.T) {
		_, _, err := DecodeMessage([]byte(""))
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

func TestSplit(t *testing.T) {
	t.Run("splits to length in header", func(t *testing.T) {
		incomingMsg := "Content-Length: 17\r\n\r\n{\"Method\":\"test\"}extrastuff"

		length, data, err := Split([]byte(incomingMsg), false)
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if length != 39 {
			t.Errorf("expected length of 39, got: %d", length)
		}

		if slices.Compare(data, []byte(incomingMsg)[:39]) != 0 {
			t.Errorf("expected %s, got: %s", []byte(incomingMsg)[:39], data)
		}
	})

	t.Run("Waits (0, nil, nil) if in no input", func(t *testing.T) {
		length, data, err := Split([]byte(""), false)
		if err != nil {
			t.Errorf("did not expect an error, got: %s", err)
		}

		if length != 0 {
			t.Errorf("expected a length of 0, got: %d", length)
		}

		if data != nil {
			t.Errorf("did not expect data, got: %s", data)
		}
	})

	t.Run("Waits (0, nil, nil) if input is shorter than expected length", func(t *testing.T) {
		length, data, err := Split([]byte("Content-Length: 17\r\n\r\n{\""), false)
		if err != nil {
			t.Errorf("did not expect an error, got: %s", err)
		}

		if length != 0 {
			t.Errorf("expected a length of 0, got: %d", length)
		}

		if data != nil {
			t.Errorf("did not expect data, got: %s", data)
		}
	})
}
