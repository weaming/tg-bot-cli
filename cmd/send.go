package cmd

import (
	"fmt"
	"os"

	"tg/internal/api"
	"tg/internal/util"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送消息或文件",
	RunE:  runSend,
}

var (
	sendTo             string
	sendText           string
	sendParseMode      string
	sendFile           string
	sendCaption        string
	sendReplyTo        int
	sendDisablePreview bool
	sendSilent         bool
	sendProtect        bool
	sendThread         int
	sendButtons        []string
)

func init() {
	rootCmd.AddCommand(sendCmd)

	f := sendCmd.Flags()
	f.StringVarP(&sendTo, "to", "t", "", "目标 chat_id 或 username")
	f.StringVarP(&sendText, "text", "m", "", "消息文本，使用 \"-\" 从 stdin 读取")
	f.StringVar(&sendParseMode, "parse-mode", "", "解析模式：HTML | Markdown | MarkdownV2")
	f.StringVarP(&sendFile, "file", "f", "", "要发送的文件路径")
	f.StringVarP(&sendCaption, "caption", "c", "", "文件说明文字")
	f.IntVarP(&sendReplyTo, "reply-to", "r", 0, "回复的消息 ID")
	f.BoolVarP(&sendDisablePreview, "disable-preview", "d", false, "禁用链接预览")
	f.BoolVarP(&sendSilent, "silent", "s", false, "静默发送（不通知）")
	f.BoolVarP(&sendProtect, "protect", "p", false, "防止转发和保存")
	f.IntVar(&sendThread, "thread", 0, "话题 ID（message_thread_id）")
	f.StringArrayVarP(&sendButtons, "button", "b", nil, "Inline 按钮行，格式：文字:URL,文字2:URL2（多次使用添加多行，| 分隔行）")
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

	replyMarkup, err := util.ParseButtons(sendButtons)
	if err != nil {
		return err
	}

	if sendFile != "" {
		return sendMedia(client, replyMarkup)
	}
	return sendTextMsg(client, replyMarkup)
}

func sendTextMsg(client *api.Client, replyMarkup *api.InlineKeyboardMarkup) error {
	text, err := util.ReadTextOrStdin(sendText)
	if err != nil {
		return fmt.Errorf("读取文本失败: %w", err)
	}
	if text == "" {
		return fmt.Errorf("--text 不能为空")
	}

	if sendParseMode == "MarkdownV2" {
		text = util.EscapeMarkdownV2(text)
	}

	msg, err := client.SendMessage(api.SendMessageParams{
		ChatID:                sendTo,
		Text:                  text,
		ParseMode:             sendParseMode,
		MessageThreadID:       sendThread,
		ReplyToMessageID:      sendReplyTo,
		DisableWebPagePreview: sendDisablePreview,
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

func sendMedia(client *api.Client, replyMarkup *api.InlineKeyboardMarkup) error {
	if _, err := os.Stat(sendFile); err != nil {
		return fmt.Errorf("文件不存在: %s", sendFile)
	}

	caption, err := util.ReadTextOrStdin(sendCaption)
	if err != nil {
		return fmt.Errorf("读取说明文字失败: %w", err)
	}

	if sendParseMode == "MarkdownV2" {
		caption = util.EscapeMarkdownV2(caption)
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
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
