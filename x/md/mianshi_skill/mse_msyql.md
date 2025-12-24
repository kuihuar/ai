

### mysql事务处理，编码实现银行转帐操作，要求使用数据库事务保证数据一致性，并处理可能的死锁情况
- 
### 问题分析
1. 银行转帐操作涉及到两个账户的余额更新，需要保证数据一致性
2. 可能出现死锁情况，需要处理
3. 查看数据库的隔离级别，是否支持事务
4. 查看数据库的锁机制，是否支持行级锁
5. 查看数据库的事务隔离级别，是否支持可重复读
6. 锁
    - 乐观锁
      1. 避免冲突：在高并发环境下，乐观锁可以减少锁的竞争，避免死锁。
      2. 更新控制：通过版本号控制，确保在更新时数据没有被其他事务修改
    - 命名锁（GET_LOCK 获取锁，操作完成后使用 RELEASE_LOCK）
      1. 避免冲突：通过SET NX命令，设置锁，避免死锁
7. 场景
选择哪种机制取决于具体场景的需求。如果操作频繁且竞争激烈，乐观锁可能更合适；如果需要确保某些关键操作的独占性，则应使用命名锁。      
    
   
### 实现方式
1. 使用数据库事务，保证数据一致性
2. 处理死锁情况，使用SET NX命令，设置锁，避免死锁
### 实现代码
```sql
-- DECIMAL(15, 4) 的范围如下：
-- -9,999,999,999,999.9999 到 9,999,999,999,999.9999
CREATE TABLE accounts (
    account_id VARCHAR(36) PRIMARY KEY,  -- 账户ID
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,  -- 余额
    version INT NOT NULL DEFAULT 0  -- 乐观锁版本号
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- 更新时间戳
) ENGINE=InnoDB;

-- 检查当前时区
SELECT @@global.time_zone, @@session.time_zone;

-- 设置时区为 UTC
SET GLOBAL time_zone = '+00:00';
SET GLOBAL time_zone = '+00:00';
SELECT account_number, CONVERT_TZ(updated_at, '+00:00', 'America/New_York') AS updated_at_ny
FROM Accounts;

可以使用bigint类型存储时间戳，然后使用CONVERT_TZ函数进行时区转换。
或者使用DATETIME或TIMESTAMP
DATETIME：

存储为一个日期和时间的组合，范围从 1000-01-01 00:00:00 到 9999-12-31 23:59:59。
不会自动进行时区转换，存储的时间是按实际输入的值保存。
TIMESTAMP：

存储为自 1970-01-01 00:00:00 UTC以来的秒数，范围从 1970-01-01 00:00:01 UTC 到 2038-01-19 03:14:07 UTC。
会自动转换为 UTC 存储，并在查询时根据连接的时区进行转换


-- 交易记录表（可选，用于审计）
CREATE TABLE transactions (
    id VARCHAR(36) PRIMARY KEY,
    from_account VARCHAR(36) NOT NULL,
    to_account VARCHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE accounts ADD INDEX idx_account_id (account_id);
命名锁（GET_LOCK(只有一个事务可以操作)） > 行级锁（FOR UPDATE） > 悲观锁
-- 银行转帐操作 
-- 1. 开启事务
-- 2. 查询账户余额
-- 3. 扣减账户余额
-- 4. 查询账户余额
-- 5. 增加账户余额
-- 6. 提交事务
-- 7. 处理死锁情况，使用SET NX命令，设置锁，避免死锁
START TRANSACTION; # 开始事务
SELECT balance,version FROM accounts WHERE account_id = 123 FOR UPDATE; # 查询账户余额
UPDATE accounts SET balance = balance - 100 version=version+1 WHERE account_id = 123 AND version=111; # 扣减账户余额
SELECT balance,version FROM accounts WHERE account_id = 456 FOR UPDATE; # 查询账户余额
UPDATE accounts SET balance = balance + 100 WHERE account_id = 456; # 增加账户余额
COMMIT; # 提交事务
```

```sql
-- 使用一个存储过程实现银行转帐操作
-- 存储过程的参数包括发送者账户、接收者账户、转帐金额
-- 存储过程的实现步骤如下：
-- 1. 开启事务
-- 2. 锁定发送者账户
-- 3. 查询发送者账户余额
-- 4. 如果余额不足，抛出异常
-- 5. 更新发送者账户余额
-- 6. 更新接收者账户余额
-- 7. 提交事务
-- 8. 处理死锁情况，使用SET NX命令，设置锁，避免死锁
DELIMITER //

CREATE PROCEDURE TransferFunds(
    IN sender_account VARCHAR(20),
    IN receiver_account VARCHAR(20),
    IN amount DECIMAL(15, 2)
)
BEGIN
    DECLARE sender_balance DECIMAL(15, 2);
    DECLARE lock_acquired INT;

    -- 尝试获取锁
    SET lock_acquired = GET_LOCK(CONCAT('transfer_', sender_account), 10);

    IF lock_acquired THEN
        -- 开始事务
        START TRANSACTION;

        -- 锁定发送者账户
        SELECT balance INTO sender_balance
        FROM Accounts
        WHERE account_number = sender_account
        FOR UPDATE;

        IF sender_balance >= amount THEN
            -- 更新发送者账户余额
            UPDATE Accounts
            SET balance = balance - amount
            WHERE account_number = sender_account;

            -- 更新接收者账户余额
            UPDATE Accounts
            SET balance = balance + amount
            WHERE account_number = receiver_account;

            -- 提交事务
            COMMIT;
        ELSE
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient funds';
            ROLLBACK;
        END IF;

        -- 释放锁
        SELECT RELEASE_LOCK(CONCAT('transfer_', sender_account));
    ELSE
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Could not acquire lock';
    END IF;
END //

DELIMITER ;
```

```sql
DELIMITER //

CREATE PROCEDURE TransferFunds(
    IN sender_account VARCHAR(20),
    IN receiver_account VARCHAR(20),
    IN amount DECIMAL(15, 2)
)
BEGIN
    DECLARE sender_balance DECIMAL(15, 2);
    DECLARE current_version INT;

    -- 开始事务
    START TRANSACTION;

    -- 获取发送者账户的余额和版本号
    SELECT balance, version INTO sender_balance, current_version
    FROM Accounts
    WHERE account_number = sender_account
    FOR UPDATE;

    IF sender_balance >= amount THEN
        -- 更新发送者账户余额，检查版本号
        UPDATE Accounts
        SET balance = balance - amount, version = version + 1
        WHERE account_number = sender_account AND version = current_version;

        IF ROW_COUNT() = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Version conflict, please retry';
            ROLLBACK;
            LEAVE;
        END IF;

        -- 更新接收者账户余额
        UPDATE Accounts
        SET balance = balance + amount, version = version + 1
        WHERE account_number = receiver_account;

        IF ROW_COUNT() = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Receiver account not found';
            ROLLBACK;
            LEAVE;
        END IF;

        -- 提交事务
        COMMIT;
    ELSE
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient funds';
        ROLLBACK;
    END IF;
END //

DELIMITER ;

CALL TransferFunds('12345', '67890', 100.00);



```

### 事务的特性
1. 原子性：事务中的所有操作要么全部成功，要么全部失败
2. 一致性：事务执行前后，数据库的状态保持一致
3. 隔离性：多个事务并发执行时，相互之间不会影响
4. 持久性：事务提交后，对数据库的修改是永久的
### 死锁处理
1. 设置锁的超时时间，避免死锁
2. 避免循环依赖，合理设计事务执行顺序（排序帐户ID）






### 分析 下Sql语句，设计合适的索引并解释优化原理
SELEDCT USER_ID,ORDER_DATE, STATUS 
FROM ORDERS
WHERE USER_ID = 123 AND ORDER_DATE > '2024-01-01' AND STATUS IN( 'PAID','SHIPPED')
ORDER BY ORDER_DATE DESC
LIMIT 10;
### 索引分析
EXPLAIN 
SELECT USER_ID,ORDER_DATE, STATUS 
FROM ORDERS
WHERE USER_ID = 123 
  AND ORDER_DATE > '2024-01-01' 
  AND STATUS IN('PAID','SHIPPED')
ORDER BY ORDER_DATE DESC, ORDER_ID DESC
LIMIT 10;
### 索引设计
1. 主键索引：USER_ID
2. 普通索引：ORDER_DATE
3. 普通索引：STATUS
<!-- -- 预期结果：
-- type=range
-- key=idx_pagination
-- Extra=Using where; Using index -->
```SQL
CREATE TABLE ORDERS (
    ORDER_ID      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,  -- 订单主键
    USER_ID       INT UNSIGNED NOT NULL,                       -- 用户ID
    ORDER_DATE    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 下单时间
    STATUS        ENUM('PAID','SHIPPED','CANCELED') NOT NULL,  -- 订单状态
    TOTAL_AMOUNT  DECIMAL(15,2) NOT NULL,                      -- 订单金额
    ADDRESS_ID    INT UNSIGNED,                                -- 收货地址
    INDEX idx_user_status_date (USER_ID, STATUS, ORDER_DATE)   -- 复合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


ALTER TABLE ORDERS 
ADD INDEX idx_optimizer (USER_ID, STATUS, ORDER_DATE DESC);
```

### 问题排查与解决
生产环境中，golang服务出现了频繁的panic，如何定位并解决？
```go
package main
import (
	"fmt"
	"time"
    "runtime/debug"
    	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)
var logger = logrus.New()
func main() {
  defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v\n%s", err, debug.Stack())
		}
        if err := recover(); err !=nil {
            logger.WithFields(logrus.Fields{
                "error": err,
            }).Errorf("Panic occurred:%v", error)
        }
	}()
}
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v\n%s", err, debug.Stack())
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
```
#### 定位步骤
    1. 查看日志，定位panic的位置
    2. 查看堆栈信息，定位panic的原因
#### 解决步骤
    3. 解决panic，恢复服务
    4. 监控服务，确保服务恢复正常 
    5. 修复代码，确保代码的稳定性
    6. 测试服务，确保服务的稳定性
    7. 部署服务，确保服务的稳定性
    8. 监控服务，确保服务的稳定性
#### 预防机制
1. 代码审查：强制检查指针、并发操作、错误处理
2. 混沌工程：定期注入故障测试恢复能力
3. 压测验证：使用 vegeta 进行并发压测 
```go 
go
// 指针安全访问
if user != nil {
    fmt.Println(user.Name)
}

// 安全类型断言
if val, ok := obj.(string); ok {
    // 处理 string 类型
}

// 带缓冲的 channel 操作
select {
case ch <- data:
default:
    log.Println("channel 阻塞，启用降级策略")
}
```
高频 panic 场景排查表
|Panic 类型|	常见原因|	定位方法|	解决方案|
|---|---|---|---|
|nil pointer dereference|	未初始化的指针/接口	|检查所有指针初始化逻辑|	使用 if x != nil 防护
|index out of range|	数组/切片越界访问|	审查索引计算逻辑	|增加边界检查|
|concurrent map read/write|	并发读写未加锁的 map	|使用 -race 编译参数检测|	加 sync.RWMutex
|type assertion failure|	接口类型断言失败|	检查断言代码|	使用 val, ok := x.(T)|
|send on closed channel|	已关闭的 channel 被发送数据	|审查 channel 生命周期管理|	使用 sync.Once 关闭|
|runtime: out of memory|	内存泄漏或超大对象分配	|使用 pprof 分析内存分布|	优化数据结构/限制资源|

### 内存泄漏问题如何排查？
#### 一、内存泄漏初步确认
1. 监控进程内存（RSS持续增长）
```bash
# 监控进程内存（RSS持续增长）
$ watch -n 1 "ps -eo pid,rss,comm | grep myapp"

# 查看GC状态（关注下次GC阈值是否持续增长）
$ curl http://localhost:6060/debug/pprof/heap?debug=1
# 关注指标：
# next_gc: 下次GC阈值（若持续上涨则泄漏）
# heap_objects: 堆对象数（异常增长）
```
#### 二、定位泄漏源工具链
1. 生成内存profile
```bash
# 通过pprof获取堆内存快照
$ go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# 对比两次内存快照（间隔10分钟）
$ curl -o heap1.pprof http://localhost:6060/debug/pprof/heap
$ sleep 600
$ curl -o heap2.pprof http://localhost:6060/debug/pprof/heap
$ go tool pprof -base heap1.pprof heap2.pprof -http=:8080
```
三、常见泄漏场景排查表
|泄漏类型	|典型特征	|排查方法	|解决方案|
|---|---|---|---|
|对象泄漏	|对象未被释放	|pprof -alloc_objects	|使用sync.Pool/手动释放|
|Goroutine泄漏	|goroutine数持续增长	|debug/pprof/goroutine	|检查未退出的循环/阻塞操作|
|内存泄漏	|RSS持续增长	|debug/pprof/heap	|优化数据结构/限制资源|
|Channel泄漏	|channel阻塞	|debug/pprof/block	|检查无界channel|
|锁泄漏	|锁未释放	|debug/pprof/mutex	|检查未释放的锁|
|文件句柄泄漏	|文件句柄未关闭	|debug/pprof/fd	|检查未关闭的文件句柄|
|缓存未清理	|大对象持有引用	|pprof -alloc_space	|添加LRU淘汰策略|
|通道阻塞	|发送/接收阻塞导致对象滞留|	pprof -inuse_objects	|使用带缓冲通道+超时机制|
|全局变量累积	|map/slice未清除旧数据	|对比多次heap快照差异	|定期清理过期数据|
|CGO资源未释放	|C内存分配未配对释放	|使用Valgrind检测	|确保C.free调用|

2. 关键分析指标
inuse_space：当前内存使用量

alloc_space：历史累计分配量

inuse_objects：存活对象数

alloc_objects：总分配对象数
CPU使用率高如何定位？
请描述排查过程和用户的工具。