# Rich Markdown 全类型样例

## Inline formatting

**bold** and __bold__

*italic* and _italic_

~~strikethrough~~

==marked text==

||spoiler text||

`inline fixed-width code`

<u>underlined text</u>, <ins>inserted text</ins>, <sub>subscript</sub>, <sup>superscript</sup>

## Links and detected entities

[URL](https://t.me/)

[email](mailto:user@example.com)

[phone](tel:+123456789)

[mention](tg://user?id=123456789)

[#chapter-1](#chapter-1)

#hashtag $USD +12345678901 4242 4242 4242 4242

https://t.me/ t.me a@t.me /command @username

![👍](tg://emoji?id=5368324170671202286)

![22:45 tomorrow](tg://time?unix=1647531900&format=wDT)

## Headings

# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6

## Paragraph and code

Paragraph text with **nested _italic_ formatting**.

```python
def hello(name):
    return f"Hello, {name}!"
```

```math
E = mc^2
```

---

## Lists

- unordered item
* unordered item
+ unordered item

1. ordered item
2. ordered item

---

1. ordered item with same seq
1. ordered item with same seq
1. ordered item with same seq

- [ ] unchecked checkbox
- [x] checked checkbox

- Nested item
  - Nested child
    1. Nested ordered child

## Block quotation

> Block quotation started
>
> Block quotation continued on the next line
> Block quotation continued on the same line
>
> The last line of the blockquote

## Media

![](https://telegram.org/img/t_logo.png)
![](https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4)
![](https://interactive-examples.mdn.mozilla.net/media/cc0-audio/t-rex-roar.mp3)
![](https://upload.wikimedia.org/wikipedia/commons/2/2c/Rotating_earth_%28large%29.gif)

![Photo caption](https://telegram.org/img/t_logo.png "Photo caption")
![Video caption](https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4 "Video caption")
![Audio caption](https://interactive-examples.mdn.mozilla.net/media/cc0-audio/t-rex-roar.mp3 "Audio caption")

## Table

| Header 1 | Header 2 | Header 3 |
|:---------|:--------:|---------:|
| left     | center   | right    |
| **bold** | <tg-spoiler>ready</tg-spoiler> | `42` |

### Wide table: 12 columns and 8 rows

| C01 | C02 | C03 | C04 | C05 | C06 | C07 | C08 | C09 | C10 | C11 | C12 |
|:----|----:|:----:|:----:|:----|:----:|:----|:----|:----|:----:|----:|:----|
| R1-01 | R1-02 | R1-03 | R1-04 | R1-05 | R1-06 | R1-07 | R1-08 | R1-09 | R1-10 | R1-11 | R1-12 |
| Alpha | 1024 | Ready | 98.76 | 2026-08-11 | HK+8 | Short | Medium length value | Very long content that should test cell wrapping | OK | 0.42 | End |
| Beta | 2048 | Pending | 12.34 | 2026-08-12 | UTC | North | East | Longer descriptive value for width testing | WAIT | 0.84 | End |
| Gamma | 4096 | Failed | 0.01 | 2026-08-13 | GMT | South | West | Another intentionally wide table cell | ERROR | 1.00 | End |
| Delta | 8192 | Ready | 456.78 | 2026-08-14 | HK+8 | One | Two | [https://t.me/](https://t.me/) | OK | 2.56 | End |
| Epsilon | 16384 | Pending | 789.01 | 2026-08-15 | UTC | Three | Four | Keep this sentence long enough to exceed a narrow column | WAIT | 3.14 | End |
| Zeta | 32768 | Ready | 999.99 | 2026-08-16 | GMT | Five | Six | Final long value for the overwidth rendering test | OK | 6.28 | End |

## Footnotes

Text with a reference[^id1] and another one[^id2].

[^id1]: Definition of the first footnote.
[^id2]: Definition with _italic text_ and <u>HTML underline</u>.

## Formula

Inline formula: $x^2 + y^2$.

Block formula:

$$
E = mc^2
$$

## Rich HTML extensions inside Markdown

<a name="chapter-1"></a>

<tg-reference name="note-1">Referenced text</tg-reference>

<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>

<tg-time unix="1647531900" format="wDT">22:45 tomorrow</tg-time>

<tg-math>x^2 + y^2</tg-math>

<aside>Pull quote<cite>The Author</cite></aside>

<tg-map lat="41.9" long="12.5" zoom="14" />

<details open>
  <summary>Summary with **bold text**</summary>

  ### Details heading
  - List item with _italic text_
  - List item with <tg-spoiler>spoiler</tg-spoiler>
</details>

<tg-collage>

![](https://telegram.org/img/t_logo.png)
![](https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4)

</tg-collage>

<tg-slideshow>

![](https://telegram.org/img/t_logo.png)
![](https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4)

</tg-slideshow>
