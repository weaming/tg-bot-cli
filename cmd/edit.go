package cmd

import (
	"tg/internal/api"
	"tg/internal/util"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑消息文本",
	RunE:  runEdit,
}

var (
	editChat           string
	editMsgID          int
	editText           string
	editParseMode      string
	editDisablePreview bool
	editButtons        []string
)

func init() {
	rootCmd.AddCommand(editCmd)

	f := editCmd.Flags()
	f.StringVarP(&editChat, "chat", "c", "", "chat_id 或 username")
	f.IntVarP(&editMsgID, "msg", "m", 0, "要编辑的消息 ID（必填）")
	f.StringVarP(&editText, "text", "t", "", "新文本，使用 \"-\" 从 stdin 读取（必填）")
	f.StringVar(&editParseMode, "parse-mode", "HTML", "解析模式：HTML | Markdown | MarkdownV2")
	f.BoolVarP(&editDisablePreview, "disable-preview", "d", false, "禁用链接预览")
	f.StringArrayVarP(&editButtons, "button", "b", nil, "Inline 按钮行，格式同 send")

	editCmd.MarkFlagRequired("msg")
	editCmd.MarkFlagRequired("text")
}

func runEdit(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	chat, err := resolveTarget(editChat)
	if err != nil {
		return err
	}

	text, err := util.ReadTextOrStdin(editText)
	if err != nil {
		return err
	}

	replyMarkup, err := util.ParseButtons(editButtons)
	if err != nil {
		return err
	}

	msg, err := client.EditMessageText(api.EditMessageTextParams{
		ChatID:                chat,
		MessageID:             editMsgID,
		Text:                  text,
		ParseMode:             editParseMode,
		DisableWebPagePreview: editDisablePreview,
		ReplyMarkup:           replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
