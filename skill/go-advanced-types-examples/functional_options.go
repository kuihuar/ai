package main

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	EnableMetrics  bool
	MaxConnections int
}

type Option func(*ServerConfig)

func WithAddr(addr string) Option             { return func(c *ServerConfig) { c.Addr = addr } }
func WithReadTimeout(d time.Duration) Option  { return func(c *ServerConfig) { c.ReadTimeout = d } }
func WithWriteTimeout(d time.Duration) Option { return func(c *ServerConfig) { c.WriteTimeout = d } }
func WithMetrics(enabled bool) Option         { return func(c *ServerConfig) { c.EnableMetrics = enabled } }
func WithMaxConnections(n int) Option         { return func(c *ServerConfig) { c.MaxConnections = n } }

func NewServerConfig(opts ...Option) *ServerConfig {
	cfg := &ServerConfig{
		Addr:           ":8080",
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		EnableMetrics:  false,
		MaxConnections: 100,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func DemoFunctionalOptions() {
	cfg := NewServerConfig(
		WithAddr(":9090"),
		WithReadTimeout(5*time.Second),
		WithMetrics(true),
		WithMaxConnections(1000),
	)
	fmt.Printf("cfg: %+v\n", *cfg)
}
