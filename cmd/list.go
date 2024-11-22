package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sshe/config"
	"strings"
)

var (
	ips                  []string
	ipStarts             []string
	ipEnds               []string
	ipContains           []string
	users                []string
	userStarts           []string
	userEnds             []string
	userContains         []string
	conditionTags        []string
	conditionTagStarts   []string
	conditionTagEnds     []string
	conditionTagContains []string
)

// get 命令
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Search the machine list according to conditions.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// 获取所有节点
		nodes := config.GlobalNode.Nodes
		var matchedNodes []config.Node

		// 遍历所有节点，进行筛选
		for _, node := range nodes {
			// 检查 IP 相关条件
			if !matchIPs(node.IP) {
				continue
			}

			// 检查用户名相关条件
			if !matchUsers(node.Username) {
				continue
			}

			// 检查标签相关条件
			if !matchTags(node.Tags) {
				continue
			}

			// 如果节点符合所有条件，则加入匹配结果
			matchedNodes = append(matchedNodes, node)
		}

		// 输出结果
		// 如果没有匹配的节点
		if len(matchedNodes) == 0 {
			fmt.Println("No matching nodes found.")
		} else {
			// 计算最大宽度
			maxIPLen := len("IP")
			maxUserLen := len("Username")
			maxTagLen := len("Tags")

			// 计算每列的最大宽度
			for _, node := range matchedNodes {
				maxIPLen = max(maxIPLen, len(node.IP))
				maxUserLen = max(maxUserLen, len(node.Username))
				maxTagLen = max(maxTagLen, len(formatTags(node.Tags)))
			}

			// 打印表头
			fmt.Printf("%-*s %-*s %-*s\n", maxIPLen, "IP", maxUserLen, "Username", maxTagLen, "Tags")

			// 打印每个节点的信息
			for _, node := range matchedNodes {
				// 使用动态宽度的格式化输出对齐每列
				fmt.Printf("%-*s %-*s %-*s\n", maxIPLen, node.IP, maxUserLen, node.Username, maxTagLen, formatTags(node.Tags))
			}
		}

		return nil
	},
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return "No tags"
	}
	var tagString string
	for _, tag := range tags {
		tagString += "#" + tag + " "
	}
	return tagString
}

// 根据 IP 条件进行匹配
func matchIPs(ip string) bool {
	if len(ips) > 0 && !contains(ips, ip) {
		return false
	}
	for _, ipStart := range ipStarts {
		if !strings.HasPrefix(ip, ipStart) {
			return false
		}
	}
	for _, ipEnd := range ipEnds {
		if !strings.HasSuffix(ip, ipEnd) {
			return false
		}
	}
	for _, ipContain := range ipContains {
		if !strings.Contains(ip, ipContain) {
			return false
		}
	}
	return true
}

// 根据用户名条件进行匹配
func matchUsers(username string) bool {
	if len(users) > 0 && !contains(users, username) {
		return false
	}
	for _, userStart := range userStarts {
		if !strings.HasPrefix(username, userStart) {
			return false
		}
	}
	for _, userEnd := range userEnds {
		if !strings.HasSuffix(username, userEnd) {
			return false
		}
	}
	for _, userContain := range userContains {
		if !strings.Contains(username, userContain) {
			return false
		}
	}
	return true
}

// 根据标签条件进行匹配
func matchTags(tags []string) bool {
	if len(conditionTags) > 0 && !hasTags(tags, conditionTags) {
		return false
	}
	for _, tagStart := range conditionTagStarts {
		if !startsWithAny(tags, tagStart) {
			return false
		}
	}
	for _, tagEnd := range conditionTagEnds {
		if !endsWithAny(tags, tagEnd) {
			return false
		}
	}
	for _, tagContain := range conditionTagContains {
		if !containsAny(tags, tagContain) {
			return false
		}
	}
	return true
}

// 判断数组是否包含某个元素
func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// 判断标签数组中是否包含指定的标签
func hasTags(tags, conditionTags []string) bool {
	for _, tag := range conditionTags {
		if !contains(tags, tag) {
			return false
		}
	}
	return true
}

// 判断标签数组中是否有以指定字符串开头的标签
func startsWithAny(tags []string, prefix string) bool {
	for _, tag := range tags {
		if strings.HasPrefix(tag, prefix) {
			return true
		}
	}
	return false
}

// 判断标签数组中是否有以指定字符串结尾的标签
func endsWithAny(tags []string, suffix string) bool {
	for _, tag := range tags {
		if strings.HasSuffix(tag, suffix) {
			return true
		}
	}
	return false
}

// 判断标签数组中是否包含指定的子字符串
func containsAny(tags []string, substr string) bool {
	for _, tag := range tags {
		if strings.Contains(tag, substr) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayVarP(&ips, "ip", "i", []string{}, "Search by IP address.")
	listCmd.Flags().StringArrayVarP(&ipStarts, "ip-start", "", []string{}, "Search by the beginning of the IP address.")
	listCmd.Flags().StringArrayVarP(&ipEnds, "ip-end", "", []string{}, "Search by the end of the IP address.")
	listCmd.Flags().StringArrayVarP(&ipContains, "ip-contain", "", []string{}, "Search by the content contained in the IP address.")

	listCmd.Flags().StringArrayVarP(&users, "user", "u", []string{}, "Search by username.")
	listCmd.Flags().StringArrayVarP(&userStarts, "user-start", "", []string{}, "Search by the beginning of the username.")
	listCmd.Flags().StringArrayVarP(&userEnds, "user-end", "", []string{}, "Search by the end of the username.")
	listCmd.Flags().StringArrayVarP(&userContains, "user-contain", "", []string{}, "Search by the content contained in the username.")

	listCmd.Flags().StringArrayVarP(&conditionTags, "tag", "t", []string{}, "Search by tag.")
	listCmd.Flags().StringArrayVarP(&conditionTagStarts, "tag-start", "", []string{}, "Search by the beginning of the tag.")
	listCmd.Flags().StringArrayVarP(&conditionTagEnds, "tag-end", "", []string{}, "Search by the end of the tag.")
	listCmd.Flags().StringArrayVarP(&conditionTagContains, "tag-contain", "", []string{}, "Search by the content contained in the tag.")
}
