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
		fmt.Printf("Version:    \x1b[32m%s\x1b[0m\n", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
