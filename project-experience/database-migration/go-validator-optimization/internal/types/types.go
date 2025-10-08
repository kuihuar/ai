// internal/types/types.go
// 共享类型定义

package types

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	User     string `json:"user" yaml:"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	Charset  string `json:"charset" yaml:"charset" mapstructure:"charset"`
}

// DatabaseInstance 数据库实例配置
type DatabaseInstance struct {
	Name     string `json:"name" yaml:"name" mapstructure:"name"`             // 实例名称
	Host     string `json:"host" yaml:"host" mapstructure:"host"`             // 实例主机地址
	User     string `json:"user" yaml:"user" mapstructure:"user"`             // 用户名
	Password string `json:"password" yaml:"password" mapstructure:"password"` // 密码
	Database string `json:"database" yaml:"database" mapstructure:"database"` // 数据库名称
	Charset  string `json:"charset" yaml:"charset" mapstructure:"charset"`    // 字符集
}

// TableComparison 表对比结果
type TableComparison struct {
	Table         string `json:"table" yaml:"table" mapstructure:"table"`
	AzureChecksum string `json:"azure_checksum" yaml:"azure_checksum" mapstructure:"azure_checksum"`
	AWSChecksum   string `json:"aws_checksum" yaml:"aws_checksum" mapstructure:"aws_checksum"`
	Match         bool   `json:"match" yaml:"match" mapstructure:"match"`
	AzureInstance string `json:"azure_instance" yaml:"azure_instance" mapstructure:"azure_instance"`
	AWSInstance   string `json:"aws_instance" yaml:"aws_instance" mapstructure:"aws_instance"`
	AzureDatabase string `json:"azure_database" yaml:"azure_database" mapstructure:"azure_database"`
	AWSDatabase   string `json:"aws_database" yaml:"aws_database" mapstructure:"aws_database"`
}

// DatabaseResult 数据库验证结果
type DatabaseResult struct {
	Database         string            `json:"database" yaml:"database" mapstructure:"database"`
	AzureInstance    string            `json:"azure_instance" yaml:"azure_instance" mapstructure:"azure_instance"` // Azure实例名称
	AWSInstance      string            `json:"aws_instance" yaml:"aws_instance" mapstructure:"aws_instance"`       // AWS实例名称
	AzureTables      int               `json:"azure_tables" yaml:"azure_tables" mapstructure:"azure_tables"`
	AWSTables        int               `json:"aws_tables" yaml:"aws_tables" mapstructure:"aws_tables"`
	TableComparisons []TableComparison `json:"table_comparisons" yaml:"table_comparisons" mapstructure:"table_comparisons"`
	Status           string            `json:"status" yaml:"status" mapstructure:"status"`
	Errors           []string          `json:"errors" yaml:"errors" mapstructure:"errors"`
	StartTime        string            `json:"start_time" yaml:"start_time" mapstructure:"start_time"`
	EndTime          string            `json:"end_time" yaml:"end_time" mapstructure:"end_time"`
}

// ValidationSummary 验证摘要
type ValidationSummary struct {
	Timestamp             string                    `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
	TotalDatabases        int                       `json:"total_databases" yaml:"total_databases" mapstructure:"total_databases"`
	SuccessfulValidations int                       `json:"successful_validations" yaml:"successful_validations" mapstructure:"successful_validations"`
	InconsistentDatabases int                       `json:"inconsistent_databases" yaml:"inconsistent_databases" mapstructure:"inconsistent_databases"`
	ErrorDatabases        int                       `json:"error_databases" yaml:"error_databases" mapstructure:"error_databases"`
	SuccessRate           string                    `json:"success_rate" yaml:"success_rate" mapstructure:"success_rate"`
	Results               map[string]DatabaseResult `json:"results" yaml:"results" mapstructure:"results"`
}

// Config 配置文件结构
type Config struct {
	Azure      []DatabaseInstance `json:"azure" yaml:"azure" mapstructure:"azure"`                   // Azure实例列表
	AWS        []DatabaseInstance `json:"aws" yaml:"aws" mapstructure:"aws"`                         // AWS实例列表
	MaxWorkers int                `json:"max_workers" yaml:"max_workers" mapstructure:"max_workers"` // 最大并发数
}

// DatabasePair 数据库对比对
type DatabasePair struct {
	AzureInstance DatabaseInstance `json:"azure_instance" yaml:"azure_instance" mapstructure:"azure_instance"`
	AWSInstance   DatabaseInstance `json:"aws_instance" yaml:"aws_instance" mapstructure:"aws_instance"`
}

// ValidationOptions 验证选项
type ValidationOptions struct {
	ConfigFile    string
	OutputFile    string
	MaxWorkers    int
	Verbose       bool
	DryRun        bool
	AzureHost     string
	AzureUser     string
	AzurePassword string
	AzureDatabase string
	AWSHost       string
	AWSUser       string
	AWSPassword   string
	AWSDatabase   string
}
