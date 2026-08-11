package cmd

import (
	"fmt"

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
	editMd2Html     bool
	editSplitTable  bool
	editRichFormat  string
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
	f.StringVarP(&editInputFile, "input-file", "i", "", "从文件或 stdin（-）读取新文本")
	f.BoolVarP(&editMd2Html, "md2html", "", false, "将 markdown 转换为 HTML（.md 文件自动转换）")
	f.BoolVarP(&editSplitTable, "split-table", "", false, "将 markdown 表格拆分成多行列表模式")
	f.StringVar(&editRichFormat, "rich", "", "使用 Rich Message：html 或 markdown")
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

	if editRichFormat != "" {
		return editRichText(client, chat, text)
	}

	if editMd2Html || api.IsMarkdownFile(editInputFile) {
		text = api.ConvertMarkdownToHTML(text, editSplitTable)
		editParseMode = "HTML"
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

func editRichText(client *api.Client, chat, text string) error {
	content, err := buildRichMessage(
		editRichFormat,
		text,
		editMd2Html || api.IsMarkdownFile(editInputFile),
	)
	if err != nil {
		return err
	}
	if editParseMode != "" {
		return fmt.Errorf("--rich 不能与 --parse-mode 一起使用")
	}
	if editLinkPreview {
		return fmt.Errorf("--rich 不支持 --link-preview")
	}

	replyMarkup, err := api.ParseButtons(editButtons)
	if err != nil {
		return err
	}

	msg, err := client.EditMessageText(api.EditMessageTextParams{
		ChatID:      chat,
		MessageID:   editMsgID,
		RichMessage: &content,
		ReplyMarkup: replyMarkup,
	})
	if err != nil {
		return err
	}

	printResult(msg, "message_id: %d", msg.MessageID)
	return nil
}
