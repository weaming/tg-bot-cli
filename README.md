# tg-bot-cli

Telegram Bot API 命令行工具。

## 安装

需要 Go 1.25+，并确保 `$(go env GOPATH)/bin` 在 PATH 中。

```bash
# 锁定 token 安装到 GOPATH/bin
make install-tg TOKEN=123:AAA... BINARY=tg

# 锁定 token + 目标 chat
make install-tg TOKEN=123:AAA... TARGET=@mychannel BINARY=tg-chan

# 多 bot / 多目标场景
make install-tg TOKEN=123:AAA... TARGET=-100123456 BINARY=tg-work
make install-tg TOKEN=456:BBB... TARGET=@personal  BINARY=tg-personal

# 单独的程序转换 std markdown -> tg HTML
make install-md2tg
```

参考 [md2tg readme](cmd/md2tg/README.md)

## 优先级

- Token：`--token` / `-T` > `TG_BOT_TOKEN` > 编译内置
- Target：`--to` / `--chat` > `TG_TARGET` > 编译内置
- Proxy：`--proxy` / `-x` > `TG_PROXY` > `HTTPS_PROXY` / `HTTP_PROXY`

## 命令

```bash
tg config                                        # 查看当前配置来源
tg send   -t <chat>  -m "内容"                  # 发文本
tg send   -t <chat>  -i <file>                  # 从文件读取
tg send   -t <chat>  -i - --md2html            # 从 stdin 读取并转 HTML
tg send   -t <chat>  -f ./photo.jpg -c "说明"  # 发媒体文件
tg edit   -c <chat> -m <id> -t "新内容"
tg delete -c <chat> -m <id>
tg forward -f <from> -t <to> -m <id>
tg copy    -f <from> -t <to> -m <id>
tg pin     -c <chat> -m <id>
tg unpin   -c <chat> [-m <id>]
tg info    -c <chat>
tg me
```

## 消息格式

普通文本通过 Telegram 的 `sendMessage` 发送，支持原有的 3 种格式：

| 格式       | 说明                                                 |
| ---------- | ---------------------------------------------------- |
| Markdown   | 旧版 Markdown，语法简单但功能有限                    |
| MarkdownV2 | 功能更多，但需要严格处理转义                         |
| HTML       | 使用 HTML 标签表达格式，`--md2html` 默认转换为此格式 |

默认不指定 `--rich` 时，继续使用 legacy 的普通消息发送路径，以保持已有脚本和调用方式向后兼容。使用 Rich Message 时，Rich Markdown 已足够渲染大多数 Markdown 内容，建议默认使用 `--rich=md`；需要更强的结构化语法或媒体能力时再使用 `--rich=html`。

此外，`tg send` 支持 Telegram 的 2 种 Rich Message 格式。它们通过 `sendRichMessage` 发送，不使用 `--parse-mode`：

| 选项          | 格式          | 说明                                                                 |
| ------------- | ------------- | -------------------------------------------------------------------- |
| `--rich=md`   | Rich Markdown | 适合大多数 Markdown 内容，并将上下标等扩展转换为可嵌入的 HTML 标签   |
| `--rich=html` | Rich HTML     | 输出结构化内容，支持原生标题、表格、列表、公式和媒体，表达能力最完整 |

Rich Message 还支持标记文本、剧透、任务列表、脚注、Telegram 媒体链接，以及 `<tg-map>`、`<details>`、`<tg-collage>` 等 Rich HTML 标签。输入 `.md` 文件或显式指定 `--md2html` 时，会先进行 Markdown 转换；`--rich` 不能与 `--parse-mode` 同时使用。

参考 [Telegram Rich Messages 官方文档](https://core.telegram.org/bots/api#rich-messages)。

## 常用选项

| 选项           | 简写 | 说明                                                |
| -------------- | ---- | --------------------------------------------------- |
| `--to`         | `-t` | 目标 chat，username 可省略 `@`                      |
| `--chat`       | `-c` | 目标 chat（edit/delete/pin/info）                   |
| `--text`       | `-m` | 消息文本                                            |
| `--input-file` | `-i` | 从文件或 stdin（-）读取                             |
| `--md2html`    | —    | 将 Markdown 转换为普通 HTML 或 Rich Message 内容    |
| `--rich`       | —    | 使用 Rich Message：`html` 或 `md`                   |
| `--file`       | `-f` | 文件路径（自动识别类型）                            |
| `--caption`    | `-c` | 文件说明                                            |
| `--msg`        | `-m` | 消息 ID                                             |
| `--from`       | `-f` | 来源 chat（forward/copy）                           |
| `--reply-to`   | `-r` | 回复的消息 ID                                       |
| `--silent`     | `-s` | 静默发送                                            |
| `--protect`    | `-p` | 防止转发/保存                                       |
| `--thread`     | —    | 话题群 message_thread_id                            |
| `--button`     | `-b` | Inline 按钮（多次使用 = 多行，`,` 同行，`\|` 换行） |
| `--parse-mode` | —    | 默认不传（纯文本）；可选 HTML \| MarkdownV2         |
| `--json`       | `-j` | 输出完整 JSON 响应                                  |

## 示例

```bash
# 发送 HTML 消息
tg send -t mychannel -m "<b>加粗</b> 内容" --parse-mode HTML

# 发送 MarkdownV2 消息
tg send -t mychannel -m "这是 *粗体* 和 _斜体_" --parse-mode MarkdownV2

# 从文件读取并自动转 HTML（.md 文件）
tg send -t mychannel -i ./readme.md

# 从 stdin 读取并转 HTML
cat readme.md | tg send -t mychannel -i - --md2html

# 发送 Rich Markdown（.md 文件自动转换）
tg send -t mychannel -i ./rich.md --rich=md

# 发送 Rich HTML
tg send -t mychannel -i ./rich.md --rich=html

# 带按钮
tg send -t mychannel -m "点击" \
  -b "官网:https://example.com,文档:https://docs.example.com" \
  -b "GitHub:https://github.com"

# 发图片
tg send -t 123456789 -f ./image.png -c "说明"

# 从管道发送
echo "定时任务完成" | tg send -t alerts -i -

# 获取 JSON 格式结果
tg info -c mychannel -j

# 检查编译内置情况
tg config
```
