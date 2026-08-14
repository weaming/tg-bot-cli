package cmd

import "testing"

func TestBuildRichMarkdownWithConversion(t *testing.T) {
	content, err := buildRichMessage("md", "收益 ^2^ ~n~", true)
	if err != nil {
		t.Fatalf("build rich markdown: %v", err)
	}

	if content.Markdown != "收益 <sup>2</sup> <sub>n</sub>" {
		t.Errorf("unexpected rich markdown: %q", content.Markdown)
	}
	if content.HTML != "" {
		t.Errorf("rich markdown should not populate HTML: %q", content.HTML)
	}
}

func TestBuildRichMarkdownWithoutConversion(t *testing.T) {
	content, err := buildRichMessage("md", "收益 ^2^ ~n~", false)
	if err != nil {
		t.Fatalf("build raw rich markdown: %v", err)
	}

	if content.Markdown != "收益 ^2^ ~n~" {
		t.Errorf("unexpected raw rich markdown: %q", content.Markdown)
	}
}

func TestShouldConvertMarkdown(t *testing.T) {
	tests := []struct {
		name               string
		noParse            bool
		explicitConversion bool
		inputFile          string
		expected           bool
	}{
		{
			name:      "markdown file enables conversion",
			inputFile: "message.md",
			expected:  true,
		},
		{
			name:               "explicit conversion enables conversion",
			explicitConversion: true,
			expected:           true,
		},
		{
			name:      "no parse disables automatic conversion",
			noParse:   true,
			inputFile: "message.md",
			expected:  false,
		},
		{
			name:               "no parse overrides explicit conversion",
			noParse:            true,
			explicitConversion: true,
			expected:           false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := shouldConvertMarkdown(
				testCase.noParse,
				testCase.explicitConversion,
				testCase.inputFile,
			)
			if actual != testCase.expected {
				t.Errorf("shouldConvertMarkdown() = %v, want %v", actual, testCase.expected)
			}
		})
	}
}
