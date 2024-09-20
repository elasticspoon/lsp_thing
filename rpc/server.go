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
	Input    io.Reader
	Output   io.Writer
	Context  context.Context
	Log      *log.Logger
	Data     map[string][][]byte
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

	server := Server{
		Log:      log,
		Input:    input,
		Output:   output,
		Timeout:  DefaultTimeout,
		handlers: handlers,
	}
	server.Context = NewServerContext(context.Background(), &server)
	server.Log.Println(server.Context)
	return &server, nil
}

func (server *Server) Serve() {
	scanner := bufio.NewScanner(server.Input)
	scanner.Split(Split)

	for scanner.Scan() {
		input := scanner.Bytes()
		request, err := DecodeMessage(input)
		if err != nil {
			server.Log.Printf("got error: %s", err)
		}

		server.Log.Println(request.Method)
		server.Log.Println(string(*request.Params))
		msg := Message{
			Method: request.Method,
			Params: *request.Params,
		}
		if !request.NullID {
			msg.ID = *request.ID
			msg.Context = NewContext(server.Context, request.ID)
		}

		server.handlers.Handle(msg)
	}
}

func WriteReponse(writer io.Writer, msg any, log *log.Logger) {
	encodedMsg := EncodeMessage(msg)
	log.Printf("writing: %s", string(encodedMsg))
	writer.Write([]byte(encodedMsg))
}

type handlers struct {
	HoverResponse       lsp.HoverResponseFunc
	TextDocumentDidOpen lsp.DocumentDidOpenFunc
	Initialize          lsp.InitializeResponseFunc
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
	case "textDocument/didOpen":
		if handlers.TextDocumentDidOpen != nil {
			var params lsp.TextDocumentOpenParams
			if err := json.Unmarshal(msg.Params, &params); err == nil {
				err = handlers.TextDocumentDidOpen(msg.Context, &params)
				return nil, err
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

func WithDidOpenHandler(hoverFunc lsp.DocumentDidOpenFunc) HandlerOption {
	return func(handlers *handlers) error {
		handlers.TextDocumentDidOpen = hoverFunc
		return nil
	}
}

type Message struct {
	Context context.Context
	Method  string
	Params  json.RawMessage
	ID
}

type Handler interface {
	Hander(ctx *context.Context) (result any, err error)
}

var serverKey key = 2

func NewServerContext(ctx context.Context, server *Server) context.Context {
	return context.WithValue(ctx, serverKey, server)
}

func ServerFromContext(ctx context.Context) (*Server, bool) {
	id, ok := ctx.Value(serverKey).(*Server)
	return id, ok
}
