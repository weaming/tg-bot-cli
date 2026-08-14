import { readFileSync } from "node:fs";
import { resolve } from "node:path";

import { describe, expect, test } from "bun:test";

import { createRichMarkdownConverter } from "../src/index";

function getNativeLibraryPath(): string {
	const libraryName = process.platform === "darwin"
		? "libtgmarkdown.dylib"
		: process.platform === "win32"
			? "tgmarkdown.dll"
			: "libtgmarkdown.so";

	return resolve(import.meta.dir, "../native", libraryName);
}

function readSampleFile(fileName: string): string {
	return readFileSync(resolve(import.meta.dir, "../../../tests/samples", fileName), "utf8");
}

describe("tg-rich-markdown Koffi binding", () => {
	const convert = createRichMarkdownConverter({ nativeLibraryPath: getNativeLibraryPath() });

	test("converts superscript and subscript", () => {
		expect(convert("H^2^ + H~2~")).toBe("H<sup>2</sup> + H<sub>2</sub>");
	});

	test("resolves the native library without an explicit path", () => {
		const automaticConverter = createRichMarkdownConverter();

		expect(automaticConverter("H^2^")).toBe("H<sup>2</sup>");
	});

	test("preserves code blocks", () => {
		const input = "```go\nvalue := H^2^\n```";
		expect(convert(input)).toBe(input);
	});

	test("preserves the official rich markdown sample", () => {
		const input = readSampleFile("official-rich-markdown.md");
		expect(convert(input)).toBe(input);
	});

	test("converts the CommonMark GFM sample to expected Rich Markdown", () => {
		const input = readSampleFile("CommonMark_GFM_Extensions.md");
		const expected = readSampleFile("CommonMark_GFM_Extensions.expected-rich.md");

		expect(convert(input)).toBe(expected);
	});
});
