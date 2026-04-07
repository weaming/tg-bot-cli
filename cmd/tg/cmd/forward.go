package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weaming/tg-bot-cli/api"
)

var forwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "转发消息",
	RunE:  runForward,
}

var (
	forwardFrom    string
	forwardTo      string
	forwardMsgID   int
	forwardThread  int
	forwardSilent  bool
	forwardProtect bool
)

func init() {
	rootCmd.AddCommand(forwardCmd)

	f := forwardCmd.Flags()
	f.StringVarP(&forwardFrom, "from", "f", "", "来源 chat_id（必填）")
	f.StringVarP(&forwardTo, "to", "t", "", "目标 chat_id 或 username")
	f.IntVarP(&forwardMsgID, "msg", "m", 0, "消息 ID（必填）")
	f.IntVar(&forwardThread, "thread", 0, "目标话题 ID")
	f.BoolVarP(&forwardSilent, "silent", "s", false, "静默发送")
	f.BoolVarP(&forwardProtect, "protect", "p", false, "防止转发和保存")

	forwardCmd.MarkFlagRequired("from")
	forwardCmd.MarkFlagRequired("msg")
}

func runForward(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	target, err := resolveTarget(forwardTo)
	if err != nil {
		return err
	}

	msg, err := client.ForwardMessage(api.ForwardMessageParams{
		ChatID:              target,
		FromChatID:          forwardFrom,
		MessageID:           forwardMsgID,
		MessageThreadID:     forwardThread,
		DisableNotification: forwardSilent,
		ProtectContent:      forwardProtect,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
