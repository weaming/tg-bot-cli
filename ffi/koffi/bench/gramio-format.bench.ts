import { readFileSync } from "node:fs";
import { resolve } from "node:path";

import { markdownToFormattable } from "@gramio/format/markdown";
import { rich } from "@gramio/format/rich";

const ITERATIONS = 10_000;
const WARMUP_ITERATIONS = 1_000;
const BYTES_PER_MEBIBYTE = 1024 * 1024;

interface BenchmarkCase {
	name: string;
	path: string;
}

interface BenchmarkResult {
	name: string;
	bytes: number;
	nanosecondsPerOperation: number;
	mebibytesPerSecond: number;
	outputLength: number;
}

const benchmarkCases: BenchmarkCase[] = [
	{
		name: "CommonMark_GFM_Extensions",
		path: resolve(import.meta.dir, "../../../tests/samples/CommonMark_GFM_Extensions.md"),
	},
	{
		name: "official-rich-markdown",
		path: resolve(import.meta.dir, "../../../tests/samples/official-rich-markdown.md"),
	},
];

function convertMarkdownToRichMarkdown(markdown: string): string {
	const formattable = markdownToFormattable(markdown);
	return rich([formattable]).markdown;
}

function runBenchmarkCase(caseDefinition: BenchmarkCase): BenchmarkResult {
	const input = readFileSync(caseDefinition.path, "utf8");

	for (let iteration = 0; iteration < WARMUP_ITERATIONS; iteration++) {
		convertMarkdownToRichMarkdown(input);
	}

	let outputLength = 0;
	const startTime = Bun.nanoseconds();
	for (let iteration = 0; iteration < ITERATIONS; iteration++) {
		outputLength += convertMarkdownToRichMarkdown(input).length;
	}
	const elapsedNanoseconds = Bun.nanoseconds() - startTime;
	const nanosecondsPerOperation = elapsedNanoseconds / ITERATIONS;
	const mebibytesPerSecond = input.length / (nanosecondsPerOperation / 1e9) / BYTES_PER_MEBIBYTE;

	return {
		name: caseDefinition.name,
		bytes: input.length,
		nanosecondsPerOperation,
		mebibytesPerSecond,
		outputLength,
	};
}

function printResult(result: BenchmarkResult): void {
	console.log(
		`${result.name}: ${result.nanosecondsPerOperation.toFixed(2)} ns/op, `
			+ `${result.mebibytesPerSecond.toFixed(2)} MiB/s, `
			+ `${result.bytes} input bytes, ${result.outputLength} output bytes`,
	);
}

for (const caseDefinition of benchmarkCases) {
	printResult(runBenchmarkCase(caseDefinition));
}
