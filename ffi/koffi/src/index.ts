import { existsSync } from "node:fs";
import { createRequire } from "node:module";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import koffi from "koffi";

export interface RichMarkdownConverterOptions {
	nativeLibraryPath?: string;
}

export type RichMarkdownConverter = (markdown: string) => string;

type LoadedLibrary = ReturnType<typeof koffi.load>;

interface NativeBindings {
	convert: (markdown: string) => unknown;
	library: LoadedLibrary;
}

const nativeBindingsCache = new Map<string, NativeBindings>();
const nodeRequire = createRequire(import.meta.url);

const NATIVE_PACKAGE_BY_TARGET: Record<string, string> = {
	"darwin-arm64": "@weaming/tg-rich-markdown-darwin-arm64",
	"darwin-x64": "@weaming/tg-rich-markdown-darwin-x64",
};

function getNativeLibraryNames(): string[] {
	switch (process.platform) {
		case "darwin":
			return ["libtgmarkdown.dylib"];
		case "win32":
			return ["tgmarkdown.dll"];
		default:
			return ["libtgmarkdown.so"];
	}
}

function getRuntimeTarget(): string {
	return `${process.platform}-${process.arch}`;
}

function resolvePackagedLibraryPath(): string | undefined {
	const packageName = NATIVE_PACKAGE_BY_TARGET[getRuntimeTarget()];
	if (!packageName) {
		return undefined;
	}

	try {
		const nativeLibraryPath = nodeRequire.resolve(packageName);
		return existsSync(nativeLibraryPath) ? nativeLibraryPath : undefined;
	} catch {
		return undefined;
	}
}

function resolveNativeLibraryPath(nativeLibraryPath?: string): string {
	const configuredPath = nativeLibraryPath?.trim() || process.env.TG_MARKDOWN_LIBRARY?.trim();
	if (configuredPath) {
		return resolve(configuredPath);
	}

	const packagedLibraryPath = resolvePackagedLibraryPath();
	if (packagedLibraryPath) {
		return packagedLibraryPath;
	}

	const packageDirectory = resolve(dirname(fileURLToPath(import.meta.url)), "..");
	const candidatePaths = getNativeLibraryNames().flatMap((libraryName) => [
		resolve(process.cwd(), libraryName),
		resolve(packageDirectory, "native", libraryName),
	]);
	const existingPath = candidatePaths.find((candidatePath) => existsSync(candidatePath));
	if (existingPath) {
		return existingPath;
	}

	throw new Error(
		`tg-rich-markdown native library not found for ${getRuntimeTarget()}; `
			+ "install a supported platform package, set TG_MARKDOWN_LIBRARY, "
			+ "or run make ffi-rich-markdown first",
	);
}

function loadNativeBindings(nativeLibraryPath: string): NativeBindings {
	const cachedBindings = nativeBindingsCache.get(nativeLibraryPath);
	if (cachedBindings) {
		return cachedBindings;
	}

	let library: LoadedLibrary;
	try {
		library = koffi.load(nativeLibraryPath);
	} catch (error) {
		throw new Error(`Failed to load tg-rich-markdown native library: ${nativeLibraryPath}`, {
			cause: error,
		});
	}

	const freeNativeString = library.func("void FreeMarkdownString(char *)");
	const richMarkdownString = koffi.disposable(
		"RichMarkdownString",
		"char *",
		freeNativeString,
	);
	const convertNative = library.func(
		"RichMarkdownString ConvertMarkdownToRichMarkdown(const char *)",
	);
	const bindings: NativeBindings = {
		convert: (markdown: string) => convertNative(markdown),
		library,
	};

	nativeBindingsCache.set(nativeLibraryPath, bindings);
	return bindings;
}

export function createRichMarkdownConverter(
	options: RichMarkdownConverterOptions = {},
): RichMarkdownConverter {
	const nativeLibraryPath = resolveNativeLibraryPath(options.nativeLibraryPath);
	const bindings = loadNativeBindings(nativeLibraryPath);

	return (markdown: string): string => {
		if (typeof markdown !== "string") {
			throw new TypeError("markdown must be a string");
		}

		const result = bindings.convert(markdown);
		if (typeof result !== "string") {
			throw new Error("tg-rich-markdown native library returned an invalid result");
		}

		return result;
	};
}

export function convertMarkdownToRichMarkdown(
	markdown: string,
	options: RichMarkdownConverterOptions = {},
): string {
	return createRichMarkdownConverter(options)(markdown);
}
