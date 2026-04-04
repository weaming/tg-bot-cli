package util

import (
	"io"
	"os"
)

// ReadTextOrStdin 当 text 为 "-" 时从 stdin 读取，否则直接返回
func ReadTextOrStdin(text string) (string, error) {
	if text != "-" {
		return text, nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
