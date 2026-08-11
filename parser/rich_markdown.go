package parser

import (
	"strings"

	"github.com/yuin/goldmark"
	gmAst "github.com/yuin/goldmark/ast"
	gmParser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type sourceRange struct {
	start int
	stop  int
}

type richSourceModel struct {
	textRanges []sourceRange
}

// parseRichSource 收集可以安全转换的源码文本区间，代码和 HTML 内容保持原样。
func parseRichSource(input string) richSourceModel {
	source := []byte(input)
	md := goldmark.New(
		goldmark.WithParserOptions(
			gmParser.WithAutoHeadingID(),
		),
	)
	document := md.Parser().Parse(text.NewReader(source))
	model := richSourceModel{
		textRanges: make([]sourceRange, 0),
	}
	openHTMLTags := make([]richHTMLTagFrame, 0)
	collectRichTextRanges(document, &openHTMLTags, &model, source)
	return model
}

// collectRichTextRanges 遍历 AST，排除代码、HTML 和代码块中的文本。
func collectRichTextRanges(node gmAst.Node, openHTMLTags *[]richHTMLTagFrame, model *richSourceModel, source []byte) {
	switch node.Kind() {
	case gmAst.KindCodeSpan,
		gmAst.KindCodeBlock,
		gmAst.KindFencedCodeBlock,
		gmAst.KindHTMLBlock:
		return
	case gmAst.KindRawHTML:
		rawHTMLNode := node.(*gmAst.RawHTML)
		updateRichHTMLTagStack(string(rawHTMLNode.Segments.Value(source)), openHTMLTags)
		return
	case gmAst.KindText:
		if len(*openHTMLTags) > 0 {
			return
		}
		textNode := node.(*gmAst.Text)
		model.textRanges = append(model.textRanges, sourceRange{
			start: textNode.Segment.Start,
			stop:  textNode.Segment.Stop,
		})
	}

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		collectRichTextRanges(child, openHTMLTags, model, source)
	}
}

// rewriteRichTextRanges 只替换 AST 标记出的安全文本区间，并保留其余源码。
func rewriteRichTextRanges(input string, transform func(string) string) string {
	model := parseRichSource(input)
	if len(model.textRanges) == 0 {
		return input
	}
	return rewriteRichSource(input, model,
		func(text string, _ sourceRange) string { return transform(text) },
	)
}

func rewriteRichSource(input string, model richSourceModel, transformText func(string, sourceRange) string) string {
	if len(model.textRanges) == 0 {
		return input
	}

	var builder strings.Builder
	builder.Grow(len(input))
	textIndex := 0
	lastEnd := 0
	for textIndex < len(model.textRanges) {
		currentRange := model.textRanges[textIndex]
		textIndex++

		if currentRange.start < lastEnd || currentRange.start < 0 || currentRange.stop > len(input) {
			continue
		}
		builder.WriteString(input[lastEnd:currentRange.start])
		builder.WriteString(transformText(input[currentRange.start:currentRange.stop], currentRange))
		lastEnd = currentRange.stop
	}
	builder.WriteString(input[lastEnd:])
	return builder.String()
}

func findRichClosingDelimiter(input string, start int, delimiter byte) int {
	for index := start; index < len(input); index++ {
		if input[index] == '\n' {
			return -1
		}
		if input[index] != delimiter || (index > 0 && input[index-1] == '\\') {
			continue
		}
		return index
	}
	return -1
}

// convertRichMarkdownText 将 Rich Markdown 不支持的上下标扩展转换为 HTML 标签。
func convertRichMarkdownText(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	for index := 0; index < len(input); {
		if input[index] == '\\' && index+1 < len(input) {
			builder.WriteString(input[index : index+2])
			index += 2
			continue
		}

		if strings.HasPrefix(input[index:], "$$") {
			closingIndex := strings.Index(input[index+2:], "$$")
			if closingIndex >= 0 {
				closingIndex += index + 2
				builder.WriteString(input[index : closingIndex+2])
				index = closingIndex + 2
				continue
			}
		}

		if input[index] == '$' && (index == 0 || input[index-1] != '$') {
			closingIndex := findRichClosingDelimiter(input, index+1, '$')
			if closingIndex > index+1 && (closingIndex+1 == len(input) || input[closingIndex+1] != '$') {
				builder.WriteString(input[index : closingIndex+1])
				index = closingIndex + 1
				continue
			}
		}

		if input[index] == '^' && (index == 0 || input[index-1] != '^') {
			closingIndex := findRichClosingDelimiter(input, index+1, '^')
			if closingIndex > index+1 {
				content := input[index+1 : closingIndex]
				if !strings.ContainsAny(content, "^[] \t\r\n") {
					builder.WriteString("<sup>")
					builder.WriteString(content)
					builder.WriteString("</sup>")
					index = closingIndex + 1
					continue
				}
			}
		}

		if input[index] == '~' &&
			(index == 0 || input[index-1] != '~') &&
			!strings.HasPrefix(input[index:], "~~") {
			closingIndex := findRichClosingDelimiter(input, index+1, '~')
			if closingIndex > index+1 && (closingIndex+1 == len(input) || input[closingIndex+1] != '~') {
				content := input[index+1 : closingIndex]
				if !strings.ContainsAny(content, "~ \t\r\n") {
					builder.WriteString("<sub>")
					builder.WriteString(content)
					builder.WriteString("</sub>")
					index = closingIndex + 1
					continue
				}
			}
		}

		builder.WriteByte(input[index])
		index++
	}

	return builder.String()
}
