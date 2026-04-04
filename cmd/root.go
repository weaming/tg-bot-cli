package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tg/internal/api"

	"github.com/spf13/cobra"
)

var (
	CompiledToken  string // 编译时注入：-X tg/cmd.CompiledToken=xxx
	CompiledTarget string // 编译时注入：-X tg/cmd.CompiledTarget=xxx

	flagToken string
	flagProxy string
	flagJSON  bool
)

var rootCmd = &cobra.Command{
	Use:   "tg",
	Short: "Telegram Bot CLI",
	Long:  "通过 Telegram Bot API 发送和管理消息",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "T", "", "Bot Token（优先于环境变量和编译内置）")
	rootCmd.PersistentFlags().StringVarP(&flagProxy, "proxy", "x", "", "HTTP/SOCKS5 代理，如 http://127.0.0.1:7890（优先于 TG_PROXY 和系统环境变量）")
	rootCmd.PersistentFlags().BoolVarP(&flagJSON, "json", "j", false, "以 JSON 格式输出完整结果")
}

// resolveToken 按优先级返回 token：--token > TG_BOT_TOKEN > 编译内置
func resolveToken() (string, error) {
	if flagToken != "" {
		return flagToken, nil
	}
	if envToken := os.Getenv("TG_BOT_TOKEN"); envToken != "" {
		return envToken, nil
	}
	if CompiledToken != "" {
		return CompiledToken, nil
	}
	return "", fmt.Errorf("未找到 Bot Token，请通过 --token、TG_BOT_TOKEN 环境变量或编译时注入提供")
}

// resolveTarget 按优先级返回目标 chat_id：flag 值 > TG_TARGET > 编译内置
// 非数字 ID 自动补全 @ 前缀
func resolveTarget(flagValue string) (string, error) {
	target := flagValue
	if target == "" {
		target = os.Getenv("TG_TARGET")
	}
	if target == "" {
		target = CompiledTarget
	}
	if target == "" {
		return "", fmt.Errorf("未指定目标，请通过 --to/--chat、TG_TARGET 环境变量或编译时注入提供")
	}
	return normalizeTarget(target), nil
}

// normalizeTarget 对不带 @ 的用户名自动补全前缀
func normalizeTarget(target string) string {
	if strings.HasPrefix(target, "@") {
		return target
	}
	if _, err := strconv.ParseInt(target, 10, 64); err == nil {
		return target // 纯数字 ID，不处理
	}
	return "@" + target
}

// resolveProxy 按优先级返回代理地址：--proxy > TG_PROXY > 空（由系统环境变量接管）
func resolveProxy() string {
	if flagProxy != "" {
		return flagProxy
	}
	return os.Getenv("TG_PROXY")
}

// newClient 创建已验证 token 的 API 客户端
func newClient() (*api.Client, error) {
	token, err := resolveToken()
	if err != nil {
		return nil, err
	}
	return api.NewClient(token, resolveProxy())
}

// printResult 根据 --json flag 决定输出格式
func printResult(v any, shortFmt string, args ...any) {
	if flagJSON {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(data))
		return
	}
	fmt.Printf(shortFmt+"\n", args...)
}
