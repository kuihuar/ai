选项模式（Functional Options）
作用：优雅地构造复杂对象，支持默认值和可选参数。
```go
type Server struct {
    host string
    port int
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) { s.host = host }
}

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080} // 默认值
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// 使用
server := NewServer(WithHost("0.0.0.0"), WithPort(80))

```