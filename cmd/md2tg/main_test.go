package main

import (
	"strings"
	"testing"
)

func TestConvertMarkdownRichHTML(t *testing.T) {
	result, err := convertMarkdown("# Title\n\n| A | B |\n|---|---|\n| 1 | 2 |", false, true)
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
	if _, err := convertMarkdown("| A |\n|---|\n| 1 |", true, true); err == nil {
		t.Fatal("expected --split-table and --rich to be rejected together")
	}
}
