// internal/config/config.go
// 配置文件处理 - 使用Viper

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"multi-database-validator-optimization/internal/types"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitViper 初始化Viper配置
func InitViper() error {
	// 设置配置文件名称（不包含扩展名）
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // 默认类型

	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./examples")
	viper.AddConfigPath("$HOME/.multi-database-validator")
	viper.AddConfigPath("/etc/multi-database-validator")

	// 绑定环境变量
	bindEnvVars()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置
			fmt.Println("⚠️  未找到配置文件，使用默认配置")
			setDefaults()
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	} else {
		fmt.Printf("使用配置文件: %s\n", viper.ConfigFileUsed())
	}

	// 解析配置到结构体
	var globalConfig types.Config
	if err := viper.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 验证配置（只在有实际配置时验证）
	if len(globalConfig.Azure) > 0 || len(globalConfig.AWS) > 0 {
		if err := validateConfig(globalConfig); err != nil {
			return fmt.Errorf("配置验证失败: %v", err)
		}
	}

	return nil
}

// bindEnvVars 绑定环境变量
func bindEnvVars() {
	// 设置环境变量前缀
	viper.SetEnvPrefix("MDV")
	viper.AutomaticEnv()

	// 注释掉环境变量绑定，避免干扰YAML数组解析
	// 如果需要使用环境变量，请使用命令行参数替代
	// viper.BindEnv("azure.0.host", "MDV_AZURE_HOST")
	// viper.BindEnv("azure.0.user", "MDV_AZURE_USER")
	// viper.BindEnv("azure.0.password", "MDV_AZURE_PASSWORD")
	// viper.BindEnv("azure.0.database", "MDV_AZURE_DATABASE")

	// viper.BindEnv("aws.0.host", "MDV_AWS_HOST")
	// viper.BindEnv("aws.0.user", "MDV_AWS_USER")
	// viper.BindEnv("aws.0.password", "MDV_AWS_PASSWORD")
	// viper.BindEnv("aws.0.database", "MDV_AWS_DATABASE")

	viper.BindEnv("max_workers", "MDV_MAX_WORKERS")
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 设置默认的Azure配置
	viper.SetDefault("azure", []map[string]interface{}{
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

	// 设置默认的AWS配置
	viper.SetDefault("aws", []map[string]interface{}{
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

	viper.SetDefault("max_workers", 3)
	viper.SetDefault("output_dir", "output")
	viper.SetDefault("output", "consistency_report.json")
}

// validateConfig 验证配置
func validateConfig(config types.Config) error {
	if len(config.Azure) == 0 {
		return fmt.Errorf("Azure实例列表不能为空")
	}

	if len(config.AWS) == 0 {
		return fmt.Errorf("AWS实例列表不能为空")
	}

	if len(config.Azure) != len(config.AWS) {
		return fmt.Errorf("Azure和AWS实例数量不匹配: Azure=%d, AWS=%d", len(config.Azure), len(config.AWS))
	}

	if config.MaxWorkers <= 0 {
		return fmt.Errorf("最大并发数必须大于0")
	}

	return nil
}

// CreateDefaultConfig 创建默认配置文件
func CreateDefaultConfig(filename string) error {
	// 根据文件扩展名设置配置类型
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		viper.SetConfigType("json")
	case ".yaml", ".yml":
		viper.SetConfigType("yaml")
	case ".toml":
		viper.SetConfigType("toml")
	default:
		viper.SetConfigType("yaml")
	}

	// 设置默认配置
	setDefaults()

	// 写入配置文件
	if err := viper.WriteConfigAs(filename); err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}

	return nil
}

// GetConfigValue 获取配置值
func GetConfigValue(key string) interface{} {
	return viper.Get(key)
}

// SetConfigValue 设置配置值
func SetConfigValue(key string, value interface{}) {
	viper.Set(key, value)
}

// GetConfig 获取完整配置
func GetConfig() (*types.Config, error) {
	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}
	return &config, nil
}

// WatchConfig 监听配置文件变化
func WatchConfig() {
	viper.WatchConfig()
}

// OnConfigChange 配置文件变化回调
func OnConfigChange(fn func(in fsnotify.Event)) {
	viper.OnConfigChange(fn)
}

// GetOutputDir 获取输出目录
func GetOutputDir() string {
	outputDir := viper.GetString("output_dir")
	if outputDir == "" {
		outputDir = "output"
	}
	return outputDir
}

// GetReportsDir 获取报告目录
func GetReportsDir() string {
	reportsDir := filepath.Join(GetOutputDir(), "reports")
	ensureDirExists(reportsDir)
	return reportsDir
}

// GetLogsDir 获取日志目录
func GetLogsDir() string {
	logsDir := filepath.Join(GetOutputDir(), "logs")
	ensureDirExists(logsDir)
	return logsDir
}

// GetTempDir 获取临时目录
func GetTempDir() string {
	tempDir := filepath.Join(GetOutputDir(), "temp")
	ensureDirExists(tempDir)
	return tempDir
}

// GetReportPath 获取报告文件路径
func GetReportPath(filename string) string {
	if filename == "" {
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("consistency_report_%s.json", timestamp)
	}

	// 如果文件名包含路径，直接返回
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return filename
	}

	// 否则放在reports目录下
	return filepath.Join(GetReportsDir(), filename)
}

// ensureDirExists 确保目录存在
func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// InitOutputDirs 初始化输出目录
func InitOutputDirs() error {
	dirs := []string{
		GetOutputDir(),
		GetReportsDir(),
		GetLogsDir(),
		GetTempDir(),
	}

	for _, dir := range dirs {
		if err := ensureDirExists(dir); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}
	}

	return nil
}
