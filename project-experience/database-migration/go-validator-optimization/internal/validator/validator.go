// internal/validator/validator.go
// 验证器包

package validator

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"multi-database-validator-optimization/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

// MultiDatabaseValidator 多数据库一致性验证器
type MultiDatabaseValidator struct {
	config  *types.Config
	results map[string]types.DatabaseResult
	mu      sync.RWMutex
}

// NewMultiDatabaseValidator 创建新的验证器
func NewMultiDatabaseValidator(config *types.Config) *MultiDatabaseValidator {
	return &MultiDatabaseValidator{
		config:  config,
		results: make(map[string]types.DatabaseResult),
	}
}

// ValidateAllDatabases 并行验证所有数据库对比对
func (v *MultiDatabaseValidator) ValidateAllDatabases() error {
	log.Printf("开始验证 %d 个数据库对比对，最大并发数: %d", len(v.config.Azure), v.config.MaxWorkers)

	// 创建数据库对比对
	databasePairs := make([]types.DatabasePair, len(v.config.Azure))
	for i := range v.config.Azure {
		databasePairs[i] = types.DatabasePair{
			AzureInstance: v.config.Azure[i],
			AWSInstance:   v.config.AWS[i],
		}
	}

	// 使用goroutine和channel进行并发控制
	semaphore := make(chan struct{}, v.config.MaxWorkers)
	var wg sync.WaitGroup
	resultsChan := make(chan types.DatabaseResult, len(databasePairs))

	// 启动验证任务
	for _, pair := range databasePairs {
		wg.Add(1)
		go func(p types.DatabasePair) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			result := v.validateDatabase(p)
			resultsChan <- result
		}(pair)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果
	for result := range resultsChan {
		v.mu.Lock()
		v.results[result.Database] = result
		v.mu.Unlock()
		log.Printf("数据库对比 %s vs %s 验证完成，状态: %s", result.AzureInstance, result.AWSInstance, result.Status)
	}

	log.Println("所有数据库验证完成")
	return nil
}

// validateDatabase 验证单个数据库对比对的一致性
func (v *MultiDatabaseValidator) validateDatabase(pair types.DatabasePair) types.DatabaseResult {
	azureInstance := pair.AzureInstance
	awsInstance := pair.AWSInstance

	log.Printf("开始验证数据库对比: %s (Azure: %s) vs %s (AWS: %s)",
		azureInstance.Database, azureInstance.Name, awsInstance.Database, awsInstance.Name)

	// 初始化结果
	result := types.DatabaseResult{
		Database:         azureInstance.Database,
		AzureInstance:    azureInstance.Name,
		AWSInstance:      awsInstance.Name,
		TableComparisons: []types.TableComparison{},
		Status:           "SUCCESS",
		Errors:           []string{},
		StartTime:        time.Now().Format(time.RFC3339),
	}

	azureConn, awsConn, err := v.connectDatabases(azureInstance, awsInstance)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, err.Error())
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}
	defer azureConn.Close()
	defer awsConn.Close()

	// 获取表列表
	azureTables, err := v.getTableList(azureConn, azureInstance.Database)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("获取Azure表列表失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	awsTables, err := v.getTableList(awsConn, awsInstance.Database)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("获取AWS表列表失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	result.AzureTables = len(azureTables)
	result.AWSTables = len(awsTables)

	// 检查表数量一致性
	if len(azureTables) != len(awsTables) {
		result.Status = "WARNING"
		errorMsg := fmt.Sprintf("表数量不一致: Azure(%d) vs AWS(%d)", len(azureTables), len(awsTables))
		result.Errors = append(result.Errors, errorMsg)
		log.Printf("数据库 %s: %s", azureInstance.Database, errorMsg)
	}

	// 对比每个表的数据一致性
	log.Printf("开始验证数据库 %s 中的 %d 个表", azureInstance.Database, len(azureTables))

	for i, table := range azureTables {
		log.Printf("验证表 %d/%d: %s", i+1, len(azureTables), table)

		if !contains(awsTables, table) {
			errorMsg := fmt.Sprintf("表 %s 在AWS中不存在", table)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "INCONSISTENT"
			log.Printf("数据库 %s: %s", azureInstance.Database, errorMsg)
			continue
		}

		// 计算校验和
		azureChecksum, err := v.calculateTableChecksum(azureConn, azureInstance.Database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("表 %s Azure校验和计算失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", azureInstance.Database, errorMsg)
			continue
		}

		awsChecksum, err := v.calculateTableChecksum(awsConn, awsInstance.Database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("表 %s AWS校验和计算失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", azureInstance.Database, errorMsg)
			continue
		}

		// 记录对比结果
		tableComparison := types.TableComparison{
			Table:         table,
			AzureChecksum: azureChecksum,
			AWSChecksum:   awsChecksum,
			Match:         azureChecksum == awsChecksum,
			AzureInstance: azureInstance.Name,
			AWSInstance:   awsInstance.Name,
			AzureDatabase: azureInstance.Database,
			AWSDatabase:   awsInstance.Database,
		}

		result.TableComparisons = append(result.TableComparisons, tableComparison)

		// 检查是否一致
		if azureChecksum != awsChecksum {
			result.Status = "INCONSISTENT"
			log.Printf("数据不一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s",
				azureInstance.Name, azureInstance.Database, table,
				awsInstance.Name, awsInstance.Database, table)
		} else {
			log.Printf("数据一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s",
				azureInstance.Name, azureInstance.Database, table,
				awsInstance.Name, awsInstance.Database, table)
		}
	}

	// 记录结束时间
	result.EndTime = time.Now().Format(time.RFC3339)

	// 统计验证结果
	consistentTables := 0
	for _, comparison := range result.TableComparisons {
		if comparison.Match {
			consistentTables++
		}
	}
	totalTables := len(result.TableComparisons)

	log.Printf("数据库 %s 验证完成:", azureInstance.Database)
	log.Printf("  状态: %s", result.Status)
	log.Printf("  表数量: Azure(%d) vs AWS(%d)", result.AzureTables, result.AWSTables)
	log.Printf("  数据一致性: %d/%d", consistentTables, totalTables)
	log.Printf("  错误数量: %d", len(result.Errors))

	return result
}

// connectDatabases 连接数据库
func (v *MultiDatabaseValidator) connectDatabases(azureInstance, awsInstance types.DatabaseInstance) (*sql.DB, *sql.DB, error) {
	// 连接Azure数据库
	azureDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		azureInstance.User, azureInstance.Password, azureInstance.Host, azureInstance.Database, azureInstance.Charset)

	azureConn, err := sql.Open("mysql", azureDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("Azure数据库连接失败: %v", err)
	}

	if err := azureConn.Ping(); err != nil {
		azureConn.Close()
		return nil, nil, fmt.Errorf("Azure数据库连接测试失败: %v", err)
	}

	// 连接AWS数据库
	awsDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		awsInstance.User, awsInstance.Password, awsInstance.Host, awsInstance.Database, awsInstance.Charset)

	awsConn, err := sql.Open("mysql", awsDSN)
	if err != nil {
		azureConn.Close()
		return nil, nil, fmt.Errorf("AWS数据库连接失败: %v", err)
	}

	if err := awsConn.Ping(); err != nil {
		azureConn.Close()
		awsConn.Close()
		return nil, nil, fmt.Errorf("AWS数据库连接测试失败: %v", err)
	}

	return azureConn, awsConn, nil
}

// getTableList 获取指定数据库的表列表
func (v *MultiDatabaseValidator) getTableList(conn *sql.DB, database string) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name"
	rows, err := conn.Query(query, database)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}

// calculateTableChecksum 计算表的校验和
func (v *MultiDatabaseValidator) calculateTableChecksum(conn *sql.DB, database, tableName string) (string, error) {
	// 获取表的行数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, tableName)
	var rowCount int
	if err := conn.QueryRow(countQuery).Scan(&rowCount); err != nil {
		return "", err
	}

	// 空表处理
	if rowCount == 0 {
		return "empty_table", nil
	}

	// 大表分批处理
	if rowCount > 100000 {
		return v.calculateLargeTableChecksum(conn, database, tableName, rowCount)
	}

	// 小表直接计算
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s` ORDER BY 1", database, tableName)
	rows, err := conn.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// 读取所有数据
	var data []string
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}

		rowData := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = "NULL"
			}
		}
		data = append(data, strings.Join(rowData, "|"))
	}

	// 计算MD5校验和
	dataStr := strings.Join(data, "\n")
	hash := md5.Sum([]byte(dataStr))
	return fmt.Sprintf("%x", hash), nil
}

// calculateLargeTableChecksum 大表分批计算校验和
func (v *MultiDatabaseValidator) calculateLargeTableChecksum(conn *sql.DB, database, tableName string, totalRows int) (string, error) {
	batchSize := 10000
	var checksums []string

	log.Printf("开始分批计算表 %s.%s，总行数: %d", database, tableName, totalRows)

	for offset := 0; offset < totalRows; offset += batchSize {
		query := fmt.Sprintf("SELECT * FROM `%s`.`%s` ORDER BY 1 LIMIT %d OFFSET %d",
			database, tableName, batchSize, offset)

		rows, err := conn.Query(query)
		if err != nil {
			return "", err
		}

		// 获取列信息
		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			return "", err
		}

		// 读取批次数据
		var batchData []string
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				rows.Close()
				return "", err
			}

			rowData := make([]string, len(columns))
			for i, val := range values {
				if val != nil {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = "NULL"
				}
			}
			batchData = append(batchData, strings.Join(rowData, "|"))
		}
		rows.Close()

		// 计算批次校验和
		batchStr := strings.Join(batchData, "\n")
		hash := md5.Sum([]byte(batchStr))
		checksums = append(checksums, fmt.Sprintf("%x", hash))
	}

	// 合并所有批次的校验和
	combinedStr := strings.Join(checksums, "")
	hash := md5.Sum([]byte(combinedStr))
	log.Printf("表 %s.%s 分批计算完成，共 %d 个批次", database, tableName, len(checksums))

	return fmt.Sprintf("%x", hash), nil
}

// GenerateReport 生成验证报告
func (v *MultiDatabaseValidator) GenerateReport(outputFile string) (*types.ValidationSummary, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// 统计验证结果
	totalDatabases := len(v.config.Azure)
	successfulValidations := 0
	inconsistentDatabases := 0
	errorDatabases := 0

	for _, result := range v.results {
		switch result.Status {
		case "SUCCESS":
			successfulValidations++
		case "INCONSISTENT":
			inconsistentDatabases++
		case "ERROR":
			errorDatabases++
		}
	}

	// 生成报告摘要
	summary := &types.ValidationSummary{
		Timestamp:             time.Now().Format(time.RFC3339),
		TotalDatabases:        totalDatabases,
		SuccessfulValidations: successfulValidations,
		InconsistentDatabases: inconsistentDatabases,
		ErrorDatabases:        errorDatabases,
		SuccessRate:           fmt.Sprintf("%.2f%%", float64(successfulValidations)/float64(totalDatabases)*100),
		Results:               v.results,
	}

	// 保存报告到文件
	if err := saveReportToFile(summary, outputFile); err != nil {
		return nil, fmt.Errorf("保存报告失败: %v", err)
	}

	log.Printf("验证报告已生成: %s", outputFile)
	return summary, nil
}

// saveReportToFile 保存报告到文件
func saveReportToFile(summary *types.ValidationSummary, filename string) error {
	// 使用JSON格式保存完整报告
	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入JSON格式的报告
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("写入报告文件失败: %v", err)
	}

	return nil
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
