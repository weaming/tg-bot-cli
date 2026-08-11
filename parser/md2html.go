package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
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

func richMediaTag(dest string) string {
	mediaPath := strings.ToLower(dest)
	if queryIndex := strings.IndexAny(mediaPath, "?#"); queryIndex >= 0 {
		mediaPath = mediaPath[:queryIndex]
	}

	switch {
	case strings.HasSuffix(mediaPath, ".mp4"),
		strings.HasSuffix(mediaPath, ".mov"),
		strings.HasSuffix(mediaPath, ".webm"),
		strings.HasSuffix(mediaPath, ".mkv"),
		strings.HasSuffix(mediaPath, ".gif"):
		return "video"
	case strings.HasSuffix(mediaPath, ".mp3"),
		strings.HasSuffix(mediaPath, ".ogg"),
		strings.HasSuffix(mediaPath, ".m4a"),
		strings.HasSuffix(mediaPath, ".wav"),
		strings.HasSuffix(mediaPath, ".flac"):
		return "audio"
	default:
		return "img"
	}
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
	richTable  bool
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Reset() {
	c.listDepth = 0
	c.listStack = c.listStack[:0]
}

func (c *Converter) preProcess(input string) string {
	if c.richTable {
		return c.preProcessRich(input)
	}

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

var (
	richBlockMathPattern  = regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
	richInlineMathPattern = regexp.MustCompile(`(^|[^$])\$([^$\n]+)\$([^$]|$)`)
	richMarkPattern       = regexp.MustCompile(`==([^=\n]+)==`)
	richSpoilerPattern    = regexp.MustCompile(`\|\|([^|\n]+)\|\|`)
)

func (c *Converter) preProcessRich(input string) string {
	input = richBlockMathPattern.ReplaceAllString(input, "\n\n<tg-math-block>$1</tg-math-block>\n\n")
	input = richInlineMathPattern.ReplaceAllString(input, `$1<tg-math>$2</tg-math>$3`)
	input = richMarkPattern.ReplaceAllString(input, `<mark>$1</mark>`)
	input = richSpoilerPattern.ReplaceAllString(input, `<tg-spoiler>$1</tg-spoiler>`)

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
		if c.richTable {
			heading := node.(*gmAst.Heading)
			tagName := fmt.Sprintf("h%d", heading.Level)
			buf.WriteString("<")
			buf.WriteString(tagName)
			buf.WriteString(">")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</")
			buf.WriteString(tagName)
			buf.WriteString(">\n")
			return
		}
		buf.WriteString("\n<b>")
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</b>\n\n")

	case gmAst.KindParagraph:
		if c.richTable && c.listDepth == 0 {
			if c.renderRichMediaParagraph(buf, node, source) {
				return
			}

			var paragraphBuffer bytes.Buffer
			c.renderInlineChildrenTo(&paragraphBuffer, node, source)
			paragraphContent := paragraphBuffer.String()
			if strings.HasPrefix(paragraphContent, "<tg-math-block>") {
				buf.WriteString(paragraphContent)
				buf.WriteByte('\n')
			} else {
				buf.WriteString("<p>")
				buf.WriteString(paragraphContent)
				buf.WriteString("</p>\n")
			}
			return
		}
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
		if c.richTable {
			if fencedCodeBlock, ok := node.(*gmAst.FencedCodeBlock); ok && strings.EqualFold(string(fencedCodeBlock.Language(source)), "math") {
				buf.WriteString("<tg-math-block>")
				c.renderCodeLines(buf, node, source)
				buf.WriteString("</tg-math-block>\n")
				return
			}
		}
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
		c.renderCodeLines(buf, node, source)
		buf.WriteString("</code></pre>\n\n")

	case gmAst.KindBlockquote:
		if c.richTable {
			buf.WriteString("<blockquote>")
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				c.renderTo(buf, child, source)
			}
			buf.WriteString("</blockquote>\n")
			return
		}
		buf.WriteString("\n")
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			c.renderTo(buf, child, source)
		}
		buf.WriteString("\n\n")

	case gmAst.KindList:
		n := node.(*gmAst.List)
		if c.richTable {
			if n.IsOrdered() {
				buf.WriteString("<ol")
				if n.Start > 1 {
					fmt.Fprintf(buf, ` start="%d"`, n.Start)
				}
				buf.WriteString(">\n")
			} else {
				buf.WriteString("<ul>\n")
			}
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
			if n.IsOrdered() {
				buf.WriteString("</ol>\n")
			} else {
				buf.WriteString("</ul>\n")
			}
			return
		}
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
		if c.richTable {
			buf.WriteString("<li>")
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				c.renderTo(buf, child, source)
			}
			buf.WriteString("</li>\n")
			return
		}
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
		if c.richTable {
			buf.WriteString("<hr/>\n")
			return
		}
		buf.WriteString("\n-------------------\n\n")

	case extAst.KindFootnoteList:
		if c.richTable {
			buf.WriteString("<hr/><ol>\n")
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				c.renderTo(buf, child, source)
			}
			buf.WriteString("</ol>\n")
			return
		}

	case extAst.KindFootnote:
		if c.richTable {
			footnote := node.(*extAst.Footnote)
			fmt.Fprintf(buf, `<li><a name="fn-%d"></a>`, footnote.Index)
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				c.renderTo(buf, child, source)
			}
			buf.WriteString("</li>\n")
			return
		}

	case gmAst.KindHTMLBlock:
		if c.richTable {
			c.renderRawBlockHTML(buf, node, source)
			return
		}

	case extAst.KindTable:
		c.renderTableTo(buf, node, source)

	default:
		if node.Type() == gmAst.TypeInline {
			c.renderInlineTo(buf, node, source)
			return
		}
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

func (c *Converter) renderRichMediaParagraph(buf *bytes.Buffer, node gmAst.Node, source []byte) bool {
	mediaNode := node.FirstChild()
	if mediaNode == nil {
		return false
	}

	for currentNode := mediaNode; currentNode != nil; currentNode = currentNode.NextSibling() {
		switch typedNode := currentNode.(type) {
		case *gmAst.Image:
		case *gmAst.Text:
			if strings.TrimSpace(string(typedNode.Segment.Value(source))) != "" {
				return false
			}
		default:
			return false
		}
	}

	for currentNode := mediaNode; currentNode != nil; currentNode = currentNode.NextSibling() {
		if imageNode, ok := currentNode.(*gmAst.Image); ok {
			c.renderRichImage(buf, imageNode, source)
		}
	}
	return true
}

func (c *Converter) renderRichImage(buf *bytes.Buffer, imageNode *gmAst.Image, source []byte) {
	if c.renderRichTime(buf, imageNode, source) {
		return
	}

	destination := string(imageNode.Destination)
	title := string(imageNode.Title)
	tagName := richMediaTag(destination)
	mediaAttributes := fmt.Sprintf(` src="%s"`, escapeHTML(destination))

	if title != "" {
		buf.WriteString("<figure>")
		c.renderRichMediaElement(buf, tagName, mediaAttributes)
		buf.WriteString("<figcaption>")
		buf.WriteString(escapeHTML(title))
		buf.WriteString("</figcaption></figure>\n")
		return
	}

	c.renderRichMediaElement(buf, tagName, mediaAttributes)
	buf.WriteByte('\n')
}

func (c *Converter) renderRichTime(buf *bytes.Buffer, imageNode *gmAst.Image, source []byte) bool {
	destination := string(imageNode.Destination)
	parsedURL, err := url.Parse(destination)
	if err != nil || parsedURL.Scheme != "tg" || parsedURL.Host != "time" {
		return false
	}

	unixTimestamp := parsedURL.Query().Get("unix")
	format := parsedURL.Query().Get("format")
	if unixTimestamp == "" || format == "" {
		return false
	}

	var altText bytes.Buffer
	c.renderInlineChildrenTo(&altText, imageNode, source)

	buf.WriteString(`<tg-time unix="`)
	buf.WriteString(escapeHTML(unixTimestamp))
	buf.WriteString(`" format="`)
	buf.WriteString(escapeHTML(format))
	buf.WriteString(`">`)
	buf.WriteString(altText.String())
	buf.WriteString("</tg-time>\n")
	return true
}

func (c *Converter) renderRichMediaElement(buf *bytes.Buffer, tagName, mediaAttributes string) {
	buf.WriteString("<")
	buf.WriteString(tagName)
	buf.WriteString(mediaAttributes)
	if tagName == "img" {
		buf.WriteString("/>")
		return
	}
	buf.WriteString("></")
	buf.WriteString(tagName)
	buf.WriteString(">")
}

func (c *Converter) renderCodeLines(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.WriteString(escapeHTML(string(line.Value(source))))
	}
}

func (c *Converter) renderRawBlockHTML(buf *bytes.Buffer, node gmAst.Node, source []byte) {
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.WriteString(string(line.Value(source)))
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

	case *extAst.Strikethrough:
		if c.richTable {
			buf.WriteString("<del>")
			c.renderInlineChildrenTo(buf, node, source)
			buf.WriteString("</del>")
			return
		}
		c.renderInlineChildrenTo(buf, node, source)

	case *extAst.TaskCheckBox:
		if c.richTable {
			buf.WriteString(`<input type="checkbox"`)
			if n.IsChecked {
				buf.WriteString(" checked")
			}
			buf.WriteString("> ")
			return
		}
		buf.WriteString("☐ ")

	case *gmAst.AutoLink:
		if c.richTable {
			url := string(n.URL(source))
			buf.WriteString(`<a href="`)
			buf.WriteString(escapeHTML(url))
			buf.WriteString(`">`)
			buf.WriteString(escapeHTML(string(n.Label(source))))
			buf.WriteString("</a>")
			return
		}
		buf.WriteString(escapeHTML(string(n.Label(source))))

	case *gmAst.RawHTML:
		if c.richTable {
			buf.WriteString(string(n.Text(source)))
			return
		}
		buf.WriteString(escapeHTML(string(n.Text(source))))

	case *extAst.FootnoteLink:
		if c.richTable {
			fmt.Fprintf(buf, `<a name="fnref-%d"></a><a href="#fn-%d">[%d]</a>`, n.Index, n.Index, n.Index)
			return
		}
		fmt.Fprintf(buf, "[%d]", n.Index)

	case *extAst.FootnoteBacklink:
		if c.richTable {
			fmt.Fprintf(buf, `<a href="#fnref-%d">↩</a>`, n.Index)
			return
		}
		buf.WriteString("↩")

	case *gmAst.CodeSpan:
		buf.WriteString("<code>")
		c.renderInlineChildrenTo(buf, node, source)
		buf.WriteString("</code>")

	case *gmAst.Link:
		dest := string(n.Destination)
		if !c.richTable && !isAllowedScheme(dest) {
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
		if c.richTable {
			c.renderRichImage(buf, n, source)
			return
		}
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
	if c.richTable {
		c.renderRichTableTo(buf, table, source)
		return
	}

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

func (c *Converter) renderRichTableTo(buf *bytes.Buffer, table *extAst.Table, source []byte) {
	buf.WriteString("<table bordered striped>\n")

	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		if row.Kind() != extAst.KindTableHeader && row.Kind() != extAst.KindTableRow {
			continue
		}

		buf.WriteString("<tr>")
		cellTag := "td"
		if row.Kind() == extAst.KindTableHeader {
			cellTag = "th"
		}

		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			if cell.Kind() != extAst.KindTableCell {
				continue
			}
			buf.WriteString("<")
			buf.WriteString(cellTag)
			if tableCell, ok := cell.(*extAst.TableCell); ok && tableCell.Alignment != extAst.AlignNone {
				fmt.Fprintf(buf, ` align="%s"`, tableCell.Alignment.String())
			}
			buf.WriteString(">")
			c.renderInlineTo(buf, cell, source)
			buf.WriteString("</")
			buf.WriteString(cellTag)
			buf.WriteString(">")
		}

		buf.WriteString("</tr>\n")
	}

	buf.WriteString("</table>\n\n")
}

func Convert(input string, splitTable bool) string {
	c := NewConverter()
	c.splitTable = splitTable
	return c.convert(input)
}

// ConvertRichHTML 将 Markdown 转换为 Rich Message HTML，并保留原生表格结构。
func ConvertRichHTML(input string) string {
	c := NewConverter()
	c.richTable = true
	return c.convert(input)
}

func (c *Converter) convert(input string) string {
	input = c.preProcess(input)

	extensions := []goldmark.Extender{extension.GFM}
	if c.richTable {
		extensions = append(extensions, extension.Footnote)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
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
