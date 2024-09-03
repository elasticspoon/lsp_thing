package rpc

import (
	"babylsp/lsp"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (int, []byte, error) {
	header, content, found := bytes.Cut(data, []byte("\r\n\r\n"))
	if !found {
		return 0, nil, nil // we are waitng
	}

	contentLengthBytes := header[HEADER_PREAMBLE_LEN:]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}

	// Wait for more
	if len(content) < contentLength {
		return 0, nil, nil
	}

	// header + /r/n/r/n + content
	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}

var DefaultTimeout = time.Second * 10

type Server struct {
	handlers handlers
	Input    *io.Reader
	Output   *io.Writer
	Log      *log.Logger
	Timeout  time.Duration
	Debug    bool
}

func NewServer(log *log.Logger, input io.Reader, output io.Writer, opts ...HandlerOption) (*Server, error) {
	handlers := handlers{}
	for _, handler := range opts {
		if err := handler(&handlers); err != nil {
			return nil, err
		}
	}

	return &Server{
		Log:      log,
		Input:    &input,
		Output:   &output,
		Timeout:  DefaultTimeout,
		handlers: handlers,
	}, nil
}

func (server *Server) Serve() {
	scanner := bufio.NewScanner(*server.Input)
	scanner.Split(Split)

	for scanner.Scan() {
		input := scanner.Bytes()
		request, err := DecodeMessage(input)
		if err != nil {
			server.Log.Printf("got error: %s", err)
		}

		server.Log.Println(request.Method)
		msg := Message{
			Context: nil,
			Method:  request.Method,
			Params:  *request.Params,
		}
		if request.ID != nil {
			msg.ID = *request.ID
		}
		server.handlers.Handle(msg)
	}
}

type handlers struct {
	HoverResponse lsp.HoverResponseFunc
	Initialize    lsp.InitializeResponseFunc
}

func (handlers *handlers) Handle(msg Message) (any, error) {
	switch msg.Method {
	case "initialize":
		if handlers.Initialize != nil {
			var params lsp.InitializeRequestParams
			if err := json.Unmarshal(msg.Params, &params); err == nil {
				return handlers.Initialize(msg.Context, &params)
			}
		}
	case "textDocument/hover":
		if handlers.HoverResponse != nil {
			var params lsp.HoverParams
			if err := json.Unmarshal(msg.Params, &params); err == nil {
				return handlers.HoverResponse(msg.Context, &params)
			}
		}
	}
	return nil, nil
}

type HandlerOption func(handlers *handlers) error

func WithInitializeResponse(initFunc lsp.InitializeResponseFunc) HandlerOption {
	return func(handlers *handlers) error {
		handlers.Initialize = initFunc
		return nil
	}
}

func WithHoverReponse(hoverFunc lsp.HoverResponseFunc) HandlerOption {
	return func(handlers *handlers) error {
		handlers.HoverResponse = hoverFunc
		return nil
	}
}

type Message struct {
	Context context.Context
	Method  string
	Params  json.RawMessage
	ID      int
}

type Handler interface {
	Hander(ctx *context.Context) (result any, err error)
}
