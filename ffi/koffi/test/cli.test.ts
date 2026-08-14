import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

import { describe, expect, test } from "bun:test";

const cliPath = resolve(import.meta.dir, "../src/cli.ts");
const nativeLibraryPath = resolve(import.meta.dir, "../native", "libtgmarkdown.dylib");

function runCli(args: string[], input = "", nativeLibraryPath?: string) {
	const environment = {
		...process.env,
		...(nativeLibraryPath ? { TG_MARKDOWN_LIBRARY: nativeLibraryPath } : {}),
	};

	return spawnSync(process.execPath, [cliPath, ...args], {
		encoding: "utf8",
		input,
		env: environment,
	});
}

describe("tg-rich-markdown CLI", () => {
	test("converts stdin to Rich Markdown", () => {
		const result = runCli([], "H^2^ + H~2~", nativeLibraryPath);

		expect(result.status).toBe(0);
		expect(result.stdout).toBe("H<sup>2</sup> + H<sub>2</sub>");
		expect(result.stderr).toBe("");
	});

	test("prints help without loading the native library", () => {
		const result = runCli(["--help"]);

		expect(result.status).toBe(0);
		expect(result.stdout).toContain("Usage: tg-rich-markdown");
	});
});
