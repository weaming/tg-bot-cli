# tg-bot-cli

Telegram Bot API 命令行工具。

## 安装

```bash
# 开发用（通过环境变量传 token）
make build

# 锁定 token 编译
make install TOKEN=123:AAA... BINARY=tg

# 锁定 token + 目标 chat
make install TOKEN=123:AAA... TARGET=@mychannel BINARY=tg-chan

# 多 bot / 多目标场景
make install TOKEN=123:AAA... TARGET=-100123456 BINARY=tg-work
make install TOKEN=456:BBB... TARGET=@personal  BINARY=tg-personal
```

## 优先级

Token：`--token` > `TG_BOT_TOKEN` > 编译内置

Target：`--to`/`--chat` > `TG_TARGET` > 编译内置

## 命令

```bash
tg send   --to <chat>  --text "内容"            # 发文本
tg send   --to <chat>  --text -                 # 从 stdin 读取
tg send   --to <chat>  --file ./photo.jpg       # 发媒体文件
tg edit   --chat <chat> --msg <id> --text "新内容"
tg delete --chat <chat> --msg <id>
tg forward --from <chat> --to <chat> --msg <id>
tg copy    --from <chat> --to <chat> --msg <id>
tg pin     --chat <chat> --msg <id>
tg unpin   --chat <chat> [--msg <id>]
tg info    --chat <chat>
tg me
```

## 常用选项

| 选项                  | 说明                                                |
| --------------------- | --------------------------------------------------- |
| `--parse-mode`        | HTML \| Markdown \| MarkdownV2（默认 HTML）         |
| `--reply-to <id>`     | 回复某条消息                                        |
| `--silent`            | 静默发送                                            |
| `--protect`           | 防止转发/保存                                       |
| `--thread <id>`       | 话题群消息                                          |
| `--button "文字:URL"` | Inline 按钮（多次使用 = 多行，`,` 同行，`\|` 换行） |
| `--json`              | 输出完整 JSON 响应                                  |

## 示例

```bash
# 发送 HTML 消息
tg send --to @mychannel --text "<b>加粗</b> 内容"

# 带按钮
tg send --to @mychannel --text "点击" \
  --button "官网:https://example.com,文档:https://docs.example.com" \
  --button "GitHub:https://github.com"

# 发图片
tg send --to 123456789 --file ./image.png --caption "说明"

# 从管道发送
echo "定时任务完成" | tg send --to @alerts --text -

# 获取 JSON 格式结果
tg info --chat @mychannel --json
```
