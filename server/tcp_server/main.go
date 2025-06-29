package main

import (
	"fmt"
	"log"
	"net"
	"os"
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