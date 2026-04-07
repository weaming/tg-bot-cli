package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/weaming/tg-bot-cli/parser"
)

func main() {
	splitTable := flag.Bool("split-table", false, "Split table into key:value format")
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	result := parser.Convert(string(input), *splitTable)
	fmt.Print(result)
}
