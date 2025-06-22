package main

import (
	"fmt"
	"net"
)

func main() {
	PORT := 8080
	address := fmt.Sprintf(":%d", PORT)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error creating listener: ", err)
		return
	}

	defer listener.Close()
	
	fmt.Printf("Server is listening on port: %d", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection accepted")

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Received data:", string(buffer))
}