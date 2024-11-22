package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sshe/config"
)

var user string
var tags []string

var rootCmd = &cobra.Command{
	Use:   "sshe",
	Short: "Simple SSH management tool",
	Long:  `SSHE is a simple tool used to manage and connect to the SSH password information of remote machines, providing basic functions such as adding, deleting, querying and connecting.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute 命令执行函数
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 初始化根命令
func init() {
	// 读取配置文件
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("加载配置失败:", err)
		os.Exit(1)
	}
}
