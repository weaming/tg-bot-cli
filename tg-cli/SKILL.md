---
name: tg-cli
description: 通过 Telegram Bot CLI 发送消息、文件，编辑/删除/转发/置顶消息，查询 chat 信息
allowed-tools: [Bash]
---

# Telegram Bot CLI 技能

使用本地 `tg` 二进制与 Telegram Bot API 交互。

## Token / 目标配置

优先级（由高到低）：

| 项目      | CLI flag                         | 环境变量                   | 编译内置 |
| --------- | -------------------------------- | -------------------------- | -------- |
| Bot Token | `-T` / `--token`                 | `TG_BOT_TOKEN`             | ldflags  |
| 目标 chat | `-t` / `--to` 或 `-c` / `--chat` | `TG_TARGET`                | ldflags  |
| 代理      | `-x` / `--proxy`                 | `TG_PROXY` → `HTTPS_PROXY` | —        |

username 可省略 `@`，自动补全。

## 检查当前配置

```bash
tg config          # 显示 token/target/proxy 来源及脱敏值
tg -T <tok> config # 临时指定 token 后查看
```

## 命令速查

```bash
# 发文本
tg send -t <chat> -m "内容"

# 从 stdin 发送
echo "内容" | tg send -t <chat> -m -

# 发文件（自动识别类型）
tg send -t <chat> -f ./photo.jpg -c "说明"

# 带 inline 按钮（, 同行 | 换行）
tg send -t <chat> -m "点击" -b "文字:https://url"

# 编辑消息
tg edit -c <chat> -m <msg_id> -t "新内容"

# 删除消息
tg delete -c <chat> -m <msg_id>

# 转发（保留来源）
tg forward -f <from_chat> -t <to_chat> -m <msg_id>

# 复制（不显示来源）
tg copy -f <from_chat> -t <to_chat> -m <msg_id>

# 置顶 / 取消置顶
tg pin   -c <chat> -m <msg_id>
tg unpin -c <chat> [-m <msg_id>]

# 查询 chat 信息
tg info -c <chat>

# 查询 bot 自身信息
tg me
```

## 常用 flag

| flag           | 简写 | 说明                                  |
| -------------- | ---- | ------------------------------------- |
| `--to`         | `-t` | 目标 chat（send/forward/copy）        |
| `--chat`       | `-c` | 目标 chat（edit/delete/pin/info）     |
| `--text`       | `-m` | 消息文本，`-` 读 stdin                |
| `--file`       | `-f` | 文件路径                              |
| `--caption`    | `-c` | 文件说明（send/copy）                 |
| `--msg`        | `-m` | 消息 ID                               |
| `--from`       | `-f` | 来源 chat（forward/copy）             |
| `--reply-to`   | `-r` | 回复的消息 ID                         |
| `--silent`     | `-s` | 静默发送                              |
| `--protect`    | `-p` | 防转发/保存                           |
| `--button`     | `-b` | Inline 按钮，可多次使用               |
| `--parse-mode` | —    | 默认不传（纯文本）；可选 HTML \| Markdown \| MarkdownV2 |
| `--thread`     | —    | 话题群 message_thread_id              |
| `--json`       | `-j` | 输出完整 JSON                         |

## 文件类型自动识别

| 扩展名          | API 方法      |
| --------------- | ------------- |
| jpg/png/webp 等 | sendPhoto     |
| mp4/mov 等      | sendVideo     |
| mp3/ogg 等      | sendAudio     |
| gif             | sendAnimation |
| 其他            | sendDocument  |
