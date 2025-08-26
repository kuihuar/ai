package skill

import (
	"fmt"
	"os/exec"
)

// ==================== Go 命令详解 ====================

// Go 命令是 Go 语言的核心工具链，提供了从开发到部署的完整功能

// 1. 基础命令

// 1.1 go version - 显示 Go 版本信息
func GoVersionExample() {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Go Version: %s", string(output))
}

// 1.2 go env - 显示 Go 环境变量
func GoEnvExample() {
	cmd := exec.Command("go", "env")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Go Environment:\n%s", string(output))
}

// 1.3 go help - 显示帮助信息
func GoHelpExample() {
	cmd := exec.Command("go", "help")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Go Help:\n%s", string(output))
}

// 2. 构建和运行命令
