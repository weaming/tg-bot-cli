package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "显示当前 token/target/proxy 来源",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	printConfigRow("token", resolveTokenSource())
	printConfigRow("target", resolveTargetSource())
	printConfigRow("proxy", resolveProxySource())
	return nil
}

func printConfigRow(key, value string) {
	fmt.Printf("%-8s %s\n", key, value)
}

func resolveTokenSource() string {
	if flagToken != "" {
		return fmt.Sprintf("flag      %s", maskToken(flagToken))
	}
	if v := os.Getenv("TG_BOT_TOKEN"); v != "" {
		return fmt.Sprintf("env       %s", maskToken(v))
	}
	if CompiledToken != "" {
		return fmt.Sprintf("compiled  %s", maskToken(CompiledToken))
	}
	return "not set"
}

func resolveTargetSource() string {
	if v := os.Getenv("TG_TARGET"); v != "" {
		return fmt.Sprintf("env       %s", v)
	}
	if CompiledTarget != "" {
		return fmt.Sprintf("compiled  %s", CompiledTarget)
	}
	return "not set"
}

func resolveProxySource() string {
	if flagProxy != "" {
		return fmt.Sprintf("flag      %s", flagProxy)
	}
	if v := os.Getenv("TG_PROXY"); v != "" {
		return fmt.Sprintf("env(TG_PROXY)   %s", v)
	}
	for _, key := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"} {
		if v := os.Getenv(key); v != "" {
			return fmt.Sprintf("env(%-11s) %s", key, v)
		}
	}
	return "not set"
}

// maskToken 保留 bot id 和 secret 首尾各4位，中间替换为 ***
func maskToken(token string) string {
	colonIdx := strings.Index(token, ":")
	if colonIdx < 0 {
		return "***"
	}
	id := token[:colonIdx]
	secret := token[colonIdx+1:]
	if len(secret) <= 8 {
		return id + ":***"
	}
	return fmt.Sprintf("%s:%s***%s", id, secret[:4], secret[len(secret)-4:])
}
