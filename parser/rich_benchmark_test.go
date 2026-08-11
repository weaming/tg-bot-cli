package parser

import (
	"strings"
	"testing"
)

func BenchmarkConvertRichHTMLPlainExtensions(b *testing.B) {
	input := strings.Repeat("==plain text== ||plain spoiler|| ", 100)
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichHTML(input)
	}
}

func BenchmarkConvertRichHTMLNestedExtensions(b *testing.B) {
	input := strings.Repeat("==**bold text**== ||_italic text_|| ", 100)
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichHTML(input)
	}
}

func FuzzConvertRichFormatsDoNotPanic(f *testing.F) {
	for _, seed := range []string{
		"plain text",
		"^sup^ ~sub~ $x^2$ $$x^2$$",
		"==mark== ||spoiler|| **bold** `code`",
		"<u>HTML</u> <code>^code^</code>",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 8192 {
			t.Skip()
		}
		ConvertRichMarkdown(input)
		ConvertRichHTML(input)
	})
}

func FuzzRichCodeProtection(f *testing.F) {
	for _, seed := range []string{
		"plain code",
		"^x^ ~x~ ==mark== ||spoiler||",
		"<tag> & text",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, content string) {
		if len(content) > 4096 {
			t.Skip()
		}

		safeContent := strings.NewReplacer(
			"<", "&lt;",
			"\r", " ",
			"\n", " ",
		).Replace(content)
		input := "<code>" + safeContent + "</code>"
		result := ConvertRichHTML(input)
		if !strings.Contains(result, safeContent) {
			t.Fatalf("code content was lost: got %q, want content %q", result, safeContent)
		}
		for _, tagName := range []string{"<sup>", "<sub>", "<mark>", "<tg-spoiler>", "<tg-math>", "<tg-math-block>"} {
			if strings.Contains(result, tagName) {
				t.Fatalf("code content was converted to %s: got %q", tagName, result)
			}
		}
	})
}
