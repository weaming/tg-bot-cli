package cmd

import (
	"fmt"
	"strings"

	"github.com/weaming/tg-bot-cli/api"
)

const (
	richFormatHTML     = "html"
	richFormatMarkdown = "markdown"
)

func buildRichMessage(format, text string, convertMarkdown bool) (api.RichMessageContent, error) {
	switch strings.ToLower(format) {
	case richFormatHTML:
		if convertMarkdown {
			text = api.ConvertMarkdownToRichHTML(text)
		}
		return api.RichMessageContent{HTML: text}, nil
	case richFormatMarkdown:
		return api.RichMessageContent{Markdown: text}, nil
	default:
		return api.RichMessageContent{}, fmt.Errorf("--rich 只支持 html 或 markdown，收到: %q", format)
	}
}
