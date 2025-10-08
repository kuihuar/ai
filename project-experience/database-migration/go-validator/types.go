// types.go
// 定义数据结构和类型

package main

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

// DatabaseInstance 数据库实例配置
type DatabaseInstance struct {
	Name     string `json:"name" yaml:"name"`         // 实例名称
	Host     string `json:"host" yaml:"host"`         // 实例主机地址
	User     string `json:"user" yaml:"user"`         // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名称
	Charset  string `json:"charset" yaml:"charset"`   // 字符集
}

// TableComparison 表对比结果
type TableComparison struct {
	Table         string `json:"table"`
	AzureChecksum string `json:"azure_checksum"`
	AWSChecksum   string `json:"aws_checksum"`
	Match         bool   `json:"match"`
	AzureInstance string `json:"azure_instance"`
	AWSInstance   string `json:"aws_instance"`
	AzureDatabase string `json:"azure_database"`
	AWSDatabase   string `json:"aws_database"`
}

// DatabaseResult 数据库验证结果
type DatabaseResult struct {
	Database         string            `json:"database"`
	AzureInstance    string            `json:"azure_instance"` // Azure实例名称
	AWSInstance      string            `json:"aws_instance"`   // AWS实例名称
	AzureTables      int               `json:"azure_tables"`
	AWSTables        int               `json:"aws_tables"`
	TableComparisons []TableComparison `json:"table_comparisons"`
	Status           string            `json:"status"`
	Errors           []string          `json:"errors"`
	StartTime        string            `json:"start_time"`
	EndTime          string            `json:"end_time"`
}

// ValidationSummary 验证摘要
type ValidationSummary struct {
	Timestamp             string                    `json:"timestamp"`
	TotalDatabases        int                       `json:"total_databases"`
	SuccessfulValidations int                       `json:"successful_validations"`
	InconsistentDatabases int                       `json:"inconsistent_databases"`
	ErrorDatabases        int                       `json:"error_databases"`
	SuccessRate           string                    `json:"success_rate"`
	Results               map[string]DatabaseResult `json:"results"`
}

// Config 配置文件结构
type Config struct {
	Azure      []DatabaseInstance `json:"azure" yaml:"azure"`             // Azure实例列表
	AWS        []DatabaseInstance `json:"aws" yaml:"aws"`                 // AWS实例列表
	MaxWorkers int                `json:"max_workers" yaml:"max_workers"` // 最大并发数
}

// DatabasePair 数据库对比对
type DatabasePair struct {
	AzureInstance DatabaseInstance `json:"azure_instance"`
	AWSInstance   DatabaseInstance `json:"aws_instance"`
}
