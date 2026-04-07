package main

import "github.com/weaming/tg-bot-cli/cmd/tg/cmd"

// 编译时注入示例：
//
//	go build -ldflags "-X 'main.CompiledToken=<token>' -X 'main.CompiledTarget=<chat_id>'" -o tg .
var (
	CompiledToken  string
	CompiledTarget string
)

func main() {
	cmd.CompiledToken = CompiledToken
	cmd.CompiledTarget = CompiledTarget
	cmd.Execute()
}
