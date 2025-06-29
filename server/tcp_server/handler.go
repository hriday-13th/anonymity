package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)


func handleConnection(clientConn net.Conn, serverLog *log.Logger) {
	defer clientConn.Close()
	serverLog.Println("New connection accepted")

	reader := bufio.NewReader(clientConn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		serverLog.Println("Error reading from request line:", err)
		return
	}
	serverLog.Println("Request line:", strings.TrimSpace(requestLine))

	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		serverLog.Println("Malformed request line")
		return
	}

	method := parts[0]
	target := parts[1]

	if strings.ToUpper(method) != "CONNECT" {
		serverLog.Printf("Blocking non-CONNECT request: %s %s", method, target)
		resp := "HTTP/1.1 403 Forbidden\r\nContent-Type: text/plain\r\n\r\nThis proxy only allows HTTPS over CONNECT tunnels.\r\n"
		clientConn.Write([]byte(resp))
		return
	}

	serverLog.Printf("CONNECT request for %s", target)

	serverConn, err := net.Dial("tcp", target)
	if err != nil {
		serverLog.Printf("Failed to connect to target %s: %v", target, err)
		resp := "HTTP/1.1 502 Bad Gateway\r\nContent-Type: text/plain\r\n\r\nFailed to connect to target.\r\n"
		clientConn.Write([]byte(resp))
		return
	}
	defer serverConn.Close()

	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		serverLog.Println("Error sending 200 response:", err)
		return
	}
	serverLog.Printf("Tunnel established to %s", target)

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(serverConn, clientConn)
		if err != nil {
			serverLog.Printf("Error copying from client -> server: %v", err)
		}
		done <- struct{}{}
	}()
	
	go func() {
		_, err = io.Copy(clientConn, serverConn)
		if err != nil && !isClosedConnError(err) {
			serverLog.Printf("Error copying from server -> client: %v", err)
		}
		done <- struct{}{}
	}()

	<-done

	serverLog.Printf("Closed tunnel to %s", target)
}

func isClosedConnError(err error) bool {
	if err == nil {
		return false
	}
	if ne, ok := err.(*net.OpError); ok && ne.Err.Error() == "use of closed network connection" {
		return true
	}
	return false
}