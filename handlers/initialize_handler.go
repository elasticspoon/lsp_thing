package handlers

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"context"
	"fmt"
)

func IntializeHandler(ctx context.Context, params *lsp.InitializeRequestParams) (*lsp.InitializeResponse, error) {
	id, ok := rpc.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to obtain id from context")
	}

	return lsp.NewInitializeResponse(id.ID), nil
}
