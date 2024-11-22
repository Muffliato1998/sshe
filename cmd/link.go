package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"os"
	"sshe/config"
	"sshe/utils"
	"time"
)

// link 命令
var linkCmd = &cobra.Command{
	Use:   "link <IP>",
	Short: "Connect to a matching node.",
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

		password, err := utils.DecryptAES(node.Password, config.GlobalConfig.SecretKey)
		if err != nil {
			return err
		}

		err = sshConnect(node, password)
		if err != nil {
			return err
		}

		return nil
	},
}

// 使用 SSH 连接到节点
func sshConnect(node config.Node, password string) error {
	// 创建 SSH 客户端配置
	clientConfig := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接到远程节点
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", node.IP), clientConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil && err.Error() != "EOF" {
			fmt.Printf("Failed to close SSH client: %v\n", err)
		}
	}(client)

	// 创建一个新的 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			fmt.Printf("Failed to close SSH session: %v\n", err)
		}
	}(session)

	// 设置会话的输入和输出，连接到本地终端
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 请求伪终端（Pty）模拟交互式 shell 环境
	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return fmt.Errorf("failed to request pseudo terminal: %w", err)
	}

	// 启动交互式 shell 会话
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// 等待会话结束
	if err := session.Wait(); err != nil {
		return fmt.Errorf("session exited with error: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// 设置 link 命令的标志
	linkCmd.Flags().StringVarP(&user, "user", "u", "", "Specifies the username for connection.")
}
