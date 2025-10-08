// main.go
// 主程序入口 - Cobra + Viper优化版本

package main

import (
	"log"

	"multi-database-validator-optimization/cmd"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 执行Cobra命令
	cmd.Execute()
}
