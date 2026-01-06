# Express 风格详解

## 一、什么是 Express？

### 1.1 Express.js 简介

**Express.js** 是 Node.js 最流行的 Web 框架，以其**简洁优雅的 API** 而闻名。

```javascript
// Express.js 示例
const express = require('express');
const app = express();

// 路由定义
app.get('/users/:id', (req, res) => {
    res.json({ id: req.params.id });
});

// 中间件
app.use((req, res, next) => {
    console.log('Request:', req.method, req.path);
    next();
});

// 启动服务器
app.listen(3000);
```

### 1.2 Express 风格的特点

1. **简洁的 API**: 方法名直观（get, post, use）
2. **链式调用**: 可以链式调用多个方法
3. **中间件模式**: 强大的中间件系统
4. **路由分组**: 支持路由分组和嵌套

## 二、Fiber 的 Express 风格

### 2.1 API 对比

#### Express.js

```javascript
// Express.js
const express = require('express');
const app = express();

app.get('/users/:id', (req, res) => {
    res.json({ id: req.params.id });
});

app.post('/users', (req, res) => {
    res.json({ message: 'created' });
});

app.use('/api', middleware);
```

#### Fiber（Express 风格）

```go
// Fiber（Express 风格）
app := fiber.New()

app.Get("/users/:id", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"id": c.Params("id")})
})

app.Post("/users", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "created"})
})

app.Use("/api", middleware)
```

**相似度**: 几乎 **100%**！

### 2.2 路由定义对比

#### Express.js

```javascript
// Express.js 路由定义
app.get('/users/:id', handler);
app.post('/users', handler);
app.put('/users/:id', handler);
app.delete('/users/:id', handler);

// 路由分组
const router = express.Router();
router.get('/:id', handler);
app.use('/users', router);
```

#### Fiber

```go
// Fiber 路由定义（Express 风格）
app.Get("/users/:id", handler)
app.Post("/users", handler)
app.Put("/users/:id", handler)
app.Delete("/users/:id", handler)

// 路由分组
group := app.Group("/users")
group.Get("/:id", handler)
```

**几乎完全一致**！

### 2.3 中间件对比

#### Express.js

```javascript
// Express.js 中间件
app.use((req, res, next) => {
    console.log('Request:', req.method, req.path);
    next();  // 继续下一个中间件
});

// 错误处理中间件
app.use((err, req, res, next) => {
    res.status(500).json({ error: err.message });
});
```

#### Fiber

```go
// Fiber 中间件（Express 风格）
app.Use(func(c fiber.Ctx) error {
    fmt.Println("Request:", c.Method(), c.Path())
    return c.Next()  // 继续下一个中间件
})

// 错误处理中间件
app.Use(func(c fiber.Ctx, err error) error {
    return c.Status(500).JSON(fiber.Map{"error": err.Error()})
})
```

**几乎完全一致**！

## 三、Express 风格的核心特性

### 3.1 简洁的方法名

#### Express.js

```javascript
app.get('/path', handler);
app.post('/path', handler);
app.put('/path', handler);
app.delete('/path', handler);
app.use('/path', middleware);
```

#### Fiber

```go
app.Get("/path", handler)
app.Post("/path", handler)
app.Put("/path", handler)
app.Delete("/path", handler)
app.Use("/path", middleware)
```

**特点**:
- ✅ 方法名直观（get, post, use）
- ✅ 符合 HTTP 方法命名
- ✅ 易于记忆和使用

### 3.2 链式调用

#### Express.js

```javascript
// Express.js 链式调用
app
    .use(middleware1)
    .use(middleware2)
    .get('/path', handler1, handler2);
```

#### Fiber

```go
// Fiber 链式调用（Express 风格）
app.
    Use(middleware1).
    Use(middleware2).
    Get("/path", handler1, handler2)
```

**特点**:
- ✅ 可以链式调用多个方法
- ✅ 代码更简洁
- ✅ 易于阅读

### 3.3 路由参数

#### Express.js

```javascript
// Express.js 路由参数
app.get('/users/:id', (req, res) => {
    const id = req.params.id;
    res.json({ id });
});

app.get('/users/:id/posts/:postId', (req, res) => {
    const { id, postId } = req.params;
    res.json({ id, postId });
});
```

#### Fiber

```go
// Fiber 路由参数（Express 风格）
app.Get("/users/:id", func(c fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"id": id})
})

app.Get("/users/:id/posts/:postId", func(c fiber.Ctx) error {
    id := c.Params("id")
    postId := c.Params("postId")
    return c.JSON(fiber.Map{"id": id, "postId": postId})
})
```

**几乎完全一致**！

### 3.4 路由分组

#### Express.js

```javascript
// Express.js 路由分组
const api = express.Router();

api.get('/users', getUsers);
api.post('/users', createUser);

app.use('/api/v1', api);
```

#### Fiber

```go
// Fiber 路由分组（Express 风格）
api := app.Group("/api/v1")

api.Get("/users", getUsers)
api.Post("/users", createUser)
```

**几乎完全一致**！

### 3.5 中间件执行

#### Express.js

```javascript
// Express.js 中间件执行
app.use((req, res, next) => {
    console.log('Before');
    next();  // 继续下一个中间件
    console.log('After');
});
```

#### Fiber

```go
// Fiber 中间件执行（Express 风格）
app.Use(func(c fiber.Ctx) error {
    fmt.Println("Before")
    err := c.Next()  // 继续下一个中间件
    fmt.Println("After")
    return err
})
```

**几乎完全一致**！

## 四、Express 风格的优势

### 4.1 降低学习成本

对于熟悉 Express.js 的开发者：

```javascript
// Express.js（JavaScript）
app.get('/users/:id', (req, res) => {
    res.json({ id: req.params.id });
});
```

```go
// Fiber（Go，Express 风格）
app.Get("/users/:id", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"id": c.Params("id")})
})
```

**学习成本**: 几乎为 **0**！

### 4.2 代码风格一致

```go
// Fiber 的 Express 风格让代码更一致
app.
    Use(logger.New()).
    Use(recovery.New()).
    Get("/", handler).
    Post("/users", createUser).
    Put("/users/:id", updateUser).
    Delete("/users/:id", deleteUser)
```

### 4.3 易于迁移

从 Express.js 迁移到 Fiber 非常简单：

```javascript
// Express.js
app.get('/users/:id', (req, res) => {
    res.json({ id: req.params.id });
});
```

```go
// Fiber（几乎可以直接翻译）
app.Get("/users/:id", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"id": c.Params("id")})
})
```

## 五、与其他框架的对比

### 5.1 Gin（非 Express 风格）

```go
// Gin（非 Express 风格）
r := gin.Default()

r.GET("/users/:id", func(c *gin.Context) {
    c.JSON(200, gin.H{"id": c.Param("id")})
})
```

**特点**:
- 使用 `*gin.Context` 而不是接口
- 方法名大写（GET 而不是 Get）
- 使用 `gin.H` 而不是 `fiber.Map`

### 5.2 GoFrame（非 Express 风格）

```go
// GoFrame（非 Express 风格）
s := g.Server()

s.Group("/", func(group *ghttp.RouterGroup) {
    group.GET("/users/:id", func(r *ghttp.Request) {
        r.Response.WriteJson(g.Map{"id": r.Get("id")})
    })
})
```

**特点**:
- 使用分组函数而不是链式调用
- 使用 `*ghttp.Request` 而不是接口
- API 风格不同

### 5.3 Fiber（Express 风格）

```go
// Fiber（Express 风格）
app := fiber.New()

app.Get("/users/:id", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"id": c.Params("id")})
})
```

**特点**:
- ✅ 与 Express.js 几乎完全一致
- ✅ 易于学习和使用
- ✅ 降低学习成本

## 六、Express 风格的实际应用

### 6.1 完整的 Express 风格示例

```go
package main

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
    app := fiber.New()
    
    // 全局中间件（Express 风格）
    app.Use(logger.New())
    app.Use(recover.New())
    
    // 路由定义（Express 风格）
    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
    
    // 路由分组（Express 风格）
    api := app.Group("/api/v1")
    
    api.Get("/users", getUsers)
    api.Get("/users/:id", getUser)
    api.Post("/users", createUser)
    api.Put("/users/:id", updateUser)
    api.Delete("/users/:id", deleteUser)
    
    // 启动服务器
    app.Listen(":3000")
}

// 处理器（Express 风格）
func getUsers(c fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "users": []fiber.Map{
            {"id": 1, "name": "Alice"},
            {"id": 2, "name": "Bob"},
        },
    })
}

func getUser(c fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"id": id, "name": "Alice"})
}
```

### 6.2 中间件使用（Express 风格）

```go
// 自定义中间件（Express 风格）
func authMiddleware(c fiber.Ctx) error {
    token := c.Get("Authorization")
    if token == "" {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }
    // 验证 token
    return c.Next()
}

// 使用中间件（Express 风格）
app.Use("/api", authMiddleware)
app.Get("/api/users", getUsers)
```

## 七、总结

### 7.1 Express 风格的核心特点

1. **简洁的 API**: 方法名直观（get, post, use）
2. **链式调用**: 可以链式调用多个方法
3. **中间件模式**: 强大的中间件系统
4. **路由分组**: 支持路由分组和嵌套
5. **参数获取**: 简洁的参数获取方式

### 7.2 Express 风格的优势

- ✅ **降低学习成本**: 熟悉 Express.js 的开发者可以快速上手
- ✅ **代码风格一致**: 代码更易读、易维护
- ✅ **易于迁移**: 从 Express.js 迁移到 Fiber 很简单
- ✅ **社区熟悉**: 大多数 Web 开发者都熟悉 Express 风格

### 7.3 Express 风格的适用场景

- ✅ **熟悉 Express.js**: 团队熟悉 Express.js
- ✅ **快速开发**: 需要快速开发 Web 应用
- ✅ **降低学习成本**: 希望降低 Go Web 开发的学习成本
- ✅ **代码一致性**: 希望代码风格与 Express.js 一致

Express 风格是 Fiber 的重要特色，它让 Go Web 开发变得更加简单和熟悉！

