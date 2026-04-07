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

Token：`--token` / `-T` > `TG_BOT_TOKEN` > 编译内置

Target：`--to` / `--chat` > `TG_TARGET` > 编译内置

Proxy：`--proxy` / `-x` > `TG_PROXY` > `HTTPS_PROXY` / `HTTP_PROXY`

## 命令

```bash
tg config                                        # 查看当前配置来源
tg send   -t <chat>  -m "内容"                  # 发文本
tg send   -t <chat>  -m -                        # 从 stdin 读取
tg send   -t <chat>  -f ./photo.jpg -c "说明"   # 发媒体文件
tg edit   -c <chat> -m <id> -t "新内容"
tg delete -c <chat> -m <id>
tg forward -f <from> -t <to> -m <id>
tg copy    -f <from> -t <to> -m <id>
tg pin     -c <chat> -m <id>
tg unpin   -c <chat> [-m <id>]
tg info    -c <chat>
tg me
```

## 常用选项

| 选项           | 简写 | 说明                                                |
| -------------- | ---- | --------------------------------------------------- |
| `--to`         | `-t` | 目标 chat，username 可省略 `@`                      |
| `--chat`       | `-c` | 目标 chat（edit/delete/pin/info）                   |
| `--text`       | `-m` | 消息文本，`-` 读 stdin                              |
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

# 带按钮
tg send -t mychannel -m "点击" \
  -b "官网:https://example.com,文档:https://docs.example.com" \
  -b "GitHub:https://github.com"

# 发图片
tg send -t 123456789 -f ./image.png -c "说明"

# 从管道发送
echo "定时任务完成" | tg send -t alerts -m -

# 获取 JSON 格式结果
tg info -c mychannel -j

# 检查编译内置情况
tg config
```
