package parser

import (
	"testing"
)

func TestHeading(t *testing.T) {
	input := "# Heading 标题"
	expected := "<b>Heading 标题</b>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("Heading failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestBoldItalic(t *testing.T) {
	input := "**bold** and *italic* and ***both***"
	expected := "<b>bold</b> and <i>italic</i> and <i><b>both</b></i>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("BoldItalic failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestCodeBlock(t *testing.T) {
	input := "```python\ndef hello():\n    print(\"world\")\n```"
	expected := "<pre><code class=\"language-python\">def hello():\n    print(\"world\")\n</code></pre>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("CodeBlock failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestInlineCode(t *testing.T) {
	input := "inline code: `const x = 1`"
	expected := "inline code: <code>const x = 1</code>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("InlineCode failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestBulletList(t *testing.T) {
	input := "- bullet 1\n- bullet 2"
	expected := "● bullet 1\n● bullet 2"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("BulletList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestNestedList(t *testing.T) {
	input := "- bullet 1\n  - nested 1\n  - nested 2"
	// Non-breaking spaces (\u00a0) are used for indentation
	// Nested items render on same line as parent text when parent has no own content
	expected := "● bullet 1\u00a0\u00a0○ nested 1\n\u00a0\u00a0○ nested 2"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("NestedList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. first\n2. second\n3. third"
	expected := "1. first\n2. second\n3. third"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("OrderedList failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestLink(t *testing.T) {
	input := "[click here](https://example.com)"
	expected := `<a href="https://example.com">click here</a>`
	result := Convert(input, false)
	if result != expected {
		t.Errorf("Link failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestImage(t *testing.T) {
	input := "![cat](https://example.com/cat.png?size=100)"
	expected := "<image_url>https://example.com/cat.png?size=100</image_url>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("Image failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestHR(t *testing.T) {
	input := "---"
	expected := "-------------------"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("HR failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestLaTeX(t *testing.T) {
	input := "$$E = mc^2$$"
	expected := "<pre><code>E = mc^2\n</code></pre>"
	result := Convert(input, false)
	if result != expected {
		t.Errorf("LaTeX failed:\n  got:  %q\n  want: %q", result, expected)
	}
}

func TestTable(t *testing.T) {
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |"
	result := Convert(input, false)
	if len(result) == 0 {
		t.Errorf("Table failed: empty result")
	}
}
