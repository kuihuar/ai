// validator.go
// 多数据库验证器核心逻辑

package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MultiDatabaseValidator 多数据库验证器
type MultiDatabaseValidator struct {
	databasePairs []DatabasePair
	results       map[string]DatabaseResult
	mutex         sync.RWMutex
}

// NewMultiDatabaseValidator 创建新的验证器实例
func NewMultiDatabaseValidator(databasePairs []DatabasePair) *MultiDatabaseValidator {
	// 验证配置参数
	if err := validateDatabasePairs(databasePairs); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	return &MultiDatabaseValidator{
		databasePairs: databasePairs,
		results:       make(map[string]DatabaseResult),
	}
}

// validateDatabasePairs 验证数据库对比对配置的有效性
func validateDatabasePairs(databasePairs []DatabasePair) error {
	if len(databasePairs) == 0 {
		return fmt.Errorf("数据库对比对列表不能为空")
	}

	for i, pair := range databasePairs {
		// 检查Azure实例配置
		if pair.AzureInstance.Host == "" || pair.AzureInstance.User == "" || pair.AzureInstance.Password == "" || pair.AzureInstance.Database == "" {
			return fmt.Errorf("第%d个对比对的Azure实例配置缺少必要参数", i+1)
		}

		// 检查AWS实例配置
		if pair.AWSInstance.Host == "" || pair.AWSInstance.User == "" || pair.AWSInstance.Password == "" || pair.AWSInstance.Database == "" {
			return fmt.Errorf("第%d个对比对的AWS实例配置缺少必要参数", i+1)
		}
	}

	return nil
}

// getConnectionString 生成数据库连接字符串
func (v *MultiDatabaseValidator) getConnectionString(instance DatabaseInstance) string {
	charset := instance.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=%s&parseTime=true&loc=Local",
		instance.User, instance.Password, instance.Host, instance.Database, charset)
}

// getTableList 获取指定数据库的表列表
func (v *MultiDatabaseValidator) getTableList(db *sql.DB, database string) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = ? 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.Query(query, database)
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("扫描表名失败: %v", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历表列表失败: %v", err)
	}

	log.Printf("数据库 %s 包含 %d 个表", database, len(tables))
	return tables, nil
}

// calculateTableChecksum 计算表的校验和
func (v *MultiDatabaseValidator) calculateTableChecksum(db *sql.DB, database, tableName string) (string, error) {
	// 获取表的行数
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, tableName)
	if err := db.QueryRow(countQuery).Scan(&rowCount); err != nil {
		return "", fmt.Errorf("获取表行数失败: %v", err)
	}

	// 空表处理
	if rowCount == 0 {
		log.Printf("表 %s.%s 为空表", database, tableName)
		return "empty_table", nil
	}

	// 大表分批处理，避免内存溢出
	if rowCount > 100000 {
		log.Printf("表 %s.%s 行数较多(%d)，使用分批计算", database, tableName, rowCount)
		return v.calculateLargeTableChecksum(db, database, tableName)
	}

	// 小表直接计算
	log.Printf("计算表 %s.%s 的校验和，行数: %d", database, tableName, rowCount)
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s` ORDER BY 1", database, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("查询表数据失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("获取列信息失败: %v", err)
	}

	// 读取所有数据
	var allData []string
	for rows.Next() {
		// 创建扫描目标
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", fmt.Errorf("扫描行数据失败: %v", err)
		}

		// 将数据转换为字符串
		rowData := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = "NULL"
			}
		}

		// 将行数据添加到总数据中
		rowStr := fmt.Sprintf("%v", rowData)
		allData = append(allData, rowStr)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("遍历表数据失败: %v", err)
	}

	// 计算MD5校验和
	dataStr := fmt.Sprintf("%v", allData)
	hash := md5.Sum([]byte(dataStr))
	return fmt.Sprintf("%x", hash), nil
}

// calculateLargeTableChecksum 大表分批计算校验和
func (v *MultiDatabaseValidator) calculateLargeTableChecksum(db *sql.DB, database, tableName string) (string, error) {
	// 获取总行数
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, tableName)
	if err := db.QueryRow(countQuery).Scan(&totalRows); err != nil {
		return "", fmt.Errorf("获取表总行数失败: %v", err)
	}

	batchSize := 10000 // 每批处理10000行
	var checksums []string

	log.Printf("开始分批计算表 %s.%s，总行数: %d", database, tableName, totalRows)

	// 分批读取数据
	for offset := 0; offset < totalRows; offset += batchSize {
		log.Printf("处理批次: %d-%d", offset, min(offset+batchSize, totalRows))

		query := fmt.Sprintf(`
			SELECT * FROM `+"`%s`"+`.`+"`%s`"+` 
			ORDER BY 1 
			LIMIT %d OFFSET %d
		`, database, tableName, batchSize, offset)

		rows, err := db.Query(query)
		if err != nil {
			return "", fmt.Errorf("查询批次数据失败: %v", err)
		}

		// 获取列信息
		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			return "", fmt.Errorf("获取列信息失败: %v", err)
		}

		// 读取批次数据
		var batchData []string
		for rows.Next() {
			// 创建扫描目标
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			// 扫描行数据
			if err := rows.Scan(valuePtrs...); err != nil {
				rows.Close()
				return "", fmt.Errorf("扫描批次数据失败: %v", err)
			}

			// 将数据转换为字符串
			rowData := make([]string, len(columns))
			for i, val := range values {
				if val != nil {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = "NULL"
				}
			}

			// 将行数据添加到批次数据中
			rowStr := fmt.Sprintf("%v", rowData)
			batchData = append(batchData, rowStr)
		}

		rows.Close()

		if err := rows.Err(); err != nil {
			return "", fmt.Errorf("遍历批次数据失败: %v", err)
		}

		// 计算批次数据的校验和
		batchStr := fmt.Sprintf("%v", batchData)
		hash := md5.Sum([]byte(batchStr))
		checksums = append(checksums, fmt.Sprintf("%x", hash))
	}

	// 合并所有批次的校验和
	combinedStr := ""
	for _, checksum := range checksums {
		combinedStr += checksum
	}

	combinedHash := md5.Sum([]byte(combinedStr))
	log.Printf("表 %s.%s 分批计算完成，共 %d 个批次", database, tableName, len(checksums))

	return fmt.Sprintf("%x", combinedHash), nil
}

// validateDatabase 验证单个数据库的一致性
func (v *MultiDatabaseValidator) validateDatabase(pair DatabasePair) DatabaseResult {
	log.Printf("开始验证数据库对比: %s (Azure: %s) vs %s (AWS: %s)",
		pair.AzureInstance.Database, pair.AzureInstance.Name,
		pair.AWSInstance.Database, pair.AWSInstance.Name)

	startTime := time.Now()
	result := DatabaseResult{
		Database:         pair.AzureInstance.Database,
		AzureInstance:    pair.AzureInstance.Name,
		AWSInstance:      pair.AWSInstance.Name,
		Status:           "SUCCESS",
		Errors:           []string{},
		TableComparisons: []TableComparison{},
		StartTime:        startTime.Format(time.RFC3339),
	}

	// 连接数据库
	azureConnStr := v.getConnectionString(pair.AzureInstance)
	awsConnStr := v.getConnectionString(pair.AWSInstance)

	azureDB, err := sql.Open("mysql", azureConnStr)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("连接Azure数据库失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}
	defer azureDB.Close()

	awsDB, err := sql.Open("mysql", awsConnStr)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("连接AWS数据库失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}
	defer awsDB.Close()

	// 测试连接
	if err := azureDB.Ping(); err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("Azure数据库连接测试失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	if err := awsDB.Ping(); err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("AWS数据库连接测试失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	// 获取表列表
	azureTables, err := v.getTableList(azureDB, pair.AzureInstance.Database)
	if err != nil {
		result.Status = "ERROR"
		result.Errors = append(result.Errors, fmt.Sprintf("获取Azure表列表失败: %v", err))
		result.EndTime = time.Now().Format(time.RFC3339)
		return result
	}

	awsTables, err := v.getTableList(awsDB, pair.AWSInstance.Database)
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
		log.Printf("数据库 %s: %s", pair.AzureInstance.Database, errorMsg)
	}

	// 对比每个表的数据一致性
	log.Printf("开始验证数据库 %s 中的 %d 个表", pair.AzureInstance.Database, len(azureTables))

	for i, table := range azureTables {
		log.Printf("验证表 %d/%d: %s", i+1, len(azureTables), table)

		// 检查表是否在AWS中存在
		tableExists := false
		for _, awsTable := range awsTables {
			if awsTable == table {
				tableExists = true
				break
			}
		}

		if !tableExists {
			errorMsg := fmt.Sprintf("表 %s 在AWS中不存在", table)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "INCONSISTENT"
			log.Printf("数据库 %s: %s", pair.AzureInstance.Database, errorMsg)
			continue
		}

		// 计算Azure表的校验和
		azureChecksum, err := v.calculateTableChecksum(azureDB, pair.AzureInstance.Database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("计算Azure表 %s 校验和失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", pair.AzureInstance.Database, errorMsg)
			continue
		}

		// 计算AWS表的校验和
		awsChecksum, err := v.calculateTableChecksum(awsDB, pair.AWSInstance.Database, table)
		if err != nil {
			errorMsg := fmt.Sprintf("计算AWS表 %s 校验和失败: %v", table, err)
			result.Errors = append(result.Errors, errorMsg)
			result.Status = "ERROR"
			log.Printf("数据库 %s: %s", pair.AzureInstance.Database, errorMsg)
			continue
		}

		// 记录对比结果
		tableComparison := TableComparison{
			Table:         table,
			AzureChecksum: azureChecksum,
			AWSChecksum:   awsChecksum,
			Match:         azureChecksum == awsChecksum,
			AzureInstance: pair.AzureInstance.Name,
			AWSInstance:   pair.AWSInstance.Name,
			AzureDatabase: pair.AzureInstance.Database,
			AWSDatabase:   pair.AWSInstance.Database,
		}

		result.TableComparisons = append(result.TableComparisons, tableComparison)

		// 检查是否一致
		if azureChecksum != awsChecksum {
			result.Status = "INCONSISTENT"
			log.Printf("数据不一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s",
				pair.AzureInstance.Name, pair.AzureInstance.Database, table,
				pair.AWSInstance.Name, pair.AWSInstance.Database, table)
		} else {
			log.Printf("数据一致 - Azure实例: %s 数据库: %s 表: %s vs AWS实例: %s 数据库: %s 表: %s",
				pair.AzureInstance.Name, pair.AzureInstance.Database, table,
				pair.AWSInstance.Name, pair.AWSInstance.Database, table)
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

	log.Printf("数据库 %s 验证完成:", pair.AzureInstance.Database)
	log.Printf("  状态: %s", result.Status)
	log.Printf("  表数量: Azure(%d) vs AWS(%d)", result.AzureTables, result.AWSTables)
	log.Printf("  数据一致性: %d/%d", consistentTables, totalTables)
	log.Printf("  错误数量: %d", len(result.Errors))

	return result
}

// validateAllDatabases 并行验证所有数据库
func (v *MultiDatabaseValidator) validateAllDatabases(maxWorkers int) {
	log.Printf("开始验证 %d 个数据库对比对，最大并发数: %d", len(v.databasePairs), maxWorkers)

	// 使用goroutine和channel实现并发控制
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, pair := range v.databasePairs {
		wg.Add(1)
		go func(p DatabasePair) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行验证
			result := v.validateDatabase(p)

			// 保存结果，使用数据库名称作为key
			v.mutex.Lock()
			v.results[p.AzureInstance.Database] = result
			v.mutex.Unlock()

			log.Printf("数据库对比 %s vs %s 验证完成，状态: %s",
				p.AzureInstance.Database, p.AWSInstance.Database, result.Status)
		}(pair)
	}

	// 等待所有验证完成
	wg.Wait()
	log.Printf("所有数据库验证完成")
}

// generateReport 生成验证报告
func (v *MultiDatabaseValidator) generateReport(outputFile string) (*ValidationSummary, error) {
	// 统计验证结果
	totalDatabases := len(v.databasePairs)
	successfulValidations := 0
	inconsistentDatabases := 0
	errorDatabases := 0

	v.mutex.RLock()
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
	v.mutex.RUnlock()

	// 计算成功率
	successRate := "0%"
	if totalDatabases > 0 {
		successRate = fmt.Sprintf("%.2f%%", float64(successfulValidations)/float64(totalDatabases)*100)
	}

	// 生成报告摘要
	summary := &ValidationSummary{
		Timestamp:             time.Now().Format(time.RFC3339),
		TotalDatabases:        totalDatabases,
		SuccessfulValidations: successfulValidations,
		InconsistentDatabases: inconsistentDatabases,
		ErrorDatabases:        errorDatabases,
		SuccessRate:           successRate,
		Results:               make(map[string]DatabaseResult),
	}

	// 复制结果
	v.mutex.RLock()
	for db, result := range v.results {
		summary.Results[db] = result
	}
	v.mutex.RUnlock()

	// 保存报告到文件
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("创建报告文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		return nil, fmt.Errorf("写入报告文件失败: %v", err)
	}

	log.Printf("验证报告已生成: %s", outputFile)
	return summary, nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
