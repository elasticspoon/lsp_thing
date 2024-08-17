package main

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	logger := getLogger("/home/bandito/Projects/lsp_thing/log.txt")
	logger.Println("logging started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("got error: %s", err)
		}
		handleMessage(logger, writer, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, method string, contents []byte) {
	logger.Printf("got method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("unmarshal error: %s", err)
		}
		logger.Printf("Connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		writeReponse(writer, msg)

		logger.Println("Sent bytes")
	case "textDocument/hover":
		logger.Println("Sent bytes")

		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("unmarshal error: %s", err)
		}

		logger.Printf("Hovering over line: %d, char: %d", request.Params.Position.Line, request.Params.Position.Character)

		msg := lsp.NewHoverResponse(request.ID)
		writeReponse(writer, msg)
	}
}

// func initializeHandler(id int) lsp.InitializeResponse {
// 	return lsp.NewInitializeResponse(id)
// }
//
// func hoverHandler(id int) lsp.HoverResponse {
// 	return lsp.NewHoverResponse(id)
// }

func writeReponse(writer io.Writer, msg any) {
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
