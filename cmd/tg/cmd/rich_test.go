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
