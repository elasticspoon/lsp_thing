package handlers

import (
	"babylsp/lsp"
	"context"
)

func HoverHandler(ctx context.Context, params *lsp.HoverParams) (*lsp.HoverResponse, error) {
	// i should add ID as a struct and then
	// i can create a FromContext function
	// to pull the id out
	id := ctx.Value("id").(*int)
	return &lsp.HoverResponse{
		Response: lsp.Response{
			ID: id,
			Message: lsp.Message{
				RPC: "2.0",
			},
		},
		Result: lsp.HoverResult{
			Contents: lsp.MarkupContent{
				Kind:  "plaintext",
				Value: "howdy",
			},
		},
	}, nil
}
