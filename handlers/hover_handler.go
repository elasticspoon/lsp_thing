package handlers

import "babylsp/lsp"

func HoverHandler(params *lsp.HoverParams) (*lsp.HoverResponse, error) {
	return &lsp.HoverResponse{
		Response: lsp.Response{
			ID: 0,
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
