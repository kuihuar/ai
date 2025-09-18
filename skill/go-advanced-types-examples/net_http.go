package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// DemoHTTPNetworking: start a tuned HTTP server, call it with a tuned client, then graceful shutdown.
func DemoHTTPNetworking() {
	// Server with sensible timeouts
	srv := &http.Server{
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("ok"))
		}),
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	// Start serving
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Println("serve error:", err)
		}
	}()

	addr := "http://" + ln.Addr().String()

	// Tuned HTTP client Transport
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	resp, err := client.Get(addr)
	if err != nil {
		log.Println("client get error:", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		fmt.Println("HTTP resp:", string(body))
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
