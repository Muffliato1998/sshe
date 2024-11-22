package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sshe/config"
)

// get 命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// 打印标题
		fmt.Print("\n========================================\n")
		fmt.Print("              SSHE Version             \n")
		fmt.Print("========================================\n")

		// 打印版本信息，带颜色
		fmt.Printf("Version:    \x1b[32m%s\x1b[0m\n", config.Version)
		if config.Commit != "" {
			fmt.Printf("Commit:     \x1b[33m%s\x1b[0m\n", config.Commit)
		}

		// 打印分隔符
		fmt.Print("\n========================================\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
