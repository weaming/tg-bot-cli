package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/weaming/tg-bot-cli/parser"
)

func convertMarkdown(input string, splitTable bool, richMode string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(richMode)) {
	case "":
		return parser.Convert(input, splitTable), nil
	case "html":
		if splitTable {
			return "", fmt.Errorf("--split-table 不能与 --rich 一起使用")
		}
		return parser.ConvertRichHTML(input), nil
	case "md":
		if splitTable {
			return "", fmt.Errorf("--split-table 不能与 --rich 一起使用")
		}
		return parser.ConvertRichMarkdown(input), nil
	default:
		return "", fmt.Errorf("--rich 只支持 html 或 md，收到: %q", richMode)
	}
}

func main() {
	splitTable := flag.Bool("split-table", false, "Split table into key:value format")
	richMode := flag.String("rich", "", "Output Telegram Rich format: html or md")
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
