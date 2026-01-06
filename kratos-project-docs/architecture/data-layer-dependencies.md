# 数据访问层外部依赖管理

## 概述

在 Kratos 分层架构中，Data 层负责所有外部数据源的访问，包括：
- **数据库**：MySQL、PostgreSQL、MongoDB 等
- **缓存**：Redis、Memcached 等
- **消息队列**：Kafka、RabbitMQ、RocketMQ 等
- **外部服务**：gRPC、HTTP API 等

本文档介绍如何在 Data 层组织和管理这些外部依赖。

## 依赖分类和组织

### 目录结构

```
internal/data/
├── data.go              # Data 结构体和初始化
├── greeter.go           # 业务相关的 Repository 实现
├── clients/             # 外部客户端封装
│   ├── redis.go        # Redis 客户端
│   ├── kafka.go        # Kafka 客户端
│   ├── mq.go           # 消息队列客户端（RabbitMQ/RocketMQ）
│   ├── grpc.go         # gRPC 客户端
│   └── http.go         # HTTP 客户端
├── cache/               # 缓存相关实现
│   └── cache.go
└── queue/               # 消息队列相关实现
    └── producer.go
```

### 依赖分类

#### 1. 存储类依赖
- **数据库**：MySQL、PostgreSQL、MongoDB
- **缓存**：Redis、Memcached

#### 2. 消息类依赖
- **消息队列**：Kafka、RabbitMQ、RocketMQ、NATS

#### 3. 服务类依赖
- **gRPC 服务**：其他微服务的 gRPC 接口
- **HTTP 服务**：RESTful API、第三方服务

## Data 结构体设计

### 统一管理所有依赖

```go
// internal/data/data.go
package data

import (
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewRedisClient,
	NewKafkaProducer,
	NewMQClient,
	NewGRPCClients,
	NewHTTPClients,
	NewGreeterRepo,
)

// Data 统一管理所有数据访问依赖
type Data struct {
	db          *gorm.DB           // 数据库
	redis       *redis.Client      // Redis 客户端
	kafkaWriter *kafka.Writer      // Kafka 生产者
	mqClient    MQClient           // 消息队列客户端（接口）
	grpcClients map[string]*grpc.ClientConn  // gRPC 客户端池
	httpClients map[string]HTTPClient        // HTTP 客户端池
	logger      log.Logger
}

// NewData 初始化所有数据访问依赖
func NewData(
	c *conf.Data,
	db *gorm.DB,
	redisClient *redis.Client,
	kafkaWriter *kafka.Writer,
	mqClient MQClient,
	grpcClients map[string]*grpc.ClientConn,
	httpClients map[string]HTTPClient,
	logger log.Logger,
) (*Data, func(), error) {
	d := &Data{
		db:          db,
		redis:       redisClient,
		kafkaWriter: kafkaWriter,
		mqClient:    mqClient,
		grpcClients: grpcClients,
		httpClients: httpClients,
		logger:      logger,
	}
	
	cleanup := func() {
		log.NewHelper(logger).Info("closing data resources")
		
		// 关闭数据库连接
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		
		// 关闭 Redis 连接
		if redisClient != nil {
			redisClient.Close()
		}
		
		// 关闭 Kafka Writer
		if kafkaWriter != nil {
			kafkaWriter.Close()
		}
		
		// 关闭 MQ 连接
		if mqClient != nil {
			mqClient.Close()
		}
		
		// 关闭所有 gRPC 连接
		for name, conn := range grpcClients {
			if conn != nil {
				conn.Close()
				log.NewHelper(logger).Infof("closed gRPC connection: %s", name)
			}
		}
		
		// 关闭所有 HTTP 客户端
		for name, client := range httpClients {
			if client != nil {
				client.Close()
				log.NewHelper(logger).Infof("closed HTTP client: %s", name)
			}
		}
	}
	
	return d, cleanup, nil
}
```

## 客户端封装

### 1. Redis 客户端

```go
// internal/data/clients/redis.go
package data

import (
	"context"
	"time"
	
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	if c.Redis == nil {
		return nil, nil
	}
	
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Network:      c.Redis.Network,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
	})
	
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	
	log.NewHelper(logger).Info("Redis client initialized")
	return rdb, nil
}
```

### 2. Kafka 客户端

```go
// internal/data/clients/kafka.go
package data

import (
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
)

// NewKafkaProducer 创建 Kafka 生产者
func NewKafkaProducer(c *conf.Data, logger log.Logger) (*kafka.Writer, error) {
	if c.Kafka == nil {
		return nil, nil
	}
	
	writer := &kafka.Writer{
		Addr:     kafka.TCP(c.Kafka.Brokers...),
		Balancer: &kafka.LeastBytes{},
		Async:    true,
		BatchSize: 100,
		BatchTimeout: 10 * time.Millisecond,
	}
	
	log.NewHelper(logger).Info("Kafka producer initialized")
	return writer, nil
}

// NewKafkaConsumer 创建 Kafka 消费者（如果需要）
func NewKafkaConsumer(c *conf.Data, logger log.Logger) (*kafka.Reader, error) {
	if c.Kafka == nil {
		return nil, nil
	}
	
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Kafka.Brokers,
		Topic:    c.Kafka.Topic,
		GroupID:  c.Kafka.GroupId,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	
	log.NewHelper(logger).Info("Kafka consumer initialized")
	return reader, nil
}
```

### 3. 消息队列客户端（接口抽象）

```go
// internal/data/clients/mq.go
package data

import (
	"context"
	
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
)

// MQClient 消息队列客户端接口
type MQClient interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	Close() error
}

// RabbitMQClient RabbitMQ 实现
type RabbitMQClient struct {
	// RabbitMQ 连接
	logger log.Logger
}

func NewRabbitMQClient(c *conf.Data, logger log.Logger) (MQClient, error) {
	// 实现 RabbitMQ 客户端
	return &RabbitMQClient{logger: logger}, nil
}

func (c *RabbitMQClient) Publish(ctx context.Context, topic string, message []byte) error {
	// 实现发布逻辑
	return nil
}

func (c *RabbitMQClient) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	// 实现订阅逻辑
	return nil
}

func (c *RabbitMQClient) Close() error {
	// 实现关闭逻辑
	return nil
}

// RocketMQClient RocketMQ 实现
type RocketMQClient struct {
	// RocketMQ 连接
	logger log.Logger
}

func NewRocketMQClient(c *conf.Data, logger log.Logger) (MQClient, error) {
	// 实现 RocketMQ 客户端
	return &RocketMQClient{logger: logger}, nil
}

// NewMQClient 根据配置创建对应的 MQ 客户端
func NewMQClient(c *conf.Data, logger log.Logger) (MQClient, error) {
	if c.Mq == nil {
		return nil, nil
	}
	
	switch c.Mq.Type {
	case "rabbitmq":
		return NewRabbitMQClient(c, logger)
	case "rocketmq":
		return NewRocketMQClient(c, logger)
	default:
		return nil, fmt.Errorf("unsupported MQ type: %s", c.Mq.Type)
	}
}
```

### 4. gRPC 客户端池

```go
// internal/data/clients/grpc.go
package data

import (
	"context"
	"time"
	
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGRPCClients 创建 gRPC 客户端池
func NewGRPCClients(c *conf.Data, logger log.Logger) (map[string]*grpc.ClientConn, error) {
	if c.Grpc == nil || len(c.Grpc.Clients) == 0 {
		return make(map[string]*grpc.ClientConn), nil
	}
	
	clients := make(map[string]*grpc.ClientConn)
	
	for name, endpoint := range c.Grpc.Clients {
		conn, err := grpc.Dial(
			endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(5*time.Second),
			grpc.WithBlock(),
		)
		if err != nil {
			log.NewHelper(logger).Errorf("failed to connect to gRPC service %s: %v", name, err)
			// 继续初始化其他客户端，不中断
			continue
		}
		
		clients[name] = conn
		log.NewHelper(logger).Infof("gRPC client initialized: %s -> %s", name, endpoint)
	}
	
	return clients, nil
}

// GetGRPCClient 获取指定的 gRPC 客户端
func (d *Data) GetGRPCClient(name string) (*grpc.ClientConn, error) {
	conn, ok := d.grpcClients[name]
	if !ok {
		return nil, fmt.Errorf("gRPC client not found: %s", name)
	}
	return conn, nil
}
```

### 5. HTTP 客户端池

```go
// internal/data/clients/http.go
package data

import (
	"context"
	"io"
	"net/http"
	"time"
	
	"sre/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPClient HTTP 客户端接口
type HTTPClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	Close() error
}

// KratosHTTPClient Kratos HTTP 客户端实现
type KratosHTTPClient struct {
	client *http.Client
	logger log.Logger
}

func NewKratosHTTPClient(endpoint string, logger log.Logger) (HTTPClient, error) {
	client, err := http.NewClient(
		context.Background(),
		http.WithEndpoint(endpoint),
		http.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}
	
	return &KratosHTTPClient{
		client: client,
		logger: logger,
	}, nil
}

func (c *KratosHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	// 实现 GET 请求
	return nil, nil
}

func (c *KratosHTTPClient) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	// 实现 POST 请求
	return nil, nil
}

func (c *KratosHTTPClient) Close() error {
	// 实现关闭逻辑
	return nil
}

// NewHTTPClients 创建 HTTP 客户端池
func NewHTTPClients(c *conf.Data, logger log.Logger) (map[string]HTTPClient, error) {
	if c.Http == nil || len(c.Http.Clients) == 0 {
		return make(map[string]HTTPClient), nil
	}
	
	clients := make(map[string]HTTPClient)
	
	for name, endpoint := range c.Http.Clients {
		client, err := NewKratosHTTPClient(endpoint, logger)
		if err != nil {
			log.NewHelper(logger).Errorf("failed to create HTTP client %s: %v", name, err)
			continue
		}
		
		clients[name] = client
		log.NewHelper(logger).Infof("HTTP client initialized: %s -> %s", name, endpoint)
	}
	
	return clients, nil
}

// GetHTTPClient 获取指定的 HTTP 客户端
func (d *Data) GetHTTPClient(name string) (HTTPClient, error) {
	client, ok := d.httpClients[name]
	if !ok {
		return nil, fmt.Errorf("HTTP client not found: %s", name)
	}
	return client, nil
}
```

## 配置管理

### 扩展配置定义

```protobuf
// internal/conf/conf.proto
message Data {
  message Database {
    string driver = 1;
    string source = 2;
  }
  message Redis {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration read_timeout = 3;
    google.protobuf.Duration write_timeout = 4;
  }
  message Kafka {
    repeated string brokers = 1;
    string topic = 2;
    string group_id = 3;
  }
  message MQ {
    string type = 1;  // rabbitmq, rocketmq
    string endpoint = 2;
    string username = 3;
    string password = 4;
  }
  message GRPC {
    map<string, string> clients = 1;  // name -> endpoint
  }
  message HTTP {
    map<string, string> clients = 1;  // name -> endpoint
  }
  
  Database database = 1;
  Redis redis = 2;
  Kafka kafka = 3;
  MQ mq = 4;
  GRPC grpc = 5;
  HTTP http = 6;
}
```

### 配置文件示例

```yaml
# configs/config.yaml
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/test?parseTime=True&loc=Local
  redis:
    network: tcp
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
  kafka:
    brokers:
      - 127.0.0.1:9092
    topic: user-events
    group_id: user-service
  mq:
    type: rabbitmq
    endpoint: amqp://guest:guest@localhost:5672/
  grpc:
    clients:
      user-service: 127.0.0.1:9001
      order-service: 127.0.0.1:9002
  http:
    clients:
      payment-api: https://api.payment.com
      notification-api: https://api.notification.com
```

## Repository 使用示例

### 在 Repository 中使用多种依赖

```go
// internal/data/greeter.go
package data

import (
	"context"
	"encoding/json"
	
	"sre/internal/biz"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

type greeterRepo struct {
	data   *Data
	log    *log.Helper
}

func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	// 1. 保存到数据库
	if err := r.data.db.Create(g).Error; err != nil {
		return nil, err
	}
	
	// 2. 更新缓存
	key := fmt.Sprintf("greeter:%d", g.Id)
	if r.data.redis != nil {
		data, _ := json.Marshal(g)
		r.data.redis.Set(ctx, key, data, time.Hour)
	}
	
	// 3. 发送消息到 Kafka
	if r.data.kafkaWriter != nil {
		message, _ := json.Marshal(map[string]interface{}{
			"event": "greeter.created",
			"data":  g,
		})
		r.data.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Topic: "greeter-events",
			Value: message,
		})
	}
	
	// 4. 调用外部 gRPC 服务
	if conn, err := r.data.GetGRPCClient("notification-service"); err == nil {
		// 调用通知服务
		// client := v1.NewNotificationServiceClient(conn)
		// client.SendNotification(ctx, &v1.SendNotificationRequest{...})
	}
	
	return g, nil
}

func (r *greeterRepo) FindByID(ctx context.Context, id int64) (*biz.Greeter, error) {
	// 1. 先查缓存
	key := fmt.Sprintf("greeter:%d", id)
	if r.data.redis != nil {
		val, err := r.data.redis.Get(ctx, key).Result()
		if err == nil {
			var g biz.Greeter
			if json.Unmarshal([]byte(val), &g) == nil {
				return &g, nil
			}
		}
	}
	
	// 2. 查数据库
	var g biz.Greeter
	if err := r.data.db.Where("id = ?", id).First(&g).Error; err != nil {
		return nil, err
	}
	
	// 3. 回写缓存
	if r.data.redis != nil {
		data, _ := json.Marshal(g)
		r.data.redis.Set(ctx, key, data, time.Hour)
	}
	
	return &g, nil
}
```

## 最佳实践

### 1. 依赖注入顺序

使用 Wire 时，确保依赖按正确顺序初始化：

```go
// internal/data/data.go
var ProviderSet = wire.NewSet(
	// 先初始化客户端
	NewRedisClient,
	NewKafkaProducer,
	NewMQClient,
	NewGRPCClients,
	NewHTTPClients,
	
	// 再初始化 Data
	NewData,
	
	// 最后初始化 Repository
	NewGreeterRepo,
)
```

### 2. 优雅关闭

确保所有资源都能正确关闭：

```go
func (d *Data) Close() error {
	var errs []error
	
	// 按依赖顺序关闭
	if d.redis != nil {
		if err := d.redis.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	
	if d.kafkaWriter != nil {
		if err := d.kafkaWriter.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	
	// ... 关闭其他资源
	
	if len(errs) > 0 {
		return fmt.Errorf("errors closing resources: %v", errs)
	}
	return nil
}
```

### 3. 连接池管理

合理配置连接池参数：

```go
// Redis 连接池
redis.NewClient(&redis.Options{
	PoolSize:     10,      // 连接池大小
	MinIdleConns: 5,       // 最小空闲连接
	MaxRetries:   3,       // 最大重试次数
})

// gRPC 连接池
grpc.WithKeepaliveParams(keepalive.ClientParameters{
	Time:                10 * time.Second,
	Timeout:             3 * time.Second,
	PermitWithoutStream: true,
})
```

### 4. 错误处理和重试

实现统一的错误处理和重试机制：

```go
// internal/data/utils/retry.go
func Retry(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return lastErr
}
```

### 5. 健康检查

为每个依赖实现健康检查：

```go
// internal/data/health.go
func (d *Data) HealthCheck(ctx context.Context) map[string]error {
	health := make(map[string]error)
	
	// Redis 健康检查
	if d.redis != nil {
		health["redis"] = d.redis.Ping(ctx).Err()
	}
	
	// 数据库健康检查
	if d.db != nil {
		if sqlDB, err := d.db.DB(); err == nil {
			health["database"] = sqlDB.Ping()
		}
	}
	
	// gRPC 健康检查
	for name, conn := range d.grpcClients {
		health[fmt.Sprintf("grpc.%s", name)] = conn.WaitForReady(ctx)
	}
	
	return health
}
```

### 6. 监控和日志

为每个依赖添加监控和日志：

```go
func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	start := time.Now()
	defer func() {
		r.log.Infof("Save greeter took %v", time.Since(start))
	}()
	
	// 实现逻辑
}
```

## 测试支持

### Mock 客户端

为测试创建 Mock 客户端：

```go
// internal/data/clients/mock.go
type MockRedisClient struct {
	data map[string]string
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	// Mock 实现
}

// 在测试中使用
func TestGreeterRepo_Save(t *testing.T) {
	mockRedis := &MockRedisClient{}
	data := &Data{
		redis: mockRedis,
		// ...
	}
	repo := NewGreeterRepo(data, log.NewStdLogger(os.Stdout))
	// 测试逻辑
}
```

## 总结

1. **统一管理**：在 `Data` 结构体中统一管理所有外部依赖
2. **接口抽象**：使用接口抽象不同类型的客户端（如 MQClient）
3. **配置驱动**：通过配置文件管理所有依赖的连接信息
4. **优雅关闭**：确保所有资源都能正确关闭
5. **错误处理**：实现统一的错误处理和重试机制
6. **健康检查**：为每个依赖实现健康检查
7. **测试支持**：提供 Mock 客户端便于单元测试

通过以上方式，可以清晰地组织和管理 Data 层的所有外部依赖，保持代码的可维护性和可测试性。

