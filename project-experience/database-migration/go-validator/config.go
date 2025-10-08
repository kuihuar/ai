// config.go
// 配置文件处理

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// loadConfig 从配置文件加载配置
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("解析YAML配置文件失败: %v", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("解析JSON配置文件失败: %v", err)
		}
	default:
		// 默认尝试JSON格式
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败，请确保文件格式为JSON或YAML: %v", err)
		}
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	// 设置默认值
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 3
	}

	return &config, nil
}

// validateConfig 验证配置
func validateConfig(config Config) error {
	if len(config.Azure) == 0 {
		return fmt.Errorf("Azure实例列表不能为空")
	}

	if len(config.AWS) == 0 {
		return fmt.Errorf("AWS实例列表不能为空")
	}

	if len(config.Azure) != len(config.AWS) {
		return fmt.Errorf("Azure和AWS实例数量不匹配")
	}

	// 验证每个实例配置
	for i, azureInstance := range config.Azure {
		if azureInstance.Name == "" || azureInstance.Host == "" || azureInstance.User == "" || azureInstance.Password == "" || azureInstance.Database == "" {
			return fmt.Errorf("Azure实例%d配置不完整", i+1)
		}
	}

	for i, awsInstance := range config.AWS {
		if awsInstance.Name == "" || awsInstance.Host == "" || awsInstance.User == "" || awsInstance.Password == "" || awsInstance.Database == "" {
			return fmt.Errorf("AWS实例%d配置不完整", i+1)
		}
	}

	return nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(filename string) error {
	config := Config{
		Azure: []DatabaseInstance{
			{
				Name:     "azure-db1",
				Host:     "your-azure-mysql1.mysql.database.azure.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db1",
				Charset:  "utf8mb4",
			},
			{
				Name:     "azure-db2",
				Host:     "your-azure-mysql2.mysql.database.azure.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db2",
				Charset:  "utf8mb4",
			},
		},
		AWS: []DatabaseInstance{
			{
				Name:     "aws-db1",
				Host:     "your-aws-rds1.region.rds.amazonaws.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db1",
				Charset:  "utf8mb4",
			},
			{
				Name:     "aws-db2",
				Host:     "your-aws-rds2.region.rds.amazonaws.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db2",
				Charset:  "utf8mb4",
			},
		},
		MaxWorkers: 3,
	}

	var data []byte
	var err error

	// 根据文件扩展名选择序列化方式
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("序列化YAML配置失败: %v", err)
		}
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化JSON配置失败: %v", err)
		}
	default:
		// 默认使用JSON格式
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化配置失败: %v", err)
		}
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Printf("默认配置文件已创建: %s\n", filename)
	return nil
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *Config {
	return &Config{
		Azure: []DatabaseInstance{
			{
				Name:     "azure-db1",
				Host:     "your-azure-mysql1.mysql.database.azure.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db1",
				Charset:  "utf8mb4",
			},
		},
		AWS: []DatabaseInstance{
			{
				Name:     "aws-db1",
				Host:     "your-aws-rds1.region.rds.amazonaws.com",
				User:     "your_username",
				Password: "your_password",
				Database: "db1",
				Charset:  "utf8mb4",
			},
		},
		MaxWorkers: 3,
	}
}
