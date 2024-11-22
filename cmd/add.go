package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"sshe/config"
	"sshe/utils"
	"strings"
)

// add 命令
var addCmd = &cobra.Command{
	Use:   "add <IP>",
	Short: "Add node connection information.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAddCommand,
}

func runAddCommand(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	ip := args[0]

	// 检查 IP 地址合法性
	if err := utils.AssertIpAddressValid(ip); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	// 查询现有记录
	existUsernames := getExistingUsernames(ip)

	// 获取用户名
	username, err := getUsername(existUsernames, ip)
	if err != nil {
		return err
	}

	// 处理标签输入
	tags, err := handleTagsInput(tags)
	if err != nil {
		return err
	}

	// 获取密码
	password, err := getPassword()
	if err != nil {
		return err
	}

	// 加密密码并添加节点
	cipherText, err := utils.EncryptAES(password, config.GlobalConfig.SecretKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	if err := config.AddNode(ip, username, cipherText, tags); err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}

	fmt.Println("Node added successfully!")
	return nil
}

// 查询现有用户名
func getExistingUsernames(ip string) []string {
	var existUsernames []string
	for _, node := range config.GlobalNode.Nodes {
		if node.IP == ip {
			existUsernames = append(existUsernames, node.Username)
		}
	}
	return existUsernames
}

// 获取用户名
func getUsername(existUsernames []string, ip string) (string, error) {
	if user != "" {
		for _, existUser := range existUsernames {
			if existUser == user {
				return "", fmt.Errorf("the username %s already exists for IP %s", user, ip)
			}
		}
		return user, nil
	}

	if len(existUsernames) > 0 {
		fmt.Printf("\nThe IP-recorded usernames are: %s.", strings.Join(existUsernames, ", "))
	}
	fmt.Print("\nInput username (default: root): ")

	var inputUser string
	_, err := fmt.Scanln(&inputUser)
	if err != nil && err.Error() != "unexpected newline" {
		return "", fmt.Errorf("error reading the username: %w", err)
	}
	if inputUser == "" {
		inputUser = "root"
	}

	for _, existUser := range existUsernames {
		if existUser == inputUser {
			return "", fmt.Errorf("the username %s already exists for IP %s", inputUser, ip)
		}
	}
	return inputUser, nil
}

// 处理标签输入
func handleTagsInput(existingTags []string) ([]string, error) {
	if len(existingTags) > 0 {
		fmt.Printf("\nTags can facilitate searching. Current tags: %v", existingTags)
	} else {
		fmt.Print("\nTags can facilitate searching. No tags added.")
	}
	fmt.Print("\nAdd new tags (format: #tag1#tag2, or enter to skip): ")

	var tagStr string
	_, err := fmt.Scanln(&tagStr)
	if err != nil && err.Error() != "unexpected newline" {
		return nil, fmt.Errorf("error reading tags: %w", err)
	}
	if tagStr != "" {
		for _, tag := range strings.Split(tagStr, "#") {
			if tag != "" {
				existingTags = append(existingTags, tag)
			}
		}
	}
	return existingTags, nil
}

// 获取密码
func getPassword() ([]byte, error) {
	fmt.Print("\nInput password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	fmt.Println() // 换行以避免后续输出和密码提示混在一起
	return password, nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&user, "user", "u", "", "Specifies the username for connection.")
	addCmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Specify the tags for connection. Multiple tags are supported.")
}
