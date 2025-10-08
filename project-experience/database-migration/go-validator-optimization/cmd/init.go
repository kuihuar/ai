// cmd/init.go
// init命令定义

package cmd

import (
	"fmt"
	"strings"

	"multi-database-validator-optimization/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	format string
	output string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "创建默认配置文件",
	Long: `创建默认配置文件

支持多种格式的配置文件:
- JSON (.json)
- YAML (.yaml, .yml)
- TOML (.toml)

使用示例:
  multi-database-validator init                    # 创建默认YAML配置文件
  multi-database-validator init --format json     # 创建JSON格式配置文件
  multi-database-validator init --output my-config.yaml  # 指定输出文件名`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// 添加标志
	initCmd.Flags().StringVarP(&format, "format", "f", "yaml", "配置文件格式 (json, yaml, toml)")
	initCmd.Flags().StringVarP(&output, "output", "o", "config.yaml", "输出文件名")

	// 验证标志
	initCmd.MarkFlagRequired("output")
}

func runInit(cmd *cobra.Command, args []string) error {
	// 验证格式
	validFormats := []string{"json", "yaml", "yml", "toml"}
	format = strings.ToLower(format)

	validFormat := false
	for _, f := range validFormats {
		if format == f {
			validFormat = true
			break
		}
	}

	if !validFormat {
		return fmt.Errorf("不支持的格式: %s，支持的格式: %v", format, validFormats)
	}

	// 根据格式设置文件扩展名
	if !strings.Contains(output, ".") {
		switch format {
		case "json":
			output += ".json"
		case "yaml", "yml":
			output += ".yaml"
		case "toml":
			output += ".toml"
		}
	}

	// 设置配置类型
	switch format {
	case "json":
		viper.SetConfigType("json")
	case "yaml", "yml":
		viper.SetConfigType("yaml")
	case "toml":
		viper.SetConfigType("toml")
	}

	// 设置默认配置
	setDefaultConfig()

	// 写入配置文件
	if err := config.CreateDefaultConfig(output); err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}

	fmt.Printf("✅ 默认配置文件已创建: %s\n", output)
	fmt.Printf("📝 请编辑配置文件设置正确的数据库连接信息\n")

	return nil
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	// 设置默认Azure配置
	viper.Set("azure", []map[string]interface{}{
		{
			"name":     "azure-db1",
			"host":     "your-azure-mysql1.mysql.database.azure.com",
			"user":     "your_username",
			"password": "your_password",
			"database": "db1",
			"charset":  "utf8mb4",
		},
		{
			"name":     "azure-db2",
			"host":     "your-azure-mysql2.mysql.database.azure.com",
			"user":     "your_username",
			"password": "your_password",
			"database": "db2",
			"charset":  "utf8mb4",
		},
	})

	// 设置默认AWS配置
	viper.Set("aws", []map[string]interface{}{
		{
			"name":     "aws-db1",
			"host":     "your-aws-rds1.region.rds.amazonaws.com",
			"user":     "your_username",
			"password": "your_password",
			"database": "db1",
			"charset":  "utf8mb4",
		},
		{
			"name":     "aws-db2",
			"host":     "your-aws-rds2.region.rds.amazonaws.com",
			"user":     "your_username",
			"password": "your_password",
			"database": "db2",
			"charset":  "utf8mb4",
		},
	})

	// 设置默认并发数
	viper.Set("max_workers", 3)
}
