package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/weaming/tg-bot-cli/api"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送消息或文件",
	RunE:  runSend,
}

var (
	sendTo          string
	sendText        string
	sendInputFile   string
	sendMd2Html     bool
	sendNoParse     bool
	sendSplitTable  bool
	sendRichFormat  string
	sendParseMode   string
	sendFile        string
	sendCaption     string
	sendReplyTo     int
	sendLinkPreview bool
	sendSilent      bool
	sendProtect     bool
	sendThread      int
	sendButtons     []string
	sendAsDoc       bool
)

func init() {
	rootCmd.AddCommand(sendCmd)

	f := sendCmd.Flags()
	f.StringVarP(&sendTo, "to", "t", "", "目标 chat_id 或 username")
	f.StringVarP(&sendText, "text", "m", "", "消息文本")
	f.StringVarP(&sendInputFile, "input-file", "i", "", "从文件或 stdin（-）读取消息文本")
	f.BoolVarP(&sendMd2Html, "md2html", "", false, "转换 Markdown 扩展（Rich Markdown 保留 Markdown，普通消息转换为 HTML；.md 文件自动启用）")
	f.BoolVar(&sendNoParse, "no-parse", false, "禁用本地 Markdown/MarkdownV2 预处理，原样交给 Telegram")
	f.BoolVarP(&sendSplitTable, "split-table", "", false, "将 markdown 表格拆分成多行列表模式")
	f.StringVar(&sendRichFormat, "rich", "", "使用 Rich Message：md 或 html")
	f.StringVar(&sendParseMode, "parse-mode", "", "解析模式：HTML | MarkdownV2")
	f.StringVarP(&sendFile, "file", "f", "", "要发送的文件路径")
	f.StringVarP(&sendCaption, "caption", "c", "", "文件说明文字")
	f.IntVarP(&sendReplyTo, "reply-to", "r", 0, "回复的消息 ID")
	f.BoolVarP(&sendLinkPreview, "link-preview", "l", false, "启用链接预览")
	f.BoolVarP(&sendSilent, "silent", "s", false, "静默发送（不通知）")
	f.BoolVarP(&sendProtect, "protect", "p", false, "防止转发和保存")
	f.IntVar(&sendThread, "thread", 0, "话题 ID（message_thread_id）")
	f.StringArrayVarP(&sendButtons, "button", "b", nil, "Inline 按钮行，格式：文字:URL,文字2:URL2（多次使用添加多行，| 分隔行）")
	f.BoolVar(&sendAsDoc, "as-doc", false, "强制作为文档发送（绕过图片尺寸限制）")
}

func runSend(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	target, err := resolveTarget(sendTo)
	if err != nil {
		return err
	}
	sendTo = target

	replyMarkup, err := api.ParseButtons(sendButtons)
	if err != nil {
		return err
	}

	if sendFile != "" {
		if sendRichFormat != "" {
			return fmt.Errorf("--rich 不能与 --file 一起使用")
		}
		return sendMedia(client, replyMarkup)
	}
	return sendTextMsg(client, replyMarkup)
}

func sendTextMsg(client *api.Client, replyMarkup *api.InlineKeyboardMarkup) error {
	text, err := api.ReadFromInput(sendInputFile)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	if sendText != "" {
		text = sendText
	}
	if text == "" {
		return fmt.Errorf("--text 不能为空")
	}
	if sendRichFormat != "" {
		return sendRichTextMsg(client, replyMarkup, text)
	}

	if shouldConvertMarkdown(sendNoParse, sendMd2Html, sendInputFile) {
		text = api.ConvertMarkdownToHTML(text, sendSplitTable)
		sendParseMode = "HTML"
	} else if !sendNoParse && sendParseMode == "MarkdownV2" {
		text = api.EscapeMarkdownV2(text)
	}

	msg, err := client.SendMessage(api.SendMessageParams{
		ChatID:                sendTo,
		Text:                  text,
		ParseMode:             sendParseMode,
		MessageThreadID:       sendThread,
		ReplyToMessageID:      sendReplyTo,
		DisableWebPagePreview: !sendLinkPreview,
		DisableNotification:   sendSilent,
		ProtectContent:        sendProtect,
		ReplyMarkup:           replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}

func sendRichTextMsg(client *api.Client, replyMarkup *api.InlineKeyboardMarkup, text string) error {
	content, err := buildRichMessage(
		sendRichFormat,
		text,
		shouldConvertMarkdown(sendNoParse, sendMd2Html, sendInputFile),
	)
	if err != nil {
		return err
	}
	if sendParseMode != "" {
		return fmt.Errorf("--rich 不能与 --parse-mode 一起使用")
	}
	if sendLinkPreview {
		return fmt.Errorf("--rich 不支持 --link-preview")
	}

	var replyParameters *api.ReplyParameters
	if sendReplyTo > 0 {
		replyParameters = &api.ReplyParameters{MessageID: sendReplyTo}
	}

	msg, err := client.SendRichMessage(api.SendRichMessageParams{
		ChatID:              sendTo,
		RichMessage:         content,
		MessageThreadID:     sendThread,
		ReplyParameters:     replyParameters,
		DisableNotification: sendSilent,
		ProtectContent:      sendProtect,
		ReplyMarkup:         replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}

func shouldConvertMarkdown(noParse, explicitConversion bool, inputFile string) bool {
	return !noParse && (explicitConversion || api.IsMarkdownFile(inputFile))
}

func sendMedia(client *api.Client, replyMarkup *api.InlineKeyboardMarkup) error {
	if _, err := os.Stat(sendFile); err != nil {
		return fmt.Errorf("文件不存在: %s", sendFile)
	}

	caption, err := api.ReadTextOrStdin(sendCaption)
	if err != nil {
		return fmt.Errorf("读取说明文字失败: %w", err)
	}

	if sendParseMode == "MarkdownV2" {
		caption = api.EscapeMarkdownV2(caption)
	}

	msg, err := client.SendMedia(api.SendMediaParams{
		ChatID:              sendTo,
		FilePath:            sendFile,
		Caption:             caption,
		ParseMode:           sendParseMode,
		MessageThreadID:     sendThread,
		ReplyToMessageID:    sendReplyTo,
		DisableNotification: sendSilent,
		ProtectContent:      sendProtect,
		ReplyMarkup:         replyMarkup,
		ForceDocument:       sendAsDoc,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
