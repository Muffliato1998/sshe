package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sshe/config"
	"sshe/utils"
	"strings"
)

// get 命令
var getCmd = &cobra.Command{
	Use:   "get <IP>",
	Short: "Get info of a specific node.",
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
			return fmt.Errorf("no data matching IP %s was found", ip)
		}

		// 选择节点
		node, err := selectNode(ip, nodes)
		if err != nil {
			return err
		}
		return printNodeInfo(node, true)
	},
}

// 选择节点
func selectNode(ip string, nodes []config.Node) (config.Node, error) {
	if len(nodes) == 1 {
		return nodes[0], nil
	}

	// 多个用户记录时提示用户选择
	fmt.Printf("\nIP %s has multiple user records:\n", ip)
	userMap := map[string]config.Node{}
	for _, node := range nodes {
		fmt.Printf("%s, ", node.Username)
		userMap[node.Username] = node
	}
	fmt.Print("\nPlease select one of them (default: root): ")

	var inputUser string
	_, err := fmt.Scanln(&inputUser)
	if err != nil && err.Error() != "unexpected newline" {
		return config.Node{}, fmt.Errorf("error reading the username: %w", err)
	}

	// 默认值处理
	if inputUser == "" {
		inputUser = "root"
	}

	// 查找用户对应的节点
	node, exists := userMap[inputUser]
	if !exists {
		return config.Node{}, fmt.Errorf("no data matching IP %s and username %s was found", ip, inputUser)
	}

	return node, nil
}

// 打印节点信息
func printNodeInfo(node config.Node, printPassword bool) error {
	fmt.Println("\nFound Node!")
	fmt.Printf("IP: %s\n", node.IP)
	fmt.Printf("Username: %s\n", node.Username)

	if printPassword {
		password, err := utils.DecryptAES(node.Password, config.GlobalConfig.SecretKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt password: %w", err)
		}
		fmt.Printf("Password: %s\n", password)
	}

	if len(node.Tags) > 0 {
		fmt.Printf("Tags: %s\n", "#"+strings.Join(node.Tags, " #"))
	}
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&user, "user", "u", "", "Specifies the username for connection.")
}
