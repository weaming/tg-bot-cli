package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	gmAst "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extAst "github.com/yuin/goldmark/extension/ast"
	gmParser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func escapeHTML(text string) string {
	var sb strings.Builder
	for _, r := range text {
		switch r {
		case '&':
			sb.WriteString("&amp;")
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func isAllowedScheme(dest string) bool {
	return strings.HasPrefix(dest, "http://") ||
		strings.HasPrefix(dest, "https://") ||
		strings.HasPrefix(dest, "tg://")
}

func isImageURL(dest string) bool {
	if dest == "" {
		return false
	}
	path := dest
	if idx := strings.IndexByte(path, '?'); idx >= 0 {
		path = path[:idx]
	}
	path = strings.ToLower(path)

	imgExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"}
	for _, ext := range imgExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return strings.HasPrefix(dest, "data:image/")
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// displayWidth returns the display width of text in Telegram <pre>.
// CJK characters render at ~1 column width, same as ASCII.
func displayWidth(text string) int {
	w := 0
	for _, r := range text {
		cp := int(r)
		if (0x1100 <= cp && cp <= 0x115F) ||
			(0x2E80 <= cp && cp <= 0x303F) ||
			(0x3040 <= cp && cp <= 0x33FF) ||
			(0x3400 <= cp && cp <= 0x4DBF) ||
			(0x4E00 <= cp && cp <= 0xA4FF) ||
			(0xAC00 <= cp && cp <= 0xD7FF) ||
			(0xF900 <= cp && cp <= 0xFAFF) ||
			(0xFE10 <= cp && cp <= 0xFE6F) ||
			(0xFF01 <= cp && cp <= 0xFF60) ||
			(0xFFE0 <= cp && cp <= 0xFFE6) {
			w += 1
		} else {
			w += 1
		}
	}
	return w
}

func ljust(text string, width int) string {
	return text + strings.Repeat(" ", maxInt(0, width-displayWidth(text)))
}

type listState struct {
	ordered bool
	start   int
	index   int
}

type Converter struct {
	listDepth  int
	listStack  []listState
	splitTable bool
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Reset() {
	c.listDepth = 0
	c.listStack = c.listStack[:0]
}

func (c *Converter) preProcess(input string) string {
	var result strings.Builder
	i := 0
	for {
		if i >= len(input) {
			break
		}
		idx := strings.Index(input[i:], "$$")
		if idx < 0 {
			result.WriteString(input[i:])
			break
		}
		pos := i + idx

		if pos > 0 && input[pos-1] == '`' {
			closeIdx := strings.Index(input[pos+2:], "`")
			if closeIdx >= 0 {
				result.WriteString(input[i : pos+2+closeIdx+1])
				i = pos + 2 + closeIdx + 1
			} else {
				result.WriteString(input[i:])
				break
			}
			continue
		}

		endIdx := strings.Index(input[pos+2:], "$$")
		if endIdx < 0 {
			result.WriteString(input[i:])
			break
		}
		endPos := pos + 2 + endIdx

		if endPos+2 < len(input) && input[endPos+2] == '`' {
			result.WriteString(input[i : endPos+4])
			i = endPos + 4
			continue
		}

		content := input[pos+2 : endPos]
		before := input[i:pos]
		if before != "" {
			if !strings.HasSuffix(before, "\n") {
				result.WriteString(before + "\n\n")
			} else if !strings.HasSuffix(before, "\n\n") {
				result.WriteString(before + "\n")
			} else {
				result.WriteString(before)
			}
		}
		result.WriteString("```\n" + content + "\n```\n\n")
		i = endPos + 2
	}
	input = result.String()

	input = strings.ReplaceAll(input, "****", "** **")
	input = strings.ReplaceAll(input, "____", "__ __")
	input = strings.ReplaceAll(input, "**__", "** __")
	input = strings.ReplaceAll(input, "__**", "__ **")
	return input
}

func (c *Converter) render(node gmAst.Node, source []byte) string {
	var buf bytes.Buffer
	c.renderTo(&buf, node, source)
	return buf.String()
}

func (c *Converter) renderTo(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	switch node.Kind() {
	case gmAst.KindDocument:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}

	case gmAst.KindHeading:
		buf.WriteString("\n<b>")
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</b>\n\n")

	case gmAst.KindParagraph:
		if c.listDepth > 0 {
			c.renderInlineTo(buf, node, source)
			buf.WriteByte('\n')
		} else {
			c.renderInlineTo(buf, node, source)
			buf.WriteString("\n\n")
		}

	case gmAst.KindText:
		n := node.(*gmAst.Text)
		textVal := string(n.Segment.Value(source))
		buf.WriteString(escapeHTML(textVal))
		if n.HardLineBreak() {
			buf.WriteByte('\n')
		}

	case gmAst.KindEmphasis:
		n := node.(*gmAst.Emphasis)
		if n.Level >= 2 {
			buf.WriteString("<b>")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</b>")
		} else {
			buf.WriteString("<i>")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</i>")
		}

	case gmAst.KindCodeSpan:
		buf.WriteString("<code>")
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</code>")

	case gmAst.KindLink:
		n := node.(*gmAst.Link)
		dest := string(n.Destination)
		if !isAllowedScheme(dest) {
			c.renderInlineChildrenTo(buf, node, source)
			return
		}
		buf.WriteString(`<a href="`)
		buf.WriteString(escapeHTML(dest))
		buf.WriteString(`">`)
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</a>")

	case gmAst.KindImage:
		n := node.(*gmAst.Image)
		dest := string(n.Destination)
		title := string(n.Title)
		if isImageURL(dest) {
			buf.WriteString("<image_url>")
			buf.WriteString(escapeHTML(dest))
			buf.WriteString("</image_url>")
		} else {
			buf.WriteString("<a href=\"")
			buf.WriteString(escapeHTML(dest))
			buf.WriteString("\">[图片: ")
			buf.WriteString(escapeHTML(title))
			buf.WriteString("]</a>")
		}

	case gmAst.KindCodeBlock, gmAst.KindFencedCodeBlock:
		buf.WriteString("<pre><code")
		if fcb, ok := node.(*gmAst.FencedCodeBlock); ok {
			lang := fcb.Language(source)
			if len(lang) > 0 {
				buf.WriteString(` class="language-`)
				buf.WriteString(escapeHTML(string(lang)))
				buf.WriteByte('"')
			}
		}
		buf.WriteString(">")
		lines := node.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			buf.WriteString(escapeHTML(string(line.Value(source))))
		}
		buf.WriteString("</code></pre>\n\n")

	case gmAst.KindBlockquote:
		buf.WriteString("\n")
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}
		buf.WriteString("\n\n")

	case gmAst.KindList:
		n := node.(*gmAst.List)
		c.listDepth++
		c.listStack = append(c.listStack, listState{
			ordered: n.IsOrdered(),
			start:   int(n.Start),
			index:   0,
		})

		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}

		c.listStack = c.listStack[:len(c.listStack)-1]
		c.listDepth--

	case gmAst.KindListItem:
		state := &c.listStack[len(c.listStack)-1]
		state.index++
		indent := strings.Repeat("\u00A0\u00A0", maxInt(0, c.listDepth-1))

		var bullet string
		if state.ordered {
			bullet = fmt.Sprintf("%d.", state.start+state.index-1)
		} else {
			if c.listDepth == 1 {
				bullet = "●"
			} else if c.listDepth == 2 {
				bullet = "○"
			} else {
				bullet = "▪"
			}
		}

		buf.WriteString(indent + bullet + " ")
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}
		buf.WriteByte('\n')

	case gmAst.KindThematicBreak:
		buf.WriteString("\n-------------------\n\n")

	case extAst.KindTable:
		c.renderTableTo(buf, node, source)

	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}
	}
}

func (c *Converter) renderInlineChildrenTo(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		c.renderInlineTo(buf, child, source)
	}
}

func (c *Converter) renderInlineTo(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	switch n := node.(type) {
	case *gmAst.Text:
		textVal := string(n.Segment.Value(source))
		buf.WriteString(escapeHTML(textVal))
		if n.HardLineBreak() {
			buf.WriteByte('\n')
		}

	case *gmAst.Emphasis:
		if n.Level >= 2 {
			buf.WriteString("<b>")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</b>")
		} else {
			buf.WriteString("<i>")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</i>")
		}

	case *gmAst.CodeSpan:
		buf.WriteString("<code>")
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</code>")

	case *gmAst.Link:
		dest := string(n.Destination)
		if !isAllowedScheme(dest) {
			c.renderInlineChildrenTo(buf, node, source)
			return
		}
		buf.WriteString(`<a href="`)
		buf.WriteString(escapeHTML(dest))
		buf.WriteString(`">`)
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</a>")

	case *gmAst.Image:
		dest := string(n.Destination)
		title := string(n.Title)
		if isImageURL(dest) {
			buf.WriteString("<image_url>")
			buf.WriteString(escapeHTML(dest))
			buf.WriteString("</image_url>")
		} else {
			buf.WriteString("<a href=\"")
			buf.WriteString(escapeHTML(dest))
			buf.WriteString("\">[图片: ")
			buf.WriteString(escapeHTML(title))
			buf.WriteString("]</a>")
		}

	default:
		c.renderInlineChildrenTo(buf, node, source)
	}
}

func (c *Converter) renderTableTo(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	table := node.(*extAst.Table)

	headers := make([]string, 0)
	dataRows := make([][]string, 0)

	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		if row.Kind() == extAst.KindTableHeader {
			for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
				if cell.Kind() == extAst.KindTableCell {
					var cellBuf bytes.Buffer
					c.renderInlineTo(&cellBuf, cell, source)
					headers = append(headers, cellBuf.String())
				}
			}
		} else if row.Kind() == extAst.KindTableRow {
			rowData := make([]string, 0)
			for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
				if cell.Kind() == extAst.KindTableCell {
					var cellBuf bytes.Buffer
					c.renderInlineTo(&cellBuf, cell, source)
					rowData = append(rowData, cellBuf.String())
				}
			}
			dataRows = append(dataRows, rowData)
		}
	}

	if c.splitTable {
		buf.WriteString("\n")
		var rows []string
		for _, row := range dataRows {
			var cells []string
			for i, cell := range row {
				header := headers[i]
				if header == "" {
					cells = append(cells, fmt.Sprintf("<b>%s</b>: %s", headers[i], cell))
				} else {
					cells = append(cells, fmt.Sprintf("<b>%s</b>: %s", headers[i], cell))
				}
			}
			rows = append(rows, strings.Join(cells, "\n"))
		}
		buf.WriteString(strings.Join(rows, "\n───────────────\n") + "\n")
		return
	}

	// Grid table format (gentleman style)
	allRows := append([][]string{headers}, dataRows...)
	numCols := 0
	for _, row := range allRows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}
	if numCols == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, numCols)
	for _, row := range allRows {
		for i := 0; i < numCols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			w := displayWidth(cell)
			if w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Build output lines
	buf.WriteString("<pre>")
	for idx, row := range allRows {
		for i := 0; i < numCols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			buf.WriteString(ljust(cell, colWidths[i]))
			if i < numCols-1 {
				buf.WriteString("  ")
			}
		}
		buf.WriteByte('\n')
		if idx == 0 {
			for i := 0; i < numCols; i++ {
				buf.WriteString(strings.Repeat("─", colWidths[i]))
				if i < numCols-1 {
					buf.WriteString("")
				}
			}
			buf.WriteByte('\n')
		}
	}
	buf.WriteString("</pre>\n\n")
}

func Convert(input string, splitTable bool) string {
	c := NewConverter()
	c.splitTable = splitTable
	input = c.preProcess(input)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			gmParser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader([]byte(input))
	doc := md.Parser().Parse(reader)

	result := c.render(doc, []byte(input))

	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")

	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return strings.Trim(result, "\n ")
}
