package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Config 配置文件
type Config struct {
	Language       string `yaml:"language"`
	PasswordEncode string `yaml:"password_encode"`
	SecretKey      string `yaml:"secret_key"`
}

// Node 节点
type Node struct {
	IP       string   `yaml:"ip"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Tags     []string `yaml:"tag"`
}

// NodesFile 存储节点的文件结构
type NodesFile struct {
	Nodes    []Node              `yaml:"nodes"`
	TagIndex map[string][]string `yaml:"tag_index"`
}

var (
	Version     = "v2024.11.22"
	Commit      = ""
	configPath  = filepath.Join(os.Getenv("HOME"), ".sshe", "sshe.conf")
	nodesPath   = filepath.Join(os.Getenv("HOME"), ".sshe", "node.yaml")
	defaultConf = Config{
		Language:       "zh-cn",
		PasswordEncode: "aes",
		SecretKey:      "sshe2024",
	}
	GlobalConfig = Config{}
	GlobalNode   = NodesFile{}
)

// loadYAMLFile 用于读取 YAML 文件并解码
func loadYAMLFile(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode YAML file %s: %v", filePath, err)
	}
	return nil
}

// writeYAMLFile 用于将数据编码并写入 YAML 文件
func writeYAMLFile(filePath string, v interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}(file)

	encoder := yaml.NewEncoder(file)
	defer func(encoder *yaml.Encoder) {
		err := encoder.Close()
		if err != nil {
			fmt.Printf("Failed to close encoder: %v\n", err)
		}
	}(encoder)

	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to write YAML file %s: %v", filePath, err)
	}
	return nil
}

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 创建 ~/.sshe/ 目录，若已存在则不影响
	home := filepath.Join(os.Getenv("HOME"), ".sshe")
	if err := os.MkdirAll(home, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 检查并初始化 sshe.conf 配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建并写入默认配置
		GlobalConfig = defaultConf
		if err := writeYAMLFile(configPath, GlobalConfig); err != nil {
			return err
		}
	} else {
		// 配置文件存在，读取并补充缺失的字段
		if err := loadYAMLFile(configPath, &GlobalConfig); err != nil {
			return err
		}

		// 如果某些配置项缺失，补充默认值
		if GlobalConfig.Language == "" {
			GlobalConfig.Language = defaultConf.Language
		}
		if GlobalConfig.PasswordEncode == "" {
			GlobalConfig.PasswordEncode = defaultConf.PasswordEncode
		}
		if GlobalConfig.SecretKey == "" {
			GlobalConfig.SecretKey = defaultConf.SecretKey
		}

		// 重新写回补充后的配置
		if err := writeYAMLFile(configPath, GlobalConfig); err != nil {
			return err
		}
	}

	// 检查并初始化 node.yaml
	if _, err := os.Stat(nodesPath); os.IsNotExist(err) {
		// 节点文件不存在，创建并初始化
		initialNodes := NodesFile{
			Nodes:    []Node{},
			TagIndex: map[string][]string{},
		}
		if err := writeYAMLFile(nodesPath, initialNodes); err != nil {
			return err
		}
	} else {
		// 节点文件存在，读取并补充缺失的部分
		if err := loadYAMLFile(nodesPath, &GlobalNode); err != nil {
			return err
		}

		// 补充缺失的字段
		if GlobalNode.Nodes == nil {
			GlobalNode.Nodes = []Node{}
		}
		if GlobalNode.TagIndex == nil {
			GlobalNode.TagIndex = map[string][]string{}
		}

		// 重新写回更新后的节点数据
		if err := writeYAMLFile(nodesPath, GlobalNode); err != nil {
			return err
		}
	}

	return nil
}

// AddNode 将节点信息添加到配置文件
func AddNode(ip, username, encryptedPassword string, tags []string) error {
	node := Node{
		IP:       ip,
		Username: username,
		Password: encryptedPassword,
		Tags:     tags,
	}

	// 添加到 GlobalNode
	GlobalNode.Nodes = append(GlobalNode.Nodes, node)
	for _, tag := range tags {
		GlobalNode.TagIndex[tag] = append(GlobalNode.TagIndex[tag], fmt.Sprintf("%s@%s", node.IP, node.Username))
	}

	// 重新写回节点文件
	if err := writeYAMLFile(nodesPath, GlobalNode); err != nil {
		return fmt.Errorf("failed to update nodes file: %v", err)
	}

	return nil
}

// GetNode 根据 IP 和用户名获取节点信息
func GetNode(ip, username string) ([]Node, error) {
	var matchedNodes []Node

	for _, node := range GlobalNode.Nodes {
		if (username == "" && node.IP == ip) || (username != "" && node.IP == ip && node.Username == username) {
			matchedNodes = append(matchedNodes, node)
		}
	}

	return matchedNodes, nil
}

// DeleteNode 根据 IP 和用户名删除节点
func DeleteNode(ip, username string) error {
	var nodeFound bool
	for i, node := range GlobalNode.Nodes {
		if node.IP == ip && node.Username == username {
			// 删除该节点
			GlobalNode.Nodes = append(GlobalNode.Nodes[:i], GlobalNode.Nodes[i+1:]...)

			// 删除 tag_index 中的记录
			for _, tag := range node.Tags {
				taggedNodes := GlobalNode.TagIndex[tag]
				for j, taggedNode := range taggedNodes {
					if taggedNode == node.IP+"@"+node.Username {
						GlobalNode.TagIndex[tag] = append(taggedNodes[:j], taggedNodes[j+1:]...)
						break
					}
				}
			}
			// 标记已找到
			nodeFound = true
			break
		}
	}

	if !nodeFound {
		return nil
	}

	// 重新写回节点文件
	if err := writeYAMLFile(nodesPath, GlobalNode); err != nil {
		return fmt.Errorf("failed to update nodes file: %v", err)
	}

	return nil
}
