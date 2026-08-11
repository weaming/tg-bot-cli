package parser

import (
	"bytes"

	gmAst "github.com/yuin/goldmark/ast"
	gmParser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type richSyntaxKind uint8

const (
	richSyntaxMark richSyntaxKind = iota + 1
	richSyntaxSpoiler
	richSyntaxSuperscript
	richSyntaxSubscript
	richSyntaxMath
	richSyntaxBlockMath
	richSyntaxHTMLCode
)

type richSyntaxNode struct {
	gmAst.BaseInline
	kind    richSyntaxKind
	segment text.Segment
	content text.Segment
}

var richSyntaxNodeKind = gmAst.NewNodeKind("RichSyntax")

func (n *richSyntaxNode) Inline() {
}

func (n *richSyntaxNode) Dump(source []byte, level int) {
	gmAst.DumpHelper(n, source, level, nil, nil)
}

func (n *richSyntaxNode) Kind() gmAst.NodeKind {
	return richSyntaxNodeKind
}

type richSyntaxParser struct{}

func (p *richSyntaxParser) Trigger() []byte {
	return []byte{'^', '~', '=', '|', '$', '<'}
}

func (p *richSyntaxParser) Parse(parent gmAst.Node, block text.Reader, pc gmParser.Context) gmAst.Node {
	source := block.Source()
	_, lineSegment := block.PeekLine()
	start := lineSegment.Start
	if start < 0 || start >= len(source) || (start > 0 && source[start-1] == '\\') {
		return nil
	}

	switch source[start] {
	case '<':
		return p.parseHTMLCode(block, source, start)
	case '$':
		return p.parseMath(block, source, start)
	case '^':
		return p.parseSuperscript(block, source, start)
	case '~':
		return p.parseSubscript(block, source, start)
	case '=':
		return p.parseDelimited(block, source, start, "==", richSyntaxMark)
	case '|':
		return p.parseDelimited(block, source, start, "||", richSyntaxSpoiler)
	default:
		return nil
	}
}

func (p *richSyntaxParser) parseHTMLCode(block text.Reader, source []byte, start int) gmAst.Node {
	openingEnd := findRichHTMLTagEndBytes(source, start+1)
	if openingEnd < 0 || !richHTMLTagNameMatches(source[start+1:openingEnd], "code") {
		return nil
	}

	closingStart, closingEnd := findRichHTMLCodeClosingTag(source, openingEnd+1)
	if closingStart < 0 || closingEnd < 0 {
		return nil
	}

	segment := text.NewSegment(start, closingEnd+1)
	block.Advance(closingEnd + 1 - start)
	return &richSyntaxNode{
		kind:    richSyntaxHTMLCode,
		segment: segment,
	}
}

func (p *richSyntaxParser) parseMath(block text.Reader, source []byte, start int) gmAst.Node {
	if start+1 < len(source) && source[start+1] == '$' {
		closingStart := bytes.Index(source[start+2:], []byte("$$"))
		if closingStart < 0 {
			return nil
		}
		closingStart += start + 2
		return p.newNode(block, start, closingStart+2, closingStart, richSyntaxBlockMath)
	}

	closingStart := findRichSyntaxDelimiter(source, start+1, "$", false)
	if closingStart <= start+1 || (closingStart+1 < len(source) && source[closingStart+1] == '$') {
		return nil
	}
	return p.newNode(block, start, closingStart+1, closingStart, richSyntaxMath)
}

func (p *richSyntaxParser) parseSuperscript(block text.Reader, source []byte, start int) gmAst.Node {
	if (start+1 < len(source) && source[start+1] == '^') || (start > 0 && source[start-1] == '^') {
		return nil
	}
	closingStart := findRichSyntaxDelimiter(source, start+1, "^", false)
	if closingStart <= start+1 {
		return nil
	}
	content := source[start+1 : closingStart]
	if bytes.ContainsAny(content, "^[] \t\r\n") {
		return nil
	}
	return p.newNode(block, start, closingStart+1, closingStart, richSyntaxSuperscript)
}

func (p *richSyntaxParser) parseSubscript(block text.Reader, source []byte, start int) gmAst.Node {
	if (start+1 < len(source) && source[start+1] == '~') || (start > 0 && source[start-1] == '~') {
		return nil
	}
	closingStart := findRichSyntaxDelimiter(source, start+1, "~", false)
	if closingStart <= start+1 || (closingStart+1 < len(source) && source[closingStart+1] == '~') {
		return nil
	}
	content := source[start+1 : closingStart]
	if bytes.ContainsAny(content, "~ \t\r\n") {
		return nil
	}
	return p.newNode(block, start, closingStart+1, closingStart, richSyntaxSubscript)
}

func (p *richSyntaxParser) parseDelimited(block text.Reader, source []byte, start int, delimiter string, kind richSyntaxKind) gmAst.Node {
	if !bytes.HasPrefix(source[start:], []byte(delimiter)) {
		return nil
	}
	closingStart := findRichSyntaxDelimiter(source, start+len(delimiter), delimiter, false)
	if closingStart <= start+len(delimiter) {
		return nil
	}
	return p.newNode(block, start, closingStart+len(delimiter), closingStart, kind)
}

func (p *richSyntaxParser) newNode(block text.Reader, start, stop, contentStop int, kind richSyntaxKind) gmAst.Node {
	contentLength := 1
	if kind == richSyntaxMark || kind == richSyntaxSpoiler || kind == richSyntaxBlockMath {
		contentLength = 2
	}
	contentStart := start + contentLength

	block.Advance(stop - start)
	return &richSyntaxNode{
		kind:    kind,
		segment: text.NewSegment(start, stop),
		content: text.NewSegment(contentStart, contentStop),
	}
}

func findRichSyntaxDelimiter(source []byte, start int, delimiter string, allowNewline bool) int {
	delimiterBytes := []byte(delimiter)
	for index := start; index < len(source); index++ {
		if !allowNewline && source[index] == '\n' {
			return -1
		}
		if source[index] == '\\' {
			index++
			continue
		}
		if bytes.HasPrefix(source[index:], delimiterBytes) {
			return index
		}
	}
	return -1
}

func newRichSyntaxParser() gmParser.InlineParser {
	return &richSyntaxParser{}
}

func richSyntaxParserOption() gmParser.Option {
	return gmParser.WithInlineParsers(
		util.Prioritized(newRichSyntaxParser(), 50),
	)
}

func richSyntaxTag(kind richSyntaxKind) string {
	switch kind {
	case richSyntaxMark:
		return "mark"
	case richSyntaxSpoiler:
		return "tg-spoiler"
	case richSyntaxSuperscript:
		return "sup"
	case richSyntaxSubscript:
		return "sub"
	case richSyntaxMath:
		return "tg-math"
	case richSyntaxBlockMath:
		return "tg-math-block"
	default:
		return ""
	}
}

func (n *richSyntaxNode) isBlockMath() bool {
	return n.kind == richSyntaxBlockMath
}

func (n *richSyntaxNode) contentValue(source []byte) string {
	return string(n.content.Value(source))
}
