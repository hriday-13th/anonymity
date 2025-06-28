package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open server log file: %v", err)
	}
	defer logFile.Close()

	serverLog := log.New(logFile, "[SERVER] ", log.LstdFlags)

	PORT := 8080
	address := fmt.Sprintf(":%d", PORT)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		serverLog.Println("Error creating listener: ", err)
		return
	}

	defer listener.Close()
	
	serverLog.Printf("Server is listening on port: %d", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			serverLog.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn, serverLog)
	}
}

func handleConnection(conn net.Conn, serverLog *log.Logger) {
	defer conn.Close()
	serverLog.Println("New connection accepted")

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		serverLog.Println("Error reading from connection:", err)
		return
	}

	received := strings.TrimSpace(string(buffer[:n]))
	serverLog.Println("Received data:", received)

	response := "Message Received!"
	_, err = conn.Write([]byte(response))

	if err != nil {
		serverLog.Println("Error writing to connection:", err)
		return
	}
	serverLog.Println("Sent message to client")
}