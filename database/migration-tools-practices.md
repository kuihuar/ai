# 数据库表迁移管理工具与实践

## 目录
- [概述](#概述)
- [主流迁移工具](#主流迁移工具)
- [ORM 框架内置迁移](#orm-框架内置迁移)
- [云平台迁移服务](#云平台迁移服务)
- [代码生成式迁移](#代码生成式迁移)
- [最佳实践](#最佳实践)
- [工具对比](#工具对比)

---

## 概述

数据库迁移管理是软件开发中的关键环节，确保数据库 Schema 变更能够：
- **版本化控制**：跟踪每个 Schema 变更
- **可重复执行**：在多个环境一致执行
- **可回滚**：支持回退到之前的版本
- **团队协作**：多人协作时避免冲突

---

## 主流迁移工具

### 1. Flyway

**语言支持**：Java、.NET、Go、Python、Node.js 等

#### 特点
- ✅ 基于 SQL 脚本的迁移
- ✅ 支持版本控制（V1__Create_user_table.sql）
- ✅ 自动检测已执行的迁移
- ✅ 支持回滚（Pro 版本）
- ✅ 支持多种数据库

#### 使用示例

```bash
# 安装
mvn flyway:migrate

# 或使用 CLI
flyway migrate
```

**迁移文件命名规范**：
```
V1__Create_user_table.sql
V2__Add_email_to_user.sql
V3__Create_post_table.sql
```

**配置文件** `flyway.conf`：
```properties
flyway.url=jdbc:mysql://localhost:3306/mydb
flyway.user=root
flyway.password=password
flyway.locations=filesystem:./migrations
```

#### 优势
- 简单易用，学习曲线低
- 社区版本免费
- 支持多种数据库
- 迁移历史表自动管理

#### 劣势
- 社区版不支持回滚
- 需要手写 SQL
- 大型项目可能文件较多

---

### 2. Liquibase

**语言支持**：Java、.NET、Python、Node.js 等

#### 特点
- ✅ 支持多种格式（SQL、XML、YAML、JSON）
- ✅ 变更集（Changeset）概念
- ✅ 支持回滚脚本
- ✅ 数据库无关性
- ✅ 复杂的变更逻辑支持

#### 使用示例

**XML 格式** `changelog.xml`：
```xml
<?xml version="1.0" encoding="UTF-8"?>
<databaseChangeLog
    xmlns="http://www.liquibase.org/xml/ns/dbchangelog"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog
    http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-3.8.xsd">

    <changeSet id="1" author="developer">
        <createTable tableName="user">
            <column name="id" type="int" autoIncrement="true">
                <constraints primaryKey="true"/>
            </column>
            <column name="name" type="varchar(255)">
                <constraints nullable="false"/>
            </column>
            <column name="email" type="varchar(255)">
                <constraints nullable="false" unique="true"/>
            </column>
        </createTable>
    </changeSet>
    
    <changeSet id="2" author="developer">
        <addColumn tableName="user">
            <column name="age" type="int"/>
        </addColumn>
    </changeSet>
</databaseChangeLog>
```

**YAML 格式** `changelog.yaml`：
```yaml
databaseChangeLog:
  - changeSet:
      id: 1
      author: developer
      changes:
        - createTable:
            tableName: user
            columns:
              - column:
                  name: id
                  type: int
                  autoIncrement: true
                  constraints:
                    primaryKey: true
              - column:
                  name: name
                  type: varchar(255)
                  constraints:
                    nullable: false
```

**执行迁移**：
```bash
liquibase --changeLogFile=changelog.xml update
liquibase --changeLogFile=changelog.xml rollback 1
```

#### 优势
- 数据库无关的变更定义
- 支持复杂的变更逻辑
- 内置回滚支持
- 变更集（Changeset）粒度控制

#### 劣势
- 学习曲线较陡
- XML/YAML 配置可能冗长
- 对于简单变更可能过度设计

---

### 3. Alembic（Python/SQLAlchemy）

**语言支持**：Python

#### 特点
- ✅ SQLAlchemy 官方迁移工具
- ✅ 自动生成迁移脚本
- ✅ 支持版本控制
- ✅ 支持升级和降级

#### 使用示例

```bash
# 初始化迁移环境
alembic init migrations

# 自动生成迁移脚本
alembic revision --autogenerate -m "create user table"

# 执行迁移
alembic upgrade head

# 回滚
alembic downgrade -1
```

**迁移文件** `migrations/versions/001_create_user.py`：
```python
"""create user table

Revision ID: 001
Revises: 
Create Date: 2024-01-01 10:00:00.000000

"""
from alembic import op
import sqlalchemy as sa

def upgrade():
    op.create_table('user',
        sa.Column('id', sa.Integer(), nullable=False),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('email', sa.String(255), nullable=False),
        sa.PrimaryKeyConstraint('id'),
        sa.UniqueConstraint('email')
    )

def downgrade():
    op.drop_table('user')
```

#### 优势
- 与 SQLAlchemy 深度集成
- 自动生成迁移脚本
- Python 生态系统支持好
- 支持复杂的迁移逻辑

#### 劣势
- 仅支持 Python 项目
- 需要 SQLAlchemy 模型

---

### 4. TypeORM Migrations（Node.js/TypeScript）

**语言支持**：TypeScript/JavaScript

#### 特点
- ✅ TypeORM 官方迁移工具
- ✅ 基于 TypeScript 类定义
- ✅ 自动生成迁移
- ✅ 支持事务

#### 使用示例

```bash
# 生成迁移
typeorm migration:generate -n CreateUserTable

# 运行迁移
typeorm migration:run

# 回滚迁移
typeorm migration:revert
```

**迁移文件** `migrations/1234567890-CreateUserTable.ts`：
```typescript
import { MigrationInterface, QueryRunner } from "typeorm";

export class CreateUserTable1234567890 implements MigrationInterface {
    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            CREATE TABLE "user" (
                "id" SERIAL NOT NULL,
                "name" VARCHAR(255) NOT NULL,
                "email" VARCHAR(255) NOT NULL UNIQUE,
                PRIMARY KEY ("id")
            )
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.dropTable("user");
    }
}
```

#### 优势
- TypeScript 类型安全
- 与 TypeORM 深度集成
- 支持事务
- 自动生成迁移

#### 劣势
- 仅支持 TypeScript/Node.js
- 需要 TypeORM 模型定义

---

### 5. Sequelize Migrations（Node.js）

**语言支持**：JavaScript/TypeScript

#### 特点
- ✅ Sequelize ORM 官方迁移工具
- ✅ 基于 JavaScript 迁移文件
- ✅ 支持 CLI 命令
- ✅ 版本控制

#### 使用示例

```bash
# 创建迁移
sequelize migration:generate --name create-user-table

# 运行迁移
sequelize db:migrate

# 回滚
sequelize db:migrate:undo
```

**迁移文件**：
```javascript
'use strict';

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('users', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: Sequelize.INTEGER
      },
      name: {
        type: Sequelize.STRING,
        allowNull: false
      },
      email: {
        type: Sequelize.STRING,
        allowNull: false,
        unique: true
      },
      createdAt: {
        allowNull: false,
        type: Sequelize.DATE
      },
      updatedAt: {
        allowNull: false,
        type: Sequelize.DATE
      }
    });
  },

  down: async (queryInterface, Sequelize) => {
    await queryInterface.dropTable('users');
  }
};
```

---

### 6. Rails Migrations（Ruby on Rails）

**语言支持**：Ruby

#### 特点
- ✅ Rails 框架内置
- ✅ 基于 Ruby DSL
- ✅ 自动生成迁移
- ✅ 版本控制

#### 使用示例

```bash
# 生成迁移
rails generate migration CreateUsers name:string email:string

# 运行迁移
rails db:migrate

# 回滚
rails db:rollback
```

**迁移文件**：
```ruby
class CreateUsers < ActiveRecord::Migration[7.0]
  def change
    create_table :users do |t|
      t.string :name, null: false
      t.string :email, null: false
      
      t.timestamps
    end
    
    add_index :users, :email, unique: true
  end
end
```

---

### 7. Django Migrations（Python）

**语言支持**：Python

#### 特点
- ✅ Django 框架内置
- ✅ 自动生成迁移
- ✅ 基于模型定义
- ✅ 支持数据迁移

#### 使用示例

```bash
# 生成迁移
python manage.py makemigrations

# 运行迁移
python manage.py migrate

# 回滚（需要指定版本）
python manage.py migrate app_name 0001
```

**迁移文件**：
```python
from django.db import migrations, models

class Migration(migrations.Migration):
    dependencies = []

    operations = [
        migrations.CreateModel(
            name='User',
            fields=[
                ('id', models.AutoField(primary_key=True)),
                ('name', models.CharField(max_length=255)),
                ('email', models.EmailField(unique=True)),
            ],
        ),
    ]
```

---

## ORM 框架内置迁移

### 1. Ent Migrations（Go）

**特点**：
- ✅ Ent 框架内置迁移系统
- ✅ 基于 Schema 定义自动生成
- ✅ 类型安全
- ✅ 支持版本控制

**使用示例**：

```go
// entc.go - 生成迁移代码
package main

import (
    "log"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
)

func main() {
    err := entc.Generate("./schema", &gen.Config{})
    if err != nil {
        log.Fatalln(err)
    }
}

// 生成迁移
//go:generate go run entc.go
```

```go
// 使用迁移
import (
    "entgo.io/ent/dialect/sql/schema"
)

func main() {
    client, err := ent.Open("mysql", "root:password@tcp(localhost:3306)/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    // 自动迁移
    if err := client.Schema.Create(ctx); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
}
```

---

### 2. GORM AutoMigrate（Go）

**特点**：
- ✅ 自动迁移
- ✅ 基于模型定义
- ✅ 简单易用

**使用示例**：

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Email string `gorm:"uniqueIndex"`
}

func main() {
    db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    
    // 自动迁移
    db.AutoMigrate(&User{})
}
```

**注意**：GORM 的 AutoMigrate 主要用于开发环境，生产环境建议使用专门的迁移工具。

---

## 云平台迁移服务

### 1. AWS Database Migration Service (DMS)

**特点**：
- ✅ 零停机迁移
- ✅ 支持多种数据库
- ✅ 自动监控
- ✅ 数据验证

**适用场景**：
- 云平台间迁移
- 数据库升级
- 架构迁移

---

### 2. Azure Database Migration Service

**特点**：
- ✅ Azure 官方迁移服务
- ✅ 评估工具
- ✅ 实时迁移
- ✅ 数据验证

---

### 3. Google Cloud Database Migration Service

**特点**：
- ✅ Google Cloud 官方服务
- ✅ 多种数据库支持
- ✅ 实时同步
- ✅ 监控和告警

---

## 代码生成式迁移

### 1. Prisma Migrate（TypeScript/Node.js）

**特点**：
- ✅ 基于 Prisma Schema 自动生成
- ✅ 类型安全
- ✅ 开发和生产环境一致

**使用示例**：

```bash
# 生成迁移
npx prisma migrate dev --name create_user

# 应用迁移
npx prisma migrate deploy

# 重置数据库
npx prisma migrate reset
```

**Schema 定义** `schema.prisma`：
```prisma
model User {
  id    Int    @id @default(autoincrement())
  name  String
  email String @unique
  
  @@map("users")
}
```

---

### 2. Hasura Migrations

**特点**：
- ✅ Hasura GraphQL 引擎
- ✅ 基于元数据
- ✅ 版本控制

**使用示例**：

```bash
# 导出迁移
hasura migrate create "create_user_table" --from-server

# 应用迁移
hasura migrate apply
```

---

## 最佳实践

### 1. 迁移文件管理

**命名规范**：
```
V{version}__{description}.sql
或
{timestamp}_{description}.sql
```

**示例**：
```
V1__Create_user_table.sql
V2__Add_email_index.sql
20240101120000_Create_user_table.sql
```

### 2. 版本控制

- ✅ 所有迁移文件纳入 Git 版本控制
- ✅ 使用语义化版本号
- ✅ 每个迁移文件包含变更描述

### 3. 迁移脚本编写

**原则**：
- ✅ 幂等性：多次执行结果一致
- ✅ 可回滚：提供回滚脚本
- ✅ 事务性：使用事务保证原子性
- ✅ 数据安全：迁移前备份数据

**示例（幂等性）**：
```sql
-- 好的做法
CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL
);

-- 避免的做法
CREATE TABLE users (...);  -- 如果表已存在会报错
```

### 4. 环境管理

**不同环境的策略**：
- **开发环境**：自动迁移，允许重置
- **测试环境**：自动迁移，定期重置
- **生产环境**：手动审核，逐步执行

### 5. 大表迁移策略

**分阶段迁移**：
1. 创建新表结构
2. 数据迁移（批量处理）
3. 数据验证
4. 切换表名
5. 清理旧表

**示例**：
```sql
-- 阶段1：创建新表
CREATE TABLE users_new LIKE users;
ALTER TABLE users_new ADD COLUMN email VARCHAR(255);

-- 阶段2：数据迁移（分批）
INSERT INTO users_new (id, name, email)
SELECT id, name, CONCAT(name, '@example.com') 
FROM users 
WHERE id BETWEEN 1 AND 10000;

-- 阶段3：验证数据
SELECT COUNT(*) FROM users_new;
SELECT COUNT(*) FROM users;

-- 阶段4：切换（需要停机）
RENAME TABLE users TO users_old, users_new TO users;

-- 阶段5：清理
DROP TABLE users_old;
```

### 6. 迁移测试

**测试清单**：
- ✅ 单元测试：测试迁移脚本逻辑
- ✅ 集成测试：在测试环境执行迁移
- ✅ 回滚测试：验证回滚脚本
- ✅ 性能测试：大表迁移性能验证

### 7. 监控和告警

**关键指标**：
- 迁移执行时间
- 迁移成功率
- 数据一致性
- 性能影响

### 8. 团队协作

**工作流程**：
1. 开发人员在分支中创建迁移
2. 代码审查迁移脚本
3. 在测试环境验证
4. 合并到主分支
5. 生产环境执行（需要审批）

---

## 工具对比

| 工具 | 语言 | 学习曲线 | 回滚支持 | 自动生成 | 数据库支持 | 适用场景 |
|------|------|----------|----------|----------|------------|----------|
| **Flyway** | 多语言 | ⭐⭐ | Pro版 | ❌ | 广泛 | Java/Spring项目 |
| **Liquibase** | 多语言 | ⭐⭐⭐ | ✅ | ❌ | 广泛 | 复杂迁移需求 |
| **Alembic** | Python | ⭐⭐ | ✅ | ✅ | SQLAlchemy支持 | Python/SQLAlchemy |
| **TypeORM** | TypeScript | ⭐⭐ | ✅ | ✅ | 广泛 | Node.js/TypeScript |
| **Sequelize** | JavaScript | ⭐⭐ | ✅ | ✅ | 广泛 | Node.js项目 |
| **Rails** | Ruby | ⭐ | ✅ | ✅ | Rails支持 | Ruby on Rails |
| **Django** | Python | ⭐ | ✅ | ✅ | Django支持 | Django项目 |
| **Ent** | Go | ⭐⭐⭐ | ✅ | ✅ | MySQL/PostgreSQL/SQLite | Go/Ent项目 |
| **Prisma** | TypeScript | ⭐⭐ | ✅ | ✅ | 广泛 | Node.js/Prisma |
| **GORM AutoMigrate** | Go | ⭐ | ❌ | ✅ | 广泛 | 开发环境 |

---

## 选择建议

### 根据项目类型选择

1. **Java/Spring 项目**
   - 推荐：Flyway 或 Liquibase
   - 理由：生态成熟，社区支持好

2. **Python 项目**
   - SQLAlchemy：Alembic
   - Django：Django Migrations

3. **Node.js/TypeScript 项目**
   - TypeORM：TypeORM Migrations
   - Prisma：Prisma Migrate
   - Sequelize：Sequelize Migrations

4. **Go 项目**
   - Ent：Ent Migrations
   - GORM：建议配合 Flyway 或自建迁移工具

5. **Ruby 项目**
   - Rails：Rails Migrations

### 根据需求选择

1. **简单项目**
   - Flyway、Rails Migrations、Django Migrations

2. **复杂迁移需求**
   - Liquibase（支持复杂逻辑）

3. **ORM 深度集成**
   - 使用对应 ORM 的迁移工具

4. **跨数据库项目**
   - Liquibase、Prisma

---

## 总结

数据库迁移管理是软件开发的重要环节，选择合适的工具和遵循最佳实践能够：

- ✅ 提高开发效率
- ✅ 降低迁移风险
- ✅ 保证数据一致性
- ✅ 便于团队协作

**关键原则**：
1. **版本化**：所有变更纳入版本控制
2. **可重复**：迁移可以在任何环境执行
3. **可回滚**：支持回退到之前版本
4. **自动化**：尽可能自动化执行
5. **测试**：充分测试迁移和回滚脚本

选择合适的工具并遵循最佳实践，可以确保数据库迁移的顺利进行。

---

## 参考资料

- [Flyway 官方文档](https://flywaydb.org/documentation/)
- [Liquibase 官方文档](https://docs.liquibase.com/)
- [Alembic 官方文档](https://alembic.sqlalchemy.org/)
- [Ent Migrations 文档](https://entgo.io/docs/migrate/)
- [Prisma Migrate 文档](https://www.prisma.io/docs/concepts/components/prisma-migrate)

