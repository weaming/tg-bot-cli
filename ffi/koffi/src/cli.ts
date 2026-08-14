#!/usr/bin/env node

import { readFileSync, realpathSync } from "node:fs";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { convertMarkdownToRichMarkdown } from "./index";

const CLI_VERSION = "0.1.0";

interface CliOptions {
	inputPath?: string;
	showHelp: boolean;
	showVersion: boolean;
}

function getUsage(): string {
	return `Usage: tg-rich-markdown [options] [input]

Convert Markdown from a file or stdin to Telegram Rich Markdown.

Options:
  -h, --help        Show this help message
  -v, --version     Show the version

Input:
  input             Markdown file path; omit it or use - to read stdin
`;
}

export function parseArguments(args: string[]): CliOptions {
	const options: CliOptions = {
		showHelp: false,
		showVersion: false,
	};

	for (let index = 0; index < args.length; index++) {
		const argument = args[index];
		if (argument === "-h" || argument === "--help") {
			options.showHelp = true;
			continue;
		}
		if (argument === "-v" || argument === "--version") {
			options.showVersion = true;
			continue;
		}
		if (argument.startsWith("-")) {
			throw new Error(`Unknown option: ${argument}`);
		}
		if (options.inputPath) {
			throw new Error("Only one input path can be provided");
		}

		options.inputPath = argument;
	}

	return options;
}

function readInput(inputPath?: string): string {
	if (!inputPath || inputPath === "-") {
		return readFileSync(0, "utf8");
	}

	return readFileSync(resolve(inputPath), "utf8");
}

export function run(args: string[] = process.argv.slice(2)): number {
	try {
		const options = parseArguments(args);
		if (options.showHelp) {
			process.stdout.write(getUsage());
			return 0;
		}
		if (options.showVersion) {
			process.stdout.write(`${CLI_VERSION}\n`);
			return 0;
		}

		const markdown = readInput(options.inputPath);
		const richMarkdown = convertMarkdownToRichMarkdown(markdown);
		process.stdout.write(richMarkdown);
		return 0;
	} catch (error) {
		const message = error instanceof Error ? error.message : String(error);
		process.stderr.write(`Error: ${message}\n`);
		return 1;
	}
}

function isMainModule(): boolean {
	if (!process.argv[1]) {
		return false;
	}

	try {
		return realpathSync(resolve(process.argv[1])) === realpathSync(fileURLToPath(import.meta.url));
	} catch {
		return false;
	}
}

if (isMainModule()) {
	process.exitCode = run();
}
