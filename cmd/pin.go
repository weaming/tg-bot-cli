package cmd

import (
	"fmt"

	"tg/internal/api"

	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "置顶消息",
	RunE:  runPin,
}

var unpinCmd = &cobra.Command{
	Use:   "unpin",
	Short: "取消置顶消息",
	RunE:  runUnpin,
}

var (
	pinChat   string
	pinMsgID  int
	pinSilent bool

	unpinChat  string
	unpinMsgID int
)

func init() {
	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(unpinCmd)

	pf := pinCmd.Flags()
	pf.StringVarP(&pinChat, "chat", "c", "", "chat_id 或 username")
	pf.IntVarP(&pinMsgID, "msg", "m", 0, "要置顶的消息 ID（必填）")
	pf.BoolVarP(&pinSilent, "silent", "s", false, "静默置顶（不发送通知）")
	pinCmd.MarkFlagRequired("msg")

	uf := unpinCmd.Flags()
	uf.StringVarP(&unpinChat, "chat", "c", "", "chat_id 或 username")
	uf.IntVarP(&unpinMsgID, "msg", "m", 0, "要取消置顶的消息 ID（不填则取消最近一条）")
}

func runPin(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	chat, err := resolveTarget(pinChat)
	if err != nil {
		return err
	}

	ok, err := client.PinChatMessage(api.PinChatMessageParams{
		ChatID:              chat,
		MessageID:           pinMsgID,
		DisableNotification: pinSilent,
	})
	if err != nil {
		return err
	}

	printResult(map[string]any{"ok": ok}, fmt.Sprintf("pinned: %v", ok))
	return nil
}

func runUnpin(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	chat, err := resolveTarget(unpinChat)
	if err != nil {
		return err
	}

	ok, err := client.UnpinChatMessage(api.UnpinChatMessageParams{
		ChatID:    chat,
		MessageID: unpinMsgID,
	})
	if err != nil {
		return err
	}

	printResult(map[string]any{"ok": ok}, fmt.Sprintf("unpinned: %v", ok))
	return nil
}
