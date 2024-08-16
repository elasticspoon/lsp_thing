package main

import (
	"babylsp/lsp"
	"babylsp/rpc"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	logger := getLogger("/home/bandito/Projects/lsp_thing/log.txt")
	logger.Println("logging started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("got error: %s", err)
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("got method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("unmarshal error: %s", err)
		}
		logger.Printf("Connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)

		writer := os.Stdout
		writer.Write([]byte(rpc.EncodeMessage(msg)))

		logger.Println("Sent bytes")
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("file %s is not vaild", filename))
	}

	return log.New(logfile, "[baby_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
