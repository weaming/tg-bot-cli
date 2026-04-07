package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weaming/tg-bot-cli/api"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑消息文本",
	RunE:  runEdit,
}

var (
	editChat        string
	editMsgID       int
	editText        string
	editInputFile   string
	editParseMode   string
	editLinkPreview bool
	editButtons     []string
)

func init() {
	rootCmd.AddCommand(editCmd)

	f := editCmd.Flags()
	f.StringVarP(&editChat, "chat", "c", "", "chat_id 或 username")
	f.IntVarP(&editMsgID, "msg", "m", 0, "要编辑的消息 ID（必填）")
	f.StringVarP(&editText, "text", "t", "", "新文本（必填）")
	f.StringVarP(&editInputFile, "input-file", "i", "", "从文件读取新文本")
	f.StringVar(&editParseMode, "parse-mode", "", "解析模式：HTML | MarkdownV2")
	f.BoolVarP(&editLinkPreview, "link-preview", "l", false, "启用链接预览")
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

	text, err := api.ReadFromInput(editInputFile)
	if err != nil {
		return err
	}
	if editText != "" {
		text = editText
	}

	replyMarkup, err := api.ParseButtons(editButtons)
	if err != nil {
		return err
	}

	msg, err := client.EditMessageText(api.EditMessageTextParams{
		ChatID:                chat,
		MessageID:             editMsgID,
		Text:                  text,
		ParseMode:             editParseMode,
		DisableWebPagePreview: !editLinkPreview,
		ReplyMarkup:           replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
