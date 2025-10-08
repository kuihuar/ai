// main.go
// 主程序入口

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			// 创建默认配置文件
			configFile := "config.json"
			if len(os.Args) > 2 {
				configFile = os.Args[2]
			}
			if err := createDefaultConfig(configFile); err != nil {
				log.Fatalf("创建配置文件失败: %v", err)
			}
			return
		case "help", "-h", "--help":
			printUsage()
			return
		}
	}

	// 加载配置
	var config *Config
	var err error

	// 尝试从配置文件加载
	configFiles := []string{"config.yaml", "config.yml", "config.json"}
	var configFile string

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			configFile = file
			break
		}
	}

	if configFile != "" {
		config, err = loadConfig(configFile)
		if err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}
		fmt.Printf("从配置文件加载配置: %s\n", configFile)
	} else {
		// 使用默认配置
		config = getDefaultConfig()
		fmt.Println("使用默认配置")
	}

	// 创建数据库对比对
	databasePairs := make([]DatabasePair, len(config.Azure))
	for i := range config.Azure {
		databasePairs[i] = DatabasePair{
			AzureInstance: config.Azure[i],
			AWSInstance:   config.AWS[i],
		}
	}

	// 创建验证器实例
	validator := NewMultiDatabaseValidator(databasePairs)

	// 执行验证
	startTime := time.Now()
	validator.validateAllDatabases(config.MaxWorkers)
	duration := time.Since(startTime)

	// 生成报告
	summary, err := validator.generateReport("consistency_report.json")
	if err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}

	// 输出验证摘要
	printSummary(summary, duration)

	// 如果有不一致或错误，显示详细信息
	if summary.InconsistentDatabases > 0 || summary.ErrorDatabases > 0 {
		printDetailedResults(summary)
	}

	// 根据验证结果设置退出码
	if summary.ErrorDatabases > 0 {
		os.Exit(1)
	}
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("多数据库一致性验证工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run .                    # 使用默认配置或config.json/config.yaml运行验证")
	fmt.Println("  go run . init [filename]    # 创建默认配置文件 (支持.json/.yaml/.yml)")
	fmt.Println("  go run . help               # 显示帮助信息")
	fmt.Println("")
	fmt.Println("支持的配置文件格式:")
	fmt.Println("  - JSON: config.json")
	fmt.Println("  - YAML: config.yaml 或 config.yml")
	fmt.Println("")
	fmt.Println("JSON配置文件格式 (config.json):")
	fmt.Println(`{
  "azure": [
    {
      "name": "azure-db1",
      "host": "your-azure-mysql1.mysql.database.azure.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    }
  ],
  "aws": [
    {
      "name": "aws-db1",
      "host": "your-aws-rds1.region.rds.amazonaws.com",
      "user": "your_username",
      "password": "your_password",
      "database": "db1",
      "charset": "utf8mb4"
    }
  ],
  "max_workers": 3
}`)
	fmt.Println("")
	fmt.Println("YAML配置文件格式 (config.yaml):")
	fmt.Println(`azure:
  - name: azure-db1
    host: your-azure-mysql1.mysql.database.azure.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4

aws:
  - name: aws-db1
    host: your-aws-rds1.region.rds.amazonaws.com
    user: your_username
    password: your_password
    database: db1
    charset: utf8mb4

max_workers: 3`)
}

// printSummary 打印验证摘要
func printSummary(summary *ValidationSummary, duration time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("数据库一致性验证完成")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("验证时间: %s\n", summary.Timestamp)
	fmt.Printf("总耗时: %v\n", duration.Round(time.Second))
	fmt.Printf("总数据库数: %d\n", summary.TotalDatabases)
	fmt.Printf("验证成功: %d\n", summary.SuccessfulValidations)
	fmt.Printf("数据不一致: %d\n", summary.InconsistentDatabases)
	fmt.Printf("验证错误: %d\n", summary.ErrorDatabases)
	fmt.Printf("成功率: %s\n", summary.SuccessRate)
	fmt.Println(strings.Repeat("=", 60))
}

// printDetailedResults 打印详细结果
func printDetailedResults(summary *ValidationSummary) {
	fmt.Println("\n详细信息:")
	for dbName, result := range summary.Results {
		if result.Status == "INCONSISTENT" || result.Status == "ERROR" {
			fmt.Printf("\n数据库: %s\n", dbName)
			fmt.Printf("Azure实例: %s\n", result.AzureInstance)
			fmt.Printf("AWS实例: %s\n", result.AWSInstance)
			fmt.Printf("状态: %s\n", result.Status)
			fmt.Printf("开始时间: %s\n", result.StartTime)
			fmt.Printf("结束时间: %s\n", result.EndTime)
			fmt.Printf("表数量: Azure(%d) vs AWS(%d)\n", result.AzureTables, result.AWSTables)

			if len(result.Errors) > 0 {
				fmt.Println("错误信息:")
				for _, error := range result.Errors {
					fmt.Printf("  - %s\n", error)
				}
			}

			// 显示不一致的表
			if len(result.TableComparisons) > 0 {
				fmt.Println("表对比结果:")
				for _, comparison := range result.TableComparisons {
					status := "✓"
					if !comparison.Match {
						status = "✗"
					}
					fmt.Printf("  %s %s\n", status, comparison.Table)
					if !comparison.Match {
						fmt.Printf("    Azure实例: %s 数据库: %s\n", comparison.AzureInstance, comparison.AzureDatabase)
						fmt.Printf("    AWS实例: %s 数据库: %s\n", comparison.AWSInstance, comparison.AWSDatabase)
						fmt.Printf("    Azure校验和: %s\n", comparison.AzureChecksum)
						fmt.Printf("    AWS校验和: %s\n", comparison.AWSChecksum)
					}
				}
			}
		}
	}
}
