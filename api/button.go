package api

import (
	"fmt"
	"strings"
)

// ParseButtons 将按钮字符串列表解析为 InlineKeyboardMarkup。
//
// 每个字符串代表一行，行内用 "," 分隔多个按钮，
// 每个按钮格式为 "文字:URL"，也支持用 "|" 作为行分隔符。
//
// 示例：
//
//	[]string{"按钮1:https://a.com,按钮2:https://b.com", "按钮3:https://c.com"}
//	[]string{"按钮1:https://a.com|按钮2:https://b.com"}
func ParseButtons(buttonRows []string) (*InlineKeyboardMarkup, error) {
	if len(buttonRows) == 0 {
		return nil, nil
	}

	var keyboard [][]InlineKeyboardButton

	for _, rowStr := range buttonRows {
		// 支持 "|" 作为额外的行分隔符
		rows := strings.Split(rowStr, "|")

		for _, row := range rows {
			var rowButtons []InlineKeyboardButton

			cells := strings.Split(row, ",")
			for _, cell := range cells {
				cell = strings.TrimSpace(cell)
				if cell == "" {
					continue
				}

				idx := strings.Index(cell, ":")
				if idx <= 0 {
					return nil, fmt.Errorf("按钮格式错误（应为 '文字:URL'）: %q", cell)
				}

				text := cell[:idx]
				url := cell[idx+1:]
				rowButtons = append(rowButtons, InlineKeyboardButton{
					Text: text,
					URL:  url,
				})
			}

			if len(rowButtons) > 0 {
				keyboard = append(keyboard, rowButtons)
			}
		}
	}

	if len(keyboard) == 0 {
		return nil, nil
	}
	return &InlineKeyboardMarkup{InlineKeyboard: keyboard}, nil
}
