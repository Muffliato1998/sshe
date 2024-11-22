module sshe

go 1.21.13

// 直接依赖
require (
	github.com/spf13/cobra v1.8.1
	golang.org/x/crypto v0.29.0
	golang.org/x/term v0.26.0
	gopkg.in/yaml.v3 v3.0.1
)

// 间接依赖
require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.27.0 // indirect
)
