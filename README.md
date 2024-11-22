# sshe

`sshe` 是一个基于 Cobra 的命令行工具，用于管理和快速创建 SSH 连接。提供了包括添加、删除、连接、查询等功能。

## 安装

工具开发环境：go version go1.21.13 darwin/amd64
IDE：JetBrains Goland 2024.2.3

1. 克隆代码：
```bash
git clone git@github.com:Muffliato1998/sshe.git
cd sshe
```

2. 安装依赖
```bash
go mod tidy
```

3. 构建项目
```bash
go build -ldflags="-s -w" -o sshe
mv sshe /usr/local/bin/sshe
chmod +x /usr/local/bin/sshe
```

4. 测试执行
```bash
sshe version
sshe help
```
成功输出版本号和帮助信息则说明安装完成。


## 指令参数
```bash
➜ sshe help   
SSHE is a simple tool used to manage and connect to the SSH password information of remote machines, providing basic functions such as adding, deleting, querying and connecting.

Usage:
  sshe [flags]
  sshe [command]

Available Commands:
  add         Add node connection information.
  delete      Delete a matching node.
  get         Get info of a specific node.
  help        Help about any command
  link        Connect to a matching node.
  list        Search the machine list according to conditions.
  version     Show version.

Flags:
  -h, --help   help for sshe

Use "sshe [command] --help" for more information about a command.
```

## 配置文件

sshe 的配置文件存放在 `~/.sshe` 目录下，包括以下文件：

- `config.yaml`: 配置文件，目前只支持配置保存节点密码使用的加密秘钥串 `secret_key`；
- `nodes.yaml`: 存放用户保存的连接节点信息，包括节点名称、IP 地址、端口号、用户名、密码等。

```yaml
# 存放节点连接信息
nodes:
    - ip: 10.2.147.88
      username: root
      password: b028eafb2b829a2558b06a15dd6f34d6
      tag:
        - tag1
        - tag2
# 方便搜索的索引
tag_index:
    tag1:
        - 10.2.147.88@root
    tag2:
        - 10.2.147.88@root
```

## 使用例子



### 添加节点

添加一个新的节点连接信息：

```bash
sshe add 192.168.1.100 -u root -t tag1 -t tag2
```

说明：
- `add` 将添加一个名为 `192.168.1.100`，用户为 `root` 的节点连接信息到配置文件中；
- `-t` 为该节点指定标签，标签有助于后续根据条件搜索或过滤节点；
- 在添加时，IP 和用户名全局唯一，所以如果添加的节点已经存在，将会提示错误；

### 查看节点

查看指定已存储的节点信息：

```bash
ssh get 192.168.1.100
```

说明：

- `get` 命令可以查看一个指定节点的详细信息，包括 IP 地址、用户名、**密码** 和标签(会明文显示密码，需要注意防止密码泄露)；
- 可使用 `-u` 参数指定用户名，如果没有指定，在查询时发现有多个用户名相同的节点，将会提示输入用户名进一步确认。

### 删除节点

删除某个已存储的节点连接信息：

```bash
ssh delete 192.168.1.100
```

说明：

- `delete` 命令可以删除一个已存在的节点连接信息，删除后无法恢复，请谨慎操作；
- 可使用 `-u` 参数指定用户名，如果没有指定，在删除时发现有多个用户名相同的节点，将会提示输入用户名进一步确认。

### 连接节点

使用存储的节点连接信息，通过 SSH 连接到指定节点：

```bash
sshe link 192.168.1.100
```

说明：

- `link` 命令可以连接一个已存在的节点，通过 SSH 连接到指定节点；
- 可使用 `-u` 参数指定用户名，如果没有指定，在连接时发现有多个用户名相同的节点，将会提示输入用户名进一步确认；
- 如果该节点配置正确，命令执行后会自动打开一个 SSH 会话。

### 搜索节点列表

根据条件搜索节点信息，并列出所有匹配的节点：

```bash
ssh list -t tag1 --ip-start 10.2
```

说明：

- `list` 命令可以列出所有匹配的节点，支持通过标签、IP 地址、用户名等条件进行过滤，在列出节点时不会展示密码信息；
- 支持同时使用多个参数过滤，以 AND 方式组合查询。

以下是将这些参数整理成的表格：

| **参数**           | **缩写** | **参数说明**        | **示例**                         |
|------------------|--------|-----------------|--------------------------------|
| `--ip`           | `-i`   | 按完整 IP 地址搜索     | `sshe list --ip 192.168.1.100` |
| `--ip-start`     | 无      | 按 IP 地址的开头部分搜索  | `sshe list --ip-start 192.168` |
| `--ip-end`       | 无      | 按 IP 地址的结尾部分搜索  | `sshe list --ip-end .100`      |
| `--ip-contain`   | 无      | 按 IP 地址中包含的内容搜索 | `sshe list --ip-contain 168.1` |
| `--user`         | `-u`   | 按用户名搜索          | `sshe list --user root`        |
| `--user-start`   | 无      | 按用户名的开头部分搜索     | `sshe list --user-start ad`    |
| `--user-end`     | 无      | 按用户名的结尾部分搜索     | `sshe list --user-end min`     |
| `--user-contain` | 无      | 按用户名中包含的内容搜索    | `sshe list --user-contain mi`  |
| `--tag`          | `-t`   | 按标签搜索           | `sshe list --tag webserver`    |
| `--tag-start`    | 无      | 按标签的开头部分搜索      | `sshe list --tag-start prod`   |
| `--tag-end`      | 无      | 按标签的结尾部分搜索      | `sshe list --tag-end server`   |
| `--tag-contain`  | 无      | 按标签中包含的内容搜索     | `sshe list --tag-contain web`  |