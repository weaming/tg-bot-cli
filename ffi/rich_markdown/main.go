package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/weaming/tg-bot-cli/api"
)

// ConvertMarkdownToRichMarkdown 转换 Markdown，并返回由 C malloc 分配的字符串。
// 调用方必须通过 FreeMarkdownString 释放返回值。
//
//export ConvertMarkdownToRichMarkdown
func ConvertMarkdownToRichMarkdown(input *C.char) *C.char {
	if input == nil {
		return C.CString("")
	}

	return C.CString(api.ConvertMarkdownToRichMarkdown(C.GoString(input)))
}

// FreeMarkdownString 释放 ConvertMarkdownToRichMarkdown 返回的字符串。
//
//export FreeMarkdownString
func FreeMarkdownString(value *C.char) {
	if value == nil {
		return
	}

	C.free(unsafe.Pointer(value))
}

func main() {}
