package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/weaming/tg-bot-cli/parser"
)

func convertMarkdown(input string, splitTable, richMode bool) (string, error) {
	if !richMode {
		return parser.Convert(input, splitTable), nil
	}
	if splitTable {
		return "", fmt.Errorf("--split-table 不能与 --rich 一起使用")
	}
	return parser.ConvertRichHTML(input), nil
}

func main() {
	splitTable := flag.Bool("split-table", false, "Split table into key:value format")
	richMode := flag.Bool("rich", false, "Output Telegram Rich HTML")
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	result, err := convertMarkdown(string(input), *splitTable, *richMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Conversion error: %v\n", err)
		os.Exit(2)
	}

	fmt.Print(result)
}
