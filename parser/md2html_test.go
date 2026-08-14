package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func readSampleFile(t *testing.T, fileName string) string {
	t.Helper()

	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位测试文件路径")
	}

	samplePath := filepath.Join(filepath.Dir(testFilePath), "..", "tests", "samples", fileName)
	content, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("读取样例失败 %s: %v", samplePath, err)
	}

	return string(content)
}

func TestHeading(t *testing.T) {
	input := "# Heading 标题"
	expected := "<b>Heading 标题</b>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("Heading failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestBoldItalic(t *testing.T) {
	input := "**bold** and *italic* and ***both***"
	expected := "<b>bold</b> and <i>italic</i> and <i><b>both</b></i>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("BoldItalic failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestCodeBlock(t *testing.T) {
	input := "```python\ndef hello():\n    print(\"world\")\n```"
	expected := "<pre><code class=\"language-python\">def hello():\n    print(\"world\")\n</code></pre>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("CodeBlock failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestInlineCode(t *testing.T) {
	input := "inline code: `const x = 1`"
	expected := "inline code: <code>const x = 1</code>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("InlineCode failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestBulletList(t *testing.T) {
	input := "- bullet 1\n- bullet 2"
	expected := "● bullet 1\n● bullet 2"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("BulletList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestNestedList(t *testing.T) {
	input := "- bullet 1\n  - nested 1\n  - nested 2"
	// Non-breaking spaces (\u00a0) are used for indentation
	// Nested items render on same line as parent text when parent has no own content
	expected := "● bullet 1\u00a0\u00a0○ nested 1\n\u00a0\u00a0○ nested 2"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("NestedList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. first\n2. second\n3. third"
	expected := "1. first\n2. second\n3. third"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("OrderedList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestLink(t *testing.T) {
	input := "[click here](https://example.com)"
	expected := `<a href="https://example.com">click here</a>`
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("Link failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestImage(t *testing.T) {
	input := "![cat](https://example.com/cat.png?size=100)"
	expected := "<image_url>https://example.com/cat.png?size=100</image_url>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("Image failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestHR(t *testing.T) {
	input := "---"
	expected := "-------------------"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("HR failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestLaTeX(t *testing.T) {
	input := "$$E = mc^2$$"
	expected := "<pre><code>E = mc^2\n</code></pre>"
	result := ConvertLegacy(input, false)
	if result != expected {
		t.Errorf("LaTeX failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	result := ConvertLegacy(input, false)
	if len(result) == 0 {
		t.Errorf("Table failed: empty result")
	}
}

func TestRichTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	result := ConvertRichHTML(input)

	if !strings.Contains(result, "<table bordered striped>") {
		t.Errorf("Rich table missing table element: %q", result)
	}
	if !strings.Contains(result, "<th>Name</th><th>Age</th>") {
		t.Errorf("Rich table missing header cells: %q", result)
	}
	if !strings.Contains(result, "<td>Alice</td><td>30</td>") {
		t.Errorf("Rich table missing data cells: %q", result)
	}
	if strings.Contains(result, "<pre>") {
		t.Errorf("Rich table must not use preformatted layout: %q", result)
	}
}

func TestRichMarkdownFeatures(t *testing.T) {
	input := `# Heading

~~deleted~~ ==marked== ||secret|| $x^2$

- [x] completed
- [ ] pending

> quoted text

![emoji](tg://emoji?id=5368324170671202286)

[^note]

[^note]: Footnote text

$$E = mc^2$$

` + "```math\nE = mc^2\n```"
	result := ConvertRichHTML(input)

	for _, expected := range []string{
		"<h1>Heading</h1>",
		"<del>deleted</del>",
		"<mark>marked</mark>",
		"<tg-spoiler>secret</tg-spoiler>",
		"<tg-math>x^2</tg-math>",
		"<ul>",
		`<input type="checkbox" checked>`,
		`<input type="checkbox">`,
		"<blockquote>",
		`<img src="tg://emoji?id=5368324170671202286"/>`,
		`<a href="#fn-1">[1]</a>`,
		`<a name="fn-1"></a>`,
		"<tg-math-block>E = mc^2",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("Rich feature missing %q in %q", expected, result)
		}
	}
	if strings.Contains(result, "<p><tg-math-block>") {
		t.Errorf("Block formula must not be nested in paragraph: %q", result)
	}
}

func TestRichMarkdownSuperscriptAndSubscript(t *testing.T) {
	input := "`^code^` `~code~`\n\n~~deleted~~ ^sup^ ~sub~\n\nFormula: $x^2+y^2$\n\n[^note]\n\n[^note]: Footnote text"
	result := ConvertRichHTML(input)

	for _, expected := range []string{
		"<code>^code^</code>",
		"<code>~code~</code>",
		"<del>deleted</del>",
		"<sup>sup</sup>",
		"<sub>sub</sub>",
		"<tg-math>x^2+y^2</tg-math>",
		`<a href="#fn-1">[1]</a>`,
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("Rich decoration missing %q in %q", expected, result)
		}
	}

	if strings.Contains(result, "<del>sub</del>") {
		t.Errorf("Subscript was parsed as strikethrough: %q", result)
	}
	if strings.Contains(result, "<sup>2 + y</sup>") {
		t.Errorf("Formula was parsed as superscript: %q", result)
	}
}

func TestRichASTProtectsCodeAndConvertsExtensions(t *testing.T) {
	input := "`$x^2$ ==mark== ||secret|| ^sup^ ~sub~ ****`\n\n" +
		"```text\n$x^2$ ==mark== ||secret|| ^sup^ ~sub~ ****\n```\n\n" +
		"$x^2$ ==mark== ||secret|| ^sup^ ~sub~ ****"
	result := ConvertRichHTML(input)

	for _, expected := range []string{
		"<code>$x^2$ ==mark== ||secret|| ^sup^ ~sub~ ****</code>",
		"<pre><code class=\"language-text\">$x^2$ ==mark== ||secret|| ^sup^ ~sub~ ****",
		"<tg-math>x^2</tg-math>",
		"<mark>mark</mark>",
		"<tg-spoiler>secret</tg-spoiler>",
		"<sup>sup</sup>",
		"<sub>sub</sub>",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("Rich AST conversion missing %q in %q", expected, result)
		}
	}
}

func TestRichBlockMathSplitsParagraph(t *testing.T) {
	input := "before $$x^2$$ after"
	result := ConvertRichHTML(input)
	expected := "<p>before</p>\n<tg-math-block>x^2</tg-math-block>\n<p>after</p>"

	if result != expected {
		t.Errorf("block math was not split into block HTML: got %q, want %q", result, expected)
	}
}

func TestRichBlockMathInsideListStaysInList(t *testing.T) {
	input := "- before $$x^2$$ after"
	result := ConvertRichHTML(input)
	expected := "<ul>\n<li>before <tg-math-block>x^2</tg-math-block> after</li>\n</ul>"

	if result != expected {
		t.Errorf("block math left its list: got %q, want %q", result, expected)
	}
}

func TestConvertRichMarkdownPreservesHTMLCodeContent(t *testing.T) {
	input := "<code>^x^ ~x~ ==mark== ||secret||</code>"
	result := ConvertRichMarkdown(input)

	if result != input {
		t.Errorf("HTML code content was changed: got %q, want %q", result, input)
	}
}

func TestRichHTMLPreservesInlineHTMLCodeContent(t *testing.T) {
	input := "<code>^x^ ~x~ ==mark== ||secret||</code>"
	result := ConvertRichHTML(input)
	expected := "<p><code>^x^ ~x~ ==mark== ||secret||</code></p>"

	if result != expected {
		t.Errorf("inline HTML code content was parsed as Markdown: got %q, want %q", result, expected)
	}
}

func TestRichHTMLPreservesInlineHTMLCodeAttributes(t *testing.T) {
	input := `<code class="language-text">^x^ ~x~</code>`
	result := ConvertRichHTML(input)
	expected := `<p><code class="language-text">^x^ ~x~</code></p>`

	if result != expected {
		t.Errorf("inline HTML code was parsed as Markdown: got %q, want %q", result, expected)
	}
}

func TestRichHTMLCodeTagScanPreservesQuotedAttributes(t *testing.T) {
	input := `<CODE data-value="a > b">^x^</CODE>`
	result := ConvertRichHTML(input)
	expected := `<p><CODE data-value="a > b">^x^</CODE></p>`

	if result != expected {
		t.Errorf("HTML code tag scan changed quoted attributes: got %q, want %q", result, expected)
	}
}

func TestRichHTMLPreservesNestedMarkdownInExtensions(t *testing.T) {
	input := "==**bold**== ||_italic_||"
	result := ConvertRichHTML(input)
	expected := "<p><mark><b>bold</b></mark> <tg-spoiler><i>italic</i></tg-spoiler></p>"

	if result != expected {
		t.Errorf("nested Markdown in Rich extensions was not rendered: got %q, want %q", result, expected)
	}
}

func TestRichHTMLNestedMarkdownFastPath(t *testing.T) {
	for _, testCase := range []struct {
		content  string
		expected bool
	}{
		{content: "plain text & symbols", expected: false},
		{content: "**bold**", expected: true},
		{content: `escaped \* text`, expected: true},
		{content: `<u>underlined</u>`, expected: true},
	} {
		if result := hasNestedRichMarkdown(testCase.content); result != testCase.expected {
			t.Errorf("nested Markdown detection for %q: got %t, want %t", testCase.content, result, testCase.expected)
		}
	}

	input := "==plain text & symbols=="
	result := ConvertRichHTML(input)
	expected := "<p><mark>plain text &amp; symbols</mark></p>"
	if result != expected {
		t.Errorf("plain Rich extension content was rendered incorrectly: got %q, want %q", result, expected)
	}
}

func TestConvertRichMarkdownExtensions(t *testing.T) {
	input := "**bold** ^sup^ ~sub~ ~~deleted~~ `^code^` $x^2+y^2$"
	expected := "**bold** <sup>sup</sup> <sub>sub</sub> ~~deleted~~ `^code^` $x^2+y^2$"

	result := ConvertRichMarkdown(input)
	if result != expected {
		t.Errorf("Rich Markdown conversion failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestConvertRichMarkdownExpandsImageReferences(t *testing.T) {
	input := "![Reference image][logo]\n\n[logo]: https://example.com/logo.png \"Logo title\""
	expected := "![Reference image](https://example.com/logo.png \"Logo title\")\n\n[logo]: https://example.com/logo.png \"Logo title\""

	result := ConvertRichMarkdown(input)
	if result != expected {
		t.Errorf("image reference was not expanded:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestConvertRichMarkdownConvertsLinkedImages(t *testing.T) {
	input := "[![Preview](https://example.com/preview.png)](https://example.com/video)"
	expected := "[Preview](https://example.com/video)"

	result := ConvertRichMarkdown(input)
	if result != expected {
		t.Errorf("linked image was not converted:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestConvertRichMarkdownPreservesHTMLCodeAfterText(t *testing.T) {
	input := "prefix <code>^x^ ~x~</code> suffix"
	result := ConvertRichMarkdown(input)

	if result != input {
		t.Errorf("HTML code after text was changed: got %q, want %q", result, input)
	}
}

func TestConvertRichMarkdownCommonMarkGFMExtensionsMatchesExpected(t *testing.T) {
	input := readSampleFile(t, "CommonMark_GFM_Extensions.md")
	expected := readSampleFile(t, "CommonMark_GFM_Extensions.expected-rich.md")
	result := ConvertRichMarkdown(input)

	if result != expected {
		t.Errorf("CommonMark GFM Rich Markdown conversion failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestConvertRichMarkdownOfficialSampleUnchanged(t *testing.T) {
	input := readSampleFile(t, "official-rich-markdown.md")
	result := ConvertRichMarkdown(input)

	if result != input {
		t.Errorf("official Rich Markdown sample was changed:\n%s", result)
	}
}

func TestConvertRichMarkdownOfficialHTMLSampleUnchanged(t *testing.T) {
	input := readSampleFile(t, "official-rich-html.html")
	result := ConvertRichMarkdown(input)

	if result != input {
		t.Errorf("official Rich HTML sample was changed:\n%s", result)
	}
}

func TestConvertRichHTMLMatchesRichMarkdownExtensions(t *testing.T) {
	input := "**bold** ^sup^ ~sub~ ~~deleted~~"
	expected := ConvertRichHTML(ConvertRichMarkdown(input))
	result := ConvertRichHTML(input)

	if result != expected {
		t.Errorf("Rich HTML pipeline mismatch:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestRichMarkdownMedia(t *testing.T) {
	input := `![Photo caption](https://example.com/photo.jpg "Photo caption")

![](https://example.com/video.mp4)

![](https://example.com/audio.mp3)

![](https://example.com/animation.gif)`
	result := ConvertRichHTML(input)

	for _, expected := range []string{
		`<figure><img src="https://example.com/photo.jpg"/><figcaption>Photo caption</figcaption></figure>`,
		`<video src="https://example.com/video.mp4"></video>`,
		`<audio src="https://example.com/audio.mp3"></audio>`,
		`<video src="https://example.com/animation.gif"></video>`,
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("Rich media missing %q in %q", expected, result)
		}
	}
}

func TestRichMarkdownSpecialMedia(t *testing.T) {
	input := `![22:45 tomorrow](tg://time?unix=1647531900&format=wDT)

![](https://example.com/photo.jpg)
![](https://example.com/video.mp4)`
	result := ConvertRichHTML(input)

	for _, expected := range []string{
		`<tg-time unix="1647531900" format="wDT">22:45 tomorrow</tg-time>`,
		`<img src="https://example.com/photo.jpg"/>`,
		`<video src="https://example.com/video.mp4"></video>`,
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("Rich special media missing %q in %q", expected, result)
		}
	}

	if strings.Contains(result, `<p><img`) || strings.Contains(result, `<img src="https://example.com/photo.jpg"/><video`) {
		t.Errorf("Media must be rendered as separate blocks: %q", result)
	}
}
