# md2tg

Markdown to Telegram HTML converter (Go implementation)

## Why HTML instead of Markdown/MarkdownV2?

Telegram supports three message formats:

| Format     | Pros                         | Cons                                                        |
| ---------- | ---------------------------- | ----------------------------------------------------------- |
| Markdown   | Simple syntax                | Limited styling, no code blocks with syntax highlighting    |
| MarkdownV2 | More features than Markdown  | Still limited, complex escaping rules, no nested formatting |
| HTML       | Full control over formatting | Requires careful escaping of HTML tags                      |

**Why md2tg uses HTML:**

1. **Code blocks with language class** - Only HTML `<pre><code class="language-xxx">` supports syntax highlighting via Telegram's built-in code block renderer. Markdown/MarkdownV2 cannot specify the language for code blocks.
2. **Image tags** - HTML `<image_url>...</image_url>` provides a clean way to embed images that Telegram renders as inline media.
3. **Consistent escaping** - HTML escaping is straightforward (`<` → `&lt;`, `&` → `&amp;`). MarkdownV2 has complex and inconsistent escaping rules that vary by context.
4. **Nested formatting** - HTML allows clearer nesting of elements that Markdown cannot express.

## 为什么选择 HTML 而不是 Markdown/MarkdownV2？

Telegram 支持三种消息格式：

| 格式       | 优点                 | 缺点                                   |
| ---------- | -------------------- | -------------------------------------- |
| Markdown   | 语法简单             | 样式有限，不支持代码块语法高亮         |
| MarkdownV2 | 比 Markdown 功能更多 | 仍然有限，转义规则复杂，不支持嵌套格式 |
| HTML       | 完全控制格式         | 需要小心转义 HTML 标签                 |

**为什么 md2tg 使用 HTML：**

1. **带语言类的代码块** - 只有 HTML `<pre><code class="language-xxx">` 支持通过 Telegram 内置代码块渲染器进行语法高亮。Markdown/MarkdownV2 无法指定代码块的语言。
2. **图片标签** - HTML `<image_url>...</image_url>` 提供了一种干净的方式来嵌入 Telegram 渲染为内联媒体的图片。
3. **一致的转义** - HTML 转义简单直接（`<` → `&lt;`，`&` → `&amp;`）。MarkdownV2 的转义规则复杂且因上下文而异。
4. **嵌套格式** - HTML 允许更清晰地嵌套 Markdown 无法表达的格式。

## Usage

```bash
echo "# Hello" | md2tg
cat README.md | md2tg
md2tg < input.md
md2tg -split-table < table.md  # tables as key:value format
```

## Features

- Headings → `<b>` bold
- Bold/italic → `<b>`/`<i>`
- Code blocks → `<pre><code>` (with language class)
- Lists (bullet/numbered) with nested support
- **Tables** (GFM)
- Links (http/https/tg only, mailto rejected)
- Images → `<image_url>` tags (with URL query string support)
- Blockquotes
- Thematic breaks
- LaTeX (`$$...$$`) → code blocks
