package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}

	defer conn.Close()

	message := "Hello server!"
	_, err = conn.Write([]byte(message))

	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}
	fmt.Println("Sent message to the server:", message)
}