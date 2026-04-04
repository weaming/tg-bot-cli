package cmd

import (
	"fmt"

	"tg/internal/api"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除消息",
	RunE:  runDelete,
}

var (
	deleteChat  string
	deleteMsgID int
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	f := deleteCmd.Flags()
	f.StringVarP(&deleteChat, "chat", "c", "", "chat_id 或 username")
	f.IntVarP(&deleteMsgID, "msg", "m", 0, "要删除的消息 ID（必填）")

	deleteCmd.MarkFlagRequired("msg")
}

func runDelete(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	chat, err := resolveTarget(deleteChat)
	if err != nil {
		return err
	}

	ok, err := client.DeleteMessage(api.DeleteMessageParams{
		ChatID:    chat,
		MessageID: deleteMsgID,
	})
	if err != nil {
		return err
	}

	printResult(map[string]any{"ok": ok}, fmt.Sprintf("deleted: %v", ok))
	return nil
}
