package cmd

import (
	"fmt"
	"strings"

	"github.com/weaming/tg-bot-cli/api"
)

const (
	richFormatHTML     = "html"
	richFormatMarkdown = "md"
)

func buildRichMessage(format, text string, convertMarkdown bool) (api.RichMessageContent, error) {
	switch strings.ToLower(format) {
	case richFormatHTML:
		if convertMarkdown {
			text = api.ConvertMarkdownToRichHTML(text)
		}
		return api.RichMessageContent{HTML: text}, nil
	case richFormatMarkdown:
		if convertMarkdown {
			text = api.ConvertMarkdownToRichMarkdown(text)
		}
		return api.RichMessageContent{Markdown: text}, nil
	default:
		return api.RichMessageContent{}, fmt.Errorf("--rich 只支持 html 或 md，收到: %q", format)
	}
}
