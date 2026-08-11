# md2tg

将 Markdown 转换为 Telegram 消息内容。除 Telegram 原有的 3 种普通消息格式外，还支持 2 种 Rich Message 格式。

## 支持的消息格式

### Telegram 原有的 3 种普通消息格式

这 3 种格式通过 `sendMessage` 的 `parse_mode` 使用：

| 格式       | 特点                                       | 适用场景                                         |
| ---------- | ------------------------------------------ | ------------------------------------------------ |
| Markdown   | 旧版 Markdown 语法简单，但功能和兼容性有限 | 兼容已有的旧消息                                 |
| MarkdownV2 | 格式更多，但转义规则复杂                   | 需要使用 Telegram 普通消息格式，并能处理严格转义 |
| HTML       | 使用 HTML 标签表达格式，嵌套和转义更直观   | `md2tg` 默认输出格式                             |

`md2tg` 默认将 Markdown 转换为 Telegram HTML。它仍然是普通消息格式：表格会转换为文本布局，不能保留 Rich Message 的原生结构。

### 新增的 2 种 Telegram Rich Message 格式

这 2 种格式通过 Telegram 的 `sendRichMessage` 使用，不应再传给普通消息的 `parse_mode`：

| 选项          | 格式          | 说明                                                                                           |
| ------------- | ------------- | ---------------------------------------------------------------------------------------------- |
| `--rich=md`   | Rich Markdown | 适合大多数 Markdown 内容；仅将 Rich Markdown 不支持的 `^上标^` 和 `~下标~` 转换为 HTML 标签 |
| `--rich=html` | Rich HTML     | 将 Markdown 转换为结构化 Rich HTML，保留原生标题、表格、列表、公式和媒体结构，表达能力最完整   |

Rich Message 支持的内容包括：

- 标题、粗体、斜体、删除线、下划线、上标、下标和行内代码
- 标记文本 `==内容==` 与剧透 `||内容||`
- 行内公式 `$x^2$` 与块公式 `$$E = mc^2$$`
- 有序列表、无序列表、任务列表、引用和脚注
- 原生表格及列对齐；Rich HTML 还支持 `bordered`、`striped` 等表格属性
- 图片、视频、音频、动画，以及 `tg://emoji`、`tg://time` 等 Telegram 媒体链接
- Rich HTML 标签，例如 `<tg-spoiler>`、`<tg-map>`、`<details>`、`<tg-collage>` 和 `<tg-slideshow>`

Rich Markdown 已足够渲染大多数 Markdown 内容。需要更强的结构化语法和媒体能力时，再使用 Rich HTML。

## 使用方法

普通 HTML 转换：

```bash
echo '# Hello' | md2tg
cat README.md | md2tg
md2tg < input.md
```

Rich Message 转换：

```bash
md2tg --rich=md < input.md
md2tg --rich=html < input.md
```

输入始终从标准输入读取，输出到标准输出。

## 转换选项

| 选项            | 说明                                                     |
| --------------- | -------------------------------------------------------- |
| `--split-table` | 将普通消息中的 Markdown 表格转换为多行 `key: value` 布局 |
| `--rich=md`     | 输出 Rich Markdown                                       |
| `--rich=html`   | 输出 Rich HTML                                           |

`--split-table` 只适用于默认的普通 HTML 转换，不能与 `--rich` 同时使用。`--rich` 的值只能是 `md` 或 `html`。

## 普通消息支持的内容

- 标题转换为加粗文本
- 粗体、斜体和删除线
- 带语言标识的代码块：`<pre><code class="language-xxx">...</code></pre>`
- 支持嵌套的有序列表和无序列表
- GFM 表格（可使用 `--split-table`）
- `http`、`https` 和 `tg` 链接
- 图片转换为 `<image_url>` 标签
- 引用和分隔线
- LaTeX 块公式转换为代码块
