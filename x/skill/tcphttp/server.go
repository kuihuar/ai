package tcphttp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func Server() {

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accepting error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	//buffer := make([]byte, 1024)

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ReadString error:", err)
			continue
		}
		fmt.Println("Received message:", message)

		_, err = conn.Write([]byte("Received message: " + message + "\n"))
		if err != nil {
			fmt.Println("Write error:", err)
			continue
		}
	}
}

func HTTPServer() {

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	router.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
}
