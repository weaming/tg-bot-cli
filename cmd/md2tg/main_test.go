package main

import (
	"strings"
	"testing"
)

func TestConvertMarkdownRichHTML(t *testing.T) {
	result, err := convertMarkdown("# Title\n\n| A | B |\n|---|---|\n| 1 | 2 |", false, "html")
	if err != nil {
		t.Fatalf("convert rich html: %v", err)
	}

	for _, expected := range []string{"<h1>Title</h1>", "<table bordered striped>", "<th>A</th>"} {
		if !strings.Contains(result, expected) {
			t.Errorf("rich HTML missing %q in %q", expected, result)
		}
	}
}

func TestConvertMarkdownRejectsLegacyTableOption(t *testing.T) {
	if _, err := convertMarkdown("| A |\n|---|\n| 1 |", true, "html"); err == nil {
		t.Fatal("expected --split-table and --rich to be rejected together")
	}
}

func TestConvertMarkdownRichMarkdown(t *testing.T) {
	result, err := convertMarkdown("**bold** ^sup^ ~sub~ ~~deleted~~", false, "md")
	if err != nil {
		t.Fatalf("convert rich markdown: %v", err)
	}

	expected := "**bold** <sup>sup</sup> <sub>sub</sub> ~~deleted~~"
	if result != expected {
		t.Errorf("unexpected rich Markdown: got %q, want %q", result, expected)
	}
}

func TestConvertMarkdownRejectsUnknownRichMode(t *testing.T) {
	if _, err := convertMarkdown("text", false, "markdown"); err == nil {
		t.Fatal("expected unknown rich mode to be rejected")
	}
}
