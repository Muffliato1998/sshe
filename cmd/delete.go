package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sshe/config"
	"sshe/utils"
)

// delete 命令
var deleteCmd = &cobra.Command{
	Use:   "delete <IP>",
	Short: "Delete a matching node.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ip := args[0]

		// 检查 IP 地址合法性
		if err := utils.AssertIpAddressValid(ip); err != nil {
			return fmt.Errorf("invalid IP address: %w", err)
		}

		// 获取节点信息
		nodes, err := config.GetNode(ip, user)
		if err != nil {
			return fmt.Errorf("failed to retrieve nodes: %w", err)
		}
		if len(nodes) == 0 {
			return fmt.Errorf("no data matching IP %s was found, deletion is not needed", ip)
		}

		// 选择节点
		node, err := selectNode(ip, nodes)
		if err != nil {
			return err
		}

		// 确认是否删除
		if err := confirmAndDeleteNode(node); err != nil {
			return err
		}

		return nil
	},
}

// 确认并删除节点
func confirmAndDeleteNode(node config.Node) error {
	// 打印节点信息
	_ = printNodeInfo(node, false)

	// 询问确认删除
	fmt.Printf("\nAre you sure to delete the node with IP %s and username %s? [y/N]: ", node.IP, node.Username)
	var sureToDelete string
	_, err := fmt.Scanln(&sureToDelete)
	if err != nil && err.Error() != "unexpected newline" {
		return fmt.Errorf("error reading the confirmation: %w", err)
	}

	if sureToDelete != "y" && sureToDelete != "Y" {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	// 删除节点
	if err := config.DeleteNode(node.IP, node.Username); err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	fmt.Printf("Node with IP %s and username %s has been deleted successfully.\n", node.IP, node.Username)
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&user, "user", "u", "", "Specifies the username for connection.")
}
