package api

import (
	"io"
	"os"
	"strings"
)

// MarkdownV2ReservedChars 会被转义的字符
const markdownV2ReservedChars = `_*[\]()~` + "`" + `>#+=|{}.!-`

// EscapeMarkdownV2 转义 MarkdownV2 保留字符
func EscapeMarkdownV2(text string) string {
	var builder strings.Builder
	for _, char := range text {
		if strings.Contains(markdownV2ReservedChars, string(char)) {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}
	return builder.String()
}

// ReadTextOrStdin 当 text 为 "-" 时从 stdin 读取，否则直接返回
func ReadTextOrStdin(text string) (string, error) {
	if text != "-" {
		return text, nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
