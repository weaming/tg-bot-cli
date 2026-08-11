package parser

import (
	"bytes"
	"strings"
)

type richHTMLTagFrame struct {
	name string
}

var richVoidHTMLTags = map[string]struct{}{
	"area": {}, "base": {}, "br": {}, "col": {}, "embed": {}, "hr": {},
	"img": {}, "input": {}, "link": {}, "meta": {}, "param": {}, "source": {},
	"track": {}, "wbr": {},
}

func richHTMLTagName(tag string) string {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "/")
	tag = strings.TrimSpace(tag)
	endIndex := 0
	for endIndex < len(tag) {
		character := tag[endIndex]
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			character == ':' || character == '-' || character == '_' {
			endIndex++
			continue
		}
		break
	}
	return strings.ToLower(tag[:endIndex])
}

func findRichHTMLTagEnd(rawHTML string, start int) int {
	var quote byte
	for index := start; index < len(rawHTML); index++ {
		character := rawHTML[index]
		if quote != 0 {
			if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		if character == '>' {
			return index
		}
	}
	return -1
}

// updateRichHTMLTagStack 根据 HTML 标签维护当前打开的标签栈。
func updateRichHTMLTagStack(rawHTML string, openHTMLTags *[]richHTMLTagFrame) {
	for searchStart := 0; searchStart < len(rawHTML); {
		openIndex := strings.IndexByte(rawHTML[searchStart:], '<')
		if openIndex < 0 {
			return
		}
		openIndex += searchStart
		closeIndex := findRichHTMLTagEnd(rawHTML, openIndex+1)
		if closeIndex < 0 {
			return
		}

		tag := strings.TrimSpace(rawHTML[openIndex+1 : closeIndex])
		searchStart = closeIndex + 1
		if tag == "" || strings.HasPrefix(tag, "!") || strings.HasPrefix(tag, "?") {
			continue
		}

		isClosingTag := strings.HasPrefix(tag, "/")
		tagName := richHTMLTagName(tag)
		if tagName == "" {
			continue
		}

		if isClosingTag {
			for index := len(*openHTMLTags) - 1; index >= 0; index-- {
				if (*openHTMLTags)[index].name == tagName {
					*openHTMLTags = (*openHTMLTags)[:index]
					break
				}
			}
			continue
		}

		if strings.HasSuffix(tag, "/") {
			continue
		}
		if _, isVoidTag := richVoidHTMLTags[tagName]; isVoidTag {
			continue
		}
		*openHTMLTags = append(*openHTMLTags, richHTMLTagFrame{
			name: tagName,
		})
	}
}

// findRichHTMLCodeClosingTag 查找大小写不敏感且支持属性空格的 code 闭合标签。
func findRichHTMLCodeClosingTag(source []byte, start int) (int, int) {
	for searchStart := start; searchStart < len(source); {
		relativeStart := bytes.IndexByte(source[searchStart:], '<')
		if relativeStart < 0 {
			return -1, -1
		}
		closingStart := searchStart + relativeStart
		if closingStart+1 >= len(source) || source[closingStart+1] != '/' {
			searchStart = closingStart + 1
			continue
		}
		closingEnd := findRichHTMLTagEndBytes(source, closingStart+2)
		if closingEnd < 0 {
			return -1, -1
		}
		if richHTMLTagNameMatches(source[closingStart+2:closingEnd], "code") {
			return closingStart, closingEnd
		}
		searchStart = closingEnd + 1
	}
	return -1, -1
}

// findRichHTMLTagEndBytes 查找不受引号内大于号影响的标签结束位置。
func findRichHTMLTagEndBytes(rawHTML []byte, start int) int {
	var quote byte
	for index := start; index < len(rawHTML); index++ {
		character := rawHTML[index]
		if quote != 0 {
			if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		if character == '>' {
			return index
		}
	}
	return -1
}

// richHTMLTagNameMatches 比较标签名并忽略 ASCII 大小写。
func richHTMLTagNameMatches(tag []byte, expected string) bool {
	tag = bytes.TrimSpace(tag)
	endIndex := 0
	for endIndex < len(tag) {
		character := tag[endIndex]
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			character == ':' || character == '-' || character == '_' {
			endIndex++
			continue
		}
		break
	}
	if endIndex != len(expected) {
		return false
	}
	for index := 0; index < endIndex; index++ {
		character := tag[index]
		if character >= 'A' && character <= 'Z' {
			character += 'a' - 'A'
		}
		if character != expected[index] {
			return false
		}
	}
	return true
}
