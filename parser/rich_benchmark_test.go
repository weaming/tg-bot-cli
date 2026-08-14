package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func readBenchmarkSample(b *testing.B, fileName string) string {
	b.Helper()

	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		b.Fatal("无法定位 benchmark 文件路径")
	}

	samplePath := filepath.Join(filepath.Dir(testFilePath), "..", "tests", "samples", fileName)
	content, err := os.ReadFile(samplePath)
	if err != nil {
		b.Fatalf("读取 benchmark 样例失败 %s: %v", samplePath, err)
	}

	return string(content)
}

func BenchmarkConvertCommonMarkGFMExtensions(b *testing.B) {
	input := readBenchmarkSample(b, "CommonMark_GFM_Extensions.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertLegacy(input, false)
	}
}

func BenchmarkConvertOfficialRichMarkdown(b *testing.B) {
	input := readBenchmarkSample(b, "official-rich-markdown.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertLegacy(input, false)
	}
}

func BenchmarkConvertRichHTMLCommonMarkGFMExtensions(b *testing.B) {
	input := readBenchmarkSample(b, "CommonMark_GFM_Extensions.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichHTML(input)
	}
}

func BenchmarkConvertRichHTMLOfficialRichMarkdown(b *testing.B) {
	input := readBenchmarkSample(b, "official-rich-markdown.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichHTML(input)
	}
}

func BenchmarkConvertRichMarkdownCommonMarkGFMExtensions(b *testing.B) {
	input := readBenchmarkSample(b, "CommonMark_GFM_Extensions.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichMarkdown(input)
	}
}

func BenchmarkConvertRichMarkdownOfficialRichMarkdown(b *testing.B) {
	input := readBenchmarkSample(b, "official-rich-markdown.md")
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		ConvertRichMarkdown(input)
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
