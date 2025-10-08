# 数据库迁移经验

这个目录包含了各种数据库迁移项目的经验总结、工具使用和最佳实践。

## 内容概览

### 迁移案例
- [MySQL从Azure到AWS迁移数据一致性对比](./mysql-azure-to-aws-consistency.md) - 详细的MySQL云平台迁移指南

### 通用迁移工具

#### 1. 数据一致性验证工具
- **pt-table-checksum**: Percona Toolkit中的数据一致性检查工具
- **mysqldbcompare**: MySQL官方数据库对比工具
- **自建Python脚本**: 可定制的数据验证脚本
- **自建Go脚本**: 高性能的多数据库一致性验证工具

#### 2. 迁移工具
- **mysqldump**: MySQL官方备份恢复工具
- **Percona XtraBackup**: 物理备份工具
- **AWS DMS**: AWS数据迁移服务
- **Azure Database Migration Service**: Azure数据迁移服务

### 迁移策略

#### 1. 零停机迁移
- 主从复制 + 切换
- 双写策略
- 实时同步

#### 2. 停机迁移
- 全量备份恢复
- 增量数据同步
- 应用切换

### 常见问题解决

#### 字符集问题
- UTF-8 vs UTF-8MB4
- 排序规则差异
- 特殊字符处理

#### 时区问题
- 时区设置统一
- 时间戳数据迁移
- 应用程序时区处理

#### 性能问题
- 大表迁移策略
- 索引重建优化
- 网络带宽限制

### 验证清单

#### 迁移前
- [ ] 环境兼容性检查
- [ ] 数据量评估
- [ ] 网络连通性测试
- [ ] 权限配置确认

#### 迁移中
- [ ] 实时监控数据同步
- [ ] 性能指标监控
- [ ] 错误日志检查
- [ ] 业务功能验证

#### 迁移后
- [ ] 数据完整性验证
- [ ] 性能基准测试
- [ ] 应用程序测试
- [ ] 备份恢复测试

## 经验总结

### 成功要素
1. **充分的前期准备**：详细的环境评估和迁移计划
2. **分阶段执行**：降低风险，便于问题定位
3. **实时监控**：及时发现和处理异常情况
4. **回滚准备**：确保可以快速恢复到原始状态
5. **团队协作**：开发、运维、测试团队密切配合

### 常见陷阱
1. **忽略字符集差异**：导致数据乱码或查询异常
2. **时区设置不当**：影响时间相关业务逻辑
3. **大表迁移策略不当**：导致迁移时间过长或失败
4. **权限配置遗漏**：影响应用程序正常访问
5. **网络带宽限制**：影响迁移速度和成功率

### 最佳实践
1. **制定详细的迁移计划**：包括时间安排、人员分工、风险控制
2. **建立监控告警机制**：实时监控迁移进度和数据一致性
3. **准备回滚方案**：确保在出现问题时可以快速恢复
4. **充分测试验证**：在正式迁移前进行充分的测试
5. **文档记录完整**：详细记录迁移过程和经验教训

## 工具使用

### Python版本多数据库一致性验证工具

```bash
# 创建虚拟环境（推荐）
python3 -m venv venv
source venv/bin/activate

# 安装依赖
pip install pymysql

# 创建默认配置文件
python3 multi_database_validator.py --init

# 修改config.json中的数据库连接信息

# 运行验证
python3 multi_database_validator.py --config config.json

# 查看帮助
python3 multi_database_validator.py --help
```

### Go版本多数据库一致性验证工具

```bash
# 进入Go项目目录
cd go-validator

# 初始化Go模块
go mod init multi-database-validator
go mod tidy

# 编译
go build -o validator

# 创建默认配置文件
./validator init

# 修改config.json中的数据库连接信息

# 运行验证
./validator --config config.json

# 查看帮助
./validator help
```

## 相关资源

- [MySQL官方迁移指南](https://dev.mysql.com/doc/mysql-backup-excerpt/8.0/en/)
- [Percona Toolkit文档](https://www.percona.com/doc/percona-toolkit/)
- [AWS RDS迁移指南](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/)
- [Azure Database迁移指南](https://docs.microsoft.com/en-us/azure/database/)
