package handlers

import (
	"babylsp/lsp"
	"context"
)

func DocumentDidOpenHandler(ctx context.Context, params *lsp.TextDocumentOpenParams) error {
	// server, ok := rpc.ServerFromContext(ctx)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to obtain server from context")
	// }
	// server.Log.Println("did open handled")
	//
	// server.Data[params.TextDocument.URI] = [][]byte{{'a', 'b'}}

	return nil
}
