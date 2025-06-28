package main

import (
	"io"
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

	clientLog.Println("Connection with proxy server established.")

	// start go routine: stdin -> server
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			clientLog.Printf("Error copying stdin to server: %v", err)
		}
	}()

	// start go routine: server -> stdout
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					clientLog.Printf("Error reading from server: %v", err)
				}
				break
			}
			data := buffer[:n]
			clientLog.Printf("Received from server: %s", strings.TrimSpace(string(data)))
			os.Stdout.Write(data)
		}
	}()

	select {}
}