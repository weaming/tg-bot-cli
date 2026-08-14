# Markdown 与 Telegram Rich Markdown

## 语法集合

这里将普通 Markdown 约定为 CommonMark，将表格、任务列表、删除线和脚注等扩展归为 GFM。Telegram Rich Markdown 在尽可能兼容 GFM 的基础上，增加了 Telegram 专用语法。

参考：

- [GitHub Flavored Markdown 规范](https://github.github.com/gfm/)
- [Telegram Rich Markdown](https://core.telegram.org/bots/api#rich-markdown-style)

## 交集

以下 Markdown/GFM 语法同时被 Telegram Rich Markdown 支持：

- 段落和标题：`# heading`
- 粗体：`**text**`、`__text__`
- 斜体：`*text*`、`_text_`
- 删除线：`~~text~~`
- 行内代码和代码块：`` `code` ``、```` ```lang ... ``` ````
- 引用：`> quote`
- 无序列表：`- item`、`* item`、`+ item`
- 有序列表：`1. item`
- 任务列表：`- [x] done`
- 链接：`[text](https://example.com)`
- 图片和媒体：`![](https://example.com/media)`
- 自动链接：URL、电子邮件地址等
- 分隔线：`---`
- GFM 表格及列对齐
- GFM 脚注：`[^id]` 和 `[^id]: definition`

媒体必须作为独立块。Telegram 根据 MIME 类型和 URL 判断它是图片、视频、音频、语音消息或动画；媒体标题作为媒体说明文字。

## Telegram 专用语法

以下语法不是 CommonMark/GFM 的标准语法：

```markdown
==marked text==
||spoiler text||
```

```markdown
![Custom emoji](tg://emoji?id=123456789)
![Tomorrow](tg://time?unix=1647531900&format=wDT)
```

```markdown
$x^2$

$$
E = mc^2
$$
```

还可以在 Rich Markdown 中嵌入 Telegram Rich HTML 标签，例如：

```html
<u>underline</u>
<sub>subscript</sub>
<sup>superscript</sup>
<details>...</details>
<aside>...</aside>
<tg-map lat="41.9" long="12.5" zoom="14"/>
<tg-collage>...</tg-collage>
<tg-slideshow>...</tg-slideshow>
```

Telegram 还会自动识别命令、用户名、电话、银行卡号、hashtag 和 cashtag 等实体；这些不是 Markdown 语法。

## Rich Markdown 与 Rich HTML

纯 Rich Markdown 的表达能力弱于 Rich HTML：

```text
纯 Rich Markdown < Rich HTML
Rich Markdown + Telegram HTML 扩展 ≈ Rich HTML
```

Rich HTML 可以直接表达 Rich Markdown 没有对应语法的能力，例如：

- 自定义 Emoji、时间实体
- 下划线、上下标、引用和 footer
- 地图、折叠块、拼贴、幻灯片
- 表格的 `bordered`、`striped`、`rowspan`、`colspan`
- 媒体的 `tg-spoiler`、`figure`、`figcaption`、`cite`
- 有序列表的 `start`、`reversed`、`value`

Rich Markdown 可以通过嵌入 Telegram 支持的 HTML 标签表达其中大部分能力，但这时已经不是纯 Markdown。Telegram 官方也建议对没有 Markdown 语法的功能使用 HTML 标签。

在命令行中：

- `--rich=md`：发送原生 Rich Markdown，适合简洁的 Markdown 文档。
- `--rich=html`：发送 Rich HTML，表达能力最完整。
- `--rich=html --input-file xxx.md`：将 Markdown 转换为 Rich HTML，适合作为统一的完整输出路径。

## Markdown 中未被 Telegram 明确保证的语法

Telegram 文档只承诺“尽可能兼容 GFM”，以下语法不应默认视为稳定的 Rich Markdown 能力：

- Setext 标题：

  ```markdown
  Heading
  =======
  ```

- 缩进代码块
- 引用式链接和图片：

  ```markdown
  [link][id]

  [id]: https://example.com
  ```

- 任意原生 HTML 标签
- HTML 注释
- 其他 Markdown 方言的专用扩展
- 部分转义、实体和换行细节

