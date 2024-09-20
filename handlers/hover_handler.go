package handlers

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"context"
	"fmt"
)

func HoverHandler(ctx context.Context, params *lsp.HoverParams) {
	server, ok := rpc.ServerFromContext(ctx)
	if !ok {
		fmt.Errorf("failed to obtain server from context")
	}
	id, ok := rpc.FromContext(ctx)
	if !ok {
		fmt.Errorf("failed to obtain id from context")
	}

	data, ok := server.Data[params.TextDocument.URI]
	if !ok {
		server.Log.Println("no data")
	}

	response := &lsp.HoverResponse{
		Response: lsp.Response{
			ID: &id.ID,
			Message: lsp.Message{
				RPC: "2.0",
			},
		},
		Result: lsp.HoverResult{
			Contents: lsp.MarkupContent{
				Kind:  "plaintext",
				Value: "wip",
			},
		},
	}
	rpc.WriteReponse(server.Output, response, server.Log)
}

// func WordFinder(char int, line int, uri string) (string, error) {
// 	// cleanUri := filepath.Clean(uri)
// 	file, err := os.Open(uri)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	scanner := bufio.NewScanner(file)
// 	// advance to correct line
// 	for scanner.Scan() && line > 0 {
// 		line--
// 	}
// 	if line != 0 {
// 		return "", fmt.Errorf("invalid line")
// 	}
//
// 	text := scanner.Text()
//
// 	return text, nil
// }
