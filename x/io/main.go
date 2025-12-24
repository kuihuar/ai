package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading data:", err)
			break
		}
		fmt.Println("Received data:", string(buffer[:n]))

		_, err = conn.Write([]byte("Received data: " + string(buffer[:n]) + "\n"))

		if err != nil {
			fmt.Println("Error writing data:", err)
			break
		}
	}
}
func main() {

	listenAddr := "localhost:8080"

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Listening error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on", listenAddr)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Accepting error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
