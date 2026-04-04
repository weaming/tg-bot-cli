package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取 chat 信息",
	RunE:  runInfo,
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "获取当前 bot 信息",
	RunE:  runMe,
}

var infoChatID string

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(meCmd)

	infoCmd.Flags().StringVarP(&infoChatID, "chat", "c", "", "chat_id 或 username")
}

func runInfo(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	chat, err := resolveTarget(infoChatID)
	if err != nil {
		return err
	}

	result, err := client.GetChat(chat)
	if err != nil {
		return err
	}

	username := ""
	if result.Username != "" {
		username = "@" + result.Username
	}
	printResult(result, "id: %d  type: %s  title: %s  username: %s",
		result.ID, result.Type, result.Title, username)
	return nil
}

func runMe(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	user, err := client.GetMe()
	if err != nil {
		return err
	}

	printResult(user, fmt.Sprintf("id: %d  @%s  %s %s", user.ID, user.Username, user.FirstName, user.LastName))
	return nil
}
