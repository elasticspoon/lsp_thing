package handlers

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"bufio"
	"context"
	"fmt"
	"os"
)

func HoverHandler(ctx context.Context, params *lsp.HoverParams) (*lsp.HoverResponse, error) {
	id, ok := rpc.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to obtain id from context")
	}

	text, err := WordFinder(params.Position.Character, params.Position.Line, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	return &lsp.HoverResponse{
		Response: lsp.Response{
			ID: &id.ID,
			Message: lsp.Message{
				RPC: "2.0",
			},
		},
		Result: lsp.HoverResult{
			Contents: lsp.MarkupContent{
				Kind:  "plaintext",
				Value: text,
			},
		},
	}, nil
}

func WordFinder(char int, line int, uri string) (string, error) {
	file, err := os.Open(uri)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	// advance to correct line
	for scanner.Scan() && line > 0 {
		line--
	}
	if line != 0 {
		return "", fmt.Errorf("invalid line")
	}

	text := scanner.Text()

	return text, nil
}
