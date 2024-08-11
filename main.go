package main

import (
	"babylsp/rpc"
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("main ran")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		handleMessage(scanner.Bytes())
	}
}

func handleMessage(in any) {
	fmt.Println(in)
}
