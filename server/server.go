package server

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"
)

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
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		input := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(input)
		if err != nil {
			server.Log.Printf("got error: %s", err)
		}

		msg := Message{
			Context: nil,
			ID:      *1,
			Method:  method,
			Params:  contents,
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
				return handlers.Initialize(&params)
			}
		}
	case "textDocument/hover":
		if handlers.HoverResponse != nil {
			var params lsp.HoverParams
			if err := json.Unmarshal(msg.Params, &params); err == nil {
				return handlers.HoverResponse(&params)
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
	ID      int
	Params  json.RawMessage
}

type Handler interface {
	Hander(ctx *context.Context) (result any, err error)
}
