package parser

import (
	"sort"
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

type sourceReplacement struct {
	start int
	stop  int
	value string
}

type richSourceModel struct {
	textRanges             []sourceRange
	structuralReplacements []sourceReplacement
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
		textRanges:             make([]sourceRange, 0),
		structuralReplacements: make([]sourceReplacement, 0),
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
	case gmAst.KindLink:
		if len(*openHTMLTags) == 0 {
			if replacement, ok := createRichLinkedImageReplacement(node.(*gmAst.Link), source); ok {
				model.structuralReplacements = append(model.structuralReplacements, replacement)
				return
			}
		}
	case gmAst.KindImage:
		if len(*openHTMLTags) == 0 {
			if replacement, ok := createRichImageReferenceReplacement(node.(*gmAst.Image), source); ok {
				model.structuralReplacements = append(model.structuralReplacements, replacement)
				return
			}
		}
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
	if len(model.textRanges) == 0 && len(model.structuralReplacements) == 0 {
		return input
	}
	return rewriteRichSource(input, model,
		func(text string, _ sourceRange) string { return transform(text) },
	)
}

func rewriteRichSource(input string, model richSourceModel, transformText func(string, sourceRange) string) string {
	if len(model.textRanges) == 0 && len(model.structuralReplacements) == 0 {
		return input
	}

	replacements := append([]sourceReplacement{}, model.structuralReplacements...)
	for _, currentRange := range model.textRanges {
		if currentRange.start < 0 || currentRange.stop > len(input) || currentRange.start > currentRange.stop {
			continue
		}

		originalText := input[currentRange.start:currentRange.stop]
		convertedText := transformText(originalText, currentRange)
		if convertedText != originalText {
			replacements = append(replacements, sourceReplacement{
				start: currentRange.start,
				stop:  currentRange.stop,
				value: convertedText,
			})
		}
	}

	if len(replacements) == 0 {
		return input
	}

	sort.SliceStable(replacements, func(leftIndex, rightIndex int) bool {
		left := replacements[leftIndex]
		right := replacements[rightIndex]
		if left.start != right.start {
			return left.start < right.start
		}
		return left.stop > right.stop
	})

	var builder strings.Builder
	builder.Grow(len(input))
	lastEnd := 0
	for _, replacement := range replacements {
		if replacement.start < lastEnd || replacement.start < 0 || replacement.stop > len(input) || replacement.start > replacement.stop {
			continue
		}
		builder.WriteString(input[lastEnd:replacement.start])
		builder.WriteString(replacement.value)
		lastEnd = replacement.stop
	}
	builder.WriteString(input[lastEnd:])
	return builder.String()
}

func createRichLinkedImageReplacement(link *gmAst.Link, source []byte) (sourceReplacement, bool) {
	image := findRichImageDescendant(link.FirstChild())
	if image == nil {
		return sourceReplacement{}, false
	}

	rangeValue, ok := richNodeSourceRange(link, source, false)
	if !ok {
		return sourceReplacement{}, false
	}

	value := "[" + escapeRichMarkdownLabel(richImageAltText(image, source)) + "]("
	value += formatRichMarkdownDestination(link.Destination)
	value += formatRichMarkdownTitle(link.Title)
	value += ")"
	return sourceReplacement{
		start: rangeValue.start,
		stop:  rangeValue.stop,
		value: value,
	}, true
}

func createRichImageReferenceReplacement(image *gmAst.Image, source []byte) (sourceReplacement, bool) {
	if image.Reference == nil {
		return sourceReplacement{}, false
	}

	rangeValue, ok := richNodeSourceRange(image, source, true)
	if !ok {
		return sourceReplacement{}, false
	}

	value := "![" + escapeRichMarkdownLabel(richImageAltText(image, source)) + "]("
	value += formatRichMarkdownDestination(image.Destination)
	value += formatRichMarkdownTitle(image.Title)
	value += ")"
	return sourceReplacement{
		start: rangeValue.start,
		stop:  rangeValue.stop,
		value: value,
	}, true
}

func findRichImageDescendant(node gmAst.Node) *gmAst.Image {
	for currentNode := node; currentNode != nil; currentNode = currentNode.NextSibling() {
		if image, ok := currentNode.(*gmAst.Image); ok {
			return image
		}
		if image := findRichImageDescendant(currentNode.FirstChild()); image != nil {
			return image
		}
	}
	return nil
}

func richImageAltText(image *gmAst.Image, source []byte) string {
	var builder strings.Builder
	appendRichNodeText(&builder, image.FirstChild(), source)
	return builder.String()
}

func appendRichNodeText(builder *strings.Builder, node gmAst.Node, source []byte) {
	for currentNode := node; currentNode != nil; currentNode = currentNode.NextSibling() {
		switch typedNode := currentNode.(type) {
		case *gmAst.Text:
			builder.Write(typedNode.Segment.Value(source))
			if typedNode.SoftLineBreak() {
				builder.WriteByte('\n')
			}
		case *gmAst.String:
			builder.Write(typedNode.Value)
		default:
			appendRichNodeText(builder, currentNode.FirstChild(), source)
		}
	}
}

func richNodeSourceRange(node gmAst.Node, source []byte, isImage bool) (sourceRange, bool) {
	start := node.Pos()
	if start < 0 || start >= len(source) {
		return sourceRange{}, false
	}

	stop, ok := findRichMarkdownNodeEnd(source, start, isImage)
	if !ok {
		return sourceRange{}, false
	}
	return sourceRange{start: start, stop: stop}, true
}

func findRichMarkdownNodeEnd(source []byte, start int, isImage bool) (int, bool) {
	labelStart := start
	if isImage {
		if source[labelStart] != '!' {
			return 0, false
		}
		labelStart++
	}
	if labelStart >= len(source) || source[labelStart] != '[' {
		return 0, false
	}

	labelEnd, ok := findRichMarkdownClosingBracket(source, labelStart)
	if !ok {
		return 0, false
	}
	if labelEnd+1 >= len(source) {
		return labelEnd + 1, true
	}

	switch source[labelEnd+1] {
	case '(':
		return findRichMarkdownClosingParenthesis(source, labelEnd+1)
	case '[':
		return findRichMarkdownClosingBracketEnd(source, labelEnd+1)
	default:
		return labelEnd + 1, true
	}
}

func findRichMarkdownClosingBracket(source []byte, openingIndex int) (int, bool) {
	depth := 0
	for index := openingIndex; index < len(source); index++ {
		if source[index] == '\\' {
			index++
			continue
		}
		switch source[index] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return index, true
			}
		}
	}
	return 0, false
}

func findRichMarkdownClosingBracketEnd(source []byte, openingIndex int) (int, bool) {
	closingIndex, ok := findRichMarkdownClosingBracket(source, openingIndex)
	if !ok {
		return 0, false
	}
	return closingIndex + 1, true
}

func findRichMarkdownClosingParenthesis(source []byte, openingIndex int) (int, bool) {
	depth := 1
	quote := byte(0)
	inAngleDestination := false
	hasDestinationWhitespace := false

	for index := openingIndex + 1; index < len(source); index++ {
		currentByte := source[index]
		if currentByte == '\\' {
			index++
			continue
		}
		if quote != 0 {
			if currentByte == quote {
				quote = 0
			}
			continue
		}

		if currentByte == '<' {
			inAngleDestination = true
			continue
		}
		if currentByte == '>' {
			inAngleDestination = false
			continue
		}
		if !inAngleDestination && hasDestinationWhitespace && (currentByte == '\'' || currentByte == '"') {
			quote = currentByte
			continue
		}
		if currentByte == ' ' || currentByte == '\t' || currentByte == '\r' || currentByte == '\n' {
			hasDestinationWhitespace = true
			continue
		}
		if inAngleDestination {
			continue
		}

		switch currentByte {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return index + 1, true
			}
		}
	}

	return 0, false
}

func escapeRichMarkdownLabel(label string) string {
	label = strings.ReplaceAll(label, "\\", "\\\\")
	return strings.ReplaceAll(label, "]", "\\]")
}

func formatRichMarkdownDestination(destination []byte) string {
	value := string(destination)
	if strings.ContainsAny(value, " \t\r\n()") {
		return "<" + value + ">"
	}
	return value
}

func formatRichMarkdownTitle(title []byte) string {
	if len(title) == 0 {
		return ""
	}
	value := strings.ReplaceAll(string(title), "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return ` "` + value + `"`
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
