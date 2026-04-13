package api

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/weaming/tg-bot-cli/parser"
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

// ReadTextOrStdin 直接返回文本
func ReadTextOrStdin(text string) (string, error) {
	return text, nil
}

// ReadFromInput 从文件或 stdin 读取内容，- 表示 stdin
func ReadFromInput(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	if text == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	data, err := os.ReadFile(text)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsMarkdownFile 检测文件扩展名是否为 .md
func IsMarkdownFile(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".md"
}

// ConvertMarkdownToHTML 将 markdown 转换为 HTML
func ConvertMarkdownToHTML(md string, splitTable bool) string {
	return parser.Convert(md, splitTable)
}
