package main

import (
	"bufio"
	"fmt"
	"io"
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

func handleConnection(clientConn net.Conn, serverLog *log.Logger) {
	defer clientConn.Close()
	serverLog.Println("New connection accepted")

	reader := bufio.NewReader(clientConn)
	var requestBuilder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			serverLog.Println("Error reading from client:", err)
			return
		}
		requestBuilder.WriteString(line)
		if line == "\r\n" || line == "\n" {
			break
		}
	}

	requestString := requestBuilder.String()
	serverLog.Println("Received HTTP request from client:", requestString)

	// buffer := make([]byte, 8192)

	// n, err := clientConn.Read(buffer)
	// if err != nil {
	// 	serverLog.Println("Error reading from client:", err)
	// 	return
	// }

	// requestData := buffer[:n]
	// requestString := string(requestData)
	// serverLog.Println("Received request from client:\n", requestString)

	host, err := extractHost(requestString)
	if err != nil {
		serverLog.Println("Failed to parse Host header:", err)
		return
	}

	destinationAddr := host + ":80"
	serverLog.Printf("Forwarding request to destination: %s", destinationAddr)

	destConn, err := net.Dial("tcp", destinationAddr)
	if err != nil {
		serverLog.Println("Error connecting to destination:", err)
		return
	}
	defer destConn.Close()

	_, err = destConn.Write([]byte(requestString))
	if err != nil {
		serverLog.Println("Error forwarding request to destination:", err)
		return
	}

	go func() {
		_, err := io.Copy(destConn, clientConn)
		if err != nil {
			serverLog.Printf("Error forwarding client -> destination: %v", err)
		}
	}()

	_, err = io.Copy(clientConn, destConn)
	if err != nil {
		serverLog.Printf("Error forwarding destination -> client: %v", err)
	}
	serverLog.Println("Connection closed.")
}

func extractHost(request string) (string, error) {
	lines := strings.Split(request, "\n")

	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "host:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("Host header not found")
}