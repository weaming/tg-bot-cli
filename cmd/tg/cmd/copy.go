package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weaming/tg-bot-cli/api"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "复制消息（不显示来源）",
	RunE:  runCopy,
}

var (
	copyFrom      string
	copyTo        string
	copyMsgID     int
	copyCaption   string
	copyParseMode string
	copyThread    int
	copyReplyTo   int
	copySilent    bool
	copyProtect   bool
	copyButtons   []string
)

func init() {
	rootCmd.AddCommand(copyCmd)

	f := copyCmd.Flags()
	f.StringVarP(&copyFrom, "from", "f", "", "来源 chat_id（必填）")
	f.StringVarP(&copyTo, "to", "t", "", "目标 chat_id 或 username")
	f.IntVarP(&copyMsgID, "msg", "m", 0, "消息 ID（必填）")
	f.StringVarP(&copyCaption, "caption", "c", "", "覆盖原始说明文字，使用 \"-\" 从 stdin 读取")
	f.StringVar(&copyParseMode, "parse-mode", "", "解析模式")
	f.IntVar(&copyThread, "thread", 0, "目标话题 ID")
	f.IntVarP(&copyReplyTo, "reply-to", "r", 0, "回复的消息 ID")
	f.BoolVarP(&copySilent, "silent", "s", false, "静默发送")
	f.BoolVarP(&copyProtect, "protect", "p", false, "防止转发和保存")
	f.StringArrayVarP(&copyButtons, "button", "b", nil, "Inline 按钮行")

	copyCmd.MarkFlagRequired("from")
	copyCmd.MarkFlagRequired("msg")
}

func runCopy(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	target, err := resolveTarget(copyTo)
	if err != nil {
		return err
	}

	caption, err := api.ReadTextOrStdin(copyCaption)
	if err != nil {
		return err
	}

	replyMarkup, err := api.ParseButtons(copyButtons)
	if err != nil {
		return err
	}

	msg, err := client.CopyMessage(api.CopyMessageParams{
		ChatID:              target,
		FromChatID:          copyFrom,
		MessageID:           copyMsgID,
		Caption:             caption,
		ParseMode:           copyParseMode,
		MessageThreadID:     copyThread,
		ReplyToMessageID:    copyReplyTo,
		DisableNotification: copySilent,
		ProtectContent:      copyProtect,
		ReplyMarkup:         replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
