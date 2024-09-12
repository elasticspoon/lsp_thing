package handlers

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"context"
	"fmt"
)

func HoverHandler(ctx context.Context, params *lsp.HoverParams) (*lsp.HoverResponse, error) {
	id, ok := rpc.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to obtain id from context")
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
				Value: "howdy",
			},
		},
	}, nil
}
