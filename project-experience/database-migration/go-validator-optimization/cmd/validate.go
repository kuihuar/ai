// cmd/validate.go
// validate命令定义

package cmd

import (
	"fmt"
	"time"

	"multi-database-validator-optimization/internal/config"
	"multi-database-validator-optimization/internal/types"
	"multi-database-validator-optimization/internal/validator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workers    int
	outputFile string
	dryRun     bool
	azureHost  string
	azureUser  string
	azurePass  string
	azureDB    string
	awsHost    string
	awsUser    string
	awsPass    string
	awsDB      string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证数据库一致性",
	Long: `验证Azure和AWS数据库的一致性

支持多种配置方式:
1. 配置文件 (推荐)
2. 命令行参数
3. 环境变量

配置优先级: 命令行参数 > 环境变量 > 配置文件 > 默认值

使用示例:
  multi-database-validator validate                           # 使用配置文件验证
  multi-database-validator validate --workers 5              # 设置并发数
  multi-database-validator validate --dry-run                # 试运行模式
  multi-database-validator validate --azure-host azure.com   # 命令行指定Azure主机`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// 添加标志
	validateCmd.Flags().IntVarP(&workers, "workers", "w", 3, "最大并发数")
	validateCmd.Flags().StringVarP(&outputFile, "output", "o", "consistency_report.json", "输出报告文件")
	validateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "试运行模式，不执行实际验证")

	// Azure配置标志
	validateCmd.Flags().StringVar(&azureHost, "azure-host", "", "Azure数据库主机")
	validateCmd.Flags().StringVar(&azureUser, "azure-user", "", "Azure数据库用户名")
	validateCmd.Flags().StringVar(&azurePass, "azure-password", "", "Azure数据库密码")
	validateCmd.Flags().StringVar(&azureDB, "azure-database", "", "Azure数据库名称")

	// AWS配置标志
	validateCmd.Flags().StringVar(&awsHost, "aws-host", "", "AWS数据库主机")
	validateCmd.Flags().StringVar(&awsUser, "aws-user", "", "AWS数据库用户名")
	validateCmd.Flags().StringVar(&awsPass, "aws-password", "", "AWS数据库密码")
	validateCmd.Flags().StringVar(&awsDB, "aws-database", "", "AWS数据库名称")

	// 绑定环境变量
	viper.BindPFlag("workers", validateCmd.Flags().Lookup("workers"))
	viper.BindPFlag("output", validateCmd.Flags().Lookup("output"))
	viper.BindPFlag("dry_run", validateCmd.Flags().Lookup("dry-run"))

	// 注意：Azure和AWS参数不绑定到Viper，只用于命令行参数覆盖
}

func runValidate(cmd *cobra.Command, args []string) error {
	// 初始化配置
	if err := initValidationConfig(); err != nil {
		return fmt.Errorf("初始化配置失败: %v", err)
	}

	// 显示配置信息
	if viper.GetBool("verbose") {
		showConfig()
	}

	// 试运行模式
	if viper.GetBool("dry_run") {
		fmt.Println("🔍 试运行模式 - 显示配置信息，不执行实际验证")
		showConfig()
		return nil
	}

	// 开始验证
	fmt.Println("🚀 开始数据库一致性验证...")
	startTime := time.Now()

	// 创建配置对象
	cfg := &types.Config{
		MaxWorkers: viper.GetInt("workers"),
	}

	// 解析Azure和AWS配置
	if azureConfig := viper.Get("azure"); azureConfig != nil {
		azureInstances := azureConfig.([]interface{})
		cfg.Azure = make([]types.DatabaseInstance, len(azureInstances))
		for i, instance := range azureInstances {
			inst := instance.(map[string]interface{})
			cfg.Azure[i] = types.DatabaseInstance{
				Name:     inst["name"].(string),
				Host:     inst["host"].(string),
				User:     inst["user"].(string),
				Password: inst["password"].(string),
				Database: inst["database"].(string),
				Charset:  inst["charset"].(string),
			}
		}
	}

	if awsConfig := viper.Get("aws"); awsConfig != nil {
		awsInstances := awsConfig.([]interface{})
		cfg.AWS = make([]types.DatabaseInstance, len(awsInstances))
		for i, instance := range awsInstances {
			inst := instance.(map[string]interface{})
			cfg.AWS[i] = types.DatabaseInstance{
				Name:     inst["name"].(string),
				Host:     inst["host"].(string),
				User:     inst["user"].(string),
				Password: inst["password"].(string),
				Database: inst["database"].(string),
				Charset:  inst["charset"].(string),
			}
		}
	}

	// 创建验证器并执行验证
	validatorInstance := validator.NewMultiDatabaseValidator(cfg)
	if err := validatorInstance.ValidateAllDatabases(); err != nil {
		return fmt.Errorf("验证失败: %v", err)
	}

	// 生成报告
	outputFile := config.GetReportPath(viper.GetString("output"))
	summary, err := validatorInstance.GenerateReport(outputFile)
	if err != nil {
		return fmt.Errorf("生成报告失败: %v", err)
	}

	duration := time.Since(startTime)

	// 显示验证结果
	fmt.Printf("✅ 验证完成，耗时: %v\n", duration)
	fmt.Printf("📊 验证结果:\n")
	fmt.Printf("  - 总数据库数: %d\n", summary.TotalDatabases)
	fmt.Printf("  - 验证成功: %d\n", summary.SuccessfulValidations)
	fmt.Printf("  - 数据不一致: %d\n", summary.InconsistentDatabases)
	fmt.Printf("  - 验证错误: %d\n", summary.ErrorDatabases)
	fmt.Printf("  - 成功率: %s\n", summary.SuccessRate)

	return nil
}

// initValidationConfig 初始化验证配置
func initValidationConfig() error {
	// 设置默认值
	viper.SetDefault("workers", 3)
	viper.SetDefault("output", "consistency_report.json")
	viper.SetDefault("dry_run", false)

	// 如果命令行指定了单实例配置，覆盖配置文件
	if azureHost != "" || awsHost != "" {
		// 创建单实例配置
		azureConfig := map[string]interface{}{
			"name":     "azure-single",
			"host":     azureHost,
			"user":     azureUser,
			"password": azurePass,
			"database": azureDB,
			"charset":  "utf8mb4",
		}

		awsConfig := map[string]interface{}{
			"name":     "aws-single",
			"host":     awsHost,
			"user":     awsUser,
			"password": awsPass,
			"database": awsDB,
			"charset":  "utf8mb4",
		}

		viper.Set("azure", []map[string]interface{}{azureConfig})
		viper.Set("aws", []map[string]interface{}{awsConfig})
	}

	return nil
}

// showConfig 显示当前配置
func showConfig() {
	fmt.Println("📋 当前配置:")
	fmt.Printf("  - 配置文件: %s\n", viper.ConfigFileUsed())
	fmt.Printf("  - 并发数: %d\n", viper.GetInt("workers"))
	fmt.Printf("  - 输出文件: %s\n", viper.GetString("output"))
	fmt.Printf("  - 详细模式: %t\n", viper.GetBool("verbose"))
	fmt.Printf("  - 试运行: %t\n", viper.GetBool("dry_run"))

	// 显示Azure配置
	if azureInstances := viper.Get("azure"); azureInstances != nil {
		if instances, ok := azureInstances.([]interface{}); ok {
			fmt.Printf("  - Azure实例数: %d\n", len(instances))
			for i, instance := range instances {
				if inst, ok := instance.(map[string]interface{}); ok {
					fmt.Printf("    [%d] %s: %s/%s\n", i+1, inst["name"], inst["host"], inst["database"])
				}
			}
		}
	}

	// 显示AWS配置
	if awsInstances := viper.Get("aws"); awsInstances != nil {
		if instances, ok := awsInstances.([]interface{}); ok {
			fmt.Printf("  - AWS实例数: %d\n", len(instances))
			for i, instance := range instances {
				if inst, ok := instance.(map[string]interface{}); ok {
					fmt.Printf("    [%d] %s: %s/%s\n", i+1, inst["name"], inst["host"], inst["database"])
				}
			}
		}
	}
}
