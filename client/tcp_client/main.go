package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	logFile, err := os.OpenFile("client.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open client log file: %v", err)
	}
	defer logFile.Close()

	clientLog := log.New(logFile, "[CLIENT] ", log.LstdFlags)

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		clientLog.Println("Error connecting to the server:", err)
		return
	}

	defer conn.Close()

	clientLog.Println("Connection with server established.")

	message := "Hello server!"
	_, err = conn.Write([]byte(message))

	if err != nil {
		clientLog.Println("Error writing to connection:", err)
		return
	}
	clientLog.Println("Sent message to the server:", message)

	// Receive response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		clientLog.Println("Error reading response from server:", err)
		return
	}

	response := strings.TrimSpace(string(buffer[:n]))
	clientLog.Println("Received response from server:", response)
	fmt.Println("Message sent successfully!!")
}