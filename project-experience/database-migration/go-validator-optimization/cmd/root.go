// cmd/root.go
// 根命令定义

package cmd

import (
	"fmt"
	"os"

	"multi-database-validator-optimization/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multi-database-validator",
	Short: "多数据库一致性验证工具",
	Long: `多数据库一致性验证工具 - 使用Cobra + Viper优化版本

这是一个用Go语言编写的多数据库一致性验证工具，用于验证MySQL数据库从Azure迁移到AWS后的数据一致性。

功能特性:
- 支持Azure和AWS多个数据库实例的对比验证
- 支持JSON、YAML、TOML等多种配置文件格式
- 支持环境变量配置
- 支持命令行参数覆盖
- 并行验证多个数据库对比对
- 大表分批处理，避免内存溢出
- 详细的验证报告和日志记录

使用示例:
  multi-database-validator validate                    # 使用默认配置验证
  multi-database-validator validate --config config.yaml  # 指定配置文件
  multi-database-validator init --format yaml         # 创建YAML配置文件
  multi-database-validator validate --max-workers 5   # 设置并发数`,
	Version: "2.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().String("log-level", "info", "日志级别 (debug, info, warn, error)")

	// 绑定环境变量
	// 注意：config参数不绑定到Viper，只用于命令行参数
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// 初始化配置
	if err := config.InitViper(); err != nil {
		fmt.Fprintf(os.Stderr, "配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 如果指定了配置文件，使用指定的文件
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "读取指定配置文件失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 初始化输出目录
	if err := config.InitOutputDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "输出目录初始化失败: %v\n", err)
		os.Exit(1)
	}
}
