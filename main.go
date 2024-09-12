package main

import (
	"babylsp/handlers"
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

	server, _ := rpc.NewServer(
		logger,
		reader,
		writer,
		rpc.WithHoverReponse(handlers.HoverHandler),
		rpc.WithInitializeResponse(handlers.IntializeHandler),
	)

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

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("file %s is not vaild", filename))
	}

	return log.New(logfile, "[baby_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
