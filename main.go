package main

import (
	"babylsp/rpc"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
)

func run(ctx context.Context, reader io.Reader, writer io.Writer, _ []string) error {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := getLogger("/home/bandito/Projects/lsp_thing/log.txt")
	logger.Println("logging started")

	server, _ := rpc.NewServer(logger, reader, writer)

	server.Serve()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdin, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// func handleMessage(logger *log.Logger, writer io.Writer, method string, contents []byte) {
// 	logger.Printf("got method: %s", method)
// 	switch method {
// 	case "initialize":
// 		var request lsp.InitializeRequest
// 		if err := json.Unmarshal(contents, &request); err != nil {
// 			logger.Printf("unmarshal error: %s", err)
// 		}
// 		logger.Printf("Connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
//
// 		msg := lsp.NewInitializeResponse(request.ID)
// 		writeReponse(writer, msg)
//
// 		logger.Println("Sent bytes")
// 	case "textDocument/hover":
// 		logger.Println("Sent bytes")
//
// 		var request lsp.HoverRequest
// 		if err := json.Unmarshal(contents, &request); err != nil {
// 			logger.Printf("unmarshal error: %s", err)
// 		}
//
// 		logger.Printf("Hovering over line: %d, char: %d", request.Params.Position.Line, request.Params.Position.Character)
//
// 		msg := lsp.NewHoverResponse(request.ID)
// 		writeReponse(writer, msg)
// 	}
// }

// func initializeHandler(id int) lsp.InitializeResponse {
// 	return lsp.NewInitializeResponse(id)
// }
//
// func hoverHandler(id int) lsp.HoverResponse {
// 	return lsp.NewHoverResponse(id)
// }

func WriteReponse(writer io.Writer, msg any) {
	encodedMsg := rpc.EncodeMessage(msg)
	writer.Write([]byte(encodedMsg))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("file %s is not vaild", filename))
	}

	return log.New(logfile, "[baby_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
